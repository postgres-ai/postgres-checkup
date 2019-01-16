# {{ .checkId }} Autovacuum: resource usage #

## Observations ##

### Settings ###

Setting name | Value | Unit | Pretty value
-------------|-------|------|--------------
{{ range $i, $key := (index (index (index .results .hosts.master) "data") "_keys") }}
{{- $value := (index (index (index $.results $.hosts.master) "data") $key) -}}
[{{ $key }}](https://postgresqlco.nf/en/doc/param/{{ $key }})|{{ $value.setting }}|{{ $value.unit }} | {{ UnitValue $value.setting $value.unit}}
{{ end }}

### CPU ###

Cpu count you can see in report A001  

### RAM ###

Ram amount you can see in report A001  

{{- $autovacuum_work_mem := (RawIntUnitValue (index (index (index .results .hosts.master) "data") "autovacuum_work_mem").setting (index (index (index .results .hosts.master) "data") "autovacuum_work_mem").unit) -}}
{{- $maintenance_work_mem := (RawIntUnitValue (index (index (index .results .hosts.master) "data") "maintenance_work_mem").setting (index (index (index .results .hosts.master) "data") "maintenance_work_mem").unit) -}}
{{- $max_connections := (RawIntUnitValue (index (index (index .results .hosts.master) "data") "max_connections").setting (index (index (index .results .hosts.master) "data") "max_connections").unit) }}

{{ if eq $autovacuum_work_mem -1 -}}
Max workers memory: {{ ByteFormat ( Mul $maintenance_work_mem $max_connections ) 0 }}
{{- else -}}
Max workers memory: {{ ByteFormat ( Mul $autovacuum_work_mem $max_connections ) 0 }}
{{- end }}

### DISK ###

:warning: Warning: collection of current impact on disks is not yet implemented. Please refer to Postgres logs and see current read and write IO bandwidth caused by autovacuum.  

## Conclusions ##


## Recommendations ##

