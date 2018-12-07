# System information #

## Current values ##

### Master DB server is `{{.hosts.master}}` ###

**System information**

{{ Code ((index (index .results .hosts.master) "data").system.raw) false }}

**Cpu information**

{{ Code ((index (index .results .hosts.master) "data").cpu.raw) false }}

**Memory information**

{{ Code (index (index .results .hosts.master) "data").ram.raw false}}

**Disk information**

{{ Code (index (index .results .hosts.master) "data").disk.raw false }}

**Virtualization information**

{{ Code (index (index .results .hosts.master) "data").virtualization.raw false }}

{{ if gt (len .hosts.replicas) 0 }}
### Slave DB servers: ###
    {{ range $key, $value := .hosts.replicas }}
#### DB slave server: `{{ $value }}` ####

        {{ if (index $.results $value) }}
**System information**

{{ Code ((index (index $.results $value) "data").system.raw) false }}

**Cpu information**

{{ Code ((index (index $.results $value) "data").cpu.raw) false }}

**Memory information**

{{ Code (index (index $.results $value) "data").ram.raw false}}

**Disk information**

{{ Code (index (index $.results $value) "data").disk.raw false }}

**Virtualization information**

{{ Code (index (index $.results $value) "data").virtualization.raw false }}
        {{ else }}
`No data`
{{ end}}{{ end }}{{ end }}

## Conclusions ##

{{.Conclusion}}

## Recommendations ##

{{.Recommended}}
