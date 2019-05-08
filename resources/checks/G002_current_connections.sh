# Current Activity: count of current connections grouped by database, user name, state
${CHECK_HOST_CMD} "${_PSQL} -f -" <<SQL
with data as (
  select
    coalesce(usename, '** ALL users **') as "User",
    coalesce(datname, '** ALL databases **') as "DB",
    coalesce(state, '** ALL states **') as "Current State",
    count(*) as "Count",
    count(*) filter (where state_change < now() - interval '1 minute') as "State changed >1m ago",
    count(*) filter (where state_change < now() - interval '1 hour') as "State changed >1h ago",
    count(*) filter (where xact_start < now() - interval '1 minute') as "Tx age >1m",
    count(*) filter (where xact_start < now() - interval '1 hour') as "Tx age >1h"
  from pg_stat_activity
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
  limit ${ROWS_LIMIT}
)
select json_object_agg(num_data.num, num_data) from num_data
SQL
