settings=$(${CHECK_HOST_CMD} "${_PSQL} -f - " <<SQL
with global_settings as (
  select
    json_object_agg(s.name, s)
  from (
    select *
    from pg_settings
    where (
      name ~ e'^(auto)?vacuum'
      or name in (
        'hot_standby_feedback',
        'maintenance_work_mem'
      )
    )
    order by name
  ) s
), table_settings as (
  select
    json_object_agg(s.namespace || '.' || s.relname, s)
  from
    (select
        (select nspname from pg_namespace where oid = relnamespace)
        namespace,
        relname,
        reloptions
    from pg_class
    where reloptions::text ~ 'autovacuum'
    order by namespace, relname
    ) s
)
select json_build_object('global_settings', (select * from global_settings), 'table_settings', (select * from table_settings));
SQL
)
data="{\"settings\": $settings}"
data=$(jq -n "$data")
echo "$data"

