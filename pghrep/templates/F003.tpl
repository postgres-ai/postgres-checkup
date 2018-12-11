# Autovacuum info #

## Current values ##

### Master DB server is `{{.hosts.master}}` ###
Setting name | Value | Unit
-------------|-------|------
{{ range $key, $value := (index (index (index .results .hosts.master) "data") "settings") }}{{$key}}|{{ $value.setting}}|{{ $value.unit }}
{{ end }}
#### iotop information ####
Command: `{{ (index (index (index .results .hosts.master) "data") "iotop").cmd }}`

Result:
```
{{ (index (index (index .results .hosts.master) "data") "iotop").data }}
```

{{ if gt (len .hosts.replicas) 0 }}
### Slave DB servers: ###
  {{ range $skey, $host := .hosts.replicas }}
#### DB slave server: `{{ $host }}` ####
    {{ if (index $.results $host) }}

Setting name | Value | Unit
-------------|-------|------
{{ range $key, $value := (index (index (index $.results $host) "data") "settings") }}{{$key}}|{{ $value.setting}}|{{ $value.unit }}
{{ end }}
#### iotop information ####
Command: `{{ (index (index (index $.results $host) "data") "iotop").cmd }}`

Result: 
```
{{ (index (index (index $.results $host) "data") "iotop").data }}
```
    {{ else}}
      No data
    {{ end }}
  {{ end }}
{{ end }}

{{/* ## Conclusions ##

{{.Conclusion}}

## Recommendations ##

{{.Recommended}}
*/}}