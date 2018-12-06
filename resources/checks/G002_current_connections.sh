# Current Activity: count of current connections grouped by database, user name, state

sql=$(curl -s -L https://raw.githubusercontent.com/NikolayS/postgres_dba/4.0/sql/a1_activity.sql | awk '{gsub("; *$", "", $0); print $0}')

${CHECK_HOST_CMD} "${_PSQL} ${PSQL_CONN_OPTIONS} -f -" <<SQL
with data as (
  $sql
),
num_data as (
  select row_number() over () num, data.* from data
)
select json_object_agg(num, num_data) as json from num_data
SQL

