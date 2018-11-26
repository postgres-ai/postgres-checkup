# Collect postgres version info
dbg "PSQL_CONN_OPTIONS: ${PSQL_CONN_OPTIONS}"
ssh ${HOST} "${_PSQL} ${PSQL_CONN_OPTIONS} -f - " <<SQL
select
    json_build_object('version', version(),
        'server_version_num', current_setting('server_version_num'),
        'server_major_ver', SPLIT_PART(current_setting('server_version'), '.', 1),
        'server_minor_ver', SPLIT_PART(current_setting('server_version'), '.', 2));
SQL
