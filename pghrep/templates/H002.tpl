# {{ .checkId }} Unused/Rarely Used Indexes #

## Observations ##
Data collected: {{ DtFormat .timestamptz }}  
Current database: {{ .database }}  
{{ if .resultData }}
{{ if .hosts.master }}
{{ if and (index .results .hosts.master) (index (index .results .hosts.master) "data") }}
Stats reset: {{ (index (index (index .results .hosts.master) "data") "database_stat").stats_age }} ago ({{ DtFormat (index (index (index .results .hosts.master) "data") "database_stat").stats_reset }})  

{{ if .resultData.unused_indexes }}
### Never Used Indexes ###
Index | {{.hosts.master}} usage {{ range $skey, $host := .hosts.replicas }}| {{ $host }} usage {{ end }}| &#9660;&nbsp;Index size | Usage
--------|-------{{ range $skey, $host := .hosts.replicas }}|--------{{ end }}|-----|-----
{{ range $i, $key := (index (index .resultData "unused_indexes") "_keys") }}
{{- $value := (index (index $.resultData "unused_indexes") $key) -}}
{{- if ne $key "_keys" -}}
{{- if eq $value.master.reason "Never Used Indexes" -}}
{{- if $value.usage -}}
{{- else -}}
{{- $key }} |
{{- $value.master.idx_scan }}{{ range $skey, $host := $.hosts.replicas }}|
{{- (index $value $host).idx_scan }}{{- end -}} |
{{- "Index&nbsp;size:"}}&nbsp;{{ Nobr $value.master.index_size }}<br/>Table&nbsp;size:&nbsp;{{ Nobr $value.master.table_size }} |
{{- if $value.usage }} Used{{ else }}Not used {{ end }}
{{/* new line */}}
{{- end -}}{{/* value.usage */}}
{{- end -}}{{/**/}}
{{- end }}{{/* in ! _keys */}}
{{- end }}{{/* range unused_indexes */}}

### Other unused indexes ###
Index | Reason |{{.hosts.master}} {{ range $skey, $host := .hosts.replicas }}| {{ $host }} {{ end }}| Usage
------|--------|-------{{ range $skey, $host := .hosts.replicas }}|--------{{ end }}|-----
{{ range $i, $key := (index (index .resultData "unused_indexes") "_keys") }}
{{- $value := (index (index $.resultData "unused_indexes") $key) -}}
{{- if ne $key "_keys" -}}
{{- if ne $value.master.reason "Never Used Indexes" -}}
{{ $key }} | {{ $value.master.reason }} | Usage:&nbsp;{{ $value.master.idx_scan }}<br/>Index&nbsp;size:{{ Nobr $value.master.index_size }}<br/>Table&nbsp;size:{{ Nobr $value.master.table_size }} {{ range $skey, $host := $.hosts.replicas }} | Usage:&nbsp;{{ (index $value $host).idx_scan }}<br/>Index&nbsp;size:{{ Nobr (index $value $host).index_size }}<br/>Table&nbsp;size:{{ Nobr (index $value $host).table_size }}{{- end -}} | {{ if $value.usage }} Used{{ else }}Not used {{ end }}
{{/* new line */}}
{{- end -}}
{{- end }}{{/* ! "_keys" */}}
{{- end }}{{/* range unused_indexes */}}
{{ end }}{{/* if unused_indexes */}}

{{- if .resultData.redundant_indexes -}}

### Redundant indexes ###

Index | {{.hosts.master}} usage {{ range $skey, $host := .hosts.replicas }}| {{ $host }} usage {{ end }}| Usage | Index size
--------|-------{{ range $skey, $host := .hosts.replicas }}|--------{{ end }}|-----|-----
{{ range $i, $key := (index (index .resultData "redundant_indexes") "_keys") }}
{{- $value := (index (index $.resultData "redundant_indexes") $key) -}}
{{- if ne $key "_keys" -}}
{{- if $value.usage -}}
{{- else -}}
{{ $key }} | {{ $value.master.index_usage }}{{ range $skey, $host := $.hosts.replicas }}|{{ (index $value $host).index_usage }}{{- end -}} | {{ if $value.usage }} Used{{ else }}Not used {{ end }} | {{ $value.master.index_size }}
{{/* new line */}}
{{- end -}}{{/* value.usage */}}
{{- end }}{{/* in ! _keys */}}
{{- end }}{{/* range redundant_indexes */}}
{{end}}{{/* if redundant_indexes */}}

{{end}}{{/* master data */}}
{{end}}{{/* master */}}

## Conclusions ##


## Recommendations ##
{{ if .resultData.drop_code }}
#### "DO" database migration code ####
```
{{ range $i, $drop_code := (index .resultData  "drop_code") }}{{ $drop_code }}
{{ end }}
```
{{ end }}

{{ if .resultData.revert_code }}
#### "UNDO" database migration code ####
```
{{ range $i, $revert_code := (index .resultData  "revert_code") }}{{ $revert_code }}
{{ end }}
```
{{ end }}

{{- else -}}
No data
{{- end }}
