sql=$(curl -s -L https://raw.githubusercontent.com/NikolayS/postgres_dba/56eded56e7ebfba3f20d7d455d8f84bd725affaf/sql/b1_table_estimation.sql | awk '{gsub("; *$", "", $0); print $0}')

${CHECK_HOST_CMD} "${_PSQL} -f -" <<SQL
with data as (
$sql
)
select json_object_agg(data."Table", data) as json from data;

SQL
