# Collect pg_settings artifact
#${CHECK_HOST_CMD} "${_PSQL} ${PSQL_CONN_OPTIONS} -f -" <<SQL
psql -U postila_ru -t -0 -f - <<SQL
with timeouts as (
  select json_object_agg(s.name,s ) from pg_settings s where name in ('statement_timeout', 'idle_in_transaction_session_timeout', 'authentication_timeout')
), locks as (
  select json_object_agg(s.name,s ) from pg_settings s where name in ('deadlock_timeout', 'lock_timeout', 'max_locks_per_transaction', 'max_pred_locks_per_page', 'max_pred_locks_per_relation', 'max_pred_locks_per_transaction')
)
select json_build_object('timeouts', (select * from timeouts), 'locks', (select * from locks));
SQL
