Current values
===

Some collected
data with help of checker

{{ range $host, $data := .rawData }}
    Host: {{ (index $data "host") }}
    ``` {{ (index $data "data") }} ```
{{ end }}



