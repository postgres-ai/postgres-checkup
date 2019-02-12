# {{ .checkId }} Disk usage and file system type
Output of `df -TPh` (follows symlinks)

## Observations ##
{{ if .hosts.master }}
{{ if (index .results .hosts.master) }}
### Master (`{{.hosts.master}}`) ###
Name | FS Type | Size | Available | Use | Used | Mount Point | Path | Device
-----|---------|------|-----------|-----|------|-------------|------|-------
{{ range $i, $name := (index (index (index .results .hosts.master) "data") "_keys") -}}
{{ $name }} {{ range $k, $val_name := (index (index (index (index $.results $.hosts.master) "data") $name) "_keys") -}}
 | {{ (index (index (index (index $.results $.hosts.master) "data") $name) $val_name) }} {{ end }}{{/* end of range $k, $val_name */}}
{{ end }}{{/* end of range $i, $name := */}}
{{ end }}{{/* end of if .hosts.master data */}}
{{ end }}{{/* end of if .hosts.master */}}

{{ if gt (len .hosts.replicas) 0 }}
### Replica servers: ###
  {{ range $skey, $host := .hosts.replicas }}
#### Replica (`{{ $host }}`) ####
Name | FS Type | Size | Available | Use | Used | Mount Point | Path | Device
-----|---------|------|-----------|-----|------|-------------|------|-------
{{- if (index $.results $host) }}
{{ range $i, $name := (index (index (index $.results $host) "data") "_keys") -}}
{{ $name }} {{ range $k, $val_name := (index (index (index (index $.results $host) "data") $name) "_keys") -}}
 | {{ (index (index (index (index $.results $host) "data") $name) $val_name) }}{{ " " }}
{{- end }} {{/* range $k, $val_name : */}}
{{ end }}{{/* if (index $.results $host) */}}
{{ end }}{{/* range $i, $name := */}}
{{ end }}{{/* range $skey, $host := */}}
{{ end }}{{/* if gt (len .hosts.replicas) 0 */}}

## Conclusions ##

## Recommendations ##

