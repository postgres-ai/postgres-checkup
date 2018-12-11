# {{ .checkId }} Version information #

## Observations ##

### Master (`{{.hosts.master}}`) ###
```
{{ (index (index .results .hosts.master) "data").version }}
```

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

