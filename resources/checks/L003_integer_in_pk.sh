if [[ ! -z ${IS_LARGE_DB+x} ]] && [[ ${IS_LARGE_DB} == "1" ]]; then
  MIN_RELPAGES=100
else
  MIN_RELPAGES=0
fi

f_stdout=$(mktemp)
f_stderr=$(mktemp)

(${CHECK_HOST_CMD} "${_PSQL} -f - " <<SQL
do \$$
declare
  MIN_RELPAGES int8 = ${MIN_RELPAGES}; -- skip tables with small number of pages
  rec record;
  out text;
  out1 json;
  i numeric;
  val int8;
  ratio numeric;
  sql text;
begin
  out := '';
  i := 0;
  for rec in
    select
      c.oid,
      (select spcname from pg_tablespace where oid = reltablespace) as tblspace,
      nspname as schema_name,
      relname as table_name,
      t.typname,
      (select pg_get_serial_sequence(quote_ident(nspname) || '.' || quote_ident(relname), attname)) as seq,
      min(attname) as attname
    from pg_index i
    join pg_class c on c.oid = i.indrelid
    left join pg_namespace n on n.oid = c.relnamespace
    join pg_attribute a on
      a.attrelid = i.indrelid
      and a.attnum = any(i.indkey)
    join pg_type t on t.oid = atttypid
    where
      i.indisprimary
      and (c.relpages >  or (select pg_get_serial_sequence(quote_ident(nspname) || '.' || quote_ident(relname), attname)) is not null)
      and t.typname in ('int2', 'int4')
      and nspname <> 'pg_toast'
      group by 1, 2, 3, 4, 5, 6
      having count(*) = 1 -- skip PKs with 2+ columns
  loop
    raise debug 'table: %', rec.table_name;

    if rec.seq is null then
        sql := format('select max(%I) from %I.%I;', rec.attname, rec.schema_name, rec.table_name);
    else
        sql := format('select last_value from %s;', rec.seq);
    end if;

    raise debug 'sql: %', sql;
    execute sql into val;

    if rec.typname = 'int4' then
      ratio := (val::numeric / 2^31)::numeric;
    elsif rec.typname = 'int2' then
      ratio := (val::numeric / 2^15)::numeric;
    else
      assert false, 'unreachable point';
    end if;

    if ratio > 0.1 then -- report only if > 10% of capacity is reached
      i := i + 1;

      out1 := json_build_object(
          'table',
          coalesce(nullif(quote_ident(rec.schema_name), 'public') || '.', '') || quote_ident(rec.table_name),
          'pk',
          rec.attname,
          'type',
          rec.typname,
          'current_max_value',
          val,
          'capacity_used_percent',
          round(100 * ratio, 2)
      );

      raise debug 'cur: %', out1;

      if out <> '' then out := out || ', '; end if;

      out := out || '"' || rec.table_name || '":' || out1 || '';
    end if;
  end loop;

  out := '{' || out || '}';

  raise info '%', out;
end;
\$$ language plpgsql;
SQL
) >$f_stdout 2>$f_stderr

result=$(cat $f_stderr)
result=${result:23:$((${#result}))}
tables_data=$(echo "$result" | jq -cs 'sort_by(-(.[]."capacity_used_percent"|tonumber)) | .[]' | jq -s add)
min_table_size_bytes=$((MIN_RELPAGES * 8192))

echo "{\"tables\": $tables_data, \"min_table_size_bytes\": $min_table_size_bytes }"

rm -f "$f_stderr" "$f_stdout"
