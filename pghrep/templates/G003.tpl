# {{ .checkId }} Timeouts, locks, deadlocks #

## Observations ##

### Master (`{{.hosts.master}}`) ###

#### Timeouts ####
Setting name | Value | Unit
-------------|-------|------
{{ range $i, $key := (index (index (index .results .hosts.master) "data") "timeouts_keys") }}[{{ $key }}](https://postgresqlco.nf/en/doc/param/{{ $key }}){{ $value := (index (index (index (index $.results $.hosts.master) "data") "timeouts") $key)}}|{{ $value.setting}}|{{ $value.unit }}
{{ end }}
#### Locks ####
Setting name | Value | Unit
-------------|-------|------
{{ range $i, $key := (index (index (index .results .hosts.master) "data") "locks_keys") }}[{{ $key }}](https://postgresqlco.nf/en/doc/param/{{ $key }}){{ $value := (index (index (index (index $.results $.hosts.master) "data") "locks") $key) }}|{{ $value.setting}}|{{ $value.unit }}
{{ end }}
#### Databases data ####
Database | Conflicts | Deadlocks | Stats reset at | Stat reset
-------------|-------|-----------|----------------|------------
{{ range $i, $key := (index (index (index .results .hosts.master) "data") "databases_stat_keys") }}{{$key}}{{ $value:= (index (index (index (index $.results $.hosts.master) "data") "databases_stat") $key) }}|{{ $value.conflicts}}|{{ $value.deadlocks }}|{{ $value.stats_reset }}|{{ $value.stats_reset_age }}
{{ end }}
{{ if gt (len .hosts.replicas) 0 }}
### Replica servers: ###
{{ range $skey, $host := .hosts.replicas }}
#### Replica (`{{ $host }}`) ####
{{ if (index $.results $host) }}
#### Timeouts ####
Setting name | Value | Unit
-------------|-------|------
{{ range $i, $key := (index (index (index $.results $host) "data") "timeouts_keys") }}[{{ $key }}](https://postgresqlco.nf/en/doc/param/{{ $key }}){{ $value := (index (index (index (index $.results $host) "data") "timeouts") $key)}}|{{ $value.setting}}|{{ $value.unit }}
{{ end }}
#### Locks ####
Setting name | Value | Unit
-------------|-------|------
{{ range $i, $key := (index (index (index $.results $host) "data") "locks_keys") }}[{{ $key }}](https://postgresqlco.nf/en/doc/param/{{ $key }}){{ $value := (index (index (index (index $.results $host) "data") "locks") $key)}}|{{ $value.setting}}|{{ $value.unit }}
{{ end }}
#### Databases data ####
Database | Conflicts | Deadlocks | Stats reset at | Stat reset
-------------|-------|-----------|----------------|------------
{{ range $i, $key := (index (index (index $.results $host) "data") "databases_stat_keys") }}{{$key}}{{ $value:= (index (index (index (index $.results $host) "data") "databases_stat") $key) }}|{{ $value.conflicts}}|{{ $value.deadlocks }}|{{ $value.stats_reset }}|{{ $value.stats_reset_age }}
{{ end }}
{{ else }}
No data
{{ end}}{{ end }}{{ end }}
## Conclusions ##


## Recommendations ##

