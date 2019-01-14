# {{ .checkId }} Unused/Rarely Used Indexes #

## Observations ##
{{ if .resultData }}

    {{- if .resultData.indexes -}}
### Never Used Indexes ###
Index | {{.hosts.master}} {{ range $skey, $host := .hosts.replicas }}| {{ $host }} {{ end }}| Usage
--------|-------{{ range $skey, $host := .hosts.replicas }}|--------{{ end }}|-----
{{ range $key, $value := (index .resultData "indexes") }}
{{- if ne $key "_keys" -}}
{{- if eq $value.master.reason "Never Used Indexes" -}}
{{- if $value.usage -}}
{{- else -}}
{{ $key }} | {{ $value.master.idx_scan }}{{ range $skey, $host := $.hosts.replicas }}|{{ (index $value $host).idx_scan }}{{- end -}} | {{ if $value.usage }} Used{{ else }}Not used {{ end }}
{{/* new line */}}
{{- end -}}
{{- end -}}
{{- end }}
{{- end }}

### Other unused indexes ###
Index | Reason |{{.hosts.master}} {{ range $skey, $host := .hosts.replicas }}| {{ $host }} {{ end }}| Usage
------|--------|-------{{ range $skey, $host := .hosts.replicas }}|--------{{ end }}|-----
{{ range $key, $value := (index .resultData "indexes") }}
{{- if ne $key "_keys" -}}
{{- if ne $value.master.reason "Never Used Indexes" -}}
{{ $key }} | {{ $value.master.reason }} | Index&nbsp;size:{{ Nobr $value.master.index_size }} Table&nbsp;size:{{ Nobr $value.master.table_size }} {{ range $skey, $host := $.hosts.replicas }} |Index&nbsp;size:{{ Nobr (index $value $host).index_size }} Table&nbsp;size:{{ Nobr (index $value $host).table_size }}{{- end -}} | {{ if $value.usage }} Used{{ else }}Not used {{ end }}
{{/* new line */}}
{{- end -}}
{{- end }}
{{- end }}
{{ end }}

## Conclusions ##


## Recommendations ##
{{ if .resultData.dropCode }}
#### Drop code ####
```
{{ range $i, $drop_code := (index .resultData  "dropCode") }}{{ $drop_code }}
{{ end }}
```
{{ end }}

{{ if .resultData.revertCode }}
#### Revert code ####
```
{{ range $i, $revert_code := (index .resultData  "revertCode") }}{{ $revert_code }}
{{ end }}
```
{{ end }}


{{- else -}}
No data
{{- end }}
