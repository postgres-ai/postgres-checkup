# {{ .checkId }} Unused Indexes #
## Observations ##
Data collected: {{ DtFormat .timestamptz }}  
Current database: {{ .database }}  

{{ if (index .results .reorderedHosts.master)}}
{{- if (index (index .results .reorderedHosts.master) "data") }}  
Stats reset: {{ (index (index (index .results .reorderedHosts.master) "data") "database_stat").stats_age }} ago ({{ DtFormat (index (index (index .results .reorderedHosts.master) "data") "database_stat").stats_reset }})  
{{ if le (Int (index (index (index .results .reorderedHosts.master) "data") "database_stat").days) 30 }}  
:warning: Statistics age is less than 30 days. Make decisions on index cleanup with caution!
{{- end }}  
{{ if gt (Int (index (index (index .results .reorderedHosts.master) "data") "min_index_size_bytes")) 0 }}NOTICE: only indexes larger than {{ ByteFormat (index (index (index .results .reorderedHosts.master) "data") "min_index_size_bytes") 0 }} are analyzed.  {{end}}
### Never Used Indexes ###
{{ if (index (index (index .results .reorderedHosts.master) "data") "never_used_indexes") }}
{{ if ge (len (index (index (index .results .reorderedHosts.master) "data") "never_used_indexes")) .LISTLIMIT }}The list is limited to {{.LISTLIMIT}} items. Total: {{ Sub (len (index (index (index .results .reorderedHosts.master) "data") "never_used_indexes")) 1 }}.{{ end }}  

| \#| Table | Index | {{.reorderedHosts.master}} usage {{ range $skey, $host := .reorderedHosts.replicas }}| {{ $host }} usage {{ end }}| &#9660;&nbsp;Index size | Table size | Supports FK |
|--|-------|-------|----{{ range $skey, $host := .reorderedHosts.replicas }}|--------{{ end }}|-----|-----|-----|
|&nbsp;|=====TOTAL=====||{{ range $skey, $host := .reorderedHosts.replicas }}|{{ end }}|{{ ByteFormat ((index (index (index .results .reorderedHosts.master) "data") "never_used_indexes_total").index_size_bytes_sum) 2 }}|{{ ByteFormat ((index (index (index .results .reorderedHosts.master) "data") "never_used_indexes_total").table_size_bytes_sum) 2 }}||
{{ range $i, $key := (index (index (index (index .results .reorderedHosts.master) "data") "never_used_indexes") "_keys") }}
{{- if lt $i $.LISTLIMIT -}}
{{- $value:=(index (index (index (index $.results $.reorderedHosts.master) "data") "never_used_indexes") $key) -}}
| {{- $value.num}}|`{{ $value.formated_relation_name}}`|`{{ $value.formated_index_name}}`|
{{- RawIntFormat $value.idx_scan }}{{ range $skey, $host := $.reorderedHosts.replicas }}|{{ if (index $.results $host) }}{{ if (index (index $.results $host) "data") }}{{- if (index (index (index $.results $host) "data") "never_used_indexes") }}{{- if (index (index (index (index $.results $host) "data") "never_used_indexes") $key) }}{{ RawIntFormat ((index (index (index (index $.results $host) "data") "never_used_indexes") $key).idx_scan) }}{{end}}{{ end }}{{ end }}{{ end }}{{end}}|
{{- ByteFormat $value.index_size_bytes 2}}|
{{- ByteFormat $value.table_size_bytes 2}}|
{{- if $value.supports_fk }}Yes{{end}}|
{{/* if limit list */}}{{ end -}}
{{ end }}{{/* range */}}
{{else}}
Nothing found.
{{ end }}{{/*never used indexes found*/}}

### Rarely Used Indexes ###
{{ if (index (index (index .results .reorderedHosts.master) "data") "rarely_used_indexes") }}
{{ if ge (len (index (index (index .results .reorderedHosts.master) "data") "rarely_used_indexes")) .LISTLIMIT }}The list is limited to {{.LISTLIMIT}} items. Total: {{ Sub (len (index (index (index .results .reorderedHosts.master) "data") "rarely_used_indexes")) 1 }}.{{ end }}

|\#| Table | Index | {{.reorderedHosts.master}} usage {{ range $skey, $host := .reorderedHosts.replicas }}| {{ $host }} usage {{ end }}| &#9660;&nbsp;Index size | Table size | Comment | Supports FK|
|--|-------|-------|-----{{ range $skey, $host := .reorderedHosts.replicas }}|--------{{ end }}|-----|-----|----|-----|
{{ if (index (index (index .results .reorderedHosts.master) "data") "rarely_used_indexes_total") }}|&nbsp;|=====TOTAL=====||{{ range $skey, $host := .reorderedHosts.replicas }}|{{ end }}|{{ ByteFormat ((index (index (index .results .reorderedHosts.master) "data") "rarely_used_indexes_total").index_size_bytes_sum) 2 }}|{{ ByteFormat ((index (index (index .results .reorderedHosts.master) "data") "rarely_used_indexes_total").table_size_bytes_sum) 2 }}|||{{ end }}
{{ range $i, $key := (index (index (index (index .results .reorderedHosts.master) "data") "rarely_used_indexes") "_keys") }}
{{- if lt $i $.LISTLIMIT -}}
{{- $value:=(index (index (index (index $.results $.reorderedHosts.master) "data") "rarely_used_indexes") $key) -}}
|{{- $value.num}}|`{{ $value.formated_relation_name}}`|`{{ $value.formated_index_name}}`|
{{- "scans:" }} {{ RawIntFormat $value.idx_scan }}\/hour, writes: {{ RawIntFormat $value.writes }}\/hour{{ range $skey, $host := $.reorderedHosts.replicas }}|{{- if (index $.results $host) }}{{- if (index (index $.results $host) "data")}}{{- if (index (index (index $.results $host) "data") "rarely_used_indexes") }}{{- if (index (index (index (index $.results $host) "data") "rarely_used_indexes") $key) }}scans: {{ RawIntFormat ((index (index (index (index $.results $host) "data") "rarely_used_indexes") $key).idx_scan) }}\/hour, writes: {{ RawIntFormat ((index (index (index (index $.results $host) "data") "rarely_used_indexes") $key).writes) }}\/hour{{ end }}{{ end }}{{ end }}{{ end }}{{ end }}|
{{- ByteFormat $value.index_size_bytes 2}}|
{{- ByteFormat $value.table_size_bytes 2}}|
{{- $value.reason}}|
{{- if $value.supports_fk }}Yes{{end}}|
{{/* if limit list */}}{{ end -}}
{{ end }}{{/* range */}}
{{else}}
Nothing found.
{{ end }}{{/* rarely used indexes found */}}

{{- else -}}{{/* end if master*/}}
Nothing found
{{end}}{{/* end if master*/}}
{{- else -}}{{/* end if master data*/}}
Nothing found
{{end}}{{/* end if master data*/}}

## Conclusions ##

{{- if .processed }}
 {{- if .conclusions }}
  {{ range $conclusion := .conclusions -}}
   - {{ $conclusion.Message }}
  {{ end }}
 {{else}}
 {{end}}
{{ end }}

## Recommendations ##

{{- if .processed }}
 {{- if .recommendations }}
  {{ range $recommendation := .recommendations -}}
   - {{ $recommendation.Message }}
  {{ end }}
 {{else}}
  All good, no recommendations here.
 {{end}}
{{ end }}
