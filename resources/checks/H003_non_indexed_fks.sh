# Foreign keys with Missing/Bad Indexes
if [[ ! -z ${IS_LARGE_DB+x} ]] && [[ ${IS_LARGE_DB} == "1" ]]; then
  INDEX_MIN_RELPAGES=10
  TABLE_MIN_RELPAGES=100
else
  INDEX_MIN_RELPAGES=0
  TABLE_MIN_RELPAGES=0
fi

${CHECK_HOST_CMD} "${_PSQL} -f -" <<SQL
with data as (
  with fk_actions ( code, action ) as (
    values ('a', 'error'),
           ('r', 'restrict'),
           ('c', 'cascade'),
           ('n', 'set null'),
           ('d', 'set default')
  ), fk_list as (
    select
      pg_constraint.oid as fkoid, conrelid, confrelid as parentid,
      conname,
      relname,
      nspname,
      fk_actions_update.action as update_action,
      fk_actions_delete.action as delete_action,
      conkey as key_cols,
      pg_class.relpages
    from pg_constraint
    join pg_class on conrelid = pg_class.oid
    join pg_namespace on pg_class.relnamespace = pg_namespace.oid
    join fk_actions as fk_actions_update on confupdtype = fk_actions_update.code
    join fk_actions as fk_actions_delete on confdeltype = fk_actions_delete.code
    where contype = 'f'
  ), fk_attributes as (
    select fkoid, conrelid, attname, attnum
    from fk_list
    join pg_attribute on conrelid = attrelid and attnum = any(key_cols)
    order by fkoid, attnum
  ), fk_cols_list as (
    select fkoid, array_agg(attname) as cols_list
    from fk_attributes
    group by fkoid
  ), index_list as (
    select
      indexrelid as indexid,
      pg_class.relname as indexname,
      indrelid,
      indkey,
      indpred is not null as has_predicate
    from pg_index
    join pg_class on indexrelid = pg_class.oid
    where indisvalid and pg_class.relkind = 'i' and pg_class.relpages > ${INDEX_MIN_RELPAGES}
  ), fk_index_match as (
    select
      fk_list.*,
      indexid,
      indexname,
      indkey::int[] as indexatts,
      has_predicate,
      array_length(key_cols, 1) as fk_colcount,
      array_length(indkey,1) as index_colcount,
      relpages,
      cols_list
    from fk_list
    join fk_cols_list using (fkoid)
    left join index_list on
      conrelid = indrelid
      and (indkey::int2[])[0:(array_length(key_cols,1) -1)] operator(pg_catalog.@>) key_cols
  ), fk_perfect_match as (
    select fkoid
    from fk_index_match
    where
      (index_colcount - 1) <= fk_colcount
      and not has_predicate
  ), fk_index_check as (
    select 'no index' as issue, *, 1 as issue_sort
    from fk_index_match
    where indexid is null
    union all
    select 'questionable index' as issue, *, 2
    from fk_index_match
    where
      indexid is not null
      and fkoid not in (select fkoid from fk_perfect_match)
  ), parent_table_stats as (
    select
      fkoid,
      tabstats.relname as parent_name,
      (n_tup_ins + n_tup_upd + n_tup_del + n_tup_hot_upd) as parent_writes,
      fk_list.relpages
    from pg_stat_user_tables as tabstats
    join fk_list on relid = parentid
  ), fk_table_stats as (
    select
      fkoid,
      (n_tup_ins + n_tup_upd + n_tup_del + n_tup_hot_upd) as writes,
      seq_scan as table_scans,
      relpages
    from pg_stat_user_tables as tabstats
    join fk_list on relid = conrelid
  )
  select
    nspname as schema_name,
    relname as table_name,
    conname as fk_name,
    issue,
    conrelid,
    writes,
    table_scans,
    parent_name,
    parentid,
    parent_writes,
    cols_list,
    indexid
  from fk_index_check
  join parent_table_stats using (fkoid)
  join fk_table_stats using (fkoid)
  where
    fk_table_stats.relpages > ${TABLE_MIN_RELPAGES}
    and (
      writes > 1000
      or parent_writes > 1000
       or parent_table_stats.relpages > ${TABLE_MIN_RELPAGES}
    )
  order by issue_sort, fk_table_stats.relpages desc, table_name, fk_name
),
num_data as (
  select row_number() over () num,
    schema_name,
    table_name,
    fk_name,
    issue,
    round(pg_relation_size(conrelid)/(1024^2)::numeric) as table_mb,
    writes,
    table_scans,
    parent_name,
    round(pg_relation_size(parentid)/(1024^2)::numeric) as parent_mb,
    parent_writes,
    cols_list,
    pg_get_indexdef(indexid) as indexdef
  from data
)
select json_object_agg(num_data.num, num_data) from num_data
SQL

# Based on https://github.com/pgexperts/pgx_scripts/blob/master/indexes/fk_no_index.sql

#Copyright (c) 2014, PostgreSQL Experts, Inc.
#    and Additional Contributors (see README)
#All rights reserved.
#
#Redistribution and use in source and binary forms, with or without
#modification, are permitted provided that the following conditions are met:
#
#* Redistributions of source code must retain the above copyright notice, this
#  list of conditions and the following disclaimer.
#
#* Redistributions in binary form must reproduce the above copyright notice,
#  this list of conditions and the following disclaimer in the documentation
#  and/or other materials provided with the distribution.
#
#* Neither the name of pgx_scripts nor the names of its
#  contributors may be used to endorse or promote products derived from
#  this software without specific prior written permission.
#
#THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
#AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
#IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
#DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
#FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
#DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
#SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
#CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
#OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
#OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
