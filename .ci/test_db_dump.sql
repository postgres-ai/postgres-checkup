-- G003 Table lock_timeout 
create database checkup_test_db;
create user checkup_test_user with encrypted password 'mypass';
grant all privileges on database checkup_test_db to checkup_test_user;
alter user checkup_test_user set lock_timeout to '3s';
alter database checkup_test_db set lock_timeout = '4s';
--alter database checkup_test_db RESET configuration_parameter

-- Fillfactor
create table t_fillfactor (i int) with (fillfactor=60);

-- H002 Unused and redundant indexes
create table t_with_unused_index as select i from generate_series(1, 1000000) _(i);
create index concurrently i_unused on t_with_unused_index(i);
create table t_with_redundant_index as select i from generate_series(1, 1000000) _(i);
create index concurrently i_redundant_1 on t_with_redundant_index(i);
create index concurrently i_redundant_2 on t_with_redundant_index(i);

-- H001 invalid indexes
create table t_with_invalid_index as select i from generate_series(1, 1000000) _(i);
set statement_timeout to '10ms';
create index concurrently i_invalid on t_with_invalid_index(i);
set statement_timeout to 0;

-- H003 non indexed fks 
create table t_fk_1 as select id::int8 from generate_series(0, 1000000) _(id);
alter table t_fk_1 add primary key (id);
create table t_fk_2 as select id, (random() * 1000000)::int8 as t1_id from generate_series(1, 1000000) _(id);
alter table t_fk_2 add constraint fk_t2_t1 foreign key (t1_id) references t_fk_1(id);

-- Bloat level
create table bloated as select i from generate_series(1, 100000) _(i); 
create index i_bloated on bloated(i); 
delete from bloated where i % 2 = 0;

-- F004
create table t_with_bloat as select i from generate_series(1, 1000000) _(i);
update t_with_bloat set id = id;
