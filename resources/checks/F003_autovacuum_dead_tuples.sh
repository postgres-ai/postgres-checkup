# Collect autovacuum dead tuples

${CHECK_HOST_CMD} "${_PSQL} -f - " <<SQL
with overrided_tables as (
  select
    pc.oid as table_id,
    pn.nspname as scheme_name,
    pc.relname as table_name,
    pc.reloptions as options
  from pg_class pc
  join pg_namespace pn on pn.oid = pc.relnamespace
  where reloptions::text ~ 'autovacuum'
), data as (
  select
    coalesce(nullif(schemaname || '.', 'public.'), '') || c.relname as "relation",
    c.relkind,
    now() - last_autovacuum as since_last_autovacuum,
    now() - last_vacuum as since_last_vacuum,
    autovacuum_count as av_count,
    vacuum_count as v_count,
    n_tup_ins, 
    n_tup_upd, 
    n_tup_del,
    reltuples::int8 as pg_class_reltuples,
    n_live_tup,
    n_dead_tup,
    round((n_dead_tup::numeric * 100 / nullif(n_dead_tup + n_live_tup, 0))::numeric, 2) as dead_ratio,
    case when ot.table_id is not null then true else false end as overrided_settings
  from pg_stat_all_tables
  join pg_class c on c.oid = relid
  left join overrided_tables ot on ot.table_id = c.oid
  where reltuples > 10000
  order by 13 desc
), num_dead_tuples as (
  select
    row_number() over () num,
    data.*
  from data
), dead_tuples as (
  select json_object_agg(num_dead_tuples."relation", num_dead_tuples) as json from num_dead_tuples
), database_stat as (
  select
    row_to_json(dbstat)
  from (
    select
      sd.stats_reset::timestamptz(0),
      age(
        date_trunc('minute',now()),
        date_trunc('minute',sd.stats_reset)
      ) as stats_age
    from pg_stat_database sd
    where datname = current_database()
  ) dbstat
)
select
  json_build_object(
    'dead_tuples',
    (select * from dead_tuples),
    'database_stat',
    (select * from database_stat),
    'overrided_settings_count',
    (select count(1) from data where overrided_settings = true)
  );
SQL
