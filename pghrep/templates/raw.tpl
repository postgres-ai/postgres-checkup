Some collected
data with help of checker

{{ range $host, $data := .rawData }}

Host: {{ (index $data "host") }}

```bash
{{ (index $data "data") }}
```

{{ end }}

