# {{ .checkId }} System Information #

## Observations ##
Data collected: {{ DtFormat .timestamptz }}  

{{ if gt (len .results) 2 }} {{/* Min 2 hosts + "_keys" item */}}
### Operating System by hosts ###

| Host| Operating System | Kernel |
|-----|------------------|--------|
{{- if (index .results .hosts.master) }}
{{- if (index (index .results .hosts.master) "data") }}
|{{ .hosts.master }}|
{{- if index (index (index (index .results .hosts.master) "data").virtualization) "Operating System" }}
{{- (index (index (index (index .results .hosts.master) "data").virtualization) "Operating System") }}
{{- else }}{{- if index (index (index (index .results .hosts.master) "data").system) "operating_system" }}{{- index (index (index (index .results .hosts.master) "data").system) "operating_system" }}{{- end }}{{- end }}|
{{- if index (index (index (index .results .hosts.master) "data").virtualization) "Kernel" }}{{- (index (index (index (index .results .hosts.master) "data").virtualization) "Kernel") }}{{- else }}{{- if index (index (index (index .results .hosts.master) "data").system) "kernel_release" }}{{- (index (index (index (index .results .hosts.master) "data").system) "kernel_release") }}{{- end }}{{ end}}|
{{- end -}}{{/* master data */}}
{{- end -}}{{/* master */}}
{{- if gt (len .hosts.replicas) 0 -}}
    {{- range $key, $host := .hosts.replicas -}}
        {{ if (index $.results $host) }}
|{{ $host }}|
{{- if and (index $.results $host) (index (index $.results $host) "data") }}{{- if (index (index (index (index $.results $host) "data").virtualization) "Operating System") }}{{- (index (index (index (index $.results $host) "data").virtualization) "Operating System") }}{{ else }}{{- if index (index (index (index $.results $host) "data").system) "operating_system" }}{{- index (index (index (index $.results $host) "data").system) "operating_system" }}{{- end }}{{- end -}}|
{{- if (index (index (index (index $.results $host) "data").virtualization) "Kernel") -}}{{- (index (index (index (index $.results $host) "data").virtualization) "Kernel") }}{{- else -}}{{- if index (index (index (index $.results $host) "data").system) "kernel_release" }}{{- (index (index (index (index $.results $host) "data").system) "kernel_release") }}{{- end -}}{{- end -}}|
{{- end -}}
        {{- end -}}
    {{- end }}
{{ end }}
{{ end }}

{{ if .hosts.master }}
{{ if (index .results .hosts.master) }}
{{ if (index (index .results .hosts.master) "data") }}
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
{{ (index (index .results .hosts.master) "data").virtualization.raw }}
```
{{ end }}{{/* virtualization */}}
{{ end }}{{/* master data */}}
{{ end }}{{/* master results */}}
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
{{ (index (index $.results $value) "data").virtualization.raw }}
```
{{ end }}
        {{ else }}
`Nothing found`
{{ end}}{{ end }}{{ end }}

## Conclusions ##


## Recommendations ##
