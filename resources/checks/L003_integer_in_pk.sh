${CHECK_HOST_CMD} "${_PSQL} -f - " <<SQL
do \$$
declare
  rec record;
  out text;
  val int8;
  ratio numeric;
begin
  out := '';
  for rec in
    select
      c.oid,
      (select spcname from pg_tablespace where oid = reltablespace) as tblspace,
      nspname as schema_name,
      relname as table_name,
      t.typname,
      attname
    from pg_index i
    join pg_class c on c.oid = i.indrelid
    left join pg_namespace n on n.oid = c.relnamespace
    join pg_attribute a on
      a.attrelid = i.indrelid
      and a.attnum = any(i.indkey)
    join pg_type t on t.oid = atttypid 
    where
      i.indisprimary
      and t.typname in ('int2', 'int4')
      and nspname <> 'pg_toast'
  loop
    execute format('select max(%I) from %I.%I;', rec.attname, rec.schema_name, rec.table_name) into val;
    if rec.typname = 'int4' then
      ratio := (val::numeric / 2^31)::numeric;
    elsif rec.typname = 'int2' then
      ratio := (val::numeric / 2^15)::numeric;
    else
      assert false, 'unreachable point';
    end if;
    if ratio > 0.00 then -- report only if > 1% of capacity is reached
      out := out || format(
        e'\nTable: %I.%I, Column: %I, Type: %s, Reached value: %s (%s%%)',
        -- e'\n%I.%I, %I, %s, %s (%s%%)',
        rec.schema_name,
        rec.table_name,
        rec.attname,
        rec.typname,
        val,
        round(100 * ratio, 2)
      );
    end if;
  end loop;
  raise info '%', out;
end;
\$$ language plpgsql;
SQL
