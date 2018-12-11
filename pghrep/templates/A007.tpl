# Altered settings #

## Current values ##

### Master DB server is `{{.hosts.master}}` ###
{{ range $key, $value := (index (index (index .results .hosts.master) "data") "changes") }}
{{ if $value.sourcefile }}**Source: `{{ $value.sourcefile }}`**{{ else}}**Source: DEFAULT**{{ end }} 

Settings count: `{{ $value.count }}`

{{ if $value.examples}} {{ if (gt (len $value.examples) 0) }}Changed settings:
  
{{ range $skey, $sname := (index $value "examples") }}    {{ $sname }} {{ end }}
    {{ end }} {{ end }}

{{ end }}

{{ if gt (len .hosts.replicas) 0 }}
### Slave DB servers: ###
  {{ range $skey, $host := .hosts.replicas }}
#### DB slave server: `{{ $host }}` ####
    {{ if (index $.results $host) }}  {{ range $key, $value := (index (index (index $.results $host) "data") "changes") }} 
{{ if $value.sourcefile }}**Source: {{ $value.sourcefile }}**{{ else}}**Source: DEFAULT**{{ end }}

Settings count: {{ $value.count }}
  {{ if $value.examples}} {{ if (gt (len $value.examples) 0) }}Changed settings:

{{ range $skey, $sname := (index $value "examples") }}     {{ $sname }} {{ end }}
    {{ end }} {{ end }}
{{ end }}
    {{ else }}
      No data
    {{ end}}
  {{ end }}
{{ end }}

## Conclusions ##

{{.Conclusion}}

## Recommendations ##

{{.Recommended}}
