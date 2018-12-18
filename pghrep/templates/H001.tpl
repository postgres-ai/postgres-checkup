# {{ .checkId }} Unused/Rarely Used Indexes #

## Observations ##

### Master (`{{.hosts.master}}`) ###

{{ if (index (index (index .results .hosts.master) "data") "indexes") -}}
#### Indexes ####

Index name | Reason | Scheme name | Table name | Index size | Table size
-----------|--------|-------------|------------|------------|------------
{{ range $i, $index_name := (index (index (index (index .results .hosts.master) "data") "indexes") "_keys") }}
{{- $index_data := (index (index (index (index $.results $.hosts.master) "data") "indexes") $index_name) -}}
{{ $index_name }} | {{ $index_data.reason }} | {{ $index_data.schemaname }} | {{ $index_data.tablename }} | {{ $index_data.index_size }} | {{ $index_data.table_size }}
{{ end }}
{{- else -}}
No data
{{- end }}

{{ if gt (len .hosts.replicas) 0 }}
### Replica servers: ###
{{ range $skey, $host := .hosts.replicas }}
#### Replica server: `{{ $host }}` ####
{{ if (index (index (index $.results $host) "data") "indexes") -}}
{{ if (index $.results $host) }}
#### Indexes ####

Index name | Reason | Scheme name | Table name | Index size | Table size
-----------|--------|-------------|------------|------------|------------
{{ range $i, $index_name := (index (index (index (index $.results $host) "data") "indexes") "_keys") }}
{{- $index_data := (index (index (index (index $.results $host) "data") "indexes") $index_name) -}}
{{ $index_name }} | {{ $index_data.reason }} | {{ $index_data.schemaname }} | {{ $index_data.tablename }} | {{ $index_data.index_size }} | {{ $index_data.table_size }}
{{ end }}
{{ end }}
{{- else }}
No data
{{- end -}}
{{- end -}}
{{ end }}

## Conclusions ##


## Recommendations ##

### Master (`{{.hosts.master}}`) ###
{{ if or (index (index (index .results .hosts.master) "data") "drop_code") (index (index (index .results .hosts.master) "data") "revert_code") -}}
{{ if (index (index (index .results .hosts.master) "data") "drop_code") }}
#### Drop code ####
```
{{ range $i, $drop_code := (index (index (index .results .hosts.master) "data") "drop_code") }}{{ $drop_code }}
{{ end }}
```
{{ end }}

{{- if (index (index (index .results .hosts.master) "data") "revert_code") -}}
#### Revert code ####
```
{{ range $i, $revert_code := (index (index (index .results .hosts.master) "data") "revert_code") }}{{ $revert_code }}
{{ end }}
```
    {{- end }}
{{ else }}
No recommendations
{{- end }}

{{ if gt (len .hosts.replicas) 0 }}
### Replica servers: ###
{{ range $skey, $host := .hosts.replicas }}
#### Replica server: `{{ $host }}` ####
{{ if or (index (index (index $.results $host) "data") "drop_code") (index (index (index $.results $host) "data") "revert_code") -}}
{{ if (index $.results $host) }}
{{ if (index (index (index $.results $host) "data") "drop_code") -}}
#### Drop code ####

```
{{ range $i, $drop_code := (index (index (index $.results $host) "data") "drop_code") }}{{ $drop_code }}
{{ end }}
```
{{- end }}

{{- if (index (index (index $.results $host) "data") "revert_code") -}}
{{/* blank  row */}}
#### Revert code ####

```
{{ range $i, $revert_code := (index (index (index $.results $host) "data") "revert_code") }}{{ $revert_code }}
{{ end }}
```
{{- end -}}
{{ end }}
{{- else }}
No recommendations
{{- end -}}
{{- end -}}
{{- end -}}
