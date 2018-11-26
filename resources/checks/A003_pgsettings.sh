# Collect pg_settings artifact
dbg "PSQL_CONN_OPTIONS: ${PSQL_CONN_OPTIONS}"
ssh ${HOST} "${_PSQL} -c \"select json_object_agg(s.name, s) from pg_settings s\""
