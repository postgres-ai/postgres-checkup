# Unused/Rarely Used Indexes #

## Current values ##

### Master DB server is `{{.hosts.master}}` ###

#### Indexes ####

Index name | Reason | Scheme name | Table name | Index size | Table size
-----------|--------|-------------|------------|------------|------------
{{ range $index_name, $index_data := (index (index (index .results .hosts.master) "data") "indexes") }}{{ $index_name }} | {{ $index_data.reason }} | {{ $index_data.schemaname }} | {{ $index_data.tablename }} | {{ $index_data.index_size }} | {{ $index_data.table_size }}
{{ end }}

#### Drop code ####
{{ range $i, $drop_code := (index (index (index .results .hosts.master) "data") "drop_code") }}
```
{{ $drop_code }}
```
{{ end }}

#### Revert code ####
{{ range $i, $revert_code := (index (index (index .results .hosts.master) "data") "revert_code") }}
```
{{ $revert_code }}
```
{{ end }}


{{ if gt (len .hosts.replicas) 0 }}
### Slave DB servers: ###
    {{ range $skey, $host := .hosts.replicas }}
#### DB slave server: `{{ $host }}` ####
        {{ if (index $.results $host) }}  
            {{ range $index_name, $index_data := (index (index (index $.results $host) "data") "indexes") }} 
**Index name: `{{ $index_name }}`**

    Reason: {{ $index_data.reason }}
    Scheme name: {{ $index_data.schemaname }}
    Table name: {{ $index_data.tablename }}
    Index size: {{ $index_data.index_size }}
    Table size: {{ $index_data.table_size }}
    {{ if $index_data.drop_code }}Drop code: {{ $index_data.drop_code }} {{ end }}
    {{ if $index_data.revert_code }}Revert code: {{ $index_data.revert_code }} {{ end }}
{{ end }}{{ end }}{{ end }}{{ end }}

## Conclusions ##

{{.Conclusion}}

## Recommendations ##

{{.Recommended}}
