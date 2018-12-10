# Timeouts, locks, deadlocks #

## Current values ##

### Master DB server is `{{.hosts.master}}` ###

#### Timeouts ####
Setting name | Value | Unit
-------------|-------|------
{{ range $key, $value := (index (index (index .results .hosts.master) "data") "timeouts") }}{{$key}}|{{ $value.setting}}|{{ $value.unit }}
{{ end }}
#### Locks ####
Setting name | Value | Unit
-------------|-------|------
{{ range $key, $value := (index (index (index .results .hosts.master) "data") "locks") }}{{$key}}|{{ $value.setting}}|{{ $value.unit }}
{{ end }}
#### Databases data ####
Database | Conflicts | Deadlocks | Stats_reset
-------------|-------|-----------|-------------
{{ range $key, $value := (index (index (index .results .hosts.master) "data") "dbsdata") }}{{$key}}|{{ $value.conflicts}}|{{ $value.deadlocks }}|{{ $value.stats_reset }}
{{ end }}
{{ if gt (len .hosts.replicas) 0 }}
### Slave DB servers: ###
{{ range $skey, $host := .hosts.replicas }}
#### DB slave server: `{{ $host }}` ####
{{ if (index $.results $host) }}
#### Timeouts ####
Setting name | Value | Unit
-------------|-------|------
{{ range $key, $value := (index (index (index $.results $host) "data") "timeouts") }}{{$key}}|{{ $value.setting}}|{{ $value.unit }}
{{ end }}
#### Locks ####
Setting name | Value | Unit
-------------|-------|------
{{ range $key, $value := (index (index (index $.results $host) "data") "locks") }}{{$key}}|{{ $value.setting}}|{{ $value.unit }}
{{ end }}
#### Databases data ####
Database | Conflicts | Deadlocks | Stats_reset
-------------|-------|-----------|-------------
{{ range $key, $value := (index (index (index $.results $host) "data") "dbsdata") }}{{$key}}|{{ $value.conflicts}}|{{ $value.deadlocks }}|{{ $value.stats_reset }}
{{ end }}
{{ else }}
No data
{{ end}}{{ end }}{{ end }}
## Conclusions ##

{{.Conclusion}}

## Recommendations ##

{{.Recommended}}
