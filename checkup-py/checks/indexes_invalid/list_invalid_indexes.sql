select 
    pn.nspname as schema_name,
	pct.relname as table_name,
    pci.relname as index_name,
    pg_size_pretty(pg_relation_size(i.indexrelid)) index_size
from pg_index i
join pg_class as pci on pci.oid = i.indexrelid
join pg_class as pct on pct.oid = i.indrelid
left join pg_namespace pn on pn.oid = pct.relnamespace
-- where i.indisvalid = false; -- disable to debug
order by pci.relname desc