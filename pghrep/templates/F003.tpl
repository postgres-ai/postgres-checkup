# {{ .checkId }} Autovacuum: Dead tuples #

## Observations ##
Data collected: {{ DtFormat .timestamptz }}  
Current database: {{ .database }}  
{{- if .hosts.master }}
{{- if (index .results .hosts.master) }}  
Stats reset: {{ (index (index (index .results .hosts.master) "data") "database_stat").stats_age }} ago ({{ DtFormat (index (index (index .results .hosts.master) "data") "database_stat").stats_reset }})  
### Master (`{{.hosts.master}}`) ###

 Relation | Type | Since last autovacuum | Since last vacuum | Autovacuum Count | Vacuum Count | n_tup_ins | n_tup_upd | n_tup_del | pg_class.reltuples | n_live_tup | n_dead_tup | &#9660;Dead Tuples Ratio, %
----------|------|-----------------------|-------------------|----------|---------|-----------|-----------|-----------|--------------------|------------|------------|-----------
{{ range $i, $key := (index (index (index (index .results .hosts.master) "data") "dead_tuples") "_keys") }}
{{- $value := (index (index (index (index $.results $.hosts.master) "data") "dead_tuples") $key) -}}
{{ index $value "relation"}} | 
{{- index $value "relkind"}} | 
{{- index $value "since_last_autovacuum"}} |
{{- index $value "since_last_vacuum"}} |
{{- NumFormat (index $value "av_count") -1 }} |
{{- NumFormat (index $value "v_count") -1 }} |
{{- NumFormat (index $value "n_tup_ins") -1 }} |
{{- NumFormat (index $value "n_tup_upd") -1 }} |
{{- NumFormat (index $value "n_tup_del") -1 }} |
{{- NumFormat (index $value "pg_class_reltuples") -1 }} |
{{- NumFormat (index $value "n_live_tup") -1 }} |
{{- NumFormat (index $value "n_dead_tup") -1 }} |
{{- if ge (Int (index $value "dead_ratio")) 10 }} **{{ (index $value "dead_ratio")}}** {{else}} {{ (index $value "dead_ratio")}} {{end}}
{{ end }}
{{- else -}}{{/*Master data*/}}
No data
{{- end }}{{/*Master data*/}}
{{- else -}}{{/*Master*/}}
No data
{{ end }}{{/*Master*/}}


## Conclusions ##


## Recommendations ##

