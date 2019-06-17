# Autovacuum: resource usage

${CHECK_HOST_CMD} "${_PSQL} -f - " <<SQL
  select
    json_object_agg(s.name, s)
  from pg_settings s
  where
    name in ('log_autovacuum_min_duration',
             'autovacuum_max_workers',
             'autovacuum_work_mem',
             'work_mem',
             'maintenance_work_mem',
             'shared_buffers',
             'max_connections',
             'maintenance_work_mem',
             'log_autovacuum_min_duration',
             'autovacuum_vacuum_cost_limit',
             'vacuum_cost_limit',
             'autovacuum_vacuum_cost_delay'
             )
SQL
