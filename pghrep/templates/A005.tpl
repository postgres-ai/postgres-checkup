# {{ .checkId }} Extensions #

## Observations ##

{{ if .hosts.master }}
### Master (`{{.hosts.master}}`) ###
Database | Extension name | Installed version | Default version | Is old
---------|----------------|-------------------|-----------------|--------
{{ range $d, $db := (index (index (index .results .hosts.master) "data") "_keys") -}}
{{- $dbData := (index (index (index $.results $.hosts.master) "data") $db) -}}
{{- range $de, $dbext := (index $dbData "_keys") -}}
{{- $extData := (index $dbData $dbext) -}}
{{ $db }} | {{ $dbext }} | {{ $extData.installed_version }} | {{ $extData.default_version }} | {{ $extData.is_old }}
{{ end -}}
{{ end -}}
{{ else }}
Extensions information not found
{{ end }}

{{/* force empty line */}}

## Conclusions ##


## Recommendations ##
