# {{ .checkId }} Heap bloat #

## Observations ##

### Master (`{{.hosts.master}}`) ###
 Table | Size | Extra | Bloat | Live | Last vacuum
-------|------|-------|-------|------|-------------
{{ range $key, $value := (index (index .results .hosts.master) "data") }}{{ $key }} | {{ ( index $value "Size") }} | {{ ( index $value "Extra") }} | {{ ( index $value "Bloat") }} | {{ ( index $value "Live") }} | {{ if (index $value "Last Vaccuum") }} {{ ( index $value "Last Vaccuum") }} {{ end }}
{{ end }}

{{ if gt (len .hosts.replicas) 0 }}
### Replica servers: ###
  {{ range $skey, $host := .hosts.replicas }}
#### Replica (`{{ $host }}`) ####
    {{ if (index $.results $host) }}
 Table | Size | Extra | Bloat | Live | Last vacuum
-------|------|-------|-------|------|-------------
{{ range $key, $value := (index (index $.results $host) "data") }}{{ $key }} | {{ ( index $value "Size") }} | {{ ( index $value "Extra") }} | {{ ( index $value "Bloat") }} | {{ ( index $value "Live") }} | {{ if (index $value "Last Vaccuum") }} {{ ( index $value "Last Vaccuum") }} {{ end }}
{{ end }}
{{ else }}
No data
{{ end}}{{ end }}{{ end }}

## Conclusions ##


## Recommendations ##

