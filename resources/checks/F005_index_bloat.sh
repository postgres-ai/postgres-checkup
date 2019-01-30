sql=$(curl -s -L https://raw.githubusercontent.com/NikolayS/postgres_dba/5.0/sql/b2_btree_estimation.sql | awk '{gsub("; *$", "", $0); print $0}')

${CHECK_HOST_CMD} "${_PSQL} -f -" <<SQL
with data as (
$sql
)
select json_object_agg(data."Index (Table)", data) as json from data
SQL

