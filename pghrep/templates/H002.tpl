# {{ .checkId }} Unused and Redundant Indexes #
## Observations ##
Data collected: {{ DtFormat .timestamptz }}  
Current database: {{ .database }}  

{{- if (index .results .hosts.master)}}
{{- if (index (index .results .hosts.master) "data") }}  
Stats reset: {{ (index (index (index .results .hosts.master) "data") "database_stat").stats_age }} ago ({{ DtFormat (index (index (index .results .hosts.master) "data") "database_stat").stats_reset }})  
{{- if le (Int (index (index (index .results .hosts.master) "data") "database_stat").days) 30 }}  
:warning: Statistics age is less than 30 days. Make decisions on index cleanup with caution!
{{- end }}
### Never Used Indexes ###
{{ if gt (len (index (index (index .results .hosts.master) "data") "never_used_indexes")) .ROWS_LIMIT }}The list is limited to {{.ROWS_LIMIT}} items.{{ end }}  

\#| Table | Index | {{.hosts.master}} usage {{ range $skey, $host := .hosts.replicas }}| {{ $host }} usage {{ end }}| &#9660;&nbsp;Index size | Table size | Supports FK
--|-------|-------|----{{ range $skey, $host := .hosts.replicas }}|--------{{ end }}|-----|-----|-----
&nbsp;|=====TOTAL=====||{{ range $skey, $host := .hosts.replicas }}|{{ end }}|{{ ByteFormat ((index (index (index .results .hosts.master) "data") "never_used_indexes_total").index_size_bytes_sum) 2 }}|{{ ByteFormat ((index (index (index .results .hosts.master) "data") "never_used_indexes_total").table_size_bytes_sum) 2 }}|
{{ range $i, $key := (index (index (index (index .results .hosts.master) "data") "never_used_indexes") "_keys") }}
{{- $value:=(index (index (index (index $.results $.hosts.master) "data") "never_used_indexes") $key) -}}
{{- $value.num}}|
{{- $value.formated_relation_name}}|
{{- $value.formated_index_name}}|
{{- RawIntFormat $value.idx_scan }}{{ range $skey, $host := $.hosts.replicas }}|{{ if (index $.results $host) }}{{- if (index (index (index (index $.results $host) "data") "never_used_indexes") $key) }}{{ RawIntFormat ((index (index (index (index $.results $host) "data") "never_used_indexes") $key).idx_scan) }}{{end}}{{ end }}{{end}}|
{{- ByteFormat $value.index_size_bytes 2}}|
{{- ByteFormat $value.table_size_bytes 2}}|
{{- if $value.supports_fk }}Yes{{end}}
{{ end }}{{/* range */}}

{{- if and (index (index (index .results .hosts.master) "data") "rarely_used_indexes") (index (index (index .results .hosts.master) "data") "rarely_used_indexes_total") }}
### Rarely Used Indexes ###
{{ if gt (len (index (index (index .results .hosts.master) "data") "rarely_used_indexes")) .ROWS_LIMIT }}The list is limited to {{.ROWS_LIMIT}} items.{{ end }}  

\#| Table | Index | {{.hosts.master}} usage {{ range $skey, $host := .hosts.replicas }}| {{ $host }} usage {{ end }}| &#9660;&nbsp;Index size | Table size | Comment | Supports FK
--|-------|-------|-----{{ range $skey, $host := .hosts.replicas }}|--------{{ end }}|-----|-----|----|-----
&nbsp;|=====TOTAL=====||{{ range $skey, $host := .hosts.replicas }}|{{ end }}|{{ ByteFormat ((index (index (index .results .hosts.master) "data") "rarely_used_indexes_total").index_size_bytes_sum) 2 }}|{{ ByteFormat ((index (index (index .results .hosts.master) "data") "rarely_used_indexes_total").table_size_bytes_sum) 2 }}||
{{ range $i, $key := (index (index (index (index .results .hosts.master) "data") "rarely_used_indexes") "_keys") }}
{{- $value:=(index (index (index (index $.results $.hosts.master) "data") "rarely_used_indexes") $key) -}}
{{- $value.num}}|
{{- $value.formated_relation_name}}|
{{- $value.formated_index_name}}|
{{- "scans:" }} {{ RawIntFormat $value.idx_scan }}\/hour, writes: {{ RawIntFormat $value.writes }}\/hour{{ range $skey, $host := $.hosts.replicas }}|{{ if (index $.results $host) }}{{ if (index (index (index (index $.results $host) "data") "rarely_used_indexes") $key) }}scans: {{ RawIntFormat ((index (index (index (index $.results $host) "data") "rarely_used_indexes") $key).idx_scan) }}\/hour, writes: {{ RawIntFormat ((index (index (index (index $.results $host) "data") "rarely_used_indexes") $key).writes) }}\/hour{{end}}{{ end }}{{end}}|
{{- ByteFormat $value.index_size_bytes 2}}|
{{- ByteFormat $value.table_size_bytes 2}}|
{{- $value.reason}}|
{{- if $value.supports_fk }}Yes{{end}}
{{ end }}{{/* range */}}
{{ end }}{{/* rarely used indexes found */}}

