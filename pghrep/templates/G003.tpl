# {{ .checkId }} Timeouts, locks, deadlocks #

## Observations ##
Data collected: {{ DtFormat .timestamptz }}  
{{ if .hosts.master }}
{{ if and (index .results .hosts.master) (index (index .results .hosts.master) "data") }}
### Master (`{{.hosts.master}}`) ###
#### Timeouts ####
Setting name | Value | Unit | Pretty value
-------------|-------|------|--------------
{{ range $i, $key := (index (index (index (index .results .hosts.master) "data") "timeouts") "_keys") }}
{{- $value := (index (index (index (index $.results $.hosts.master) "data") "timeouts") $key) -}}
[{{ $key }}](https://postgresqlco.nf/en/doc/param/{{ $key }})|{{ $value.setting }}|{{ $value.unit }}|{{ UnitValue $value.setting $value.unit }}
{{ end }}
#### Locks ####
Setting name | Value | Unit | Pretty value
-------------|-------|------|--------------
{{ range $i, $key := (index (index (index (index .results .hosts.master) "data") "locks") "_keys") }}
{{- $value := (index (index (index (index $.results $.hosts.master) "data") "locks") $key) -}}
[{{ $key }}](https://postgresqlco.nf/en/doc/param/{{ $key }})|{{ $value.setting }}|{{ $value.unit }}|{{ UnitValue $value.setting $value.unit }}
{{ end }}
{{ if (index (index (index .results .hosts.master) "data") "db_specified_settings") }}
#### Database specified settings ####
Database | Setting
---------|---------
{{ range $i, $key := (index (index (index (index .results .hosts.master) "data") "db_specified_settings") "_keys") }}
{{- $value := (index (index (index (index $.results $.hosts.master) "data") "db_specified_settings") $key) -}}
{{ $value.database }} | {{ $value.setconfig }}
{{ end }}
{{- end -}}
{{ if (index (index (index .results .hosts.master) "data") "user_specified_settings") }}
#### User specified settings ####
User | Setting
---------|---------
{{ range $i, $key := (index (index (index (index .results .hosts.master) "data") "user_specified_settings") "_keys") }}
{{- $value := (index (index (index (index $.results $.hosts.master) "data") "user_specified_settings") $key) -}}
{{ $value.rolname }} | {{ $value.rolconfig }}
{{ end }}
{{- end -}}
{{ if (index (index (index .results .hosts.master) "data") "databases_stat") }}
#### Databases data ####
{{ if gt (len (index (index (index $.results $.hosts.master) "data") "databases_stat")) .ROWS_LIMIT }}The list is limited to {{.ROWS_LIMIT}} items.{{ end }}  

\# | Database | Conflicts | &#9660;&nbsp;Deadlocks | Stats reset at | Stat reset
--|-----------|-------|-----------|----------------|------------
{{ range $i, $key := (index (index (index (index .results .hosts.master) "data") "databases_stat") "_keys") }}
{{- $value:= (index (index (index (index $.results $.hosts.master) "data") "databases_stat") $key) -}}
{{ $value.num }}|{{- $key }}|{{ $value.conflicts}}|{{ $value.deadlocks }}|{{ $value.stats_reset }}|{{ $value.stats_reset_age }}
{{ end }}
{{ end }}
{{- end -}}
{{- end -}}

{{ if gt (len .hosts.replicas) 0 }}
### Replica servers: ###
{{ range $skey, $host := .hosts.replicas }}
#### Replica (`{{ $host }}`) ####
{{ if (index $.results $host) }}
#### Timeouts ####
Setting name | Value | Unit | Pretty value
-------------|-------|------|--------------
{{ range $i, $key := (index (index (index (index $.results $host) "data") "timeouts") "_keys") }}
{{- $value := (index (index (index (index $.results $host) "data") "timeouts") $key) -}}
[{{ $key }}](https://postgresqlco.nf/en/doc/param/{{ $key }})|{{ $value.setting}}|{{ $value.unit }}|{{ UnitValue $value.setting $value.unit }}
{{ end }}
#### Locks ####
Setting name | Value | Unit | Pretty value
-------------|-------|------|--------------
{{ range $i, $key := (index (index (index (index $.results $host) "data") "locks") "_keys") }}
{{- $value := (index (index (index (index $.results $host) "data") "locks") $key) -}}
[{{ $key }}](https://postgresqlco.nf/en/doc/param/{{ $key }})|{{ $value.setting}}|{{ $value.unit }}|{{ UnitValue $value.setting $value.unit }}
{{ end }}
{{ if (index (index (index $.results $host) "data") "db_specified_settings") }}
#### Database specified settings ####
Database | Setting
---------|---------
{{ range $i, $key := (index (index (index (index $.results $host) "data") "db_specified_settings") "_keys") }}
{{- $value := (index (index (index (index $.results $host) "data") "db_specified_settings") $key) -}}
{{ $value.database }} | {{ $value.setconfig }}
{{ end }}
{{- end -}}
{{ if (index (index (index $.results $host) "data") "user_specified_settings") }}
#### User specified settings ####
User | Setting
---------|---------
{{ range $i, $key := (index (index (index (index $.results $host) "data") "user_specified_settings") "_keys") }}
{{- $value := (index (index (index (index $.results $host) "data") "user_specified_settings") $key) -}}
{{ $value.rolname }} | {{ $value.rolconfig }}
{{ end }}
{{- end -}}
{{ if (index (index (index $.results $host) "data") "databases_stat") }}
#### Databases data ####
{{ if gt (len (index (index (index $.results $host) "data") "databases_stat")) .ROWS_LIMIT }}The list is limited to {{.ROWS_LIMIT}} items.{{ end }}  

Database | Conflicts | &#9660;&nbsp;Deadlocks | Stats reset at | Stat reset
-------------|-------|-----------|----------------|------------
{{ range $i, $key := (index (index (index (index $.results $host) "data") "databases_stat") "_keys") }}
{{- $value:= (index (index (index (index $.results $host) "data") "databases_stat") $key) -}}
{{$key}}|{{ $value.conflicts}}|{{ $value.deadlocks }}|{{ $value.stats_reset }}|{{ $value.stats_reset_age }}
{{ end }}
{{ end }}
{{ else }}
No data
{{ end}}{{ end }}{{ end }}
## Conclusions ##


## Recommendations ##

