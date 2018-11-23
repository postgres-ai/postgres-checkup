# Collect pg_settings artifact
dbg "PSQL_CONN_OPTIONS: ${PSQL_CONN_OPTIONS}"
psql ${PSQL_CONN_OPTIONS} -c "select json_object_agg(s.name, s) from pg_settings s;"

