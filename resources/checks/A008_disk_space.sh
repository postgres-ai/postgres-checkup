
set -x

PG_MAJOR_VER=$(ssh "$HOST" "${_PSQL} -f -" <<EOF
  select setting::integer / 10000 from pg_settings where name = 'server_version_num'
EOF
)

PG_DATA_DIR=$(ssh "$HOST" "${_PSQL} -f -" <<EOF
  show data_directory
EOF
)
PG_DATA_DIR="data_directory: ${PG_DATA_DIR}"

PG_STATS_TEMP_DIRECTORY=$(ssh "$HOST" "${_PSQL} -f -" <<EOF
  show stats_temp_directory
EOF
)
PG_STATS_TEMP_DIRECTORY="stats_temp_directory: ${PG_STATS_TEMP_DIRECTORY}"

PG_LOG_DIRECTORY=$(ssh "$HOST" "${_PSQL} -f -" <<EOF
  show log_directory
EOF
)
PG_LOG_DIRECTORY="log_directory: ${PG_LOG_DIRECTORY}"

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

local dir key value
for dir in "$PG_DATA_DIR"; do
  while IFS=": " read -r key value; do
    # save non-empty (to not override by empty)
    [[ ! -z "$key" ]] && NAME="$key" # postgres variable
    [[ ! -z "$value" ]] && VALUE="$value"
  done <<< ${dir}
  echo "key: $NAME, value: $VALUE"

done

set +x
