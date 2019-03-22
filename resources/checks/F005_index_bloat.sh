${CHECK_HOST_CMD} "${_PSQL} -f -" <<SQL
with data as (
  with overrided_tables as (
    select
      pc.oid as table_id,
      pn.nspname as scheme_name,
      pc.relname as table_name,
      pc.reloptions as options
    from pg_class pc
    join pg_namespace pn on pn.oid = pc.relnamespace
    where reloptions::text ~ 'autovacuum'
  ), step0 as (
      select
        tbl.oid tblid, nspname, tbl.relname AS tblname, idx.relname AS idxname, idx.reltuples, idx.relpages, idx.relam,
        indrelid, indexrelid, regexp_split_to_table(indkey::text, ' ')::smallint AS attnum, --indkey::smallint[] AS attnum,
        coalesce(substring(array_to_string(idx.reloptions, ' ') from 'fillfactor=([0-9]+)')::smallint, 90) as fillfactor
      from pg_index
      join pg_class idx on idx.oid = pg_index.indexrelid
      join pg_class tbl on tbl.oid = pg_index.indrelid
      join pg_namespace on pg_namespace.oid = idx.relnamespace
      join pg_am a ON idx.relam = a.oid
      where a.amname = 'btree'
        AND pg_index.indisvalid
        AND tbl.relkind = 'r'
        AND idx.relpages > 10
        AND pg_namespace.nspname NOT IN ('pg_catalog', 'information_schema')
  ), step1 as (
    select
      i.tblid,
      i.nspname as schema_name,
      i.tblname as table_name,
      i.idxname as index_name,
      i.reltuples,
      i.relpages,
      i.relam,
      a.attrelid AS table_oid,
      current_setting('block_size')::numeric AS bs,
      fillfactor,
      -- MAXALIGN: 4 on 32bits, 8 on 64bits (and mingw32 ?)
      case when version() ~ 'mingw32|64-bit|x86_64|ppc64|ia64|amd64' then 8 else 4 end as maxalign,
      /* per page header, fixed size: 20 for 7.X, 24 for others */
      24 AS pagehdr,
      /* per page btree opaque data */
      16 AS pageopqdata,
      /* per tuple header: add IndexAttributeBitMapData if some cols are null-able */
      case
        when max(coalesce(s.null_frac,0)) = 0 then 2 -- IndexTupleData size
        else 2 + (( 32 + 8 - 1 ) / 8) -- IndexTupleData size + IndexAttributeBitMapData size ( max num filed per index + 8 - 1 /8)
      end as index_tuple_hdr_bm,
      /* data len: we remove null values save space using it fractionnal part from stats */
      sum((1 - coalesce(s.null_frac, 0)) * coalesce(s.avg_width, 1024)) as nulldatawidth,
      max(case when a.atttypid = 'pg_catalog.name'::regtype then 1 else 0 end) > 0 as is_na
    from pg_attribute as a
    join step0 as i on a.attrelid = i.indexrelid AND a.attnum = i.attnum
    join pg_stats as s on
      s.schemaname = i.nspname
      and (
        (s.tablename = i.tblname and s.attname = pg_catalog.pg_get_indexdef(a.attrelid, a.attnum, true)) -- stats from tbl
        OR (s.tablename = i.idxname AND s.attname = a.attname) -- stats from functionnal cols
      )
    join pg_type as t on a.atttypid = t.oid
    where a.attnum > 0
    group by 1, 2, 3, 4, 5, 6, 7, 8, 9, 10
  ), step2 as (
    select
      *,
      (
        index_tuple_hdr_bm + maxalign
        -- Add padding to the index tuple header to align on MAXALIGN
        - case when index_tuple_hdr_bm % maxalign = 0 THEN maxalign else index_tuple_hdr_bm % maxalign end
        + nulldatawidth + maxalign
        -- Add padding to the data to align on MAXALIGN
        - case
            when nulldatawidth = 0 then 0
            when nulldatawidth::integer % maxalign = 0 then maxalign
            else nulldatawidth::integer % maxalign
          end
      )::numeric as nulldatahdrwidth
      -- , index_tuple_hdr_bm, nulldatawidth -- (DEBUG INFO)
    from step1
  ), step3 as (
    select
      *,
      -- ItemIdData size + computed avg size of a tuple (nulldatahdrwidth)
      coalesce(1 + ceil(reltuples / floor((bs - pageopqdata - pagehdr) / (4 + nulldatahdrwidth)::float)), 0) as est_pages,
      coalesce(1 + ceil(reltuples / floor((bs - pageopqdata - pagehdr) * fillfactor / (100 * (4 + nulldatahdrwidth)::float))), 0) as est_pages_ff
      -- , stattuple.pgstatindex(quote_ident(nspname)||'.'||quote_ident(idxname)) AS pst, index_tuple_hdr_bm, maxalign, pagehdr, nulldatawidth, nulldatahdrwidth, reltuples -- (DEBUG INFO)
    from step2
    join pg_am am on step2.relam = am.oid
    where am.amname = 'btree'
  ), step4 as (
    SELECT
      *,
      bs*(relpages)::bigint AS real_size,
  -------current_database(), nspname AS schemaname, tblname, idxname, bs*(relpages)::bigint AS real_size,
      bs*(relpages-est_pages)::bigint AS extra_size,
      100 * (relpages-est_pages)::float / relpages AS extra_ratio,
      bs*(relpages-est_pages_ff) AS bloat_size,
      100 * (relpages-est_pages_ff)::float / relpages AS bloat_ratio
      -- , 100-(sub.pst).avg_leaf_density, est_pages, index_tuple_hdr_bm, maxalign, pagehdr, nulldatawidth, nulldatahdrwidth, sub.reltuples, sub.relpages -- (DEBUG INFO)
    from step3
    -- WHERE NOT is_na
  )
  select
    case is_na when true then 'TRUE' else '' end as "Is N/A",
    format(
      \$out$%s
    (%s)\$out$,
      left(index_name, 50) || case when length(index_name) > 50 then '…' else '' end,
      coalesce(nullif(step4.schema_name, 'public') || '.', '') || step4.table_name
    ) as "Index (Table)",
    real_size as "Real size bytes",
    pg_size_pretty(real_size::numeric) as "Size",
    extra_ratio as "Extra ratio",
    case
      when extra_size::numeric >= 0
        then '~' || pg_size_pretty(extra_size::numeric)::text || ' (' || round(extra_ratio::numeric, 2)::text || '%)'
      else null
    end  as "Extra",
    case
      when extra_size::numeric >= 0
        then extra_size
      else null
    end as "Extra size bytes",
    case
      when bloat_size::numeric >= 0
        then '~' || pg_size_pretty(bloat_size::numeric)::text || ' (' || round(bloat_ratio::numeric, 2)::text || '%)'
      else null
    end as "Bloat",
    case
      when (bloat_size)::numeric >=0
        then bloat_size
        else null
      end as "Bloat size bytes",
    case
      when (bloat_ratio)::numeric >=0
        then bloat_ratio
        else null
      end as "Bloat ratio",
    case
      when (real_size - bloat_size)::numeric >=0
        then '~' || pg_size_pretty((real_size - bloat_size)::numeric)
        else null
     end as "Live",
    case
      when (real_size - bloat_size)::numeric >=0
        then (real_size - bloat_size)::numeric
        else null
     end as "Live bytes",     
    fillfactor,
    case when ot.table_id is not null then true else false end as overrided_settings
  from step4
  left join overrided_tables ot on ot.table_id = step4.tblid
  order by real_size desc nulls last
), limited_data as (
  select * from data limit ${ROWS_LIMIT}
), num_limited_data as (
  select
    row_number() over () num,
    limited_data.*
  from limited_data
), limited_json_data as (
  select json_object_agg(ld."Index (Table)", ld) as json from num_limited_data ld
), total_data as (
  select
    sum(1) as count,
    sum("Extra size bytes") as "Extra size bytes sum",
    sum("Real size bytes") as "Real size bytes sum",
    sum("Bloat size bytes") as "Bloat size bytes sum",
   (sum("Bloat size bytes")::numeric/sum("Real size bytes")::numeric * 100) as "Bloat ratio",
    sum("Extra size bytes") as "Extra size bytes sum"
  from data
)
select
  json_build_object(
    'index_bloat',
    (select * from limited_json_data),
    'index_bloat_total',
    (select row_to_json(total_data) from total_data),
    'overrided_settings_count',
    (select count(1) from limited_data where overrided_settings = true)

  )
SQL

