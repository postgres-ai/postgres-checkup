# Collect pg_settings artifact
sql=$(curl -s -L https://github.com/NikolayS/postgres_dba/raw/master/sql/e1_extensions.sql)
sql=${sql%;} #remove last ;

ssh ${HOST} "${_PSQL} ${PSQL_CONN_OPTIONS} -f -" <<SQL
with data as (
$sql
)
select json_object_agg(data.name, data) as json from data;

SQL

