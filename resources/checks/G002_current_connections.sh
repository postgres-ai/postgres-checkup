# Current Activity: count of current connections grouped by database, user name, state
${CHECK_HOST_CMD} "${_PSQL} -f -" <<SQL
with data as (
  select
    coalesce(usename, '** ALL users **') as "user",
    coalesce(datname, '** ALL databases **') as "database",
    coalesce(state, '** ALL states **') as "current_state",
    count(*) as "count",
    count(*) filter (where state_change < now() - interval '1 minute') as "state_changed_more_1m_ago",
    count(*) filter (where state_change < now() - interval '1 hour') as "state_changed_more_1h_ago",
    count(*) filter (where xact_start < now() - interval '1 minute') as "tx_age_more_1m",
    count(*) filter (where xact_start < now() - interval '1 hour') as "tx_age_more_1h"
  from pg_stat_activity
  where query not like 'autovacuum: %'
  group by grouping sets ((datname, usename, state), (usename, state), ())
  order by
    usename is null desc,
    datname is null desc,
    2 asc,
    3 asc,
    count(*) desc
),
num_data as (
  select row_number() over () num, data.*
  from data
)
select json_object_agg(num_data.num, num_data) from num_data
SQL
