About
===
PgHealth is the ultimate PostgreSQL database healthcheck utility.

It checks PostgreSQL settings, configs and Linux
postgres-related environment, collecting data
into convinient formats (json, md, ...) with
a series of checks.

Example
===

```bash
./check -h db1 -p 5432 --username postgres --dbname postgres --project my-site_org -e 1
```

```bash
./check -h db2 -p 5432 --username postgres --dbname postgres --project my-site_org -e 1
```

Which literaly means: "connect to server with given credentials, save data into `my-site_org`
project directory as epoch of check `1`.

As a result of health-check we have got a two directories with .json files and .md files:

```bash
./artifacts/my-site_org/json_reports/1_2018_12_06T14_12_36_+0300/
./artifacts/my-site_org/md_reports/1_2018_12_06T14_12_36_+0300/
```

`1` corresponds to an `epoch` given by argument `-e`.

Requirements
===

* bash
* coreutils
* jq >= 1.5,
* golang >= 1.8
* awk


