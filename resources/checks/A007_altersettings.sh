# Collect pg_settings artifact
#dbg "PSQL_CONN_OPTIONS: ${PSQL_CONN_OPTIONS}"
psql ${PSQL_CONN_OPTIONS} -t -A <<SQL
with settings_count as (
    select json_object_agg(coalesce(s.sourcefile, 'default'), s.count) from (select sourcefile, count(ps.*) as count from pg_settings ps group by 1) s
), changes as (
select
        json_agg(json_build_object(
            'sourcefile', s.sourcefile,
            'count', s.count,
            'examples', s.examples
        ))
    from (select sourcefile, count(ps.*) as count, (json_agg(name order by name) filter (where sourcefile is not null)) as examples from pg_settings ps group by 1) s
)
select json_build_object('settings_count', (select * from settings_count), 'changes', (select * from changes));
SQL
