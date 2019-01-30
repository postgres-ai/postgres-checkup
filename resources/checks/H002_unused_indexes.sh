unusedSql=$(curl -s -L https://raw.githubusercontent.com/NikolayS/postgres_dba/5.0/sql/i1_rare_indexes.sql | awk '{gsub("; *$", "", $0); print $0}')
redundantSql=$(curl -s -L https://raw.githubusercontent.com/NikolayS/postgres_dba/5.0/sql/i2_redundant_indexes.sql | awk '{gsub("; *$", "", $0); print $0}')
migrationSql=$(curl -s -L https://raw.githubusercontent.com/NikolayS/postgres_dba/5.0/sql/i5_indexes_migration.sql | awk '{gsub("; *$", "", $0); print $0}')

${CHECK_HOST_CMD} "${_PSQL} -f -" <<SQL
with indexes as (
  $unusedSql
), redundant_indexes as (
  $redundantSql
), migrations as (
  $migrationSql
), deploy as (
  select * from (select * from migrations limit (select count(1) from migrations)/2) as docode where docode.run_in_separate_transactions not like '--%'
), revert as (
  select * from (select * from migrations offset ((select count(1) from migrations)/2 + 1)) as revertcode where revertcode.run_in_separate_transactions not like '--%' 
), deploy_code as (
  select json_agg(jsondata.json) from (select run_in_separate_transactions as json from deploy) jsondata
), revert_code as (
  select json_agg(jsondata.json) from (select run_in_separate_transactions as json from revert) jsondata
), unsed as (
  select json_object_agg(i."index_name", i) as json from indexes i
), redundant as (
  select json_object_agg(ri."index_name", ri) as json from redundant_indexes ri
), database_stat as (
  select
    row_to_json(dbstat)
  from (
    select
      sd.stats_reset::timestamptz(0),
      age(
        date_trunc('minute',now()),
        date_trunc('minute',sd.stats_reset)
      ) as stats_age
    from pg_stat_database sd
    where datname = current_database()
  ) dbstat
)
select json_build_object(
        'unused_indexes',
        (select * from unsed),
        'redundant_indexes',
        (select * from redundant),
        'drop_code',
        (select * from deploy_code),
        'revert_code',
        (select * from revert_code),
        'database_stat',
        (select * from database_stat)
    );
SQL
