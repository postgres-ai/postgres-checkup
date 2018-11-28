sql=$(curl -s -L https://github.com/NikolayS/postgres_dba/raw/master/sql/b1_table_estimation.sql | sed -r ':a; s%(.*)/\*.*\*/%\1%; ta; /\/\*/ !b; N; ba' | sed '/^--/d' | sed '/^$/d')
sql="$sql"
sql=${sql%;} #remove last ;

ssh ${HOST} "psql -U postila_ru -t -f - " <<SQL
with data as (
$sql
)
select json_object_agg(data."Table", data) as json from data;

SQL