select 
    format('DROP INDEX CONCURRENTLY %s; -- %s, table %s', i.indexrelid::regclass::text, 'Invalid index', pct.relname) as drop_code,
    replace(
      format('%s; -- table %s', pg_get_indexdef(i.indexrelid), pct.relname),
      'CREATE INDEX',
      'CREATE INDEX CONCURRENTLY'
    ) as revert_code
from pg_index i
join pg_class as pci on pci.oid = i.indexrelid
join pg_class as pct on pct.oid = i.indrelid
left join pg_namespace pn on pn.oid = pct.relnamespace
-- where i.indisvalid = false; -- disable to debug