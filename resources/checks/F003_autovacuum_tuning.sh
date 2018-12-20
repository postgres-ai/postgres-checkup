settings=$(${CHECK_HOST_CMD} "${_PSQL} -f - " <<SQL
with global_settings as (
  select json_object_agg(s.name, s) from pg_settings s where name like 'autovacuum%'
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
#iotop_cmd="sudo iotop -o -d 5 -n 12 -k -b"
#iotop_result=$($iotop_cmd)
#iotop_data="{\"cmd\": \"$iotop_cmd\", \"data\": \"$iotop_result\"}"
#data="{\"settings\": $settings, \"iotop\": $iotop_data}"
data="{\"settings\": $settings}"
data=$(jq -n "$data")
echo "$data"

