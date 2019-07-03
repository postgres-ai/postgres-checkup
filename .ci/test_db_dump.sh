#!/bin/bash

# Generate more than 100 tables with index to check limitation of table rows count
for ((i=1; i < 110; i++))
do
  psql -d dbname -U test_user -c "create table t_$i as select i from generate_series(1, 1000) _(i);"
  psql -d dbname -U test_user -c "create index concurrently i_$i on t_$i(i);"
done