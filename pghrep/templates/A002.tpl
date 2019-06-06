# {{ .checkId }} Version Information #

## Observations ##
Data collected: {{ DtFormat .timestamptz }}  
{{ if .hosts.master }}
{{ if (index .results .hosts.master) }}
{{ if (index (index .results .hosts.master) "data") }}
### Master (`{{.hosts.master}}`) ###

```
{{ (index (index .results .hosts.master) "data").version }}
```
{{ end }}{{/*master data*/}}
{{ end }}{{/*master results*/}}
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
Nothing found
{{ end}}{{ end }}{{ end }}

## Conclusions ##

{{- if .processed }}
 {{- if .conclusions }}
  {{ range $conclusion := .conclusions -}}
   - {{ $conclusion.Message }}
  {{ end }}
 {{else}}
 {{end}}
{{ end }}

## Recommendations ##

{{- if .processed }}
 {{- if .recommendations }}
  {{ range $recommendation := .recommendations -}}
   - {{ $recommendation.Message }}
  {{ end }}
 {{else}}
  All good, no recommendations here.
 {{end}}
{{ end }}
