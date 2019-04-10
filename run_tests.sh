db_name='my_test_db'
current_dir=$( cd "$(dirname "${BASH_SOURCE[0]}")" ; pwd -P )

su - postgres -c "psql -A -t -d postgres -c 'SELECT version()'"

su - postgres -c "psql -A -t -d postgres -c \"
	select pg_terminate_backend(pid)
	from pg_stat_activity
	where datname <> current_database() and datname = '"${db_name}"'\""

su - postgres -c "psql -A -t -d postgres -c 'DROP DATABASE if exists ${db_name}'"
su - postgres -c "psql -A -t -d postgres -c 'CREATE DATABASE ${db_name}'"
su - postgres -c "psql -A -t -d ${db_name} -f "${current_dir}"/.ci/test_db_dump.sql"

echo "=======> Test started: H003 Non indexed FKs"
test_resut=1
su - postgres -c "psql -A -t -d ${db_name} -f "${current_dir}"/.ci/h003_step_1.sql"

rm -Rf ./artifacts
./checkup -h 127.0.0.1 --username postgres --project test --dbname ${db_name} -e 1 \
	--file ./resources/checks/H003_non_indexed_fks.sh

# one record must be exists
data_dir=$(cat ./artifacts/test/nodes.json | jq -r '.last_check | .dir') \
	&& result=$(cat ./artifacts/test/json_reports/$data_dir/H003_non_indexed_fks.json | jq '.results ."127.0.0.1" .data') \
	&& ([[ "$result" == "[]" ]] || [[ "$result" == "null" ]]) \
	&& echo "ERROR in H003: ${result} in '.results .\"127.0.0.1\" .data'" \
	&& echo $(cat ./artifacts/test/json_reports/$data_dir/H003_non_indexed_fks.json | jq '.') \
	&& test_resut=0

su - postgres -c "psql -A -t -d ${db_name} -f "${current_dir}"/.ci/h003_step_2.sql"
rm -Rf ./artifacts

./checkup -h 127.0.0.1 --username postgres --project test --dbname ${db_name} -e 1 \
	--file ./resources/checks/H003_non_indexed_fks.sh

# must be no records
data_dir=$(cat ./artifacts/test/nodes.json | jq -r '.last_check | .dir') \
	&& result=$(cat ./artifacts/test/json_reports/$data_dir/H003_non_indexed_fks.json | jq '.results ."127.0.0.1" .data') \
	&& cat ./artifacts/test/json_reports/$data_dir/H003_non_indexed_fks.json \
	&& (! [[ "$result" == "[]" ]]) \
	&& echo "ERROR in H003: found ${result} in '.results .\"127.0.0.1\" .data'" \
	&& test_resut=0

if [ "$test_resut" -eq "1" ]; then
	echo "<======= Test finished: H003"
else
	echo "<======= Test failed: H003"
fi
