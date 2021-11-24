#!/bin/bash
PG_VER=$1

echo "deb http://apt.postgresql.org/pub/repos/apt $(lsb_release -cs)-pgdg main" > /etc/apt/sources.list.d/pgdg.list
wget --quiet -O - https://www.postgresql.org/media/keys/ACCC4CF8.asc | apt-key add -
apt-get update
apt-get -y upgrade
apt-get -y install postgresql-${PG_VER} postgresql-contrib-${PG_VER} postgresql-client-${PG_VER} postgresql-server-dev-${PG_VER} && apt-get install -y postgresql-${PG_VER}-pg-stat-kcache
psql --version
source ~/.profile
echo "127.0.0.2 postgres.test1.node" >> /etc/hosts # replica 1
echo "127.0.0.3 postgres.test2.node" >> /etc/hosts # replica 2
echo "127.0.0.4 postgres.test3.node" >> /etc/hosts # master

# Configure postgres

## Configure pg_hba.conf
cat > /etc/postgresql/${PG_VER}/main/pg_hba.conf << EOL
local   all all trust
host all  all    0.0.0.0/0  md5
host all  all    ::1/128  trust
host replication  replication    ::1/128  md5
EOL

## Configure general postgres params
echo "listen_addresses='*'" >> /etc/postgresql/${PG_VER}/main/postgresql.conf
echo "log_filename='postgresql-${PG_VER}-main.log'" >> /etc/postgresql/${PG_VER}/main/postgresql.conf
echo "shared_preload_libraries = 'pg_stat_statements,auto_explain,pg_stat_kcache'" >> /etc/postgresql/${PG_VER}/main/postgresql.conf

## Configure general postgres master node params
echo "wal_level = hot_standby" >> /etc/postgresql/${PG_VER}/main/postgresql.conf
echo "max_wal_senders = 5" >> /etc/postgresql/${PG_VER}/main/postgresql.conf
echo "wal_keep_segments = 32" >> /etc/postgresql/${PG_VER}/main/postgresql.conf
echo "archive_mode    = on" >> /etc/postgresql/${PG_VER}/main/postgresql.conf
echo "archive_command = 'cp %p /path_to/archive/%f'" >> /etc/postgresql/${PG_VER}/main/postgresql.conf

## Start postgres master node
/etc/init.d/postgresql start 
psql -U postgres -c "create role replication with replication password 'rEpLpAssw' login"
psql -U postgres -c 'create database dbname;'
psql -U postgres dbname -b -c 'create extension if not exists pg_stat_statements;'
psql -U postgres dbname -b -c 'create extension if not exists pg_stat_kcache;'
psql -U postgres dbname -c "create role test_user superuser login;"
psql -U postgres -c 'show data_directory;'


#######################################
# Add and start a new Postgres replica locally
# Globals:
#   PG_VER
# Arguments:
#   (number) Replica's number, must be unique
#   (port) TCP port to be used
# Returns:
#   None
#######################################
function add_replica() {
  local num="$1"
  local port="$2"

  ## Configure data storage
  sudo -u postgres mkdir /var/lib/postgresql/${PG_VER}/data${num} && sudo -u postgres chmod 0700 /var/lib/postgresql/${PG_VER}/data${num}
  sudo -u postgres /usr/lib/postgresql/${PG_VER}/bin/initdb /var/lib/postgresql/${PG_VER}/data${num}
  sudo -u postgres cp /etc/postgresql/${PG_VER}/main/pg_hba.conf /var/lib/postgresql/${PG_VER}/data${num}/

  ## Configure general postgres settings
  echo "port = ${port}" >> /var/lib/postgresql/${PG_VER}/data${num}/postgresql.conf
  echo "listen_addresses='*'" >> /var/lib/postgresql/${PG_VER}/data${num}/postgresql.conf
  echo "shared_preload_libraries = 'pg_stat_statements,auto_explain,pg_stat_kcache'" >> /var/lib/postgresql/${PG_VER}/data${num}/postgresql.conf
  sudo -u postgres /usr/lib/postgresql/${PG_VER}/bin/pg_ctl -D /var/lib/postgresql/${PG_VER}/data${num} -l /var/log/postgresql/replica${num}.log start || cat /var/log/postgresql/replica${num}.log
  psql -U postgres -p ${port} -c 'show data_directory;'
  psql -U postgres -p ${port} -c 'create database dbname;'
  psql -U postgres -p ${port} dbname -b -c 'create extension if not exists pg_stat_statements;'
  psql -U postgres -p ${port} dbname -b -c 'create extension if not exists pg_stat_kcache;'
  psql -U postgres -p ${port} dbname -c "create role test_user superuser login;"
  sudo -u postgres /usr/lib/postgresql/${PG_VER}/bin/pg_ctl -D /var/lib/postgresql/${PG_VER}/data${num} -l /var/log/postgresql/replica${num}.log stop

  ## Configure replica postgres settings
  echo "hot_standby = on" >> /var/lib/postgresql/${PG_VER}/data${num}/postgresql.conf
  echo "standby_mode = 'on'" > /var/lib/postgresql/${PG_VER}/data${num}/recovery.conf
  echo "primary_conninfo = 'host=127.0.0.4 port=5432 user=replication password=rEpLpAssw'" >> /var/lib/postgresql/${PG_VER}/data${num}/recovery.conf
  echo "trigger_file = '/var/lib/postgresql/${PG_VER}/data${num}/trigger'" >> /var/lib/postgresql/${PG_VER}/data${num}/recovery.conf
  echo "restore_command = 'cp /path_to/archive/%f "%p"'" >> /var/lib/postgresql/${PG_VER}/data${num}/recovery.conf

  ## Start replica
  sudo -u postgres /usr/lib/postgresql/${PG_VER}/bin/pg_ctl -D /var/lib/postgresql/${PG_VER}/data${num} -l /var/log/postgresql/secondary1.log start || cat /var/log/postgresql/replica${num}.log
}

add_replica 1 5433
add_replica 2 5434
