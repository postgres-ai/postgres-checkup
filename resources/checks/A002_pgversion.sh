# Collect postgres version info
${CHECK_HOST_CMD} "${_PSQL} -f - " <<SQL
select
    json_build_object('version', version(),
        'server_version_num', current_setting('server_version_num'),
        'server_major_ver', (select regexp_replace(current_setting('server_version'), '\\.\\d+$', '')),
        'server_minor_ver', (select regexp_replace(current_setting('server_version'), '^.*\\.', '')));
SQL
