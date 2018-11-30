# Check disk space for important for postgres disk partitions

PG_MAJOR_VER=$(ssh "$HOST" "${_PSQL} -f -" <<EOF
  select setting::integer / 10000 from pg_settings where name = 'server_version_num'
EOF
)

PG_DATA_DIR=$(ssh "$HOST" "${_PSQL} -f -" <<EOF
  show data_directory
EOF
)
PG_DATA_DIR="data_directory: ${PG_DATA_DIR}"

PG_STATS_TEMP_DIR=$(ssh "$HOST" "${_PSQL} -f -" <<EOF
  show stats_temp_directory
EOF
)
PG_STATS_TEMP_DIR="stats_temp_directory: ${PG_STATS_TEMP_DIR}"

PG_LOG_DIR=$(ssh "$HOST" "${_PSQL} -f -" <<EOF
  show log_directory
EOF
)
PG_LOG_DIR="log_directory: ${PG_LOG_DIR}"

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

buf="{ "
for dir in "$PG_DATA_DIR" "$PG_STATS_TEMP_DIR" "$PG_STATS_TEMP_DIR" "$PG_LOG_DIR" ; do
  while IFS=": " read -r key value; do
    # save non-empty (to not override by empty strings)
    [[ ! -z "$key" ]] && NAME="$key" # postgres variable
    [[ ! -z "$value" ]] && PG_PATH="$value"
  done <<< ${dir}
 
  #echo "key: $NAME, value: $PG_PATH"

  cmd_raw_result=$(ssh "$HOST" "df -T \"${PG_PATH}\"")

  obj="\"${NAME}\": { path: \"${PG_PATH}\", df_raw_result: \"${cmd_raw_result}\" }"
  buf=" ${buf} ${obj}, "

done

buf="${buf} }"

jq -n "$buf"

