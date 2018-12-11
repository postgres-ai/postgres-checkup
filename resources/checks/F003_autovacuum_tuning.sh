settings_sql="select json_object_agg(s.name, s) from pg_settings s where name like 'autovacuum%';"
settings=$(${CHECK_HOST_CMD} "${_PSQL} -f - " <<SQL
$settings_sql
SQL
)
iotop_cmd="sudo iotop -o -d 5 -n 12 -k -b"
iotop_result=$($iotop_cmd)
iotop_data="{\"cmd\": \"$iotop_cmd\", \"data\": \"$iotop_result\"}"
data="{\"settings\": $settings, \"iotop\": $iotop_data}"
data=$(jq -n "$data")
echo "$data"

