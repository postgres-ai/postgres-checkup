# Collect pg_settings artifact
sql=$(wget --quiet -O - https://github.com/NikolayS/postgres_dba/raw/master/sql/e1_extensions.sql)
sql=${sql%;} #remove last ;

#ssh ${HOST} "
${_PSQL} ${PSQL_CONN_OPTIONS} -f - <<SQL
with data as (
$sql
)
select json_agg(jsondata.json) from (select row_to_json(data) as json from data) jsondata;

SQL

