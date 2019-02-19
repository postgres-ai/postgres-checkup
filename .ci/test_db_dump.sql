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
create schema test_schema;
create table test_schema.t_with_invalid_index as select i from generate_series(1, 1000000) _(i);
set statement_timeout to '5s';
create index concurrently test_schema.i_invalid on test_schema.t_with_invalid_index(i);
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
update t_with_bloat set i = i;

-- h002 Supports fk
create table t_red_fk_1 as select id::int8 from generate_series(0, 1000000) _(id);
alter table t_red_fk_1 add primary key (id);
create index r_red_fk_1_id_idx on t_red_fk_1(id);
create index r_red_fk_1_X_idx on t_red_fk_1(id);

create table t_red_fk_2 as select id, (random() * 1000000)::int8 as r_t1_id from generate_series(1, 1000000) _(id);
alter table t_red_fk_2 add constraint fk_red_fk2_t1 foreign key (r_t1_id) references t_red_fk_1(id);
create index r_red_fk_2_fk_idx on t_red_fk_2(r_t1_id);

-- altered settings
alter system set random_page_cost = 2.22;
select pg_reload_conf();

--slow query
create table t_slw_q as select id::int8 from generate_series(0, 10000000) _(id);
select * from t_slw_q where id between 2000000 and 6001600;

-- rarely used indexes
create table t_rar_q as select id, (random() * 1000000)::int8 as t_dat from generate_series(1, 1000000) _(id);
create index t_rar_q_idx on t_rar_q(id);
select * from t_rar_q where id = 23211;
update t_rar_q set t_dat=100 where id between 553432 and 1553432;
update t_rar_q set t_dat=200 where id between 1553432 and 2553432;
update t_rar_q set t_dat=300 where id between 2553432 and 3553432;
update t_rar_q set t_dat=400 where id between 3553432 and 4553432;
update t_rar_q set t_dat=500 where id between 4553432 and 5553432;
