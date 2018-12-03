sql=$(curl -s -L https://github.com/NikolayS/postgres_dba/raw/master/sql/b2_btree_estimation.sql | awk '{gsub("; *$", "", $0); print $0}')

ssh ${HOST} "psql -U postila_ru -t -f - " <<SQL
with data as (
$sql
)
select json_agg(jsondata.json) from (select row_to_json(data) as json from data) jsondata;
SQL

#For get objects change row 9 to `select json_object_agg(data."Index (Table)", data) as json from data;`
#but in this case we have indexes like `i_user_visits_postgrest_auth\n  (user_visits)`