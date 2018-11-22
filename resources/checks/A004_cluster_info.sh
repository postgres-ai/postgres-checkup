# Collect pg_settings artifact
[[ -z ${PSQL_CONN_OPTIONS+x} ]] && PSQL_CONN_OPTIONS="-U postila_ru"
#dbg "PSQL_CONN_OPTIONS: ${PSQL_CONN_OPTIONS}"
pgver=$(psql ${PSQL_CONN_OPTIONS}  -c "SHOW server_version" -t -A)
vers=(${pgver//./ })
majorVer=${vers[0]}

prepare_sql=""
main_sql=$(cat <<EOF
with data as (
  select s.*
  from pg_stat_database s
  where s.datname = current_database()
), alldata as (
select 'Postgres Version' as metric, version() as value
union all
select
  'Config file' as metric,
  (select setting from pg_settings where name = 'config_file') as value
union all
select
  'Role' as metric,
  case
  when pg_is_in_recovery()  then 'Replica' || ' (delay: '
    || ((((case
        when :postgres_dba_last_wal_receive_lsn() = :postgres_dba_last_wal_replay_lsn() then 0
        else extract (epoch from now() - pg_last_xact_replay_timestamp())
      end)::int)::text || ' second')::interval)::text
    || '; paused: ' || :postgres_dba_is_wal_replay_paused()::text || ')'
  else
    'Master'
  end as value
union all
(
  with repl_groups as (
    select sync_state, state, string_agg(host(client_addr), ', ') as hosts
    from pg_stat_replication
    group by 1, 2
  )
  select
    'Replicas',
    string_agg(sync_state || '/' || state || ': ' || hosts, e'\n')
  from repl_groups
)
union all
select 'Started At', pg_postmaster_start_time()::timestamptz(0)::text
union all
select 'Uptime', (now() - pg_postmaster_start_time())::interval(0)::text
union all
select
  'Checkpoints',
  (select (checkpoints_timed + checkpoints_req)::text from pg_stat_bgwriter)
union all
select
  'Forced Checkpoints',
  (
    select round(100.0 * checkpoints_req::numeric /
      (nullif(checkpoints_timed + checkpoints_req, 0)), 1)::text || '%'
    from pg_stat_bgwriter
  )
union all
select
  'Checkpoint MB/sec',
  (
    select round((nullif(buffers_checkpoint::numeric, 0) /
      ((1024.0 * 1024 /
        (current_setting('block_size')::numeric))
          * extract('epoch' from now() - stats_reset)
      ))::numeric, 6)::text
    from pg_stat_bgwriter
  )
union all
select 'Database Name' as metric, datname as value from data
union all
select 'Database Size', pg_size_pretty(pg_database_size(current_database()))
union all
select 'Stats Since', stats_reset::timestamptz(0)::text from data
union all
select 'Stats Age', (now() - stats_reset)::interval(0)::text from data
union all
select 'Installed Extensions', (
  with exts as (
    select extname || ' ' || extversion e, (-1 + row_number() over (order by extname)) / 5 i from pg_extension
  ), lines(l) as (
    select string_agg(e, ', ' order by i) l from exts group by i
  )
  select string_agg(l, e'\n') from lines
)
union all
select 'Cache Effectiveness', (round(blks_hit * 100::numeric / (blks_hit + blks_read), 2))::text || '%' from data -- no "/0" because we already work!
union all
select 'Successful Commits', (round(xact_commit * 100::numeric / (xact_commit + xact_rollback), 2))::text || '%' from data
union all
select 'Conflicts', conflicts::text from data
union all
select 'Temp Files: total size', pg_size_pretty(temp_bytes)::text from data
union all
select 'Temp Files: total number of files', temp_files::text from data
union all
select 'Temp Files: avg file size', pg_size_pretty(temp_bytes::numeric / nullif(temp_files, 0))::text from data
union all
select 'Deadlocks', deadlocks::text from data
)
select json_object_agg(alldata.metric, alldata.value) from alldata;
EOF
)

prepare_sql=1
if [[ $majorVer -lt 10 ]]; then
#  echo "Version less 10"
  prepare_sql=$(cat <<EOF
\set postgres_dba_last_wal_receive_lsn pg_last_xlog_receive_location
\set postgres_dba_last_wal_replay_lsn pg_last_xlog_replay_location
\set postgres_dba_is_wal_replay_paused pg_is_xlog_replay_paused

EOF
)

else 
#  echo "Version more or equal 10"
  prepare_sql=$(cat <<EOF
\set postgres_dba_last_wal_receive_lsn pg_last_wal_receive_lsn
\set postgres_dba_last_wal_replay_lsn pg_last_wal_replay_lsn
\set postgres_dba_is_wal_replay_paused pg_is_wal_replay_paused

EOF
)

fi

psql ${PSQL_CONN_OPTIONS} <<SQL
$prepare_sql 
$main_sql
SQL
