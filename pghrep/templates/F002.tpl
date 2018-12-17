# {{ .checkId }} Index bloat #

## Observations ##

### Master (`{{.hosts.master}}`) ###
 Index (Table) | Size | Extra | Bloat | Live | Fill factor
---------------|------|-------|-------|------|-------------
{{ range $i, $key := (index (index (index .results .hosts.master) "data") "_keys") }} {{ $tableIndex := Split $key "\n" }} {{ $table := Trim (index $tableIndex 1) " ()"}}{{ (index $tableIndex 0) }} ({{ $table }}){{ $value := (index (index (index $.results $.hosts.master) "data") $key) }} | {{ ( index $value "Size") }} | {{ ( index $value "Extra") }} | {{ ( index $value "Bloat") }} | {{ ( index $value "Live") }} | {{ ( index $value "fillfactor") }}
{{ end }}

{{ if gt (len .hosts.replicas) 0 }}
### Replica servers: ###
  {{ range $skey, $host := .hosts.replicas }}
#### Replica (`{{ $host }}`) ####
    {{ if (index $.results $host) }}
 Index (Table) | Size | Extra | Bloat | Live | Fill factor
---------------|------|-------|-------|------|-------------
{{ range $i, $key := (index (index (index $.results $host) "data") "_keys") }} {{ $tableIndex := Split $key "\n" }} {{ $table := Trim (index $tableIndex 1) " ()"}}{{ (index $tableIndex 0) }} ({{ $table }}) {{ $value := (index (index (index $.results $host) "data") $key) }}| {{ ( index $value "Size") }}{{ $value := (index (index (index $.results $host) "data") $key) }} | {{ ( index $value "Extra") }} | {{ ( index $value "Bloat") }} | {{ ( index $value "Live") }} | {{ ( index $value "fillfactor") }}
{{ end }}
{{- else -}}
No data
{{- end -}}{{ end }}{{ end }}

## Conclusions ##


## Recommendations ##

