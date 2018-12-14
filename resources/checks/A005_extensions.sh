# Collect pg_settings artifact
sql=$(curl -s -L https://raw.githubusercontent.com/NikolayS/postgres_dba/4.0/sql/e1_extensions.sql | awk '{gsub("; *$", "", $0); print $0}')

${CHECK_HOST_CMD} "${_PSQL} -f -" <<SQL
with data as (
$sql
), withsettins as (
    select data.*, (select json_object_agg(pg_settings.name, pg_settings) from pg_settings where name ~ data.name) as settings from data
)
select json_object_agg(withsettins.name, withsettins) as json from withsettins;
SQL

