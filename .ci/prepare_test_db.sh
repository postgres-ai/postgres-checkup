#!/bin/bash

if [[ ! -z "${1+x}" ]]; then
  hostname=$1
else
  hostname='localhost'
fi

# Generate more than 100 tables with index to check that
# the limitation of .md reports' rows counts work as expected.
for ((i=1; i < 110; i++))
do
  psql -h "${hostname}" -d dbname -U test_user -f -<<SQL
create table t_$i as select i from generate_series(1, 1000) _(i);
create index concurrently i_u_$i on t_$i(i);
create index concurrently i_r_$i on t_$i(i);
SQL
done