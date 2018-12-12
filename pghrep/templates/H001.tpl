# {{ .checkId }} Unused/Rarely Used Indexes #

## Observations ##

### Master (`{{.hosts.master}}`) ###

#### Indexes ####

Index name | Reason | Scheme name | Table name | Index size | Table size
-----------|--------|-------------|------------|------------|------------
{{ range $index_name, $index_data := (index (index (index .results .hosts.master) "data") "indexes") }}{{ $index_name }} | {{ $index_data.reason }} | {{ $index_data.schemaname }} | {{ $index_data.tablename }} | {{ $index_data.index_size }} | {{ $index_data.table_size }}
{{ end }}

#### Drop code ####
```
{{ range $i, $drop_code := (index (index (index .results .hosts.master) "data") "drop_code") }}{{ $drop_code }}
{{ end }}
```

#### Revert code ####
```
{{ range $i, $revert_code := (index (index (index .results .hosts.master) "data") "revert_code") }}{{ $revert_code }}
{{ end }}
```

{{ if gt (len .hosts.replicas) 0 }}
### Replica servers: ###
    {{ range $skey, $host := .hosts.replicas }}
#### Replica server: `{{ $host }}` ####
        {{ if (index $.results $host) }}
#### Indexes ####

Index name | Reason | Scheme name | Table name | Index size | Table size
-----------|--------|-------------|------------|------------|------------
{{ range $index_name, $index_data := (index (index (index $.results $host) "data") "indexes") }}{{ $index_name }} | {{ $index_data.reason }} | {{ $index_data.schemaname }} | {{ $index_data.tablename }} | {{ $index_data.index_size }} | {{ $index_data.table_size }}
{{ end }}

#### Drop code ####
```
{{ range $i, $drop_code := (index (index (index $.results $host) "data") "drop_code") }}{{ $drop_code }}
{{ end }}
```

#### Revert code ####
```
{{ range $i, $revert_code := (index (index (index $.results $host) "data") "revert_code") }}{{ $revert_code }}
{{ end }}
```
{{ end }}{{ end }}{{ end }}

## Conclusions ##


## Recommendations ##

