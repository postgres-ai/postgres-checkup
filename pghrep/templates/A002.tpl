# {{ .checkId }} Version information #

## Observations ##
Data collected: {{ DtFormat .timestamptz }}  
{{ if .hosts.master }}
{{ if and (index .results .hosts.master) (index (index .results .hosts.master) "data") }}
### Master (`{{.hosts.master}}`) ###

```
{{ (index (index .results .hosts.master) "data").version }}
```
{{ end }}{{/*master data*/}}
{{ end }}{{/*master*/}}

{{ if gt (len .hosts.replicas) 0 }}
### Replica servers: ###
{{ range $key, $value := .hosts.replicas }}
#### Replica (`{{ $value }}`) ####
{{ if (index $.results $value) }}

```
{{ (index (index $.results $value) "data").version }}
```
{{ else }}
No data
{{ end}}{{ end }}{{ end }}

## Conclusions ##


## Recommendations ##

