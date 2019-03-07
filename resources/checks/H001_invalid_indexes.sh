# Invalid keys
${CHECK_HOST_CMD} "${_PSQL} -f -" <<SQL
with fk_indexes as (
  select
    schemaname as schema_name,
    (indexrelid::regclass)::text as index_name,
    (relid::regclass)::text as table_name,
    (confrelid::regclass)::text as fk_table_ref,
    array_to_string(indclass, ', ') as opclasses
  from
    pg_stat_user_indexes
  join pg_index using (indexrelid)
  left join pg_constraint
    on array_to_string(indkey, ',') = array_to_string(conkey, ',')
      and schemaname = (connamespace::regnamespace)::text
      and conrelid = relid
      and contype = 'f'
  where idx_scan = 0
     and indisunique is false
     and conkey is not null --conkey is not null then true else false end as is_fk_idx
), data as (
  select
    pci.relname as index_name,
    pn.nspname as schema_name,
    pct.relname as table_name,
    quote_ident(pn.nspname) as formated_schema_name,
    quote_ident(pci.relname) as formated_index_name,
    quote_ident(pct.relname) as formated_table_name,
    coalesce(nullif(quote_ident(pn.nspname), 'public') || '.', '') || quote_ident(pct.relname) as formated_relation_name,
    pg_relation_size(pidx.indexrelid) index_size_bytes,
    format(
      'DROP INDEX CONCURRENTLY %s; -- %s, table %s',
      quote_ident(pci.relname),--pidx.indexrelid::regclass::text,
      'Invalid index',
      pct.relname) as drop_code,
    replace(
      format('%s; -- table %s', pg_get_indexdef(pidx.indexrelid), pct.relname),
      'CREATE INDEX',
      'CREATE INDEX CONCURRENTLY'
    ) as revert_code,
    case when fi.index_name is not null then true else false end as supports_fk
  from pg_index pidx
  join pg_class as pci on pci.oid = pidx.indexrelid
  join pg_class as pct on pct.oid = pidx.indrelid
  left join pg_namespace pn on pn.oid = pct.relnamespace
  left join fk_indexes fi on
    fi.fk_table_ref = pct.relname
    and fi.opclasses like (array_to_string(pidx.indclass, ', ') || '%')
  where pidx.indisvalid = false
), data_total as (
    select
      sum(index_size_bytes) as index_size_bytes_sum
    from data
), num_data as (
  select
    row_number() over () num,
    data.*
  from data
  limit ${ROWS_LIMIT}
), data_json as (
  select
    json_object_agg(d.schema_name || '.' || d.index_name, d) as json
  from num_data d
)
select
  json_build_object(
    'invalid_indexes',
    (select * from data_json),
    'invalid_indexes_total',
    (select row_to_json(dt) from data_total as dt)
  );
SQL
