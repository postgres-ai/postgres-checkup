# Collect autovacuum dead tuples

${CHECK_HOST_CMD} "${_PSQL} -f - " <<SQL
with data as (
  select
    coalesce(nullif(schemaname || '.', 'public.'), '') || c.relname as "relation",
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
    round((n_dead_tup::numeric / nullif(n_dead_tup + n_live_tup, 0))::numeric, 2) as dead_ratio
  from pg_stat_all_tables
  join pg_class c on c.oid = relid
  where reltuples > 10000
  order by 12 desc limit 50
), dead_tuples as (
  select json_object_agg(data."relation", data) as json from data
), database_stat as (
  select
    row_to_json(dbstat)
  from (
    select
      sd.stats_reset::timestamptz(0)
    from pg_stat_database sd
    where datname = current_database()
  ) dbstat
)
select
  json_build_object(
    'dead_tuples',
    (select * from dead_tuples),
    'database_stat',
    (select * from database_stat)
  );
SQL
