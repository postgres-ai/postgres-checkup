# {{ .checkId }} Cluster Information #

## Observations ##
Data collected: {{ DtFormat .timestamptz }}  

|&#9660;&nbsp;Indicator | {{.hosts.master}} {{ range $skey, $host := .hosts.replicas }}| {{ $host }} {{ end }}|
|--------|-------{{ range $skey, $host := .hosts.replicas }}|-------- {{ end }}|
{{ range $i, $key := (index (index (index (index .results .hosts.master) "data") "general_info") "_keys") }}
    {{- $value := (index (index (index (index $.results $.hosts.master) "data") "general_info") $key) -}}
    |{{ $key }} | 
    {{- Nobr (index $value "value") }}
    {{- range $skey, $host := $.hosts.replicas }}| 
        {{- if (index $.results $host) }}
            {{- if (index (index $.results $host) "data") }}
                {{- if (index (index (index $.results $host) "data") "general_info") }}
                    {{- (index (index (index (index $.results $host) "data") "general_info") $key).value }}
                {{- end}}
            {{- end}}
        {{- end}}
    {{- end }}|
{{ end }}

{{ if .hosts.master }}
{{- if (index .results .hosts.master) -}}
{{- if (index (index .results .hosts.master) "data") -}}
### Databases sizes ###

| Database | &#9660;&nbsp;Size |
|----------|--------|
{{ range $i, $key := (index (index (index (index .results .hosts.master) "data") "database_sizes") "_keys") }}
{{- $value := (index (index (index (index $.results $.hosts.master) "data") "database_sizes") $key) -}}
| `{{ $key }}` | {{ ByteFormat $value 2 }} |
{{ end }}
{{- end -}}
{{- end -}}
{{ end }}

## Conclusions ##


## Recommendations ##

