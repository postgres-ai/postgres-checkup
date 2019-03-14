# {{ .checkId }} Disk usage and file system type
Output of `df -TPh` (follows symlinks)

## Observations ##
Data collected: {{ DtFormat .timestamptz }}  
{{ if .hosts.master }}
    {{ if (index .results .hosts.master) }}
        {{ if (index (index .results .hosts.master) "data") }}
### Master (`{{.hosts.master}}`) ###

#### System directories ####
Device | FS Type | Size | Available | Usage | Used | Mount Point 
-------|---------|------|-----------|-----|------|-------------
{{ range $i, $name := (index (index (index (index .results .hosts.master) "data") "fs_data") "_keys") -}}
    {{- $value := (index (index (index (index $.results $.hosts.master) "data") "fs_data") $name) -}}
    {{ $value.device}}|
    {{- $value.fstype}}|
    {{- $value.size}}|
    {{- $value.avail}}|
    {{- $value.use_percent}}|
    {{- $value.used}}|
    {{- $value.mount_point}}
{{ end }}{{/* end of range $i, $name := */}}

#### Database directories ####
Name | FS Type | Size | Available | Usage | Used | Mount Point | Path | Device
-----|---------|------|-----------|-----|------|-------------|------|-------
{{ range $i, $name := (index (index (index (index .results .hosts.master) "data") "db_data") "_keys") -}}
    {{- $value := (index (index (index (index $.results $.hosts.master) "data") "db_data") $name) -}}
    {{ $name }}|
    {{- $value.fstype}}|
    {{- $value.size}}|
    {{- $value.avail}}|
    {{- $value.use_percent}}|
    {{- $value.used}}|
    {{- $value.mount_point}}|
    {{- $value.path}}|
    {{- $value.device}}
{{ end }}{{/* end of range $i, $name := */}}

        {{ end }}{{/* end of if .hosts.master data */}}
    {{ end }}{{/* end of if .results .hosts.master */}}
{{ end }}{{/* end of if .hosts.master */}}

{{ if gt (len .hosts.replicas) 0 }}
### Replica servers: ###
    {{ range $skey, $host := .hosts.replicas }}
        {{- if (index $.results $host) }}
#### Replica (`{{ $host }}`) ####

#### System directories ####
Device | FS Type | Size | Available | Usage | Used | Mount Point 
-------|---------|------|-----------|-----|------|-------------
{{ range $i, $name := (index (index (index (index $.results $host) "data") "fs_data") "_keys") -}}
    {{- $value := (index (index (index (index $.results $host) "data") "fs_data") $name) -}}
    {{ $value.device}}|
    {{- $value.fstype}}|
    {{- $value.size}}|
    {{- $value.avail}}|
    {{- $value.use_percent}}|
    {{- $value.used}}|
    {{- $value.mount_point}}
{{ end }}{{/* range $i, $name := */}}

#### Database directories ####
Name | FS Type | Size | Available | Usage | Used | Mount Point | Path | Device
-----|---------|------|-----------|-----|------|-------------|------|-------
{{ range $i, $name := (index (index (index (index $.results $host) "data") "db_data") "_keys") -}}
    {{- $value := (index (index (index (index $.results $host) "data") "db_data") $name) -}}
    {{ $name }}|
    {{- $value.fstype}}|
    {{- $value.size}}|
    {{- $value.avail}}|
    {{- $value.use_percent}}|
    {{- $value.used}}|
    {{- $value.mount_point}}|
    {{- $value.path}}|
    {{- $value.device}}
{{ end }}{{/* range $i, $name := */}}

        {{ end }}{{/* if (index $.results $host) */}}
    {{ end }}{{/* if (index $.results $host) */}}
{{ end }}{{/* if gt (len .hosts.replicas) 0 */}}

## Conclusions ##

## Recommendations ##
