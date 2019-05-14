# Collect pg_settings artifact

${CHECK_HOST_CMD} "${_PSQL} -f -" "reset statement_timeout"

${CHECK_HOST_CMD} "${_PSQL} -f -" <<SQL
with timeouts as (
  select json_object_agg(s.name,s ) from pg_settings s where name in ('statement_timeout', 'idle_in_transaction_session_timeout', 'authentication_timeout')
), locks as (
  select json_object_agg(s.name,s ) from pg_settings s where name in ('deadlock_timeout', 'lock_timeout', 'max_locks_per_transaction', 'max_pred_locks_per_page', 'max_pred_locks_per_relation', 'max_pred_locks_per_transaction')
), databases_stat as (
  select
    *,
    ((now() - sd.stats_reset)::interval(0)::text) as stats_reset_age
  from pg_stat_database sd
  where datname in (SELECT datname FROM pg_database WHERE datistemplate = false)
  order by deadlocks desc
  limit ${ROWS_LIMIT}
), num_dbs_data as (
  select
    row_number() over () num,
    ds.*
  from databases_stat ds
), dbs_data as (
  select json_object_agg(sd.datname, sd) from num_dbs_data sd
), db_specified_settings as (
    select json_object_agg(dbs.database, dbs) from
        (select
            (select datname from pg_database where oid = pd.setdatabase) as database,
            *
        from pg_db_role_setting pd
        where
            setconfig::text ~ '(lock_timeout|deadlock_timeout)'
            and setdatabase is not null and setdatabase <> 0
        ) dbs
), user_specified_settings as (
    select json_object_agg(pr.rolname, pr) from pg_roles pr where rolconfig::text ~ '(lock_timeout|deadlock_timeout)'
)
select
    json_build_object(
        'timeouts', (select * from timeouts), 'locks', (select * from locks),
        'databases_stat', (select * from dbs_data),
        'db_specified_settings', (select * from db_specified_settings),
        'user_specified_settings', (select * from user_specified_settings)
    );
SQL

${CHECK_HOST_CMD} "${_PSQL} -f -" "set statement_timeout to ${STIMEOUT}s"
