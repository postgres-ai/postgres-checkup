# Memory-related settings #

## Current values ##

### Master DB server is `{{.hosts.master}}` ###

Setting name | Value | Unit
-------------|-------|------
{{ range $key, $value := (index (index .results .hosts.master) "data") }}{{$key}}|{{ $value.setting}}|{{ $value.unit }}
{{ end }}
{{ if gt (len .hosts.replicas) 0 }}
### Slave DB servers: ###
  {{ range $skey, $host := .hosts.replicas }}
#### DB slave server: `{{ $host }}` ####
    {{ if (index $.results $host) }}
Setting name | Value | Unit
-------------|-------|------
{{ range $key, $value := (index (index $.results $host) "data") }}{{$key}}|{{ $value.setting}}|{{ $value.unit }}
{{ end }}{{ else }}
      No data
{{ end}}{{ end }}{{ end }}

## Conclusions ##

{{.Conclusion}}

## Recommendations ##

{{.Recommended}}
