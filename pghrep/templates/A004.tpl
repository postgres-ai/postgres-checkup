# {{ .checkId }} Cluster Information #

## Observations ##
Data collected: {{ DtFormat .timestamptz }}  

|&#9660;&nbsp;Indicator | {{.reorderedHosts.master}} {{ range $skey, $host := .reorderedHosts.replicas }}| {{ $host }} {{ end }}|
|--------|-------{{ range $skey, $host := .reorderedHosts.replicas }}|-------- {{ end }}|
{{ range $i, $key := (index (index (index (index .results .reorderedHosts.master) "data") "general_info") "_keys") }}
    {{- $value := (index (index (index (index $.results $.reorderedHosts.master) "data") "general_info") $key) -}}
    |{{ $key }} | 
    {{- Nobr (index $value "value") }}
    {{- range $skey, $host := $.reorderedHosts.replicas }}| 
        {{- if (index $.results $host) }}
            {{- if (index (index $.results $host) "data") }}
                {{- if (index (index (index $.results $host) "data") "general_info") }}
                    {{- (index (index (index (index $.results $host) "data") "general_info") $key).value }}
                {{- end}}
            {{- end}}
        {{- end}}
    {{- end }}|
{{ end }}

{{ if .reorderedHosts.master }}
{{- if (index .results .reorderedHosts.master) -}}
{{- if (index (index .results .reorderedHosts.master) "data") -}}
### Databases sizes ###

| Database | &#9660;&nbsp;Size |
|----------|--------|
{{ range $i, $key := (index (index (index (index .results .reorderedHosts.master) "data") "database_sizes") "_keys") }}
{{- $value := (index (index (index (index $.results $.reorderedHosts.master) "data") "database_sizes") $key) -}}
| `{{ $key }}` | {{ ByteFormat $value 2 }} |
{{ end }}
{{- end -}}
{{- end -}}
{{ end }}

## Conclusions ##


## Recommendations ##

