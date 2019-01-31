# how close to wraparound

${CHECK_HOST_CMD} "${_PSQL} -f - " <<SQL
with per_instance as (
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
    (greatest(age(c.relfrozenxid), age(t.relfrozenxid)) > 1200000000)::int as warning
  from pg_class c
  join pg_namespace n on c.relnamespace = n.oid
  left join pg_class t ON c.reltoastrelid = t.oid
  where c.relkind IN ('r', 'm') and not (n.nspname = 'pg_catalog' and c.relname <> 'pg_class')
    and n.nspname <> 'information_schema'
  order by 3 desc
  limit 50
)
select 
  json_build_object(
    'per_instance', 
    (select json_object_agg(i.datname, i) from per_instance i), 
    'per_database', 
    (select json_object_agg(d.relation, d) from per_database d)
  );
SQL
