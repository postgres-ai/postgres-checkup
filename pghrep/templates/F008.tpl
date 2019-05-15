# {{ .checkId }} Autovacuum: Resource Usage #

## Observations ##
Data collected: {{ DtFormat .timestamptz }}  
### Settings ###

{{ if .hosts.master }}
{{ if (index .results .hosts.master)}}
{{ if (index (index .results .hosts.master) "data") }}
| Setting name | Value | Unit | Pretty value |
|-------------|-------|------|--------------|
{{ range $i, $key := (index (index (index .results .hosts.master) "data") "_keys") }}
{{- $value := (index (index (index $.results $.hosts.master) "data") $key) -}}
| [{{ $key }}](https://postgresqlco.nf/en/doc/param/{{ $key }})|{{ $value.setting }}|{{ $value.unit }} | {{ UnitValue $value.setting $value.unit}} |
{{ end }}

### CPU ###

Cpu count you can see in report A001  

### RAM ###

Ram amount you can see in report A001  

{{- $autovacuum_work_mem := (RawIntUnitValue (index (index (index .results .hosts.master) "data") "autovacuum_work_mem").setting (index (index (index .results .hosts.master) "data") "autovacuum_work_mem").unit) -}}
{{- $maintenance_work_mem := (RawIntUnitValue (index (index (index .results .hosts.master) "data") "maintenance_work_mem").setting (index (index (index .results .hosts.master) "data") "maintenance_work_mem").unit) -}}
{{- $autovacuum_max_workers := (RawIntUnitValue (index (index (index .results .hosts.master) "data") "autovacuum_max_workers").setting (index (index (index .results .hosts.master) "data") "autovacuum_max_workers").unit) }}

{{ if eq $autovacuum_work_mem -1 -}}
Max workers memory: {{ ByteFormat ( Mul $maintenance_work_mem $autovacuum_max_workers ) 0 }}
{{- else -}}
Max workers memory: {{ ByteFormat ( Mul $autovacuum_work_mem $autovacuum_max_workers ) 0 }}
{{- end }}


### DISK ###

:warning: Warning: collection of current impact on disks is not yet implemented. Please refer to Postgres logs and see current read and write IO bandwidth caused by autovacuum.  
{{- else -}}{{/*Master data*/}}
No data
{{- end }}{{/*Master data*/}}
{{- else -}}{{/*Master results*/}}
No data
{{- end }}{{/*Master results*/}}
{{- else -}}{{/*Master*/}}
No data
{{ end }}{{/*Master*/}}

## Conclusions ##


## Recommendations ##

