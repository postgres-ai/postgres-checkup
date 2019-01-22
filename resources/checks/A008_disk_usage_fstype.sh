# Check disk space and file system type for important Postgres-related disk partitions

PG_MAJOR_VER=$(ssh "$HOST" "${_PSQL} -f -" <<EOF
  select setting::integer / 10000 from pg_settings where name = 'server_version_num'
EOF
)

PG_DATA_DIR=$(ssh "$HOST" "${_PSQL} -f -" <<EOF
  show data_directory
EOF
)

PG_STATS_TEMP_DIR=$(ssh "$HOST" "${_PSQL} -f -" <<EOF
  show stats_temp_directory
EOF
)

PG_LOG_DIR=$(ssh "$HOST" "${_PSQL} -f -" <<EOF
  show log_directory
EOF
)

# process relative paths
if ! [[ "${PG_LOG_DIR}" =~ ^/ ]]; then
  PG_LOG_DIR="${PG_DATA_DIR}/${PG_LOG_DIR}"
fi
if ! [[ "${PG_STATS_TEMP_DIR}" =~ ^/ ]]; then
  PG_STATS_TEMP_DIR="${PG_DATA_DIR}/${PG_STATS_TEMP_DIR}"
fi

PG_TABLESPSACES_DIRS=$(ssh "$HOST" "${_PSQL} -f -" <<EOF
  SELECT pg_catalog.pg_tablespace_location(oid)
  FROM pg_catalog.pg_tablespace
  WHERE pg_tablespace_location(oid) ~ '/'
EOF
)

if [[ "${PG_MAJOR_VER}" -ge "10" ]]; then
  PG_WAL_DIR="${PG_DATA_DIR}/pg_wal"
else
  PG_WAL_DIR="${PG_DATA_DIR}/pg_xlog"
fi

#######################################
# convert df output to json
# for usage inside print_df()
# Globals:
#   None
# Arguments:
#   (text) /path/to/dir
# Returns:
#   json
#######################################
df_to_json() {
  echo "{
  \"path\": \"$1\",
  \"device\": \"$2\",
  \"fstype\": \"$3\",
  \"size\": \"$4\",
  \"used\": \"$5\",
  \"avail\": \"$6\",
  \"use_percent\": \"$7\",
  \"mount_point\": \"$8\"
}"

}

#######################################
# ssh to host and invoke 'sudo df -TPh'
# for given path
# Globals:
#   HOST
# Arguments:
#   (text) /path/to/dir
# Returns:
#   path device fstype size used avail use_percent mount_point
#######################################
print_df() {
  local path="$1"
  df_to_json "${path}" $(ssh ${HOST} "sudo df -TPh \"${path}\" | tail -n +2")
}

# json output starts here
echo "{"

# print custom tablesapces
ts_cnt="0"
for ts in ${PG_TABLESPSACES_DIRS}; do
  ts_cnt=$(( ts_cnt + 1 ))
  echo "\"tablespace_${ts_cnt}\":"
  print_df "${ts}"
  echo ","
done

echo "\"data_directory\":"
print_df "$PG_DATA_DIR"
echo ","

echo "\"xlog_directory\":"
print_df "$PG_WAL_DIR"
echo ","

echo "\"stats_temp_directory\":"
print_df "$PG_STATS_TEMP_DIR"
echo ","

echo "\"log_directory\":"
print_df "$PG_LOG_DIR"

echo "}"

