# Collect pg_settings artifact
dbg "PSQL_CONN_OPTIONS: ${PSQL_CONN_OPTIONS}"
#psql ${PSQL_CONN_OPTIONS} -c "select row_to_json(row) from (select * from pg_settings) row"
psql ${PSQL_CONN_OPTIONS} -c "select json_build_object(row.name, row) from (select * from pg_settings) row"

