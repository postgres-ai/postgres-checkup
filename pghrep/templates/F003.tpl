# {{ .checkId }} Autovacuum info #

## Observations ##

### Master (`{{.hosts.master}}`) ###
Setting name | Value | Unit
-------------|-------|------
{{ range $key, $value := (index (index (index .results .hosts.master) "data") "settings") }}{{$key}}|{{ $value.setting}}|{{ $value.unit }}
{{ end }}
{{ if (index (index .results .hosts.master) "data").iotop }}
#### iotop information ####
Command: `{{ (index (index (index .results .hosts.master) "data") "iotop").cmd }}`

Result:
```
{{ (index (index (index .results .hosts.master) "data") "iotop").data }}
```
{{ end }}

{{ if gt (len .hosts.replicas) 0 }}
### Replica servers: ###
  {{ range $skey, $host := .hosts.replicas }}
#### Replica (`{{ $host }}`) ####
    {{ if (index $.results $host) }}

Setting name | Value | Unit
-------------|-------|------
{{ range $key, $value := (index (index (index $.results $host) "data") "settings") }}{{$key}}|{{ $value.setting}}|{{ $value.unit }}
{{ end }}
{{ if (index (index $.results $host) "data").iotop }}
#### iotop information ####
Command: `{{ (index (index (index $.results $host) "data") "iotop").cmd }}`

Result: 
```
{{ (index (index (index $.results $host) "data") "iotop").data }}
```
{{ end }}
    {{ else}}
      No data
{{ end }}{{ end }}{{ end }}

## Conclusions ##


## Recommendations ##

