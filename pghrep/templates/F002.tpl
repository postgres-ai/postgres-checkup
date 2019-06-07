# {{ .checkId }} Autovacuum: Transaction ID Wraparound Check #
{{ $TABLE_LIST_LIMIT := 10 }}
## Observations ##
- Data collected: {{ DtFormat .timestamptz }}
- Current database: {{ .database }}
{{ if .hosts.master }}
{{ if (index .results .hosts.master)}}
{{ if (index (index .results .hosts.master) "data") }}

### Master (`{{.hosts.master}}`) ###
{{ if index (index (index .results .hosts.master) "data") "per_instance" }}

#### Databases ####
{{ if gt (len (index (index (index .results .hosts.master) "data") "per_instance")) .LISTLIMIT }}The list is limited to {{.LISTLIMIT}} items.{{ end }}

| \# | Database | &#9660;&nbsp;Age | Capacity used, % | Warning | datfrozenxid |
|--|--------|-----|------------------|---------|--------------|
{{ range $i, $key := (index (index (index (index .results .hosts.master) "data") "per_instance") "_keys") }}
{{- if lt $i $.LISTLIMIT -}}
{{- $value := (index (index (index (index $.results $.hosts.master) "data") "per_instance") $key) -}}
| {{ $value.num }} |`{{- index $value "datname"}}`|
{{- NumFormat (index $value "age") -1 }} |
{{- index $value "capacity_used"}} |
{{- if (index $value "warning") }} &#9888; {{ else }} {{ end }} |
{{- NumFormat (index $value "datfrozenxid") -1}} |
{{/* if limit list */}}{{ end -}}
{{ end }}{{/* range */}}
{{- end -}}{{/* if per_instance exists */}}

#### Tables in the observed database ####
{{ if gt (len (index (index (index .results .hosts.master) "data") "per_database")) $TABLE_LIST_LIMIT }}The list is limited to {{ $TABLE_LIST_LIMIT }} items.{{ end }}

| \# | Relation | Age | &#9660;&nbsp;Capacity used, % | Warning |rel_relfrozenxid | toast_relfrozenxid |
|---|-------|-----|------------------|---------|-----------------|--------------------|
{{ range $i, $key := (index (index (index (index .results .hosts.master) "data") "per_database") "_keys") }}
{{- if lt $i $TABLE_LIST_LIMIT -}}
{{- $value := (index (index (index (index $.results $.hosts.master) "data") "per_database") $key) -}}
| {{ $value.num }} |`{{ index $value "relation"}}`{{if $value.overrided_settings}}\*{{ end }} |
{{- NumFormat (index $value "age") -1 }} |
{{- index $value "capacity_used"}} |
{{- if (index $value "warning") }} &#9888; {{ else }} {{ end }} |
{{- NumFormat (index $value "rel_relfrozenxid") -1}} |
{{- NumFormat (index $value "toast_relfrozenxid") -1}} |
{{/* if limit list */}}{{ end -}}
{{ end }}{{/* range */}}
{{/*- end -*/}}{{/* if per_instance exists */}}
{{- if gt (Int (index (index (index .results .hosts.master) "data") "overrided_settings_count")) 0 }}
\* This table has specific autovacuum settings. See 'F001 Autovacuum: Current settings'
{{- end }}
{{- else -}}{{/*Master data*/}}
Nothing found
{{- end }}{{/*Master data*/}}
{{- else -}}{{/*Master results*/}}
Nothing found
{{- end }}{{/*Master results*/}}
{{- else -}}{{/*Master*/}}
Nothing found
{{ end }}{{/*Master*/}}


## Conclusions ##

{{- if .conclusions }}
{{ range $conclusion := .conclusions -}}
- {{ $conclusion }}
{{ end }}
{{ end }}

## Recommendations ##

{{- if .recommendations }}
{{ range $recommendation := .recommendations -}}
- {{ $recommendation }}
{{ end }}
{{ end }}