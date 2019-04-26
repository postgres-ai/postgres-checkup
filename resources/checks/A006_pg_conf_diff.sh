# Collect pg_settings and pg_config values data to compare with same values from other domains
# Require comparation in plugin
${CHECK_HOST_CMD} "${_PSQL} -f -" <<SQL
with settings as (
  select
    json_object_agg(s.name, s)
  from (
    select *
    from pg_settings s
    where
      name not in ('linux_hz', 'pg_stat_kcache.linux_hz', 'transaction_read_only')
    order by name) s
), configs as (
  select
    json_object_agg(s.name, s)
  from (select * from pg_config s order by name) s
)
select json_build_object('pg_settings', (select * from settings), 'pg_config', (select * from configs));
SQL
