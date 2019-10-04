if [[ ! -z ${IS_LARGE_DB+x} ]] && [[ ${IS_LARGE_DB} == "1" ]]; then
  MIN_RELPAGES=100
else
  MIN_RELPAGES=0
fi

${CHECK_HOST_CMD} "${_PSQL} -f -" <<SQL
with fk_indexes as (
  select
    n.nspname as schema_name,
    ci.relname as index_name,
    cr.relname as table_name,
    (confrelid::regclass)::text as fk_table_ref,
    array_to_string(indclass, ', ') as opclasses
  from pg_index i
  join pg_class ci on ci.oid = i.indexrelid and ci.relkind = 'i'
  join pg_class cr on cr.oid = i.indrelid and cr.relkind = 'r'
  join pg_namespace n on n.oid = ci.relnamespace
  join pg_constraint cn on cn.conrelid = cr.oid
  left join pg_stat_user_indexes si on si.indexrelid = i.indexrelid
  where
     contype = 'f'
     and i.indisunique is false
     and conkey is not null
     and ci.relpages > ${MIN_RELPAGES}
     and si.idx_scan < 10
), table_scans as (
  select relid,
      tables.idx_scan + tables.seq_scan as all_scans,
      ( tables.n_tup_ins + tables.n_tup_upd + tables.n_tup_del ) as writes,
    pg_relation_size(relid) as table_size
      from pg_stat_user_tables as tables
      join pg_class c on c.oid = relid
      where c.relpages > ${MIN_RELPAGES}
), all_writes as (
  select sum(writes) as total_writes
  from table_scans
), indexes as (
  select
    i.indrelid,
    i.indexrelid,
    n.nspname as schema_name,
    cr.relname as table_name,
    ci.relname as index_name,
    quote_ident(n.nspname) as formated_schema_name,
    quote_ident(ci.relname) as formated_index_name,
    quote_ident(cr.relname) as formated_table_name,
    coalesce(nullif(quote_ident(n.nspname), 'public') || '.', '') || quote_ident(cr.relname) as formated_relation_name,
    si.idx_scan,
    pg_relation_size(i.indexrelid) as index_bytes,
    ci.relpages,
    (case when a.amname = 'btree' then true else false end) as idx_is_btree,
    pg_get_indexdef(i.indexrelid) as index_def,
    array_to_string(i.indclass, ', ') as opclasses
  from pg_index i
     join pg_class ci on ci.oid = i.indexrelid and ci.relkind = 'i'
     join pg_class cr on cr.oid = i.indrelid and cr.relkind = 'r'
     join pg_namespace n on n.oid = ci.relnamespace
     join pg_am a ON ci.relam = a.oid
     left join pg_stat_user_indexes si on si.indexrelid = i.indexrelid
  where
    i.indisunique = false
    and i.indisvalid = true
    and ci.relpages > ${MIN_RELPAGES}
), index_ratios as (
  select
    i.indexrelid as index_id,
    i.schema_name,
    i.table_name,
    i.index_name,
    idx_scan,
    all_scans,
        round(( case when all_scans = 0 then 0.0::numeric
          else idx_scan::numeric/all_scans * 100 end), 2) as index_scan_pct,
    writes,
    round((case when writes = 0 then idx_scan::numeric else idx_scan::numeric/writes end), 2)
      as scans_per_write,
    index_bytes as index_size_bytes,
    table_size as table_size_bytes,
    i.relpages,
    idx_is_btree,
    index_def,
    formated_index_name,
    formated_schema_name,
    formated_table_name,
    formated_relation_name,
    i.opclasses,
    case when fi.index_name is not null then true else false end as supports_fk
  from indexes i
  left join fk_indexes fi on
    fi.fk_table_ref = i.table_name
    and fi.schema_name = i.schema_name 
    and fi.opclasses like (i.opclasses || '%')
  join table_scans ts on ts.relid = i.indrelid
),
-- Never used indexes
never_used_indexes as (
  select
    'Never Used Indexes' as reason,
    ir.*
  from index_ratios ir
  where
    idx_scan = 0
    and idx_is_btree
  order by index_size_bytes desc
), never_used_indexes_num as (
  select row_number() over () num, nui.*
  from never_used_indexes nui
), never_used_indexes_total as (
  select
    sum(index_size_bytes) as index_size_bytes_sum,
    sum(table_size_bytes) as table_size_bytes_sum
  from never_used_indexes

), never_used_indexes_json as (
  select
    json_object_agg(coalesce(nuin.schema_name, 'public') || '.' || nuin.index_name, nuin) as json
  from never_used_indexes_num nuin
),
-- Rarely used indexes
rarely_used_indexes as (
  select
    'Low Scans, High Writes' as reason,
    *,
    1 as grp
  from index_ratios
  where
      scans_per_write <= 1
      and index_scan_pct < 10
      and idx_scan > 0
      and writes > 100
      and idx_is_btree
  union all
  select
    'Seldom Used Large Indexes' as reason,
    *,
    2 as grp
  from index_ratios
  where
      index_scan_pct < 5
      and scans_per_write > 1
      and idx_scan > 0
      and idx_is_btree
      and index_size_bytes > 100000000
  union all
  select
    'High-Write Large Non-Btree' as reason,
    index_ratios.*,
    3 as grp
  from index_ratios, all_writes
  where
      ( writes::numeric / ( total_writes + 1 ) ) > 0.02
      and not idx_is_btree
      and index_size_bytes > 100000000
  order by grp, index_size_bytes desc
), rarely_used_indexes_num as (
  select row_number() over () num, rui.*
  from rarely_used_indexes rui
), rarely_used_indexes_total as (
  select
    sum(index_size_bytes) as index_size_bytes_sum,
    sum(table_size_bytes) as table_size_bytes_sum
  from rarely_used_indexes
), rarely_used_indexes_json as (
  select
    json_object_agg(coalesce(ruin.schema_name, 'public') || '.' || ruin.index_name, ruin) as json
  from rarely_used_indexes_num ruin
), database_stat as (
  select
    row_to_json(dbstat)
  from (
    select
      sd.stats_reset::timestamptz(0),
      age(
        date_trunc('minute',now()),
        date_trunc('minute',sd.stats_reset)
      ) as stats_age,
      ((extract(epoch from now()) - extract(epoch from sd.stats_reset))/86400)::int as days
    from pg_stat_database sd
    where datname = current_database()
  ) dbstat
)
-- summarize data
select
  json_build_object(
    'never_used_indexes',
    (select * from never_used_indexes_json),
    'never_used_indexes_total',
    (select row_to_json(nuit) from never_used_indexes_total as nuit),
    'rarely_used_indexes',
    (select * from rarely_used_indexes_json),
    'rarely_used_indexes_total',
    (select row_to_json(ruit) from rarely_used_indexes_total as ruit),
    'database_stat',
    (select * from database_stat),
    'min_index_size_bytes',
    (select ${MIN_RELPAGES} * 8192)
  );
SQL
