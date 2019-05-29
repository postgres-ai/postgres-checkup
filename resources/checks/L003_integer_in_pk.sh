f_stdout=$(mktemp)
f_stderr=$(mktemp)

(${CHECK_HOST_CMD} "${_PSQL} -f - " <<SQL
do \$$
declare
  rec record;
  out text;
  i numeric;
  val int8;
  ratio numeric;
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
      attname,
      (select pg_get_serial_sequence(quote_ident(nspname) || '.' || quote_ident(relname), attname)) as seq
    from pg_index i
    join pg_class c on c.oid = i.indrelid
    left join pg_namespace n on n.oid = c.relnamespace
    join pg_attribute a on
      a.attrelid = i.indrelid
      and a.attnum = any(i.indkey)
    join pg_type t on t.oid = atttypid
    where
      i.indisprimary
      and (c.relpages > 1000 or (select pg_get_serial_sequence(quote_ident(nspname) || '.' || quote_ident(relname), attname)) is not null)
      and t.typname in ('int2', 'int4')
      and nspname <> 'pg_toast'
  loop
    if rec.seq is null then
        execute format('select max(%I) from %I.%I;', rec.attname, rec.schema_name, rec.table_name) into val;
    else
        execute format('SELECT last_value FROM %s;', rec.seq) into val;
    end if;
    if rec.typname = 'int4' then
      ratio := (val::numeric / 2^31)::numeric;
    elsif rec.typname = 'int2' then
      ratio := (val::numeric / 2^15)::numeric;
    else
      assert false, 'unreachable point';
    end if;
    if ratio > 0.1 then -- report only if > 10% of capacity is reached
      i := i + 1;
      out := out || '{"' || rec.table_name || '":' || json_build_object(
          'Table',
          coalesce(nullif(quote_ident(rec.schema_name), 'public') || '.', '') || quote_ident(rec.table_name),
          'PK',
          rec.attname,
          'Type',
          rec.typname,
          'Current max value',
          val,
          'Capacity used, %',
          round(100 * ratio, 2)
      ) || '}';
    end if;
  end loop;
  raise info '%', out;
end;
\$$ language plpgsql;
SQL
) >$f_stdout 2>$f_stderr

result=$(cat $f_stderr)
result=${result:23:$((${#result}))}

echo "$result" | jq -cs 'sort_by(-(.[]."Capacity used, %"|tonumber)) | .[]' | jq -s add

rm -f "$f_stderr" "$f_stdout"
