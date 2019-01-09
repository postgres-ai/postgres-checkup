# Collect pg cluster info
main_sql=$(curl -s -L https://raw.githubusercontent.com/NikolayS/postgres_dba/4.0/sql/2_table_sizes.sql | awk '{gsub("; *$", "", $0); print $0}')

${CHECK_HOST_CMD} "${_PSQL} -f - " <<SQL
with data as (
  $main_sql
)
select json_object_agg(data."Table", data) as json from data where data."Table" not like ' '
SQL
