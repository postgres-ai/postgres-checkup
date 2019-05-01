# Demo

Auto-generated demonstration based on the code in the master
branch (only single node analyzed): https://gitlab.com/postgres-ai-team/postgres-checkup-tests/tree/master/master.
Go to `md_reports/TIMESTAMP` and then open `0_Full_report.md`.

# Disclaimer: This Tool is Designed for DBA Experts

Each report consists of 3 sections: Observations, Conclusions, and Recommendations.
As of March 2019, only Observations are filled automatically. To treat the data
correctly, you need deep Postgres knowledge.

You can get Conclusions, and Recommendations from the Postgres.ai team for free,
send us your .json and .md with filled Observations sections: checkup@postgres.ai.
Limited time only. We're a small team, so "restrictions apply".

# About

Postgres Checkup ([postgres-checkup](https://gitlab.com/postgres-ai-team/postgres-checkup))
is a new-generation diagnostics tool that allows users to collect  deep analysis
of the health of a Postgres database. It aims to detect and describe all current
and potential issues in the fields of database performance, scalability, and
security, providing advices how to resolve or prevent them.

Compared to a monitoring system, postgres-checkup goes deeper into the analysis
of the database system and environment.  It combines numerous internal
characteristics of the database with data about resources and OS, producing
multiple comprehensive reports. These reports use formats which are easily
readable both by humans and machines and which are extremely oriented to DBA
problem-solving. Monitoring systems constantly collect telemetry, help to react
to issues more quickly, and are useful for post-mortem analyses. At the same
time, checkups are needed for a different purpose: detect issues at a very early
stage, advising on how to prevent them. This procedure is to be done on a
regular basis — weekly, monthly, or quarterly. Additionally, it is recommended
to run it immediately before and after any major change in the database server.

The three key principles behind postgres-checkup:

- *Unobtrusiveness*: postgres-checkup’s impact on the observing system is
close to zero. It does not use any heavy queries, keeping resource usage
very low, and avoiding having the [“observer effect.”](https://en.wikipedia.org/wiki/Observer_effect_(information_technology))

- *Zero install* (on observed machines): it is able to analyze any Linux
machine (including virtual machines), as well as Cloud Postgres instances
(such as Amazon RDs or Google Cloud SQL), not requiring any additional setup
or any changes. It does, hovewer, require a privileged access (a DBA usually
has it anyway).

- *Complex analysis*: unlike most monitoring tools, which provide raw data,
postgres-checkup combines data from various parts of the system (e.g.,
internal Postgres stats are combined with knowledge about system resources
in autovacuum setting and behavior analysis). Also, it analyzes the master
database server together with all its replicas (e.g. to build the list of
unused indexes).

# Reports Structure

Postgres-checkup produces two kinds of reports for every check:

- JSON reports (*.json) — can be consumed by any program or service, or
stores in some database.

- Markdown reports (*.md) — the main format for humans, may contain lists,
tables, pictures. Being of native format for GitLab and GitHub, such reports
are ready to be used, for instance, in their issue trackers, simplifying
workflow. Markdown reports are derived from JSON reports.

Markdown reports can be converted to different formats such as HTML or PDF.

Each report consists of three sections:

1. "Observations": automatically collected data. This is to be consumed by
an expert DBA.
1. "Conclusions": what we conclude from the Observations—what is good, what
is bad (right now, it is to be manually filled for most checks).
1. "Recommendations": action items, what to do to fix the discovered issues.
Both "Conclusions" and "Recommendations" are to be consumed by engineers who
will make decisions what, how and when to optimize, and how to react to the
findings.

# Installation and Usage

## Requirements

The supported OS of the observer machine (those from which the tool is to be
executed):

* Linux (modern RHEL/CentOS or Debian/Ubuntu; others should work as well, but
are not yet tested);
* MacOS.

The following programs must be installed on the observer machine:

* bash
* psql
* coreutils
* jq >= 1.5
* golang >= 1.8 (no binaries are shipped at the moment)
* awk
* sed

Nothing special has to be installed on the observed machines. However, these
machines must run Linux (again: modern RHEL/CentOS or Debian/Ubuntu; others
should work as well, but are not yet tested).

:warning: Only Postgres version 9.6 and higher are currently supported.

## How to Install

Use `git clone`. This is the only method of installation currently supported.

## Example of Use

Let's make a report for a project named `prod1`:
Cluster `slony` contains two servers - `db1.vpn.local` and `db1.vpn.local`.
Postgres-checkup automatically detects which one is a master:

```bash
./checkup -h db1.vpn.local -p 5432 --username postgres --dbname postgres --project prod1 -e 1
```

```bash
./checkup -h db2.vpn.local -p 5432 --username postgres --dbname postgres --project prod1 -e 1
```

Which literally means: connect to the server with given credentials, save data into `prod1`
project directory, as epoch of check `1`. Epoch is a numerical (**integer**) sign of current iteration.
For example: in half a year we can switch to "epoch number `2`".

`-h db2.vpn.local` means: try to connect to host via SSH and then use remote `psql` command to perform checks.  
If SSH is not available the local 'psql' will be used (non-psql reports will be skipped).

For comprehensive analysis, it is recommended to run the tool on the master and
all its replicas – postgres-checkup is able to combine all the information from
multiple nodes to a single report.

Some reports (such as K003) require two snapshots, to calculate "deltas" of
metrics. So, for better results, use the following example, executing it during peak working
hours, with `$DISTANCE` values from 10 min to a few hours:

```bash
$DISTANCE="1800" # 30 minutes

# Assuming that db2 is the master, db3 and db4 are its replicas
for host in db2.vpn.local db3.vpn.local db4.vpn.local; do
  ./checkup \
    -h "$host" \
    -p 5432 \
    --username postgres \
    --dbname postgres \
    --project prod1 \
    -e 1 \
    --file resources/checks/K000_query_analysis.sh # the first snapshot is needed only for reports K***
done
  
sleep "$DISTANCE"

for host in db2.vpn.local db3.vpn.local db4.vpn.local; do
  ./checkup \
    -h "$host" \
    -p 5432 \
    --username postgres \
    --dbname postgres \
    --project prod1 \
    -e 1
done
```

As a result of execution, two directories containing .json and .md files will
be created:

```bash
./artifacts/prod1/json_reports/1_2018_12_06T14_12_36_+0300/
./artifacts/prod1/md_reports/1_2018_12_06T14_12_36_+0300/
```

Each of generated files contains information about "what we check" and collected data for
all instances of the postgres cluster `prod1`.

A human-readable report can be found at:

```bash
./artifacts/prod1/md_reports/1_2018_12_06T14_12_36_+0300/Full_report.md
```

Open it with your favorite Markdown files viewer or just upload to a service such as gist.github.com.

## Docker 🐳

It's possible to use the `postgres-checkup` from a docker container.
The container will run, execute all checks and stop itself.
The check result can be found inside the `artifacts` folder in current directory (pwd).

### Usage with `docker run`

First of all we need a postgres. You can use any local or remote running instance.
For this example we run postgres in a separate docker container:

```bash
docker run \
    --name postgres \
    -e POSTGRES_PASSWORD=postgres \
    -d postgres
```

We need to know a hostname or an ip address of target database to be used with `-h` parameter:

```bash
PG_HOST=$(docker inspect --format '{{ .NetworkSettings.IPAddress }}' postgres)
```

You can use official images or build an image yourself. 
Run this command to build an image:

```bash
docker build -t postgres-checkup .
```

Then run a container with `postgres-checkup`. 
This command run the tool with access via `psql`:

```bash
docker run --rm \
    --name postgres-checkup \
    -e PGPASSWORD="postgres" \
    -v `pwd`/artifacts:/artifacts \
    postgres-checkup \
    ./checkup -h $PG_HOST -p 5432 --username postgres --dbname postgres --project docker -e 1
```

If you want to execute all supported checks you have to use `ssh` access to target host with postgres.
With docker container it's possible by mounting the ssh key file and specify the username in `-h` parameter:

```bash
docker run --rm \
    --name postgres-checkup \
    -e PGPASSWORD="postgres" \
    -v `pwd`/artifacts:/artifacts \
    -v `pwd`/ssh/key:/root/.ssh/id_rsa:ro \
    postgres-checkup \
    ./checkup -h $SSH_USER@$PG_HOST -p 5432 --username postgres --dbname postgres --project docker -e 1
```

If you try to check the local instance of postgres on your host from a container, you can't use `localhost` in `-h` parameter.
You have to use `bridge` between host OS and Docker Engine. By default host IP is `172.17.0.1` in `docker0` network, but it can be different.
More information [here](https://nickjanetakis.com/blog/docker-tip-65-get-your-docker-hosts-ip-address-from-in-a-container).

### Usage with `docker-compose`

It will run an empty `postgres` database and `postgres-checkup` application that will stop when it's done.
Local folder `artifacts` will contain `docker` subfolder with check result.

```bash
docker-compose build
docker-compose up -d

docker-compose down
```

## Credits

Some reports are based on or inspired by useful queries created and improved by
various developers, including but not limited to:
 * Jehan-Guillaume (ioguix) de Rorthais https://github.com/ioguix/pgsql-bloat-estimation
 * Alexey Lesovsky, Alexey Ermakov, Maxim Boguk, Ilya Kosmodemiansky et al. from Data Egret (aka PostgreSQL-Consulting) https://github.com/dataegret/pg-utils
 * Josh Berkus, Quinn Weaver et al. from PostgreSQL Experts, Inc. https://github.com/pgexperts/pgx_scripts

# The Full List of Reports

## А. General  / Infrastructural

- [x] A001 System, CPU, RAM, disks, virtualization #6 , #56 , #57 , #86 
- [x] A002 PostgreSQL versions (Simple) #68, #21, #86
- [x] A003 Collect pg_settings  #15, #167, #86 
- [x] A004 General cluster info  #7, #58, #59, #86, #162  
- [x] A005 Extensions #8, #60, #61, #86, #167   
- [x] A006 Config diff  #9, #62, #63, #86  
- [x] A007 ALTER SYSTEM vs postgresql.conf #18, #86  
- [x] A008 Disk usage and file system type #19, #20  
- [ ] A010 Data checksums, wal_log_hints #22  
- [ ] A011 Connection pooling. pgbouncer #23  
- [ ] A012 Anti-crash checks #177  

## B. Backups and DR  

- [ ] B001 SLO/SLA, RPO, RTO  #24  
- [ ] B002 File system, mount flags #25  
- [ ] B003 Full backups / incremental  #26  
- [ ] B004 WAL archiving (GB/day?) - #27  
- [ ] B005 Restore checks, monitoring, alerting  #28  

## C. Replication and HA

- [ ] C001 SLO/SLA  #29  
- [ ] C002 Sync/async, Streaming / wal transfer; logical decoding #30  
- [ ] C003 SPOFs; “-1 datacenter”, standby with traffic #31  
- [ ] C004 Failover #32  
- [ ] C005 Switchover #33  
- [ ] C006 Delayed replica (replay of 1 day of WALs) - #34  

## D. Monitoring / Troubleshooting   

- [ ] D001 Logging (syslog?), log_*** #35  
- [x] D002 Useful Linux tools  #36  
- [ ] D003 List of monitoring metrics #37  
- [x] D004 pg_stat_statements, tuning opts, pg_stat_kcache #38  
- [ ] D005 track_io_timing, …, auto_explain  #39  
- [ ] D006 Recommended DBA toolsets: postgres_dba, pgCenter, pgHeroother  #40  
- [ ] D007 Postgres-specific tools for troubleshooting  #137  

## E. WAL, Checkpoints

- [ ] E001 WAL/checkpoint settings, IO  #41   
- [ ] E002 Checkpoints, bgwriter, IO  #42  

## F. Autovacuum, Bloat

- [x] F001 < F003 Current autovacuum-related settings  #108, #164    
- [x] F002 < F007 Transaction wraparound check  #16, #171  
- [x] F003 < F006 Dead tuples  #164   
- [x] F004 < F001 Heap bloat estimation #87, #122  
- [x] F005 < F002 Index bloat estimation #88  
- [ ] F006 < F004 Precise heap bloat analysis 
- [ ] F007 < F005 Precise index bloat analysis 
- [x] F008 < F008 Resource usage (CPU, Memory, disk IO) #44  

## G. Performance / Connections / Memory-related Settings 

- [x] G001 Memory-related settings #45, #190  
- [x] G002 Connections #46  
- [x] G003 Timeouts, locks, deadlocks (amount) #47  
- [ ] G004 Query planner (diff) #48   
- [ ] G005 I/O settings #49   
- [ ] G006 Default_statistics_target (plus per table?) #50   

## H. Index Analysis

- [x] H001 Indexes: invalid #192, #51  
- [x] H002 Unused and redundant indexes #51, #180, #170, #168, #322  
- [x] H003 Missing FK indexes #52, #142, #173  

## J.  Capacity Planning

- [ ] J001 Capacity planning - #54  

## K. SQL query Analysis

- [x] K001 Globally aggregated query metrics #158, #178, #182, #184  
- [x] K002 Workload type ("first word" analysis) #159, #178, #179, #182, #184  
- [x] K003 Top-50 queries by total_time  #160, #172, #174, #178, #179, #182, #184, #193

## L. DB Schema Analysis
- [x] L001 (was: H003) Current sizes of DB objects (tables, indexes, mat. views)  #163  
- [ ] L002 (was: H004) Data types being used #53  
- [x] L003 Integer (int2, int4) out-of-range risks in PKs // calculate capacity remained; optional: predict when capacity will be fully used) https://gitlab.com/postgres-ai-team/postgres-checkup/issues/237

## TODO:

- [ ] DB schema, DDL, DB schema migrations

---

# Ideas :bulb: :bulb: :bulb:  :thinking\_face: 

- analyze all FKs and check if data types of referencing column and referenced one match (same thing for multi-column FKs)
- tables w/o PKs? tables not having even unique index?

## PostgreSQL:

- ready to archive WAL files (count) (need FS access) on master
- standby lag in seconds

## OS:

- FS settings (mount command parsing)
- meltdown/spectre patches
- swap settings
- memory pressure settings
- overcommit settings
- NUMA enabled?
- Huge pages?
- Transparent huge pages?
