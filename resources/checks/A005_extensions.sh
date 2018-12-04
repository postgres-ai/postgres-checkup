# Collect pg_settings artifact
ssh ${HOST} "${_PSQL} ${PSQL_CONN_OPTIONS} -f -" <<SQL
with row as (
  select
    ae.name,
    installed_version,
    default_version,
    extversion as available_version,
    case when installed_version <> extversion then 'OLD' end as actuality
  from pg_extension e
  join pg_available_extensions ae on extname = ae.name
  order by ae.name
)
select json_object_agg(row.name, row) from row;
SQL

