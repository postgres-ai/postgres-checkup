# {{ .checkId }} Unused and redundant indexes #

## Observations ##
Data collected: {{ DtFormat .timestamptz }}  
Current database: {{ .database }}  

{{- if and (index .results .hosts.master) (index (index .results .hosts.master) "data") }}  
Stats reset: {{ (index (index (index .results .hosts.master) "data") "database_stat").stats_age }} ago ({{ DtFormat (index (index (index .results .hosts.master) "data") "database_stat").stats_reset }})  
{{- if le (Int (index (index (index .results .hosts.master) "data") "database_stat").days) 30 }}  
:warning: Statistics age is less than 30 days. Make decisions on index cleanup with caution!
{{- end }}
### Never Used Indexes ###
\#| Table | Index | {{.hosts.master}} usage {{ range $skey, $host := .hosts.replicas }}| {{ $host }} usage {{ end }}| &#9660;&nbsp;Index size | Table size
--|-------|-------|----{{ range $skey, $host := .hosts.replicas }}|--------{{ end }}|-----|-----
&nbsp;|=====TOTAL=====||{{ range $skey, $host := .hosts.replicas }}|{{ end }}|{{ ByteFormat ((index (index (index .results .hosts.master) "data") "never_used_indexes_total").index_size_bytes_sum) 2 }}|{{ ByteFormat ((index (index (index .results .hosts.master) "data") "never_used_indexes_total").table_size_bytes_sum) 2 }}
{{ range $i, $key := (index (index (index (index .results .hosts.master) "data") "never_used_indexes") "_keys") }}
{{- $value:=(index (index (index (index $.results $.hosts.master) "data") "never_used_indexes") $key) -}}
{{- $value.num}}|
{{- $value.formated_table_name}}|
{{- $value.index_name}}|
{{- RawIntFormat $value.idx_scan }}{{ range $skey, $host := $.hosts.replicas }}|{{ if (index (index (index (index $.results $host) "data") "never_used_indexes") $key) }}{{ RawIntFormat ((index (index (index (index $.results $host) "data") "never_used_indexes") $key).idx_scan) }}{{end}}{{ end }}|
{{- ByteFormat $value.index_size_bytes 2}}|
{{- ByteFormat $value.table_size_bytes 2}}
{{ end }}{{/* range */}}

{{- if and (index (index (index .results .hosts.master) "data") "rarely_used_indexes") (index (index (index .results .hosts.master) "data") "rarely_used_indexes_total") }}
### Rarely Used Indexes ###
\#| Table | Index | {{.hosts.master}} usage {{ range $skey, $host := .hosts.replicas }}| {{ $host }} usage {{ end }}| &#9660;&nbsp;Index size | Table size | Comment
--|-------|-------|-----{{ range $skey, $host := .hosts.replicas }}|--------{{ end }}|-----|-----|----
&nbsp;|=====TOTAL=====||{{ range $skey, $host := .hosts.replicas }}|{{ end }}|{{ ByteFormat ((index (index (index .results .hosts.master) "data") "rarely_used_indexes_total").index_size_bytes_sum) 2 }}|{{ ByteFormat ((index (index (index .results .hosts.master) "data") "rarely_used_indexes_total").table_size_bytes_sum) 2 }}|
{{ range $i, $key := (index (index (index (index .results .hosts.master) "data") "rarely_used_indexes") "_keys") }}
{{- $value:=(index (index (index (index $.results $.hosts.master) "data") "rarely_used_indexes") $key) -}}
{{- $value.num}}|
{{- $value.formated_table_name}}|
{{- $value.index_name}}|
{{- "scans:" }} {{ RawIntFormat $value.idx_scan }}\/hour, writes: {{ RawIntFormat $value.writes }}\/hour{{ range $skey, $host := $.hosts.replicas }}|{{ if (index (index (index (index $.results $host) "data") "rarely_used_indexes") $key) }}scans: {{ RawIntFormat ((index (index (index (index $.results $host) "data") "rarely_used_indexes") $key).idx_scan) }}\/hour, writes: {{ RawIntFormat ((index (index (index (index $.results $host) "data") "rarely_used_indexes") $key).writes) }}\/hour{{end}}{{ end }}|
{{- ByteFormat $value.index_size_bytes 2}}|
{{- ByteFormat $value.table_size_bytes 2}}|
{{- $value.reason}}
{{ end }}{{/* range */}}
{{ end }}{{/* rarely used indexes found */}}

{{- if and (index (index (index .results .hosts.master) "data") "redundant_indexes") (index (index (index .results .hosts.master) "data") "redundant_indexes_total") -}}
### Redundant indexes ###
\#| Table | Index | Redundant to |{{.hosts.master}} usage {{ range $skey, $host := .hosts.replicas }}| {{ $host }} usage {{ end }}| &#9660;&nbsp;Index size | Table size
--|-------|-------|--------------|--{{ range $skey, $host := .hosts.replicas }}|--------{{ end }}|-----|-----
&nbsp;|=====TOTAL=====|||{{ range $skey, $host := .hosts.replicas }}|{{ end }}|{{ ByteFormat ((index (index (index .results .hosts.master) "data") "redundant_indexes_total").index_size_bytes_sum) 2 }}|{{ ByteFormat ((index (index (index .results .hosts.master) "data") "redundant_indexes_total").table_size_bytes_sum) 2 }}
{{ range $i, $key := (index (index (index (index .results .hosts.master) "data") "redundant_indexes") "_keys") }}
{{- $value:=(index (index (index (index $.results $.hosts.master) "data") "redundant_indexes") $key) -}}
{{- $value.num}}|
{{- $value.formated_table_name}}|
{{- $value.index_name}}|
{{- $rinexes := Split $value.reason ", " -}}{{ range $r, $rto:= $rinexes }}{{ $ridx := (Replace $rto "redundant to index: " "") }}`{{$ridx}}`<br/>{{end}}|
{{- RawIntFormat $value.idx_scan }}{{ range $skey, $host := $.hosts.replicas }}|{{ if (index (index (index (index $.results $host) "data") "never_used_indexes") $key) }}{{ RawIntFormat ((index (index (index (index $.results $host) "data") "redundant_indexes") $key).idx_scan) }}{{end}}{{ end }}|
{{- ByteFormat $value.index_size_bytes 2}}|
{{- ByteFormat $value.table_size_bytes 2}}
{{ end }}{{/* range */}}
{{ end }}{{/* redundant indexes found */}}

{{- else -}}
No data
{{end}}{{/* end if master data*/}}

## Conclusions ##


## Recommendations ##
{{ if and (index (index .results .hosts.master) "data") (index (index (index .results .hosts.master) "data") "do") (index (index (index .results .hosts.master) "data") "undo") }}
#### "DO" database migration code ####
```
{{ range $i, $drop_code := (index (index (index .results .hosts.master) "data") "do") }}{{ $drop_code }}
{{ end }}
```

#### "UNDO" database migration code ####
```
{{ range $i, $revert_code := (index (index (index .results .hosts.master) "data") "undo") }}{{ $revert_code }}
{{ end }}
```
{{ end }}{{/* do and undo found */}}
