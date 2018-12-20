# {{ .checkId }} Timeouts, locks, deadlocks #

## Observations ##
{{ if .hosts.master }}
### Master (`{{.hosts.master}}`) ###
#### Timeouts ####
Setting name | Value | Unit
-------------|-------|------
{{ range $i, $key := (index (index (index (index .results .hosts.master) "data") "timeouts") "_keys") }}
{{- $value := (index (index (index (index $.results $.hosts.master) "data") "timeouts") $key) -}}
[{{ $key }}](https://postgresqlco.nf/en/doc/param/{{ $key }})|{{ $value.setting}}|{{ $value.unit }}
{{ end }}
#### Locks ####
Setting name | Value | Unit
-------------|-------|------
{{ range $i, $key := (index (index (index (index .results .hosts.master) "data") "locks") "_keys") }}
{{- $value := (index (index (index (index $.results $.hosts.master) "data") "locks") $key) -}}
[{{ $key }}](https://postgresqlco.nf/en/doc/param/{{ $key }})|{{ $value.setting}}|{{ $value.unit }}
{{ end }}
#### Databases data ####
Database | Conflicts | Deadlocks | Stats reset at | Stat reset
-------------|-------|-----------|----------------|------------
{{ range $i, $key := (index (index (index (index .results .hosts.master) "data") "databases_stat") "_keys") }}
{{- $value:= (index (index (index (index $.results $.hosts.master) "data") "databases_stat") $key) -}}
{{$key}}|{{ $value.conflicts}}|{{ $value.deadlocks }}|{{ $value.stats_reset }}|{{ $value.stats_reset_age }}
{{ end }}
{{ end }}
{{ if gt (len .hosts.replicas) 0 }}
### Replica servers: ###
{{ range $skey, $host := .hosts.replicas }}
#### Replica (`{{ $host }}`) ####
{{ if (index $.results $host) }}
#### Timeouts ####
Setting name | Value | Unit
-------------|-------|------
{{ range $i, $key := (index (index (index (index $.results $host) "data") "timeouts") "_keys") }}
{{- $value := (index (index (index (index $.results $host) "data") "timeouts") $key) -}}
[{{ $key }}](https://postgresqlco.nf/en/doc/param/{{ $key }})|{{ $value.setting}}|{{ $value.unit }}
{{ end }}
#### Locks ####
Setting name | Value | Unit
-------------|-------|------
{{ range $i, $key := (index (index (index (index $.results $host) "data") "locks") "_keys") }}
{{- $value := (index (index (index (index $.results $host) "data") "locks") $key) -}}
[{{ $key }}](https://postgresqlco.nf/en/doc/param/{{ $key }})|{{ $value.setting}}|{{ $value.unit }}
{{ end }}
#### Databases data ####
Database | Conflicts | Deadlocks | Stats reset at | Stat reset
-------------|-------|-----------|----------------|------------
{{ range $i, $key := (index (index (index (index $.results $host) "data") "databases_stat") "_keys") }}
{{- $value:= (index (index (index (index $.results $host) "data") "databases_stat") $key) -}}
{{$key}}|{{ $value.conflicts}}|{{ $value.deadlocks }}|{{ $value.stats_reset }}|{{ $value.stats_reset_age }}
{{ end }}
{{ else }}
No data
{{ end}}{{ end }}{{ end }}
## Conclusions ##


## Recommendations ##

