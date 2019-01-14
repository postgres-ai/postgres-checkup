# {{ .checkId }} Connections and current activity #

## Observations ##
{{ if .hosts.master }}
### Master (`{{.hosts.master}}`) ###
Num | User | DB | Current state | Count | State changed >1m ago | State changed >1h ago
----|------|----|---------------|-------|-----------------------|-----------------------
{{ range $i, $key := (index (index (index .results .hosts.master) "data") "_keys") }}
    {{- $value := (index (index (index $.results $.hosts.master) "data") $key) -}}
    {{ $key }} | {{ Trim (Trim $value.User "*") " " }} | {{ Trim (Trim $value.DB "*") " " }} | {{ Trim (Trim (index $value "Current State") "*") " " }} | {{ $value.Count }} | {{ index $value "State changed >1m ago" }} | {{ index $value "State changed >1h ago" }}
{{ end }}{{/* range */}}
{{ end }}{{/* if .host.master */}}


{{- if gt (len .hosts.replicas) 0 -}}
### Replica servers: ###
{{ range $skey, $host := .hosts.replicas }}
#### Replica (`{{ $host }}`) ####
{{ if (index $.results $host) }}
Num | User | DB | Current state | Count | State changed >1m ago | State changed >1h ago
----|------|----|---------------|-------|-----------------------|-----------------------
{{ range $i, $key := (index (index (index $.results $host) "data") "_keys") }}
{{- $value := (index (index (index $.results $host) "data") $key) -}}
{{ $key }} | {{ Trim (Trim $value.User "*") " " }} | {{ Trim (Trim $value.DB "*") " " }} | {{ Trim (Trim (index $value "Current State") "*") " " }} | {{ $value.Count }} | {{ index $value "State changed >1m ago" }} | {{ index $value "State changed >1h ago" }}
{{ end }}{{/* data range */}}
{{- else -}}{{/* if $.results $host */}}
No data
{{- end -}}{{/* if $.results $host */}}
{{- end -}}{{/* replicas range*/}}
{{- end -}}{{/* if replica */}}

## Conclusions ##


## Recommendations ##

