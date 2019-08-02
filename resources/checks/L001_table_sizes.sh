# Collect pg cluster info
${CHECK_HOST_CMD} "${_PSQL} -f - " <<SQL
with data as (
  select
    c.oid,
    (select spcname from pg_tablespace where oid = reltablespace) as tblspace,
    nspname as schema_name,
    relname as table_name,
    c.reltuples as row_estimate,
    pg_total_relation_size(c.oid) as total_bytes,
    pg_indexes_size(c.oid) as index_bytes,
    pg_total_relation_size(reltoastrelid) as toast_bytes,
    pg_total_relation_size(c.oid) - pg_indexes_size(c.oid) - coalesce(pg_total_relation_size(reltoastrelid), 0) as table_bytes
  from pg_class c
  left join pg_namespace n on n.oid = c.relnamespace
  where relkind = 'r' and nspname <> 'pg_catalog'
  order by c.relpages desc
), data2 as (
  select
    null::oid as oid,
    null as tblspace,
    null as schema_name,
    '    tablespace: [' || coalesce(tblspace, 'pg_default') || ']' as table_name,
    sum(row_estimate) as row_estimate,
    sum(total_bytes) as total_bytes,
    sum(index_bytes) as index_bytes,
    sum(toast_bytes) as toast_bytes,
    sum(table_bytes) as table_bytes
  from data
  where (select count(distinct coalesce(tblspace, 'pg_default')) from data) > 1
  group by tblspace
  union all
  select * from data
), tables as (
  select
    coalesce(nullif(schema_name, 'public') || '.', '') || table_name || coalesce(' [' || tblspace || ']', '') as "table",
    row_estimate as row_estimate,
    total_bytes as total_size_bytes,
    round(
      100 * total_bytes::numeric / nullif(sum(total_bytes) over (partition by (schema_name is null), left(table_name, 3) = '***'), 0),
      2
    ) as "total_size_percent",
    table_bytes as table_size_bytes,
    round(
      100 * table_bytes::numeric / nullif(sum(table_bytes) over (partition by (schema_name is null), left(table_name, 3) = '***'), 0),
      2
    ) as "table_size_percent",
    index_bytes as indexes_size_bytes,
    round(
      100 * index_bytes::numeric / nullif(sum(index_bytes) over (partition by (schema_name is null), left(table_name, 3) = '***'), 0),
      2
    ) as "index_size_percent",
    toast_bytes as toast_size_bytes,
    round(
      100 * toast_bytes::numeric / nullif(sum(toast_bytes) over (partition by (schema_name is null), left(table_name, 3) = '***'), 0),
      2
    ) as "toast_size_percent"
  from data2
  where schema_name is distinct from 'information_schema'
  order by oid is null desc, total_bytes desc nulls last
), total_data as (
  select
    sum(1) as count,
    sum("row_estimate") as "row_estimate_sum",
    sum("total_size_bytes") as "total_size_bytes_sum",
    sum("table_size_bytes") as "table_size_bytes_sum",
    sum("indexes_size_bytes") as "indexes_size_bytes_sum",
    sum("toast_size_bytes") as "toast_size_bytes_sum"
  from tables
), num_tables as (
  select
    row_number() over () num,
    ts.*
  from tables ts
), tables_json_data as (
  select json_object_agg(nt."table", nt) as json from num_tables nt
)
select
  json_build_object(
    'tables_data',
    (select * from tables_json_data),
    'tables_data_total',
    (select row_to_json(total_data) from total_data)
  );
SQL
