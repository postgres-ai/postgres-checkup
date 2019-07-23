#!/bin/bash

if [[ ! -z "${1+x}" ]]; then
  hostname=$1
else
  hostname='localhost'
fi

# Generate more than 30000 tables with 2 index to check that
# in large db mode small tables and indexes are not analyzed.
for ((i=1; i < 3400; i++))
do
  psql -h ${hostname} -d dbname -U test_user -f -<<SQL
create table t_lg1_$i as select i from generate_series(1, 10) _(i);
create index concurrently i_lg1u_$i on t_lg1_$i(i);
create index concurrently i_lg1r_$i on t_lg1_$i(i);

create table t_lg2_$i as select i from generate_series(1, 10) _(i);
create index concurrently i_lg2u_$i on t_lg2_$i(i);
create index concurrently i_lg2r_$i on t_lg2_$i(i);

create table t_lg3_$i as select i from generate_series(1, 10) _(i);
create index concurrently i_lg3u_$i on t_lg3_$i(i);
create index concurrently i_lg3r_$i on t_lg3_$i(i);

create table t_lg4_$i as select i from generate_series(1, 10) _(i);
create index concurrently i_lg4u_$i on t_lg4_$i(i);
create index concurrently i_lg4r_$i on t_lg4_$i(i);

create table t_lg5_$i as select i from generate_series(1, 10) _(i);
create index concurrently i_lg5u_$i on t_lg5_$i(i);
create index concurrently i_lg5r_$i on t_lg5_$i(i);

create table t_lg6_$i as select i from generate_series(1, 10) _(i);
create index concurrently i_lg6u_$i on t_lg6_$i(i);
create index concurrently i_lg6r_$i on t_lg6_$i(i);

create table t_lg7_$i as select i from generate_series(1, 10) _(i);
create index concurrently i_lg7u_$i on t_lg7_$i(i);
create index concurrently i_lg7r_$i on t_lg7_$i(i);

create table t_lg8_$i as select i from generate_series(1, 10) _(i);
create index concurrently i_lg8u_$i on t_lg8_$i(i);
create index concurrently i_lg8r_$i on t_lg8_$i(i);

create table t_lg9_$i as select i from generate_series(1, 10) _(i);
create index concurrently i_lg9u_$i on t_lg9_$i(i);
create index concurrently i_lg9r_$i on t_lg9_$i(i);

create table t_lg10_$i as select i from generate_series(1, 10) _(i);
create index concurrently i_lg10u_$i on t_lg10_$i(i);
create index concurrently i_lg10r_$i on t_lg10_$i(i);
SQL
done