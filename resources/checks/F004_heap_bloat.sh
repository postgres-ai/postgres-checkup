sql=$(curl -s -L https://raw.githubusercontent.com/NikolayS/postgres_dba/5ba8c4754b8b67400434bc371197aba2137a4f6d/sql/b1_table_estimation.sql | awk '{gsub("; *$", "", $0); print $0}')

${CHECK_HOST_CMD} "${_PSQL} -f -" <<SQL
with data as (
$sql
)
select json_object_agg(data."Table", data) as json from data;

SQL
