# Invalid keys

sql=$(curl -s -L https://raw.githubusercontent.com/NikolayS/postgres_dba/8e40d2345a54e2e6ed6605c0061e3bfd76fa032f/sql/i4_invalid_indexes.sql | awk '{gsub("; *$", "", $0); print $0}')

${CHECK_HOST_CMD} "${_PSQL} -f -" <<SQL
with data as (
  $sql
),
num_data as (
  select row_number() over () num, data.* from data
)
select json_object_agg(num_data.num, num_data) from num_data
SQL

