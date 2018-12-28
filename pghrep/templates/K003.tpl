# {{ .checkId }} Top 50 queries #<a name="K003"/>

## Observations ##
{{ if .hosts.master }}
### Master (`{{.hosts.master}}`) ###
Start: {{ (index (index (index .results .hosts.master) "data") "start_timestamptz") }}  
End: {{ (index (index (index .results .hosts.master) "data") "end_timestamptz") }}  
Period, seconds: {{ (index (index (index .results .hosts.master) "data") "period_seconds") }}  
Period, age: {{ (index (index (index .results .hosts.master) "data") "period_age") }}  

\# | Calls | Total&nbsp;time | Rows | shared_blks_hit | shared_blks_read | shared_blks_dirtied | shared_blks_written | blk_read_time | blk_write_time | kcache_reads | kcache_writes | kcache_user_time_ms | kcache_system_time | Query
----|-------|------------|------|-----------------|------------------|---------------------|---------------------|---------------|----------------|--------------|---------------|---------------------|--------------------|------- 
{{ range $i, $key := (index (index (index (index .results .hosts.master) "data") "queries") "_keys") }}
{{- $value := (index (index (index (index $.results $.hosts.master) "data") "queries") $key) -}}
{{- $key}} |
{{- RoundUp $value.diff_calls 2 }}<br/>{{ RoundUp $value.per_sec_calls 2 }}/sec<br/>{{ RoundUp $value.per_call_calls 2 }}/call |
{{- RoundUp $value.diff_total_time 2 }} ms<br/>{{ RoundUp $value.per_sec_total_time 2 }} ms/sec<br/>{{ RoundUp $value.per_call_total_time 2 }} ms/call |
{{- RoundUp $value.diff_rows 2 }}<br/>{{ RoundUp $value.per_sec_rows 2 }}/sec<br/>{{ RoundUp $value.per_call_rows 2 }}/call |
{{- RoundUp $value.diff_shared_blks_hit 2 }} blks<br/>{{ RoundUp $value.per_sec_shared_blks_hit 2 }} blks/sec<br/>{{ RoundUp $value.per_call_shared_blks_hit 2 }} blks/call |
{{- RoundUp $value.diff_shared_blks_read 2 }} blks<br/>{{ RoundUp $value.per_sec_shared_blks_read 2 }} blks/sec<br/>{{ RoundUp $value.per_call_shared_blks_read 2 }} blks/call |
{{- RoundUp $value.diff_shared_blks_dirtied 2 }} blks<br/>{{ RoundUp $value.per_sec_shared_blks_dirtied 2 }} blks/sec<br/>{{ RoundUp $value.per_call_shared_blks_dirtied 2 }} blks/call |
{{- RoundUp $value.diff_shared_blks_written 2 }} blks<br/>{{ RoundUp $value.per_sec_shared_blks_written 2 }} blks/sec<br/>{{ RoundUp $value.per_call_shared_blks_written 2 }} blks/call |
{{- RoundUp $value.diff_blk_read_time 2 }} ms<br/>{{ RoundUp $value.per_sec_blk_read_time 2 }} ms/sec<br/>{{ RoundUp $value.per_call_blk_read_time 2 }} ms/call |
{{- RoundUp $value.diff_blk_write_time 2 }} ms<br/>{{ RoundUp $value.per_sec_blk_write_time 2 }} ms/sec<br/>{{ RoundUp $value.per_call_blk_write_time 2 }} ms/call |
{{- RoundUp $value.diff_kcache_reads 2 }} bytes<br/>{{ RoundUp $value.per_sec_kcache_reads 2 }} bytes/sec<br/>{{ RoundUp $value.per_call_kcache_reads 2 }} bytes/call |
{{- RoundUp $value.diff_kcache_writes 2 }} bytes<br/>{{ RoundUp $value.per_sec_kcache_writes 2 }} bytes/sec<br/>{{ RoundUp $value.per_call_kcache_writes 2 }} bytes/call |
{{- RoundUp $value.diff_kcache_user_time_ms 2 }} ms<br/>{{ RoundUp $value.per_sec_kcache_user_time_ms 2 }} ms/sec<br/>{{ RoundUp $value.per_call_kcache_user_time_ms 2 }} ms/call |
{{- RoundUp $value.diff_kcache_system_time_ms 2 }} ms<br/>{{ RoundUp $value.per_sec_kcache_system_time_ms 2 }} ms/sec<br/>{{ RoundUp $value.per_call_kcache_system_time_ms 2 }} ms/call |
{{- Nobr (LimitStr $value.query 512 ) }}
{{ end }}{{/* range */}}
{{ else }}{{/* if .host.master*/}}
No data
{{ end }}{{/* if .host.master*/}}

