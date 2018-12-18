# {{ .checkId }} Heap bloat #

## Observations ##

### Master (`{{.hosts.master}}`) ###
{{ if (index (index .results .hosts.master) "data") }}
 Table | Size | Extra | Bloat | Live | Last vacuum
-------|------|-------|-------|------|-------------
{{ range $i, $key := (index (index (index .results .hosts.master) "data") "_keys") }}
{{- $value := (index (index (index $.results $.hosts.master) "data") $key) -}}
{{ $key }} | {{ ( index $value "Size") }} | {{ ( index $value "Extra") }} | {{ ( index $value "Bloat") }} | {{ ( index $value "Live") }} | {{ if (index $value "Last Vaccuum") }} {{ ( index $value "Last Vaccuum") }} {{ end }}
{{ end }}
{{- else -}}
`No data`
{{- end -}}

{{ if gt (len .hosts.replicas) 0 }}
### Replica servers: ###
  {{ range $skey, $host := .hosts.replicas }}
#### Replica (`{{ $host }}`) ####
{{ if (index $.results $host) }}
{{- if (index (index $.results $host) "data") -}}
 Table | Size | Extra | Bloat | Live | Last vacuum
-------|------|-------|-------|------|-------------
{{ range $i, $key := (index (index (index $.results $host) "data") "_keys") }}
{{- $value := (index (index (index $.results $host) "data") $key) -}}
{{ $key }} | {{ ( index $value "Size") }} | {{ ( index $value "Extra") }} | {{ ( index $value "Bloat") }} | {{ ( index $value "Live") }} | {{ if (index $value "Last Vaccuum") }} {{ ( index $value "Last Vaccuum") }} {{ end }}
{{ end }}
{{- else -}}
`No data`
{{- end -}}
{{- else -}}
`No data`
{{- end -}}{{- end -}}{{ end }}

## Conclusions ##


## Recommendations ##

