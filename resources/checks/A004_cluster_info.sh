# Collect pg cluster info
main_sql=$(curl -s -L https://raw.githubusercontent.com/NikolayS/postgres_dba/4.0/sql/0_node.sql | awk '{gsub("; *$", "", $0); print $0}')

pgver=$(${CHECK_HOST_CMD} "${_PSQL} -c \"SHOW server_version\"")

vers=(${pgver//./ })
majorVer=${vers[0]}

prepare_sql=""

if [[ $majorVer -lt 10 ]]; then
  #  echo "Version less 10"
  prepare_sql="
\set postgres_dba_last_wal_receive_lsn pg_last_xlog_receive_location
\set postgres_dba_last_wal_replay_lsn pg_last_xlog_replay_location
\set postgres_dba_is_wal_replay_paused pg_is_xlog_replay_paused "
else 
  #  echo "Version greater or equal 10"
  prepare_sql="
\set postgres_dba_last_wal_receive_lsn pg_last_wal_receive_lsn
\set postgres_dba_last_wal_replay_lsn pg_last_wal_replay_lsn
\set postgres_dba_is_wal_replay_paused pg_is_wal_replay_paused"

fi

${CHECK_HOST_CMD} "${_PSQL} -f - " <<SQL
$prepare_sql
with data as (
$main_sql
)
select json_object_agg(data.metric, data) as json from data where data.metric not like '------%';

SQL

