#!/bin/bash

if [[ ! -z "${1+x}" ]]; then
  hostname=$1
else
  hostname='localhost'
fi

# Generate more than 3000 tables with 2 index to check that
# in large db mode small tables and indexes are not analyzed.
for ((i=1; i < 1670; i++))
do
  psql -h ${hostname} -d dbname -U test_user -f -<<SQL
create table t_lg1_$i as select i from generate_series(1, 10) _(i);
create index concurrently i_lg1u_$i on t_lg1_$i(i);
create index concurrently i_lg1r_$i on t_lg1_$i(i);
create table t_lg2_$i as select i from generate_series(1, 10) _(i);
create index concurrently i_lg2u_$i on t_lg2_$i(i);
create index concurrently i_lg2r_$i on t_lg2_$i(i);
SQL
done