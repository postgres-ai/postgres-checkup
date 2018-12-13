# Collect pg_settings artifact
${CHECK_HOST_CMD} "${_PSQL} -f -" <<SQL
with timeouts as (
  select json_object_agg(s.name,s ) from pg_settings s where name in ('statement_timeout', 'idle_in_transaction_session_timeout', 'authentication_timeout')
), locks as (
  select json_object_agg(s.name,s ) from pg_settings s where name in ('deadlock_timeout', 'lock_timeout', 'max_locks_per_transaction', 'max_pred_locks_per_page', 'max_pred_locks_per_relation', 'max_pred_locks_per_transaction')
), databases_stat as (
  select *, ((now() - sd.stats_reset)::interval(0)::text) as stats_reset_age from pg_stat_database sd where datname in (SELECT datname FROM pg_database WHERE datistemplate = false)
), dbs_data as (
  select json_object_agg(sd.datname, sd) from databases_stat sd
)
select json_build_object('timeouts', (select * from timeouts), 'locks', (select * from locks), 'databases_stat', (select * from dbs_data));
SQL
