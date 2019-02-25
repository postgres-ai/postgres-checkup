# Collect extensions and their settings

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
      select
        ae.name,
        installed_version,
        default_version,
        case when installed_version <> default_version then 'OLD' end as is_old
      from pg_extension e
      join pg_available_extensions ae on extname = ae.name
      order by ae.name
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

