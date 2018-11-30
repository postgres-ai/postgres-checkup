Current values
===

System information

Master DB server is {{.hosts.master}}
  System information
  `
  {{(index (index .results .hosts.master) "data").system.raw}}
  `
  Cpu information 
  `
  {{(index (index .results .hosts.master) "data").cpu.raw}}
  `
  Memory information
  `
  {{(index (index .results .hosts.master) "data").ram.raw}}
  `
  Disk information
  `
  {{(index (index .results .hosts.master) "data").disk.raw}}
  `
  Virtualization information
  `
  {{(index (index .results .hosts.master) "data").virtualization.raw}}
  `

{{ if gt (len .hosts.replicas) 0 }}
Slave DB servers:
    {{ range $key, $value := .hosts.replicas }}
  DB slave server: {{ $value }}
       {{ if (index $.results $value) }}
    System information
      `
{{ (index (index $.results $value) "data").system.raw }}
      `
    Cpu information 
      `
{{(index (index $.results $value) "data").cpu.raw}}
      `
    Memory information
      `
{{(index (index $.results $value) "data").ram.raw}}
      `
    Disk information
      `
{{(index (index $.results $value) "data").disk.raw}}
      `
    Virtualization information
      `
  {{(index (index $.results $value) "data").virtualization.raw}}
      `      
      {{ else }}
      No data
      {{ end}}
  {{ end }}
{{ end }}

Conclusions
===
{{.Conclusion}}

Recommendations
===
{{.Recommended}}
