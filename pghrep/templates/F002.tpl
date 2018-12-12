# {{ .checkId }} Index bloat #

## Observations ##

### Master (`{{.hosts.master}}`) ###
 Index (Table) | Size | Extra | Bloat | Live | Fill factor
---------------|------|-------|-------|------|-------------
{{ range $key, $value := (index (index .results .hosts.master) "data") }} {{ $tableIndex := Split $key "\n" }} {{ $table := Trim (index $tableIndex 1) " ()"}}{{ (index $tableIndex 0) }} ({{ $table }}) | {{ ( index $value "Size") }} | {{ ( index $value "Extra") }} | {{ ( index $value "Bloat") }} | {{ ( index $value "Live") }} | {{ ( index $value "fillfactor") }}
{{ end }}

{{ if gt (len .hosts.replicas) 0 }}
### Replica servers: ###
  {{ range $skey, $host := .hosts.replicas }}
#### Replica (`{{ $host }}`) ####
    {{ if (index $.results $host) }}
 Index (Table) | Size | Extra | Bloat | Live | Fill factor
---------------|------|-------|-------|------|-------------
{{ range $key, $value := (index (index $.results $host) "data") }} {{ $tableIndex := Split $key "\n" }} {{ $table := Trim (index $tableIndex 1) " ()"}}{{ (index $tableIndex 0) }} ({{ $table }}) | {{ ( index $value "Size") }} | {{ ( index $value "Extra") }} | {{ ( index $value "Bloat") }} | {{ ( index $value "Live") }} | {{ ( index $value "fillfactor") }}
{{ end }}
{{ else }}
No data
{{ end}}{{ end }}{{ end }}

## Conclusions ##


## Recommendations ##

