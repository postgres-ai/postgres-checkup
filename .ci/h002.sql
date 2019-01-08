create table t1 as select id::int8 from generate_series(1, 10000000) _(id);
alter table t1 add primary key (id);
create table t2 as select id, (random() * 10000000)::int8 as t1_id from generate_series(1, 10000000) _(id);
alter table t2 add constraint fk_t2_t1 foreign key (t1_id) references t1(id);
