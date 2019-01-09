# {{ .checkId }} Autovacuum: dead rows #

## Observations ##
{{ if .hosts.master }}
### Master (`{{.hosts.master}}`) ###
 Relation | Since last autovacuum | Since last vacuum | av_count | v_count | n_tup_ins | n_tup_upd | n_tup_del | pg_class_reltuples | n_live_tup | n_dead_tup | dead_ratio
----------|-----------------------|-------------------|----------|---------|-----------|-----------|-----------|--------------------|------------|------------|-----------
{{ range $i, $key := (index (index (index .results .hosts.master) "data") "_keys") }}
{{- $value := (index (index (index $.results $.hosts.master) "data") $key) -}}
{{ index $value "relation"}} | 
{{- index $value "since_last_autovacuum"}} | 
{{- index $value "since_last_vacuum"}} | 
{{- index $value "av_count"}} | 
{{- index $value "v_count"}} | 
{{- index $value "n_tup_ins"}} |
{{- index $value "n_tup_upd"}} |
{{- index $value "n_tup_del"}} |
{{- index $value "pg_class_reltuples"}} |
{{- index $value "n_live_tup"}} |
{{- index $value "n_dead_tup"}} |
{{- index $value "dead_ratio"}}
{{ end }}
{{- else }}
No data
{{- end }}

## Conclusions ##


## Recommendations ##

