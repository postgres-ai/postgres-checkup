#!/bin/bash

# Generate more than 100 tables with index to check that
# the limitation of .md reports' rows counts work as expected.
for ((i=1; i < 110; i++))
do
  psql -d dbname -U test_user -c "create table t_$i as select i from generate_series(1, 1000) _(i);"
  psql -d dbname -U test_user -c "create index concurrently i_u_$i on t_$i(i);"
  psql -d dbname -U test_user -c "create index concurrently i_r_$i on t_$i(i);"
done