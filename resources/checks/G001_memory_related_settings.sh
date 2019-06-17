# Collect pg_settings artifact
${CHECK_HOST_CMD} "${_PSQL} -f -" <<SQL
    select json_object_agg(s.name, s)
    from pg_settings s
    where name in (
        'max_connections',
        'work_mem',
        'maintenance_work_mem',
        'autovacuum_work_mem',
        'shared_buffers',
        'effective_cache_size',
        'temp_buffers',
        'autovacuum_max_workers'
    );
SQL
