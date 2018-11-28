# Collect pg cluster info
sql=$(curl -s -L https://github.com/NikolayS/postgres_dba/raw/master/sql/0_node.sql)
main_sql=${sql%;} #remove last ;

pgver=$(ssh ${HOST} "${_PSQL} -c \"SHOW server_version\"")

vers=(${pgver//./ })
majorVer=${vers[0]}

prepare_sql=""

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

ssh ${HOST} "${_PSQL} -f - " <<SQL
$prepare_sql
with data as (
$main_sql
)
select json_agg(jsondata.json) from (select row_to_json(data) as json from data) jsondata;
SQL
