#!/bin/bash

# For limit 100 rows per table will generate 110 tables/indexes
for ((i=1; i < 110; i++))
do
  psql -d dbname -U username -c "create table t_$i as select i from generate_series(1, 1000) _(i);"
  psql -d dbname -U username -c "create index concurrently i_$i on t_$i(i);"
done