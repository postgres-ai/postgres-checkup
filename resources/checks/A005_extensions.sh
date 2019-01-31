# Collect extensions and their settings

sql=$(curl -s -L https://raw.githubusercontent.com/NikolayS/postgres_dba/5.0/sql/e1_extensions.sql | awk '{gsub("; *$", "", $0); print $0}')

dbs=$(${CHECK_HOST_CMD} "${_PSQL} -f - " <<SQL
  select
    datname
  from pg_database
  where datname not in ('template0', 'template1', 'postgres')
  order by datname
SQL
)

result="{ }"

for cur_db in ${dbs}; do
  object=$(${CHECK_HOST_CMD} "${_PSQL} -d "$cur_db" -f -" <<SQL
    with data as (
      $sql
    ), withsettins as (
        select
          data.*,
          (select json_object_agg(name, setting)
        from pg_settings
        where name ~ data.name) as settings from data
        order by name
    )
    select json_object_agg(withsettins.name, withsettins) as json from withsettins;
SQL
)

  result=$(jq --arg db "${cur_db}" --argjson obj "$object" -r '. += { ($db): $obj }' <<<"${result}")
done

jq -r . <<<"$result"

