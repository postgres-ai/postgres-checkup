# {{ .checkId }} Autovacuum: Dead tuples #

## Observations ##
{{ if .hosts.master }}
### Master (`{{.hosts.master}}`) ###
 Relation | Since last autovacuum | Since last vacuum | Autovacuum Count | Vacuum Count | n_tup_ins | n_tup_upd | n_tup_del | pg_class.reltuples | n_live_tup | n_dead_tup | &#9660;&nbsp;Dead&nbsp;Tuples&nbsp;Ratio, %
----------|-----------------------|-------------------|----------|---------|-----------|-----------|-----------|--------------------|------------|------------|-----------
{{ range $i, $key := (index (index (index .results .hosts.master) "data") "_keys") }}
{{- $value := (index (index (index $.results $.hosts.master) "data") $key) -}}
{{ index $value "relation"}} | 
{{- index $value "since_last_autovacuum"}} |
{{- index $value "since_last_vacuum"}} |
{{- NumFormat (index $value "av_count") 2}} |
{{- NumFormat (index $value "v_count") 2 }} |
{{- NumFormat (index $value "n_tup_ins") 2 }} |
{{- NumFormat (index $value "n_tup_upd") 2 }} |
{{- NumFormat (index $value "n_tup_del") 2 }} |
{{- NumFormat (index $value "pg_class_reltuples") 2 }} |
{{- NumFormat (index $value "n_live_tup") 2 }} |
{{- NumFormat (index $value "n_dead_tup") 2 }} |
{{- index $value "dead_ratio"}}
{{ end }}
{{- else }}
No data
{{- end }}

## Conclusions ##


## Recommendations ##

