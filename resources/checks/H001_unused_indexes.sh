unusedSql=$(curl -s -L https://raw.githubusercontent.com/NikolayS/postgres_dba/4.0/sql/i1_rare_indexes.sql | awk '{gsub("; *$", "", $0); print $0}')
redundantSql=$(curl -s -L https://raw.githubusercontent.com/NikolayS/postgres_dba/4.0/sql/i2_redundant_indexes.sql | awk '{gsub("; *$", "", $0); print $0}')

#psql -U postila_ru -t -0 -f - <<SQL
${CHECK_HOST_CMD} "${_PSQL} ${PSQL_CONN_OPTIONS} -f -" <<SQL
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
), indexes_detail as (
  select i.*,
    (select run_in_separate_transactions from deploy where run_in_separate_transactions like '%' || i.indexname || '%') as drop_code,
    (select run_in_separate_transactions from revert where run_in_separate_transactions like '%' || i.indexname || '%') as revert_code  
    from indexes i
), unsed_indexes as (
  select json_object_agg(indexes_detail."indexname", indexes_detail) as json from indexes_detail
)
select json_build_object('indexes', (select * from unsed_indexes), 'drop_code', (select * from deploy_code), 'revert_code', (select * from revert_code));
SQL
