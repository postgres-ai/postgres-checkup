Current values
===

Master DB server is {{.hosts.master}}
  System information
  `
  {{(index .results .hosts.master).system.raw}}
  `
  Cpu information 
  `
  {{(index .results .hosts.master).cpu.raw}}
  `
  Memory information
  `
  {{(index .results .hosts.master).ram.raw}}
  `
  Disk information
  `
  {{(index .results .hosts.master).disk.raw}}
  `
  Virtualization information
  `
  {{(index .results .hosts.master).virtualization.raw}}
  `

{{ if gt (len .hosts.replicas) 0 }}
Slave DB servers:
    {{ range $key, $value := .hosts.replicas }}
  DB slave server: {{ $value }}
       {{ if (index $.results $value) }}
    System information
      `
{{(index $.results $.value).system.raw}}
      `
    Cpu information 
      `
{{(index $.results $.value).cpu.raw}}
      `
    Memory information
      `
{{(index $.results $.value).ram.raw}}
      `
    Disk information
      `
{{(index $.results $.value).disk.raw}}
      `
    Virtualization information
      `
  {{(index $.results $.value).virtualization.raw}}
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