{{ if gt (len .hosts.replicas) 0 }}
### Replica servers: ###
{{ range $key, $host := .hosts.replicas }}
#### Replica (`{{ $host }}`) ####
{{ if (index $.results $host) }}
Snapshot start time: {{ (index (index (index $.results $host) "data") "start_timestamptz") }}  
Snapshot end time: {{ (index (index (index $.results $host) "data") "end_timestamptz") }}  
Snapshot period: {{ (index (index (index $.results $host) "data") "period_age") }}  

\# | Calls | Total&nbsp;time | Rows | shared_blks_hit | shared_blks_read | shared_blks_dirtied | shared_blks_written | blk_read_time | blk_write_time | kcache_reads | kcache_writes | kcache_user_time_ms | kcache_system_time | Query
----|-------|------------|------|-----------------|------------------|---------------------|---------------------|---------------|----------------|--------------|---------------|---------------------|--------------------|------- 
{{ range $i, $key := (index (index (index (index $.results $host) "data") "queries") "_keys") }}
{{- $value := (index (index (index (index $.results $host) "data") "queries") $key) -}}
{{- $key}} |
{{- RoundUp $value.diff_calls 2 }}<br/>{{ RoundUp $value.per_sec_calls 2 }}/sec<br/>{{ RoundUp $value.per_call_calls 2 }}/call |
{{- RoundUp $value.diff_total_time 2 }} ms<br/>{{ RoundUp $value.per_sec_total_time 2 }} ms/sec<br/>{{ RoundUp $value.per_call_total_time 2 }} ms/call |
{{- RoundUp $value.diff_rows 2 }}<br/>{{ RoundUp $value.per_sec_rows 2 }}/sec<br/>{{ RoundUp $value.per_call_rows 2 }}/call |
{{- RoundUp $value.diff_shared_blks_hit 2 }} blks<br/>{{ RoundUp $value.per_sec_shared_blks_hit 2 }} blks/sec<br/>{{ RoundUp $value.per_call_shared_blks_hit 2 }} blks/call |
{{- RoundUp $value.diff_shared_blks_read 2 }} blks<br/>{{ RoundUp $value.per_sec_shared_blks_read 2 }} blks/sec<br/>{{ RoundUp $value.per_call_shared_blks_read 2 }} blks/call |
{{- RoundUp $value.diff_shared_blks_dirtied 2 }} blks<br/>{{ RoundUp $value.per_sec_shared_blks_dirtied 2 }} blks/sec<br/>{{ RoundUp $value.per_call_shared_blks_dirtied 2 }} blks/call |
{{- RoundUp $value.diff_shared_blks_written 2 }} blks<br/>{{ RoundUp $value.per_sec_shared_blks_written 2 }} blks/sec<br/>{{ RoundUp $value.per_call_shared_blks_written 2 }} blks/call |
{{- RoundUp $value.diff_blk_read_time 2 }} ms<br/>{{ RoundUp $value.per_sec_blk_read_time 2 }} ms/sec<br/>{{ RoundUp $value.per_call_blk_read_time 2 }} ms/call |
{{- RoundUp $value.diff_blk_write_time 2 }} ms<br/>{{ RoundUp $value.per_sec_blk_write_time 2 }} ms/sec<br/>{{ RoundUp $value.per_call_blk_write_time 2 }} ms/call |
{{- RoundUp $value.diff_kcache_reads 2 }} bytes<br/>{{ RoundUp $value.per_sec_kcache_reads 2 }} bytes/sec<br/>{{ RoundUp $value.per_call_kcache_reads 2 }} bytes/call |
{{- RoundUp $value.diff_kcache_writes 2 }} bytes<br/>{{ RoundUp $value.per_sec_kcache_writes 2 }} bytes/sec<br/>{{ RoundUp $value.per_call_kcache_writes 2 }} bytes/call |
{{- RoundUp $value.diff_kcache_user_time_ms 2 }} ms<br/>{{ RoundUp $value.per_sec_kcache_user_time_ms 2 }} ms/sec<br/>{{ RoundUp $value.per_call_kcache_user_time_ms 2 }} ms/call |
{{- RoundUp $value.diff_kcache_system_time_ms 2 }} ms<br/>{{ RoundUp $value.per_sec_kcache_system_time_ms 2 }} ms/sec<br/>{{ RoundUp $value.per_call_kcache_system_time_ms 2 }} ms/call |
{{- Nobr (LimitStr $value.query 512 ) }}
{{ end }}{{/* range */}}
{{- else -}}{{/* if host data */}}
No data
{{- end -}}{{/* if host data */}}
{{- end -}}{{/* hosts range */}}
{{- end -}}{{/* if replicas */}}

## Conclusions ##


## Recommendations ##

