# About
Postgres Checkup ([postgres-checkup](https://gitlab.com/postgres-ai-team/postgres-checkup)) is a new-generation diagnostics tool that allows to collect deep analysis of a Postgres database health. It aims to detect and describe all current and potential issues in the fields of database performance, scalability and security, advising how to resolve or prevent them.

Compared to a monitoring system, postgres-checkup goes deeper in database system and environment analysis, it combines numerious database internal characteristics with data about resources and OS into multiple comprehensive reports. These reports use formats which are easy readable both by humans and machines and are extremely oriented to DBA problems solving. Monitoring systems constantly collect telemetry, they help to react to issues quicker and do post-mortem analyses. Checkups are needed for a different purpose: detect issues at very early stage, advising on how to prevent them. This procedure is to be done on a regular basis — weekly, monthly, or quarterly. Additionally, it is recommended to run it right right before and right after any major change in database server.

The three key principles behind postgres-checkup:

- *Unobtrusiveness*: postgres-checkup’s impact on the observing system is close to zero. It does not use any heavy queries, keeping resource usage very low, avoids having the “observer effect”.

- *Zero install* (on observed machines): it is able to analyze any Linux machine (including virtual machines), as well as Cloud Postgres instanced (such as Amazon RDs or Google Cloud SQL), not requiring any additional setup or any changes. It does, hovewer, require the privileged access (a DBA usually has it anyway).

- *Complex analysis*: unlike most monitoring tools, which provide raw data, postgres-checkup combines data from various parts of the system (w.g.: internal Postgres stats are combined with knowledge about system resources in autovacuum setting and behavior analysis). Also, it analyzes the master database server together with all its replicas (e.g. to build the list of unused indexes).

# Reports structure
The two kinds of reports postgres-checkup produces for every check:

- JSON reports (*.json) — can be consumed by any program or service, or stores in some database.

- Markdown reports (*.md) — the main format for humans, may contain lists, tables, pictures. Being of native format for GitLab and GitHub, such reports are ready to be used, for instance, in their issue trackers, simplifying workflow. Markdown reports are derived from JSON reports.

Markdown reports can be converted to different formats such as HTML or PDF.

Each report consists of three sections:
1. Observations
1. Conclusions
1. Recommendations

## Installation

### Requirements

The supported OS of the observer machine (those from which the tool is to be executed):

* Linux (modern RHEL/CentOS or Debian/Ubuntu; others should work as well, but not yet tested);
* MacOS.

It has to have the following programs:

* bash
* psql
* coreutils
* jq >= 1.5
* golang >= 1.8
* awk
* sed

Nothing special has to be istalled on the observed machines. However, these machines must run Linux (again: modern RHEL/CentOS or Debian/Ubuntu; others should work as well, but not yet tested).


# Example of use

Let's make a report for a project named `prod1`:
Cluster `slony` contains two servers - `db1.vpn.local` and `db1.vpn.local`.
Postgres-checkup automatically detects which one is a master:

```bash
./checkup -h db1.vpn.local -p 5432 --username postgres --dbname postgres --project prod1
```

```bash
./checkup -h db2.vpn.local -p 5432 --username postgres --dbname postgres --project prod1 -e 1
```

Which literally means: "connect to the server with given credentials, save data into `prod1`
project directory as epoch of check `1`. Epoch is a numerical (**integer**) sign of current iteration.
For example: in half a year we can switch to "epoch number `2`".

At the first run we can skip `-e 1` because default epoch is `1`, but at the second argument `-e`  
must exist: we don't want to overwrite historical results.


As a result of postgres-checkup we have got two directories with .json files and .md files:

```bash
./artifacts/prod1/json_reports/1_2018_12_06T14_12_36_+0300/
./artifacts/prod1/md_reports/1_2018_12_06T14_12_36_+0300/
```

Each of generated files contains information about "what we check" and collected data for
all instances of the postgres cluster `prod1`.

Human-readable report can be found at:

```bash
./artifacts/prod1/e1_full_report.md
```

Open it with your favorite Markdown files viewer or just upload to a service such as gist.github.com.



