# {{ .checkId }} Table Sizes #

## Observations ##
Data collected: {{ DtFormat .timestamptz }}  
Current database: {{ .database }}  
{{ if .hosts.master }}
{{ if (index .results .hosts.master) }}
{{ if (index (index .results .hosts.master) "data") }}
### Master (`{{.hosts.master}}`) ###
{{ if gt (len (index (index (index .results .hosts.master) "data") "tables_data")) .LISTLIMIT }}The list is limited to {{.LISTLIMIT}} items. All items {{(len (index (index (index .results .hosts.master) "data") "tables_data"))}}.{{ end }}  

| \# | Table | Rows | &#9660;&nbsp;Total size | Table size | Index(es) Size | TOAST Size |
|---|---|------|------------|------------|----------------|------------|
|&nbsp;|===== TOTAL ===== |~{{- NumFormat (index (index (index (index $.results $.hosts.master) "data") "tables_data_total") "row_estimate_sum" ) 0 }} |
{{- ByteFormat (index (index (index (index $.results $.hosts.master) "data") "tables_data_total") "total_size_bytes_sum" ) 2 }} |
{{- ByteFormat (index (index (index (index $.results $.hosts.master) "data") "tables_data_total") "table_size_bytes_sum" ) 2 }} |
{{- ByteFormat (index (index (index (index $.results $.hosts.master) "data") "tables_data_total") "indexes_size_bytes_sum" ) 2 }} |
{{- ByteFormat (index (index (index (index $.results $.hosts.master) "data") "tables_data_total") "toast_size_bytes_sum" ) 2 }} |
{{ range $i, $key := (index (index (index (index .results .hosts.master) "data") "tables_data") "_keys") }}
{{- if lt $i $.LISTLIMIT -}}
{{- $value := (index (index (index (index $.results $.hosts.master) "data") "tables_data") $key) -}}
|{{ $value.num}} |`{{ index $value "table" }}` | ~{{ NumFormat (index $value "row_estimate") 0 }} |
{{- ByteFormat (index $value "total_size_bytes") 2 }} ({{ (RawFloatFormat (index $value "total_size_percent") 2) }}%) |
{{- ByteFormat (index $value "table_size_bytes") 2 }} ({{ (RawFloatFormat (index $value "table_size_percent") 2) }}%) |
{{- ByteFormat (index $value "indexes_size_bytes") 2 }} ({{ (RawFloatFormat (index $value "indexes_size_percent") 2) }}%) |
{{- if gt (Int (index $value "toast_size_bytes")) 0 }}{{ ByteFormat (index $value "toast_size_bytes") 2 }} ({{ (RawFloatFormat (index $value "toast_size_percent") 2) }}%){{end}} |
{{/* if limit list */}}{{ end -}}
{{ end }}
{{- else -}}{{/*Master data*/}}
Nothing found
{{- end }}{{/*Master data*/}}
{{- else -}}{{/*Master results*/}}
Nothing found
{{- end }}{{/*Master results*/}}
{{- else -}}{{/*Master*/}}
Nothing found
{{ end }}{{/*Master*/}}

## Conclusions ##


## Recommendations ##

