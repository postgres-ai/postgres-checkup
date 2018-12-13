# Collect pg_settings artifact
${CHECK_HOST_CMD} "${_PSQL} -f -" <<SQL
with stat_statements as (
  select json_object_agg(pg_settings.name, pg_settings) as json from pg_settings where name ~ 'pg_stat_statements'
), kcache as (
  select json_object_agg(pg_settings.name, pg_settings) as json from pg_settings where name ~ 'pg_stat_kcache'
)
select json_build_object('pg_stat_statements', (select * from stat_statements), 'kcache', (select * from kcache));
SQL

