sql=$(curl -s -L https://github.com/NikolayS/postgres_dba/raw/master/sql/b1_table_estimation.sql | awk '{gsub("; *$", "", $0); print $0}')

ssh ${HOST} "psql -U postila_ru -t -f - " <<SQL
with data as (
$sql
)
select json_object_agg(data."Table", data) as json from data;

SQL