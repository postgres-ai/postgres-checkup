${CHECK_HOST_CMD} "${_PSQL} -f -" <<SQL
with table_scans as (
  select relid,
      tables.idx_scan + tables.seq_scan as all_scans,
      ( tables.n_tup_ins + tables.n_tup_upd + tables.n_tup_del ) as writes,
    pg_relation_size(relid) as table_size
      from pg_stat_user_tables as tables
), all_writes as (
  select sum(writes) as total_writes
  from table_scans
), indexes as (
  select
    idx_stat.relid,
    idx_stat.indexrelid,
    idx_stat.schemaname as schema_name,
    idx_stat.relname as table_name,
    idx_stat.indexrelname as index_name,
    ( case when lower(idx_stat.schemaname) <> idx_stat.schemaname then ('"' || idx_stat.schemaname || '"') else idx_stat.schemaname::text end ) as formated_schema_name,
    ( case when lower(idx_stat.indexrelname) <> idx_stat.indexrelname then ('"' || idx_stat.indexrelname || '"') else idx_stat.indexrelname::text end ) as formated_index_name,
    ( case when lower(idx_stat.relname) <> idx_stat.relname then ('"' || idx_stat.relname || '"') else idx_stat.relname::text end ) as formated_table_name,
    idx_stat.idx_scan,
    pg_relation_size(idx_stat.indexrelid) as index_bytes,
    indexdef ~* 'using btree' as idx_is_btree,
    pg_get_indexdef(pg_index.indexrelid) as index_def
  from pg_stat_user_indexes as idx_stat
      join pg_index
          using (indexrelid)
      join pg_indexes as indexes
          on idx_stat.schemaname = indexes.schemaname
              and idx_stat.relname = indexes.tablename
              and idx_stat.indexrelname = indexes.indexname
  where pg_index.indisunique = false
), index_ratios as (
  select
    indexrelid as index_id,
    schema_name,
    table_name,
    index_name,
    idx_scan,
    all_scans,
        round(( case when all_scans = 0 then 0.0::numeric
          else idx_scan::numeric/all_scans * 100 end),2) as index_scan_pct,
    writes,
    round((case when writes = 0 then idx_scan::numeric else idx_scan::numeric/writes end),2)
      as scans_per_write,
    index_bytes as index_size_bytes,
    table_size as table_size_bytes,
    idx_is_btree,
    index_def,
    formated_index_name,
    formated_schema_name,
    formated_table_name,
    coalesce(nullif(indexes.formated_schema_name, 'public') || '.', '') || indexes.formated_table_name as formated_object_name
  from indexes
  join table_scans
  using (relid)
),
-- Never used indexes
never_used_indexes as (
    select
      'Never Used Indexes' as reason,
      index_ratios.*
    from index_ratios
    where
      idx_scan = 0
      and idx_is_btree
    order by index_size_bytes desc
), never_used_indexes_num as (
  select row_number() over () num, nui.* from never_used_indexes nui
), never_used_indexes_total as (
    select
      sum(index_size_bytes) as index_size_bytes_sum,
      sum(table_size_bytes) as table_size_bytes_sum
    from never_used_indexes
), never_used_indexes_json as (
  select
    json_object_agg(nuin.formated_object_name || '.' || nuin.index_name, nuin) as json
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
  select row_number() over () num, rui.* from rarely_used_indexes rui
), rarely_used_indexes_total as (
    select
      sum(index_size_bytes) as index_size_bytes_sum,
      sum(table_size_bytes) as table_size_bytes_sum
    from rarely_used_indexes
), rarely_used_indexes_json as (
  select
    json_object_agg(ruin.formated_object_name || '.' || ruin.index_name, ruin) as json
  from rarely_used_indexes_num ruin
),
-- Redundant indexes
index_data as (
  select
    *,
    indkey::text as columns,
    array_to_string(indclass, ', ') as opclasses
  from pg_index
), original_redundant_indexes as (
  select
    i2.indexrelid as index_id,
    tnsp.nspname AS schema_name,
    trel.relname AS table_name,
    pg_relation_size(trel.oid) as table_size_bytes,
    irel.relname AS index_name,
    am1.amname as access_method,
    (i1.indexrelid::regclass)::text as reason,
    pg_get_indexdef(i1.indexrelid) main_index_def,
    pg_size_pretty(pg_relation_size(i1.indexrelid)) main_index_size,
    pg_get_indexdef(i2.indexrelid) index_def,
    pg_relation_size(i2.indexrelid) index_size_bytes,
    s.idx_scan as index_usage,
    ( case when lower(tnsp.nspname) <> tnsp.nspname then ('"' || tnsp.nspname || '"') else tnsp.nspname::text end ) as formated_schema_name,
    ( case when lower(irel.relname) <> irel.relname then ('"' || irel.relname || '"') else irel.relname::text end ) as formated_index_name,
    ( case when lower(trel.relname) <> trel.relname then ('"' || trel.relname || '"') else trel.relname::text end ) as formated_table_name
  from
    index_data as i1
    join index_data as i2 on (
        i1.indrelid = i2.indrelid -- same table
        and i1.indexrelid <> i2.indexrelid -- NOT same index
    )
    inner join pg_opclass op1 on i1.indclass[0] = op1.oid
    inner join pg_opclass op2 on i2.indclass[0] = op2.oid
    inner join pg_am am1 on op1.opcmethod = am1.oid
    inner join pg_am am2 on op2.opcmethod = am2.oid
    join pg_stat_user_indexes as s on s.indexrelid = i2.indexrelid
    join pg_class as trel on trel.oid = i2.indrelid
    join pg_namespace as tnsp on trel.relnamespace = tnsp.oid
    join pg_class as irel on irel.oid = i2.indexrelid
  where
    not i1.indisprimary -- index 1 is not primary
    and not ( -- skip if index1 is (primary or uniq) and is NOT (primary and uniq)
        (i1.indisprimary or i1.indisunique)
        and (not i2.indisprimary or not i2.indisunique)
    )
    and  am1.amname = am2.amname -- same access type
    and (
      i2.columns like (i1.columns || '%') -- index 2 includes all columns from index 1
      or i1.columns = i2.columns -- index1 and index 2 includes same columns
    )
    and (
      i2.opclasses like (i1.opclasses || '%')
      or i1.opclasses = i2.opclasses
    )
    -- index expressions is same
    and pg_get_expr(i1.indexprs, i1.indrelid) is not distinct from pg_get_expr(i2.indexprs, i2.indrelid)
    -- index predicates is same
    and pg_get_expr(i1.indpred, i1.indrelid) is not distinct from pg_get_expr(i2.indpred, i2.indrelid)
), redundant_indexes as (
    select
      ori.*,
      coalesce(nullif(formated_schema_name, 'public') || '.', '') || formated_index_name as formated_object_name
   from original_redundant_indexes ori
), redundant_indexes_grouped as (
  select
    index_id,
    schema_name,
    table_name,
    table_size_bytes,
    index_name,
    access_method,
    string_agg(reason, ', ') as reason,
    string_agg(main_index_def, ', ') as main_index_def,
    string_agg(main_index_size, ', ') as main_index_size,
    index_def,
    index_size_bytes,
    index_usage,
    formated_index_name,
    formated_schema_name,
    formated_table_name,
    formated_object_name
  from redundant_indexes
  group by
    index_id,
    table_size_bytes,
    schema_name,
    table_name,
    index_name,
    access_method,
    index_def,
    index_size_bytes,
    index_usage,
    formated_index_name,
    formated_schema_name,
    formated_table_name,
    formated_object_name
  order by index_size_bytes desc
), redundant_indexes_num as (
  select row_number() over () num, rig.* from redundant_indexes_grouped rig
),
redundant_indexes_json as (
  select
    json_object_agg(rin.formated_object_name, rin) as json
  from redundant_indexes_num rin
), redundant_indexes_total as (
    select
      sum(index_size_bytes) as index_size_bytes_sum,
      sum(table_size_bytes) as table_size_bytes_sum
    from redundant_indexes_grouped
),
-- Do and Undo code generation
together as (
  select
    reason::text,
    index_id::text,
    schema_name::text,
    table_name::text,
    index_name::text,
    pg_size_pretty(index_size_bytes)::text as index_size,
    index_def::text,
    null as main_index_def,
    formated_index_name::text,
    formated_schema_name::text,
    formated_table_name::text,
    formated_object_name::text
  from never_used_indexes
  union all
  select
    reason::text,
    index_id::text,
    schema_name::text,
    table_name::text,
    index_name::text,
    pg_size_pretty(index_size_bytes)::text as index_size,
    index_def::text,
    main_index_def::text,
    formated_index_name::text,
    formated_schema_name::text,
    formated_table_name::text,
    formated_object_name::text
  from redundant_indexes
), do_lines as (
  select format('DROP INDEX CONCURRENTLY %s; -- %s, %s, table %s', max(formated_index_name), max(index_size), string_agg(reason, ', '), table_name) as line
  from together t1
  group by table_name, index_name
  order by table_name, index_name
), undo_lines as (
  select
    replace(
      format('%s; -- table %s', max(index_def), table_name),
      'CREATE INDEX',
      'CREATE INDEX CONCURRENTLY'
    )as line
  from together t2
  group by table_name, index_name
  order by table_name, index_name
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
    'redundant_indexes',
    (select * from redundant_indexes_json),
    'redundant_indexes_total',
    (select row_to_json(rit) from redundant_indexes_total as rit),
    'do',
    (select json_agg(dl.line) from do_lines as dl),
    'undo',
    (select json_agg(ul.line) from undo_lines as ul),
    'database_stat',
    (select * from database_stat)
  );
SQL
