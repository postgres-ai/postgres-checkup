# {{ .checkId }} Non indexed foreign keys (or with bad indexes) #

## Observations ##
{{ if .hosts.master }}
### Master (`{{.hosts.master}}`) ###
{{ if (index (index .results .hosts.master) "data")}}
Num | Schema name | Table name | FK name | Issue | Table mb | writes | Table scans | Parent name | Parent mb | Parent writes | Cols list | Indexdef
----|-------------|------------|---------|-------|----------|--------|-------------|-------------|-----------|---------------|-----------|----------
{{ range $i, $key := (index (index (index .results .hosts.master) "data") "_keys") }}
    {{- $value := (index (index (index $.results $.hosts.master) "data") $key) -}}
    {{ $key }} |
    {{- $value.schema_name }} |
    {{- $value.table_name }} |
    {{- $value.fk_name }} |
    {{- $value.issue }} |
    {{- $value.table_mb }} |
    {{- $value.writes }} |
    {{- $value.table_scans }} |
    {{- $value.parent_name }} |
    {{- $value.parent_mb}} |
    {{- NumFormat $value.parent_writes -1 }} |
    {{- $value.cols_list }} |
    {{- $value.indexdef }}
{{ end }}{{/* range */}}
{{ else }}
No data
{{- end -}}{{/* if data */}}
{{ end }}{{/* if .host.master */}}

{{- if gt (len .hosts.replicas) 0 -}}
### Replica servers: ###
{{ range $skey, $host := .hosts.replicas }}
#### Replica (`{{ $host }}`) ####
{{- if (index $.results $host) }}
{{ if (index (index $.results $host) "data")}}
Num | Schema name | Table name | FK name | Issue | Table mb | writes | Table scans | Parent name | Parent mb | Parent writes | Cols list | Indexdef
----|-------------|------------|---------|-------|----------|--------|-------------|-------------|-----------|---------------|-----------|----------
{{ range $i, $key := (index (index (index $.results $host) "data") "_keys") }}
    {{- $value := (index (index (index $.results $host) "data") $key) -}}
    {{ $key }} |
    {{- $value.schema_name }} |
    {{- $value.table_name }} |
    {{- $value.fk_name }} |
    {{- $value.issue }} |
    {{- $value.table_mb }} |
    {{- $value.writes }} |
    {{- $value.table_scans }} |
    {{- $value.parent_name }} |
    {{- $value.parent_mb}} |
    {{- $value.parent_writes}} |
    {{- $value.cols_list }} |
    {{- $value.indexdef }}
{{ end }}{{/* range */}}
{{ end }}{{/* if data */}}
{{- else -}}{{/* if $.results $host */}}
No data
{{- end -}}{{/* if $.results $host */}}
{{- end -}}{{/* replicas range*/}}
{{- end -}}{{/* if replica */}}

## Conclusions ##


## Recommendations ##

