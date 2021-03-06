# {{ .checkId }} Invalid Indexes #

## Observations ##
Data collected: {{ DtFormat .timestamptz }}  
Current database: {{ .database }}  
{{ if .hosts.master }}
{{ if (index .results .hosts.master) }}
### Master (`{{.hosts.master}}`) ###
{{ if (index (index .results .hosts.master) "data") }}
{{ if (index (index (index .results .hosts.master) "data") "invalid_indexes") }}
{{ if ge (len (index (index (index .results .hosts.master) "data") "invalid_indexes")) .LISTLIMIT }}The list is limited to {{.LISTLIMIT}} items. Total: {{ Sub (len (index (index (index .results .hosts.master) "data") "invalid_indexes")) 1 }}.{{ end }}  

| \# | Table | Index name | Index size | Supports FK |
|---|-------|------------|------------|----------|
&nbsp;|=====TOTAL=====||{{- ByteFormat (index (index (index (index $.results $.hosts.master) "data") "invalid_indexes_total") "index_size_bytes_sum") 2 }} ||
{{ range $i, $key := (index (index (index (index .results .hosts.master) "data") "invalid_indexes") "_keys") }}
{{- if lt $i $.LISTLIMIT -}}
    {{- $value := (index (index (index (index $.results $.hosts.master) "data") "invalid_indexes") $key) -}}
    | {{ $value.num }} |`{{ $value.formated_relation_name }}` | `{{ $value.formated_index_name }}` |
    {{- ByteFormat $value.index_size_bytes 2 }} |
    {{- if $value.supports_fk }}Yes{{ end }} |
{{/* if limit list */}}{{ end -}}
{{ end }}{{/* range */}}
{{- else -}}
Invalid indexes not found
{{- end -}}{{/* if data */}}
{{- else -}}
Invalid indexes not found
{{- end -}}{{/* if data */}}
{{- else -}}
Nothing found
{{- end -}}{{/* if .host.master data */}}
{{- else -}}
Nothing found
{{- end -}}{{/* if .host.master */}}


## Conclusions ##


## Recommendations ##

{{- if .hosts.master }}
{{- if (index .results .hosts.master) }}
{{- if (index (index (index .results .hosts.master) "data") "invalid_indexes") }}
#### Rebuild invalid indexes ####
```
-- Call each line separately. "CONCURRENTLY" queries cannot be
-- combined in multi-statement requests.

{{ range $i, $key := (index (index (index (index .results .hosts.master) "data") "invalid_indexes") "_keys") }}
{{- $value := (index (index (index (index $.results $.hosts.master) "data") "invalid_indexes") $key) -}}
{{ $value.drop_code }}
{{ $value.revert_code }}

{{ end }}
```
{{- end -}}{{/* if data */}}
{{- end -}}{{/* if data */}}
{{- end -}}{{/* if .host.master */}}