{{- if and (index (index (index .results .hosts.master) "data") "redundant_indexes") (index (index (index .results .hosts.master) "data") "redundant_indexes_total") -}}
### Redundant indexes ###
{{ if gt (len (index (index (index .results .hosts.master) "data") "redundant_indexes")) .ROWS_LIMIT }}The list is limited to {{.ROWS_LIMIT}} items.{{ end }}  

\#| Table | Index | Redundant to |{{.hosts.master}} usage {{ range $skey, $host := .hosts.replicas }}| {{ $host }} usage {{ end }}| &#9660;&nbsp;Index size | Table size | Supports FK
--|-------|-------|--------------|--{{ range $skey, $host := .hosts.replicas }}|--------{{ end }}|-----|-----|-----
&nbsp;|=====TOTAL=====|||{{ range $skey, $host := .hosts.replicas }}|{{ end }}|{{ ByteFormat ((index (index (index .results .hosts.master) "data") "redundant_indexes_total").index_size_bytes_sum) 2 }}|{{ ByteFormat ((index (index (index .results .hosts.master) "data") "redundant_indexes_total").table_size_bytes_sum) 2 }}|
{{ range $i, $key := (index (index (index (index .results .hosts.master) "data") "redundant_indexes") "_keys") }}
{{- $value:=(index (index (index (index $.results $.hosts.master) "data") "redundant_indexes") $key) -}}
{{- $value.num}}|
{{- $value.formated_relation_name}}|
{{- $value.formated_index_name}}|
{{- $rinexes := Split $value.reason ", " -}}{{ range $r, $rto:= $rinexes }}{{$rto}}<br/>{{end}}|
{{- RawIntFormat $value.idx_scan }}{{ range $skey, $host := $.hosts.replicas }}|{{ if (index $.results $host) }}{{ if (index (index (index (index $.results $host) "data") "never_used_indexes") $key) }}{{ RawIntFormat ((index (index (index (index $.results $host) "data") "redundant_indexes") $key).idx_scan) }}{{end}}{{ end }}{{end}}|
{{- ByteFormat $value.index_size_bytes 2}}|
{{- ByteFormat $value.table_size_bytes 2}}|
{{- if $value.supports_fk }}Yes{{end}}
{{ end }}{{/* range */}}
{{ end }}{{/* redundant indexes found */}}

{{- else -}}{{/* end if master*/}}
No data
{{end}}{{/* end if master*/}}
{{- else -}}{{/* end if master data*/}}
No data
{{end}}{{/* end if master data*/}}

## Conclusions ##


## Recommendations ##
{{ if (index .results .hosts.master) }}
{{ if (index (index .results .hosts.master) "data") }}
{{ if (index (index (index .results .hosts.master) "data") "do")}}
#### "DO" database migration code ####
```
{{ range $i, $drop_code := (index (index (index .results .hosts.master) "data") "do") }}{{ $drop_code }}
{{ end }}
```
{{end}}
{{ if (index (index (index .results .hosts.master) "data") "undo") }}
#### "UNDO" database migration code ####
```
{{ range $i, $revert_code := (index (index (index .results .hosts.master) "data") "undo") }}{{ $revert_code }}
{{ end }}
```
{{end}}
{{ end }}{{/* data found */}}
{{ end }}{{/* master */}}
