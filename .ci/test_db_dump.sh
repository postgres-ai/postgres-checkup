#!/bin/bash

for ((i=1; i < 1000; i++))
do
  psql -h localhost -d dbname -U username -c "create table t_$i as select i from generate_series(1, 1000) _(i);"
  psql -h localhost -d dbname -U username -c "create index concurrently i_$i on t_$i(i);"
done