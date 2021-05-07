# how close to wraparound

${CHECK_HOST_CMD} "${_PSQL} -f - " <<SQL
with overridden_tables as (
  select
    pc.oid as table_id,
    pn.nspname as scheme_name,
    pc.relname as table_name,
    pc.reloptions as options
  from pg_class pc
  join pg_namespace pn on pn.oid = pc.relnamespace
  where reloptions::text ~ 'autovacuum'
), per_instance as (
  select
    datname,
    age(datfrozenxid),
    round(
      age(datfrozenxid)::numeric * 100
        / (2 * 10^9 - current_setting('vacuum_freeze_min_age')::numeric
      )::numeric,
      2
    ) as capacity_used,
    datfrozenxid,
    (age(datfrozenxid) > 1200000000)::int as warning
  from pg_database
  order by 3 desc
), num_per_instance as (
  select
    row_number() over () num,
    per_instance.*
  from per_instance
), per_database as (
  select
    coalesce(nullif(n.nspname || '.', 'public.'), '') || c.relname as relation,
    greatest(age(c.relfrozenxid), age(t.relfrozenxid)) as age,
    round(
      (greatest(age(c.relfrozenxid), age(t.relfrozenxid))::numeric *
      100 / (2 * 10^9 - current_setting('vacuum_freeze_min_age')::numeric)::numeric),
      2
    ) as capacity_used,
    c.relfrozenxid as rel_relfrozenxid,
    t.relfrozenxid as toast_relfrozenxid,
    (greatest(age(c.relfrozenxid), age(t.relfrozenxid)) > 1200000000)::int as warning,
    case when ot.table_id is not null then true else false end as overridden_settings
  from pg_class c
  join pg_namespace n on c.relnamespace = n.oid
  left join pg_class t ON c.reltoastrelid = t.oid
  left join overridden_tables ot on ot.table_id = c.oid
  where c.relkind IN ('r', 'm') and not (n.nspname = 'pg_catalog' and c.relname <> 'pg_class')
    and n.nspname <> 'information_schema'
  order by 3 desc
), num_per_database as (
  select
    row_number() over () num,
    per_database.*
  from per_database
)
select
  json_build_object(
    'per_instance',
    (select json_object_agg(i.datname, i) from num_per_instance i),
    'per_database',
    (select json_object_agg(d.relation, d) from num_per_database d),
    'overridden_settings_count',
    (select count(1) from per_database where overridden_settings = true)
  );
SQL
