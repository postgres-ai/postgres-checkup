Current values
===

Available only raw data.

{{ range $host, $data := .rawData }}
    Host: {{ (index $data "host") }}
    ``` {{ (index $data "data") }} ```
{{ end }}



