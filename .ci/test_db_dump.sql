-- G003 Table lock_timeout 
create database checkup_test_db;
create user checkup_test_user with encrypted password 'mypass';
grant all privileges on database checkup_test_db to checkup_test_user;
alter user checkup_test_user set lock_timeout to '3s';
alter database checkup_test_db set lock_timeout = '4s';
--alter database checkup_test_db RESET configuration_parameter

-- rarely used indexes
create table t_rar_q as select id, (random() * 1000000)::int8 as t_dat from generate_series(1, 1000000) _(id);
create index t_rar_q_idx on t_rar_q(id);

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
set statement_timeout to '20ms';
create index concurrently i_invalid on test_schema.t_with_invalid_index(i);
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
ALTER TABLE t_with_bloat SET (autovacuum_vacuum_scale_factor=0.01);

VACUUM ANALYZE;

-- rarely used indexes
select * from t_rar_q where id = 23211;
update t_rar_q set t_dat=100 where t_dat between 553432 and 155343;
update t_rar_q set t_dat=200 where t_dat between 155343 and 255343;
update t_rar_q set t_dat=300 where t_dat between 255343 and 355343;
update t_rar_q set t_dat=400 where t_dat between 455343 and 555343;
update t_rar_q set t_dat=500 where t_dat between 555343 and 655343;
update t_rar_q set t_dat=600 where t_dat between 655343 and 755343;
update t_rar_q set t_dat=700 where t_dat between 755343 and 855343;
update t_rar_q set t_dat=800 where t_dat between 855343 and 955343;
update t_rar_q set t_dat=900 where t_dat between 955343 and 1055343;
-- F004
update t_with_bloat set i = i;

-- h002 Supports fk
select count(*) from t_slw_q;
explain select count(*) from t_slw_q;

-- L003
CREATE TABLE test_schema.orders
(
    id serial,
	cnt integer,
    CONSTRAINT orders_pk PRIMARY KEY (id)
);

INSERT INTO test_schema.orders(cnt) select id from generate_series(0, 100) _(id);
SELECT setval('test_schema.orders_id_seq'::regclass, 800000000, false);

CREATE TABLE test_schema."orders_A"
(
    id serial,
	cnt integer,
    CONSTRAINT "orders_A_pk" PRIMARY KEY (id)
);

INSERT INTO test_schema."orders_A"(cnt) select id from generate_series(0, 100) _(id);
SELECT setval('test_schema."orders_A_id_seq"'::regclass, 300000000, false);
