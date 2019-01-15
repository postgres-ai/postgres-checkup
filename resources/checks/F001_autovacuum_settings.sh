settings=$(${CHECK_HOST_CMD} "${_PSQL} -f - " <<SQL
with global_settings as (
  select json_object_agg(s.name, s) from pg_settings s
  where name like '%autovacuum%'
    or name in (
      'vacuum_cost_delay',
      'vacuum_cost_limit', 
      'hot_standby_feedback',
      'maintenance_work_mem'
    )
), table_settings as (
  select json_object_agg(s.namespace || '.' || s.relname, s) from
    (select
        (select nspname from pg_namespace where oid = relnamespace)
        namespace,
        relname,
        reloptions
    from pg_class
    where reloptions::text ~ 'autovacuum') s
)
select json_build_object('global_settings', (select * from global_settings), 'table_settings', (select * from table_settings));
SQL
)
data="{\"settings\": $settings}"
data=$(jq -n "$data")
echo "$data"

