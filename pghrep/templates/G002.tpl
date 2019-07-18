# {{ .checkId }} Connections and Current Activity #

## Observations ##
Data collected: {{ DtFormat .timestamptz }}  
{{ if .hosts.master }}
{{ if (index .results .hosts.master) }}
{{ if (index (index .results .hosts.master) "data") }}
### Master (`{{.hosts.master}}`) ###
{{ if gt (len (index (index .results .hosts.master) "data")) .LISTLIMIT }}The list is limited to {{.LISTLIMIT}} items. All items {{ (len (index (index .results .hosts.master) "data")) }}.{{ end }}  

 \# | User | DB | Current state | Count | State changed >1m ago | State changed >1h ago | Tx age >1m | Tx age >1h
|----|------|----|---------------|-------|-----------------------|-----------------------|------------|-----------
{{ range $i, $key := (index (index (index .results .hosts.master) "data") "_keys") }}
{{- if lt $i $.LISTLIMIT -}}
    {{- $value := (index (index (index $.results $.hosts.master) "data") $key) -}}
| {{ $key }} |{{ if eq (Trim (Trim $value.user "*") " ") "ALL users" }}{{ Trim (Trim $value.user "*") " " }}{{else}}`{{ Trim (Trim $value.user "*") " " }}`{{end}}|{{ if eq (Trim (Trim $value.database "*") " ") "ALL databases" }}{{ Trim (Trim $value.database "*") " " }}{{else}}`{{ Trim (Trim $value.database "*") " " }}`{{end}}| {{ Trim (Trim (index $value "current_state") "*") " " }} | {{ $value.Count }} | {{ index $value "state_changed_more_1m_ago" }} | {{ index $value "state_changed_more_1h_ago" }} | {{ index $value "tx_age_more_1m" }} | {{ index $value "tx_age_more_1h" }} | 
{{/* if limit list */}}{{ end -}}
{{ end }}{{/* range */}}
{{ end }}{{/* if .host.master data */}}
{{ end }}{{/* if .results .host.master */}}
{{ end }}{{/* if .host.master */}}

{{- if gt (len .hosts.replicas) 0 -}}
### Replica servers: ###
{{ range $skey, $host := .hosts.replicas }}
#### Replica (`{{ $host }}`) ####
{{ if (index $.results $host) }}
{{ if gt (len (index (index $.results $host) "data")) $.LISTLIMIT }}The list is limited to {{ $.LISTLIMIT }} items. All items {{ (len (index (index $.results $host) "data")) }}.{{ end }}  

| \# | User | DB | Current state | Count | State changed >1m ago | State changed >1h ago | Tx age >1m | Tx age >1h |
|----|------|----|---------------|-------|-----------------------|-----------------------|------------|------------|
{{ range $i, $key := (index (index (index $.results $host) "data") "_keys") }}
{{- if lt $i $.LISTLIMIT -}}
{{- $value := (index (index (index $.results $host) "data") $key) -}}
| {{ $key }} |{{ if eq (Trim (Trim $value.user "*") " ") "ALL users" }}{{ Trim (Trim $value.user "*") " " }}{{else}}`{{ Trim (Trim $value.user "*") " " }}`{{end}}|{{ if eq (Trim (Trim $value.database "*") " ") "ALL databases" }}{{ Trim (Trim $value.database "*") " " }}{{else}}`{{ Trim (Trim $value.database "*") " " }}`{{end}}| {{ Trim (Trim (index $value "current_state") "*") " " }} | {{ $value.count }} | {{ index $value "state_changed_more_1m_ago" }} | {{ index $value "state_changed_more_1h_ago" }} | {{ index $value "tx_age_more_1m" }} | {{ index $value "tx_age_more_1h" }} |{{/* if limit   list */}}{{ end }}
{{ end }}{{/* data range */}}
{{- else -}}{{/* if $.results $host */}}
Nothing found
{{- end -}}{{/* if $.results $host */}}
{{- end -}}{{/* replicas range*/}}
{{- end -}}{{/* if replica */}}

## Conclusions ##

{{- if .processed }}
 {{- if .conclusions }}
  {{ range $conclusion := .conclusions -}}
   - {{ $conclusion.Message }}
  {{ end }}
 {{else}}
 {{end}}
{{ end }}

## Recommendations ##

{{- if .processed }}
 {{- if .recommendations }}
  {{ range $recommendation := .recommendations -}}
   - {{ $recommendation.Message }}
  {{ end }}
 {{else}}
  All good, no recommendations here.
 {{end}}
{{ end }}
