# {{ .checkId }} Top 50 queries #

## Observations ##
{{ if .hosts.master }}
### Master (`{{.hosts.master}}`) ###
Snapshot time: {{(index (index .results .hosts.master) "data").snapshot_timestamptz }}

 Calls | Total time | Min time | Max time | Rows | shared_blks_hit | shared_blks_read | shared_blks_dirtied | shared_blks_written | local_blks_hit | local_blks_read | local_blks_dirtied | local_blks_written | temp_blks_read | temp_blks_written | blk_read_time | blk_write_time | kcache_reads | kcache_writes | kcache_user_time_ms | kcache_system_time | Query
-------|------------|----------|----------|------|-----------------|------------------|---------------------|---------------------|----------------|-----------------|--------------------|--------------------|----------------|-------------------|---------------|----------------|--------------|---------------|---------------------|--------------------|------- 
{{ range $i, $key := (index (index (index (index .results .hosts.master) "data") "queries") "_keys") }}
{{- $value := (index (index (index (index $.results $.hosts.master) "data") "queries") $key) -}}
{{- $value.diff_calls }}<br/>{{ $value.per_sec_calls }} sec<br/>{{ $value.per_call_calls }} per call |
{{- $value.diff_total_time }}<br/>{{ $value.per_sec_total_time }} sec<br/>{{ $value.per_call_total_time }} per call |
{{- $value.diff_min_time }}<br/>{{ $value.per_sec_min_time }} sec<br/>{{ $value.per_call_min_time }} per call |
{{- $value.diff_max_time }}<br/>{{ $value.per_sec_max_time }} sec<br/>{{ $value.per_call_max_time }} per call |
{{- $value.diff_rows }}<br/>{{ $value.per_sec_rows }} sec<br/>{{ $value.per_call_rows }} per call |
{{- $value.diff_shared_blks_hit }}<br/>{{ $value.per_sec_shared_blks_hit }} sec<br/>{{ $value.per_call_shared_blks_hit }} per call |
{{- $value.diff_shared_blks_read }}<br/>{{ $value.per_sec_shared_blks_read }} sec<br/>{{ $value.per_call_shared_blks_read }} per call |
{{- $value.diff_shared_blks_dirtied }}<br/>{{ $value.per_sec_shared_blks_dirtied }} sec<br/>{{ $value.per_call_shared_blks_dirtied }} per call |
{{- $value.diff_shared_blks_written }}<br/>{{ $value.per_sec_shared_blks_written }} sec<br/>{{ $value.per_call_shared_blks_written }} per call |
{{- $value.diff_local_blks_hit }}<br/>{{ $value.per_sec_local_blks_hit }} sec<br/>{{ $value.per_call_local_blks_hit }} per call |
{{- $value.diff_local_blks_read }}<br/>{{ $value.per_sec_local_blks_read }} sec<br/>{{ $value.per_call_local_blks_read }} per call |
{{- $value.diff_local_blks_dirtied }}<br/>{{ $value.per_sec_local_blks_dirtied }} sec<br/>{{ $value.per_call_local_blks_dirtied }} per call |
{{- $value.diff_local_blks_written }}<br/>{{ $value.per_sec_local_blks_written }} sec<br/>{{ $value.per_call_local_blks_written }} per call |
{{- $value.diff_temp_blks_read }}<br/>{{ $value.per_sec_temp_blks_read }} sec<br/>{{ $value.per_call_temp_blks_read }} per call |
{{- $value.diff_temp_blks_written }}<br/>{{ $value.per_sec_temp_blks_written }} sec<br/>{{ $value.per_call_temp_blks_written }} per call |
{{- $value.diff_blk_read_time }}<br/>{{ $value.per_sec_blk_read_time }} sec<br/>{{ $value.per_call_blk_read_time }} per call |
{{- $value.diff_blk_write_time }}<br/>{{ $value.per_sec_blk_write_time }} sec<br/>{{ $value.per_call_blk_write_time }} per call |
{{- $value.diff_kcache_reads }}<br/>{{ $value.per_sec_kcache_reads }} sec<br/>{{ $value.per_call_kcache_reads }} per call |
{{- $value.diff_kcache_writes }}<br/>{{ $value.per_sec_kcache_writes }} sec<br/>{{ $value.per_call_kcache_writes }} per call |
{{- $value.diff_kcache_user_time_ms }}<br/>{{ $value.per_sec_kcache_user_time_ms }} ms sec<br/>{{ $value.per_call_kcache_user_time_ms }} ms per call |
{{- $value.diff_kcache_system_time }}<br/>{{ $value.per_sec_kcache_system_time }} sec<br/>{{ $value.per_call_kcache_system_time }}  per call |
{{- Nobr $value.query }}
{{ end }}{{/* range */}}
{{ else }}{{/* if .host.master*/}}
No data
{{ end }}{{/* if .host.master*/}}

## Conclusions ##


## Recommendations ##

