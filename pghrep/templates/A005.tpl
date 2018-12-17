# {{ .checkId }} Extensions #

## Observations ##

### Master (`{{.hosts.master}}`) ###
Extension name | Installed version | Default version | Is old
---------------|-------------------|-----------------|--------
{{ range $i, $key := (index (index (index .results .hosts.master) "data") "_keys") }}
{{- $value := (index (index (index $.results $.hosts.master) "data") $key) -}}
{{ $key }} | {{ $value.installed_version }} | {{ $value.default_version }} | {{ $value.is_old }}
{{ end }}

## Conclusions ##


## Recommendations ##

