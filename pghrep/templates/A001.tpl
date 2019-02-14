# {{ .checkId }} System information #

## Observations ##
Data collected: {{ DtFormat .timestamptz }}  
{{ if .hosts.master }}
{{ if and (index .results .hosts.master) (index (index .results .hosts.master) "data") }}
### Master (`{{.hosts.master}}`) ###
{{ if (index (index .results .hosts.master) "data").system.raw}}
**System**

```
{{ (index (index .results .hosts.master) "data").system.raw}}
```
{{ end }}{{/* system */}}
{{ if (index (index .results .hosts.master) "data").cpu.raw}}
**CPU**

```
{{ (index (index .results .hosts.master) "data").cpu.raw }}
```
{{ end }}{{/* cpu */}}
{{ if (index (index .results .hosts.master) "data").ram.raw}}
**Memory**

```
{{ (index (index .results .hosts.master) "data").ram.raw }}
```
{{ end }}{{/* memory */}}
{{ if (index (index .results .hosts.master) "data").disk.raw }}
**Disk**

```
{{ (index (index .results .hosts.master) "data").disk.raw}}
```
{{ end }}{{/* disk */}}
{{ if (index (index .results .hosts.master) "data").virtualization.raw }}
**Virtualization**

```
{{ (index (index .results .hosts.master) "data").virtualization.raw }}
```
{{ end }}{{/* virtualization */}}
{{ end }}{{/* master data */}}
{{ end }}{{/* master */}}

{{ if gt (len .hosts.replicas) 0 }}
### Replica servers: ###
    {{ range $key, $value := .hosts.replicas }}
#### Replica (`{{ $value }}`) ####
        {{ if (index $.results $value) }}
{{ if (index (index $.results $value) "data").system.raw}}
**System**

```
{{ (index (index $.results $value) "data").system.raw}}
```
{{ end }}
{{ if (index (index $.results $value) "data").cpu.raw}}
**CPU**

```
{{ (index (index $.results $value) "data").cpu.raw }}
```
{{ end }}
{{ if (index (index $.results $value) "data").ram.raw}}
**Memory**

```
{{ (index (index $.results $value) "data").ram.raw }}
```
{{ end }}
{{ if (index (index $.results $value) "data").disk.raw }}
**Disk**

```
{{ (index (index $.results $value) "data").disk.raw }}
```
{{ end }}
{{ if (index (index $.results $value) "data").virtualization.raw }}
**Virtualization**

```
{{ (index (index $.results $value) "data").virtualization.raw }}
```
{{ end }}
        {{ else }}
`No data`
{{ end}}{{ end }}{{ end }}

## Conclusions ##


## Recommendations ##

