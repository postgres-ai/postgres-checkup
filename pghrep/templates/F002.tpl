# Index bloat #

## Current values ##

### Master DB server is `{{.hosts.master}}` ###
{{ range $key, $value := (index (index .results .hosts.master) "data") }} {{ $tableIndex := Split $key "\n" }} {{ $table := Trim (index $tableIndex 1) " ()"}}
Index (Table): {{ (index $tableIndex 0) }} ({{ $table }})  Size: {{ ( index $value "Size") }}  Extra: {{ ( index $value "Extra") }}  Bloat: {{ ( index $value "Bloat") }}  Live: {{ ( index $value "Live") }}  Fill factor: {{ ( index $value "fillfactor") }} 
{{ end }}

{{ if gt (len .hosts.replicas) 0 }}
### Slave DB servers: ###
  {{ range $skey, $host := .hosts.replicas }}
#### DB slave server: `{{ $host }}` ####
    {{ if (index $.results $host) }} {{ range $key, $value := (index (index $.results $host) "data") }} {{ $tableIndex := Split $key "\n" }} {{ $table := Trim (index $tableIndex 1) " ()"}} 
Index (Table): {{ (index $tableIndex 0) }} ({{ $table }})  Size: {{ ( index $value "Size") }}  Extra: {{ ( index $value "Extra") }}  Bloat: {{ ( index $value "Bloat") }}  Live: {{ ( index $value "Live") }}  Fill factor: {{ ( index $value "fillfactor") }} 
    {{ end }} {{ else }}
No data
    {{ end}}
  {{ end }}
{{ end }}

## Conclusions ##

{{.Conclusion}}

## Recommendations ##

{{.Recommended}}
