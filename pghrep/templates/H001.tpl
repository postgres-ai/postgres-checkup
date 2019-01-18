# {{ .checkId }} Invalid indexes #

## Observations ##

{{ if .hosts.master }}
### Master (`{{.hosts.master}}`) ###

{{ if (index (index .results .hosts.master) "data") }}
Num | Schema name | Table name | Index name 
----|-------------|------------|------------
{{ range $i, $key := (index (index (index .results .hosts.master) "data") "_keys") }}
    {{- $value := (index (index (index $.results $.hosts.master) "data") $key) -}}
    {{ $key }} |
    {{- $value.schema_name }} |
    {{- $value.table_name }} |
    {{- $value.index_name }} |
{{ end }}{{/* range */}}
{{- else -}}
Invalid indexes not found
{{- end -}}{{/* if data */}}
{{- end -}}{{/* if .host.master */}}

{{ if .hosts.replicas }}
### Replica servers: ###
{{ range $skey, $host := .hosts.replicas }}
#### Replica (`{{ $host }}`) ####
{{ if (index $.results $host) }}
{{- if (index (index $.results $host) "data") -}}
Num | Schema name | Table name | Index name 
----|-------------|------------|------------
{{ range $i, $key := (index (index (index $.results $host) "data") "_keys") }}
    {{- $value := (index (index (index $.results $host) "data") $key) -}}
    {{ $key }} |
    {{- $value.schema_name }} |
    {{- $value.table_name }} |
    {{- $value.index_name }} |
{{ end }}{{/* range */}}
{{ else }}
Invalid indexes not found
{{- end -}}{{/* if data */}}
{{- else -}}{{/* if $.results $host */}}
Invalid indexes not found
{{- end -}}{{/* if $.results $host */}}
{{- end -}}{{/* replicas range*/}}
{{- end -}}{{/* if replica */}}

## Conclusions ##


## Recommendations ##

{{ if (index .resultData "repair_code") }}
#### "DO" database migration code ####
```
-- Call each line separately. "CONCURRENTLY" queries cannot be
-- combined in multi-statement requests.
{{ range $i, $code := (index .resultData  "repair_code") }}
{{ $code.drop_code }}
{{ $code.revert_code }}
{{ end }}
```
{{ end }}
