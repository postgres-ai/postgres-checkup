# Current Activity: count of current connections grouped by database, user name, state

source ~/tmp/wrapper.sh

sql=$(curl -s -L https://raw.githubusercontent.com/NikolayS/postgres_dba/4.0/sql/a1_activity.sql | awk '{gsub("; *$", "", $0); print $0}')

data=$(${CHECK_HOST_CMD} "${_PSQL} ${PSQL_CONN_OPTIONS} -f -" <<SQL
with data as (
  $sql
),
num_data as (
  select row_number() over () num, data.* from data
)
select * from num_data
SQL)

data="{ \"raw\": \"${data}\" }"
jq -rn "${data}"

