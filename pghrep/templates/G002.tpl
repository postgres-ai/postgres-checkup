# {{ .checkId }} Connections and Current Activity #

## Observations ##
Data collected: {{ DtFormat .timestamptz }}  
{{ if .hosts.master }}
{{ if (index .results .hosts.master) }}
{{ if (index (index .results .hosts.master) "data") }}
### Master (`{{.hosts.master}}`) ###
{{ if gt (len (index (index .results .hosts.master) "data")) .ROWS_LIMIT }}The list is limited to {{.ROWS_LIMIT}} items.{{ end }}  

\# | User | DB | Current state | Count | State changed >1m ago | State changed >1h ago | Tx age >1m | Tx age >1h
----|------|----|---------------|-------|-----------------------|----------------------|------------|-----------
{{ range $i, $key := (index (index (index .results .hosts.master) "data") "_keys") }}
    {{- $value := (index (index (index $.results $.hosts.master) "data") $key) -}}
    {{ $key }} | {{ Trim (Trim $value.User "*") " " }} | {{ Trim (Trim $value.DB "*") " " }} | {{ Trim (Trim (index $value "Current State") "*") " " }} | {{ $value.Count }} | {{ index $value "State changed >1m ago" }} | {{ index $value "State changed >1h ago" }} | {{ index $value "Tx age >1m" }} | {{ index $value "Tx age >1h" }}
{{ end }}{{/* range */}}
{{ end }}{{/* if .host.master data */}}
{{ end }}{{/* if .results .host.master */}}
{{ end }}{{/* if .host.master */}}

{{- if gt (len .hosts.replicas) 0 -}}
### Replica servers: ###
{{ range $skey, $host := .hosts.replicas }}
#### Replica (`{{ $host }}`) ####
{{ if (index $.results $host) }}
{{ if gt (len (index (index $.results $host) "data")) $.ROWS_LIMIT }}The list is limited to {{ $.ROWS_LIMIT }} items.{{ end }}  

\# | User | DB | Current state | Count | State changed >1m ago | State changed >1h ago | Tx age >1m | Tx age >1h
----|------|----|---------------|-------|-----------------------|----------------------|------------|-----------
{{ range $i, $key := (index (index (index $.results $host) "data") "_keys") }}
{{- $value := (index (index (index $.results $host) "data") $key) -}}
{{ $key }} | {{ Trim (Trim $value.User "*") " " }} | {{ Trim (Trim $value.DB "*") " " }} | {{ Trim (Trim (index $value "Current State") "*") " " }} | {{ $value.Count }} | {{ index $value "State changed >1m ago" }} | {{ index $value "State changed >1h ago" }} | {{ index $value "Tx age >1m" }} | {{ index $value "Tx age >1h" }}
{{ end }}{{/* data range */}}
{{- else -}}{{/* if $.results $host */}}
No data
{{- end -}}{{/* if $.results $host */}}
{{- end -}}{{/* replicas range*/}}
{{- end -}}{{/* if replica */}}

## Conclusions ##


## Recommendations ##

