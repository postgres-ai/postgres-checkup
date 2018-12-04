# Collect pg_settings artifact
${CHECK_HOST_CMD} "${_PSQL} -c \"select json_object_agg(s.name, s) from pg_settings s\""
