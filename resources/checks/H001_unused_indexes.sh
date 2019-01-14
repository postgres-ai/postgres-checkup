unusedSql=$(curl -s -L https://raw.githubusercontent.com/NikolayS/postgres_dba/4.0/sql/i1_rare_indexes.sql | awk '{gsub("; *$", "", $0); print $0}')
redundantSql=$(curl -s -L https://raw.githubusercontent.com/NikolayS/postgres_dba/4.0/sql/i2_redundant_indexes.sql | awk '{gsub("; *$", "", $0); print $0}')

unusedSql=$( cat <<SQL
WITH table_scans as (
    SELECT relid,
        tables.idx_scan + tables.seq_scan as all_scans,
        ( tables.n_tup_ins + tables.n_tup_upd + tables.n_tup_del ) as writes,
                pg_relation_size(relid) as table_size
        FROM pg_stat_user_tables as tables
),
all_writes as (
    SELECT sum(writes) as total_writes
    FROM table_scans
),
indexes as (
    SELECT idx_stat.relid, idx_stat.indexrelid,
        idx_stat.schemaname, idx_stat.relname as tablename,
        idx_stat.indexrelname as indexname,
        idx_stat.idx_scan,
        pg_relation_size(idx_stat.indexrelid) as index_bytes,
        indexdef ~* 'USING btree' AS idx_is_btree
    FROM pg_stat_user_indexes as idx_stat
        JOIN pg_index
            USING (indexrelid)
        JOIN pg_indexes as indexes
            ON idx_stat.schemaname = indexes.schemaname
                AND idx_stat.relname = indexes.tablename
                AND idx_stat.indexrelname = indexes.indexname
    WHERE pg_index.indisunique = FALSE
),
index_ratios AS (
    SELECT schemaname, tablename, indexname,
        idx_scan, all_scans,
        round(( CASE WHEN all_scans = 0 THEN 0.0::NUMERIC
            ELSE idx_scan::NUMERIC/all_scans * 100 END),2) as index_scan_pct,
        writes,
        round((CASE WHEN writes = 0 THEN idx_scan::NUMERIC ELSE idx_scan::NUMERIC/writes END),2)
            as scans_per_write,
        pg_size_pretty(index_bytes) as index_size,
        pg_size_pretty(table_size) as table_size,
        idx_is_btree, index_bytes
        FROM indexes
        JOIN table_scans
        USING (relid)
),
index_groups AS (
    SELECT 'Never Used Indexes' as reason, *, 1 as grp
    FROM index_ratios
    WHERE
        idx_scan = 0
        and idx_is_btree
    UNION ALL
    SELECT 'Low Scans, High Writes' as reason, *, 2 as grp
    FROM index_ratios
    WHERE
        scans_per_write <= 1
        and index_scan_pct < 10
        and idx_scan > 0
        and writes > 100
        and idx_is_btree
    UNION ALL
    SELECT 'Seldom Used Large Indexes' as reason, *, 3 as grp
    FROM index_ratios
    WHERE
        index_scan_pct < 5
        and scans_per_write > 1
        and idx_scan > 0
        and idx_is_btree
        and index_bytes > 100000000
    UNION ALL
    SELECT 'High-Write Large Non-Btree' as reason, index_ratios.*, 4 as grp
    FROM index_ratios, all_writes
    WHERE
        ( writes::NUMERIC / ( total_writes + 1 ) ) > 0.02
        AND NOT idx_is_btree
        AND index_bytes > 100000000
    ORDER BY grp, index_bytes DESC 
)
SELECT reason, schemaname, tablename, indexname, idx_scan, all_scans,
    index_scan_pct, scans_per_write, index_size, table_size
FROM index_groups

SQL
)

#psql -U postila_ru -t -0 -f - <<SQL
${CHECK_HOST_CMD} "${_PSQL} -f -" <<SQL
with indexes as (
  $unusedSql
), migrations as (
  $redundantSql
), deploy as (
  select * from (select * from migrations limit (select count(1) from migrations)/2) as docode where docode.run_in_separate_transactions not like '--%'
), revert as (
  select * from (select * from migrations offset ((select count(1) from migrations)/2 + 1)) as revertcode where revertcode.run_in_separate_transactions not like '--%' 
), deploy_code as (
  select json_agg(jsondata.json) from (select run_in_separate_transactions as json from deploy) jsondata
), revert_code as (
  select json_agg(jsondata.json) from (select run_in_separate_transactions as json from revert) jsondata
), unsed_indexes as (
  select json_object_agg(indexes."indexname", indexes) as json from indexes
)
select json_build_object('indexes', (select * from unsed_indexes), 'drop_code', (select * from deploy_code), 'revert_code', (select * from revert_code));
SQL
