# {{ .checkId }} System information #

## Observations ##
{{ if .hosts.master }}
### Master (`{{.hosts.master}}`) ###
{{ if (index (index .results .hosts.master) "data").system.raw}}
**System**

```
{{ (index (index .results .hosts.master) "data").system.raw}}
```
{{ end }}
{{ if (index (index .results .hosts.master) "data").cpu.raw}}
**CPU**

```
{{ (index (index .results .hosts.master) "data").cpu.raw }}
```
{{ end }}
{{ if (index (index .results .hosts.master) "data").ram.raw}}
**Memory**

```
{{ (index (index .results .hosts.master) "data").ram.raw }}
```
{{ end }}
{{ if (index (index .results .hosts.master) "data").disk.raw }}
**Disk**

```
{{ (index (index .results .hosts.master) "data").disk.raw}}
```
{{ end }}
{{ if (index (index .results .hosts.master) "data").virtualization.raw }}
**Virtualization**

```
{{ (index (index .results .hosts.master) "data").virtualization.raw }}
```
{{ end }}
{{ end }}
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

