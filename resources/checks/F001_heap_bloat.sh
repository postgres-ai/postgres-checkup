sql=$(curl -s -L https://raw.githubusercontent.com/NikolayS/postgres_dba/4.0/sql/b1_table_estimation.sql | awk '{gsub("; *$", "", $0); print $0}')

ssh ${HOST} "${_PSQL} ${PSQL_CONN_OPTIONS} -f -" <<SQL
with data as (
$sql
)
select json_object_agg(data."Table", data) as json from data;

SQL