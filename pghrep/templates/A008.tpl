# {{ .checkId }} Disk usage and file system type

## Observations ##
{{ if .hosts.master }}
### Master (`{{.hosts.master}}`) ###

name | path | device | fstype | size | used | avail | use_percent | mount_point
-----|------|--------|--------|------|------|-------|-------------|------------
{{ range $i, $name := (index (index (index .results .hosts.master) "data") "_keys") -}}
{{ $name }} {{ range $k, $val_name := (index (index (index (index $.results $.hosts.master) "data") $name) "_keys") -}}
 | {{ (index (index (index (index $.results $.hosts.master) "data") $name) $val_name) }} {{ end }}{{/* end of range $k, $val_name */}}
{{ end }}{{/* end of range $i, $name := */}}
{{ end }}{{/* end of if .hosts.master */}}
