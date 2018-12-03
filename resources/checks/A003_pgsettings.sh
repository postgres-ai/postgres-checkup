# Collect pg_settings artifact
ssh ${HOST} "${_PSQL} -c \"select json_object_agg(s.name, s) from pg_settings s\""
