# Differences settings #

## Current values ##

### Settings (pg_settings) that differ ###
{{ if (index .diffData "pg_settings") }}

{{ range $key, $value := (index .diffData "pg_settings") }}
Setting {{ $key }}: {{ range $key, $value := $value }} On {{ $key }}: `{{ index $value "value" }}` {{ if (index $value "unit") }}{{ index $value "unit" }}{{ end  }}  {{ end }}
{{ end }}{{end}}
{{ if (index .diffData "pg_configs") }}
Configs(pg_config) that differ
### Configs(pg_config) that differ ###
{{ range $key, $value := (index .diffData "pg_configs") }}
Config {{ $key }}: {{ range $key, $value := $value }} On {{ $key }}: `{{ index $value "value" }}` {{ if (index $value "unit") }}{{ index $value "unit" }}{{ end  }}{{ end }}
{{ end }}{{end}}

## Conclusions ##

{{.Conclusion}}

## Recommendations ##

{{.Recommended}}