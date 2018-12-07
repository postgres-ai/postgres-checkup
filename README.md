About
===
PgHealth is the ultimate open-source PostgreSQL database healthcheck utility.

It checks PostgreSQL settings, configs and Linux postgres-related environment
with series of checks.

The observed data is saved in the form of JSON reports, ready to be consumed by machines.  
The final reports are .md files, in Markdown format, to be read by humans.

The main goal is to detect bottlenecks and prevent performance degradation.  
It also helps to detect a lot of issues with postgres instances.

Example
===

```bash
./check -h db1 -p 5432 --username postgres --dbname postgres --project my-site_org-slony -e 1
```

```bash
./check -h db2 -p 5432 --username postgres --dbname postgres --project my-site_org-slony -e 1
```

Which literaly means: "connect to server with given credentials, save data into `my-site_org-slony`
project directory as epoch of check `1`. Epoch is a numerical sign of current iteration.
For example: in half a year we can switch to "epoch number `2`".

As a result of health-check we have got two directories with .json files and .md files:

```bash
./artifacts/my-site_org-slony/json_reports/1_2018_12_06T14_12_36_+0300/
./artifacts/my-site_org-slony/md_reports/1_2018_12_06T14_12_36_+0300/
```

Each of generated files contains information about "what we check" and collected data for
all instances of the postgres cluster `my-site_org-slony`.

Requirements
===

* bash
* psql
* coreutils
* jq >= 1.5,
* golang >= 1.8
* awk


