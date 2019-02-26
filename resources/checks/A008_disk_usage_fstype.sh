# Check disk space and file system type for important Postgres-related disk partitions

if [[ "${SSH_SUPPORT}" = "false" ]]; then
  echo "SSH is not supported, skipping..." >&2
  exit 1
fi

PG_MAJOR_VER=$(${CHECK_HOST_CMD} "${_PSQL} -f -" <<EOF
  select setting::integer / 10000 from pg_settings where name = 'server_version_num'
EOF
)

PG_DATA_DIR=$(${CHECK_HOST_CMD} "${_PSQL} -f -" <<EOF
  show data_directory
EOF
)

PG_STATS_TEMP_DIR=$(${CHECK_HOST_CMD} "${_PSQL} -f -" <<EOF
  show stats_temp_directory
EOF
)

PG_LOG_DIR=$(${CHECK_HOST_CMD} "${_PSQL} -f -" <<EOF
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

PG_TABLESPSACES_DIRS=$(${CHECK_HOST_CMD} "${_PSQL} -f -" <<EOF
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
  \"fstype\": \"$3\",
  \"size\": \"$4\",
  \"avail\": \"$6\",
  \"used\": \"$5\",
  \"use_percent\": \"$7\",
  \"mount_point\": \"$8\",
  \"path\": \"$1\",
  \"device\": \"$2\"
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
  df_to_json "${path}" $(${CHECK_HOST_CMD} "sudo df -TPh \"${path}\" | tail -n +2")
}

# json output starts here
echo "{\"db_data\":{"

# print custom tablesapces
ts_cnt="0"
for ts in ${PG_TABLESPSACES_DIRS}; do
  ts_cnt=$(( ts_cnt + 1 ))
  echo "\"tablespace_${ts_cnt}\":"
  print_df "${ts}"
  echo ","
done

echo "\"PGDATA\":"
print_df "$PG_DATA_DIR"
echo ","

echo "\"WAL directory\":"
print_df "$PG_WAL_DIR"
echo ","

echo "\"stats_temp_directory\":"
print_df "$PG_STATS_TEMP_DIR"

# do not fail if log_directory does not exist
if $(${CHECK_HOST_CMD} "sudo stat \"$PG_LOG_DIR\" >/dev/null 2>&1"); then
  echo ","
  echo "\"log_directory\":"
  print_df "$PG_LOG_DIR"
fi

echo "},"
echo "\"fs_data\":{"

i=0
points=$(${CHECK_HOST_CMD} "sudo df -TPh | tail -n +2")
while read -r line; do
  if [[ $i -gt 0 ]]; then
    echo ",\"$i\":{"
  else
    echo "\"$i\":{"
  fi
  let i=$i+1
  params=($line)
  echo "  \"fstype\": \"${params[1]}\",
  \"size\": \"${params[2]}\",
  \"avail\": \"${params[4]}\",
  \"used\": \"${params[3]}\",
  \"use_percent\": \"${params[5]}\",
  \"mount_point\": \"${params[6]}\",
  \"path\": \"${params[6]}\",
  \"device\": \"${params[0]}\"
}"
done <<< "$points"
echo "}}"
