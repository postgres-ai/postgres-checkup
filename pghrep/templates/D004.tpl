# `pg_stat_statements` and `kcache` settings #

## Current values ##

### Master DB server is `{{.hosts.master}}` ###

#### `pg_stat_statements` extension settings ####
{{ range $setting_name, $setting_data := (index (index (index .results .hosts.master) "data") "pg_stat_statements") }}
**Setting: `{{ $setting_name }}`**

    Setting value: {{ $setting_data.setting }} {{ if $setting_data.unit }}{{ $setting_data.unit }} {{ end }}
    Description: {{ $setting_data.short_desc }}
    Category: {{ $setting_data.category }}
    Type: {{ $setting_data.vartype }}
    Source: {{ $setting_data.source }}
    Source file: {{ $setting_data.sourcefile }}
    {{ if $setting_data.min_val }}Min value: {{ $setting_data.min_val }} {{ end }}
    {{ if $setting_data.max_val }}Max value: {{ $setting_data.max_val }} {{ end }}
{{ end }}

#### `kcache` extension settings ####
{{ range $setting_name, $setting_data := (index (index (index .results .hosts.master) "data") "kcache") }}
**Setting: `{{ $setting_name }}`**

    Setting name: {{ $setting_data.name }}
    Setting value: {{ $setting_data.setting }} {{ if $setting_data.unit }}{{ $setting_data.unit }} {{ end }}
    Description: {{ $setting_data.short_desc }}
    Category: {{ $setting_data.category }}
    Type: {{ $setting_data.vartype }}
    Source: {{ $setting_data.source }}
    Source file: {{ $setting_data.sourcefile }}
    {{ if $setting_data.min_val }}Min value: {{ $setting_data.min_val }} {{ end }}
    {{ if $setting_data.max_val }}Max value: {{ $setting_data.max_val }} {{ end }}
{{ end }}


{{ if gt (len .hosts.replicas) 0 }}
### Slave DB servers: ###
    {{ range $skey, $host := .hosts.replicas }}
#### DB slave server: `{{ $host }}` ####
        {{ if (index $.results $host) }}  
#### `pg_stat_statements` settings ####
{{ range $setting_name, $setting_data := (index (index (index $.results $host) "data") "pg_stat_statements") }}
**Setting: `{{ $setting_name }}`**

    Setting name: {{ $setting_data.name }}
    Setting value: {{ $setting_data.setting }} {{ if $setting_data.unit }}{{ $setting_data.unit }} {{ end }}
    Description: {{ $setting_data.short_desc }}
    Category: {{ $setting_data.category }}
    Type: {{ $setting_data.vartype }}
    Source: {{ $setting_data.source }}
    Source file: {{ $setting_data.sourcefile }}
    {{ if $setting_data.min_val }}Min value: {{ $setting_data.min_val }} {{ end }}
    {{ if $setting_data.max_val }}Max value: {{ $setting_data.max_val }} {{ end }}
{{ end }}

#### `kcache` settings ####
{{ range $setting_name, $setting_data := (index (index (index $.results $host) "data") "kcache") }}
**Setting: `{{ $setting_name }}`**

    Setting name: {{ $setting_data.name }}
    Setting value: {{ $setting_data.setting }} {{ if $setting_data.unit }}{{ $setting_data.unit }} {{ end }}
    Description: {{ $setting_data.short_desc }}
    Category: {{ $setting_data.category }}
    Type: {{ $setting_data.vartype }}
    Source: {{ $setting_data.source }}
    Source file: {{ $setting_data.sourcefile }}
    {{ if $setting_data.min_val }}Min value: {{ $setting_data.min_val }} {{ end }}
    {{ if $setting_data.max_val }}Max value: {{ $setting_data.max_val }} {{ end }}
{{ end }}
        {{ end }}
    {{ end }}
{{ end }}

## Conclusions ##

{{.Conclusion}}

## Recommendations ##

{{.Recommended}}