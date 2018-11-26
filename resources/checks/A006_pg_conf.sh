# Collect pg_settings artifact
# Require comparation in plugin
pg_settings=$(ssh ${HOST} "${_PSQL} -c \"select json_object_agg(s.name, s) from pg_settings s;\" -t -A")
pg_config=$(ssh ${HOST} "${_PSQL} -c \"select json_object_agg(c.name, c) from pg_config c;\" -t -A")
echo '{ "pg_settings": '$pg_settings', "pg_config":'$pg_config'}'

