# Collect pg_settings artifact
# Require comparation in plugin
dbg "PSQL_CONN_OPTIONS: ${PSQL_CONN_OPTIONS}"
pg_settings=$(psql ${PSQL_CONN_OPTIONS} -c "select json_object_agg(s.name, s) from pg_settings s;" -t -A)
pg_config=$(psql ${PSQL_CONN_OPTIONS} -c "select json_object_agg(c.name, c) from pg_config c;" -t -A)
echo '{ "pg_settings": '$pg_settings', "pg_config":'$pg_config'}'

