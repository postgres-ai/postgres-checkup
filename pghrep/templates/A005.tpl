# {{ .checkId }} Extensions #

## Observations ##

### Master (`{{.hosts.master}}`) ###

Database | Extension name | Installed version | Default version | Is old 
---------|----------------|-------------------|-----------------|--------
{{ range $key, $value := (index (index .results .hosts.master) "data") -}}
{{ range $dkey, $dvalue := $value -}}
{{ $key }} | {{ $dkey }} | {{ $dvalue.installed_version }} | {{ $dvalue.default_version }} | {{ $dvalue.is_old }}
{{ end -}}
{{ end -}}

{{/* force empty line */}}

## Conclusions ##


## Recommendations ##

