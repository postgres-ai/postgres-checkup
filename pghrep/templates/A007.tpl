# {{ .checkId }} Altered Settings #

## Observations ##
Data collected: {{ DtFormat .timestamptz }}  
{{ if .hosts.master }}
{{ if (index .results .hosts.master) }}
{{ if (index (index .results .hosts.master) "data") }}
### Master (`{{.hosts.master}}`) ###
| Source | Settings count | Changed settings |
|-------|----------------|-----------------|
{{ range $key, $value := (index (index (index .results .hosts.master) "data") "changes") }}| {{ if $value.sourcefile }}{{ $value.sourcefile }}{{ else}}default{{ end }} | {{ $value.count }} | {{ if $value.examples}} {{ if (gt (len $value.examples) 0) }}{{ range $skey, $sname := (index $value "examples") }}{{ $sname }} {{ end }} {{ end }}
{{ end }}|{{ end }}
{{ end }}
{{ end }}
{{ end }}

{{ if gt (len .hosts.replicas) 0 }}
### Replica servers: ###
  {{ range $skey, $host := .hosts.replicas }}
#### Replica (`{{ $host }}`) ####
    {{ if (index $.results $host) }}
| Source | Settings count | Changed settings |
|-------|----------------|-----------------|
{{ range $key, $value := (index (index (index $.results $host) "data") "changes") }}|{{ if $value.sourcefile }}{{ $value.sourcefile }}{{ else}}default{{ end }} | {{ $value.count }} | {{ if $value.examples}} {{ if (gt (len $value.examples) 0) }}{{ range $skey, $sname := (index $value "examples") }}{{ $sname }} {{ end }} {{ end }}
{{ end }}|{{ end }}
{{ else }}
No data
{{ end}}{{ end }}{{ end }}

## Conclusions ##


## Recommendations ##

