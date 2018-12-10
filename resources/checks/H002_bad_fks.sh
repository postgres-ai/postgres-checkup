# Foreign keys with Missing/Bad Indexes

sql=$(curl -s -L https://raw.githubusercontent.com/NikolayS/postgres_dba/master/sql/i3_non_indexed_fks.sql | awk '{gsub("; *$", "", $0); print $0}')

${CHECK_HOST_CMD} "${_PSQL} ${PSQL_CONN_OPTIONS} -f -" <<SQL
with data as (
  $sql
),
num_data as (
  select row_number() over () num, data.* from data
)
select json_object_agg(num, num_data) as json from num_data
SQL

