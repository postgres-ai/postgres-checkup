# Foreign keys with Missing/Bad Indexes

sql=$(curl -s -L https://raw.githubusercontent.com/NikolayS/postgres_dba/master/sql/i3_non_indexed_fks.sql | awk '{gsub("; *$", "", $0); print $0}')

data=$(${CHECK_HOST_CMD} "${_PSQL} -f -" <<SQL
-- disble psql CLI options '-A -t' for this report
\pset tuples_only
\pset format aligned
\pset expanded

with data as (
  $sql
),
num_data as (
  select row_number() over () num, data.* from data
)
select * from num_data as index
SQL)

data="{ \"raw\": \"${data}\" }"
jq -n "${data}"

