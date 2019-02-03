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
* golang >= 1.8
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

A human-readable report can be found at:

```bash
./artifacts/prod1/e1_full_report.md
```

Open it with your favorite Markdown files viewer or just upload to a service such as gist.github.com.

# The Full List of Checks

## А. General  / Infrastructural

- [x] A001 System, CPU, RAM, disks, virtualization #6 , #56 , #57 , #86 
- [x] A002 PostgreSQL Versions (Simple) #68, #21, #86
- [x] A003 Collect pg_settings  #15, #167, #86 
- [x] A004 General cluster info  #7, #58, #59, #86, #162  
- [x] A005 Extensions #8, #60, #61, #86, #167   
- [x] A006 Config diff  #9, #62, #63, #86  
- [x] A007 Alter system vs postgresql.conf #18, #86  
- [x] A008 Disk usage and file system type #19, #20  
- [ ] A010 Data checksums are not enabled + wal_log_hints #22  
- [ ] A011 Connection pooling. PgBouncer #23  
- [ ] A012 Anti crash checks #177  

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
- [ ] D002 Useful Linux tools  #36  
- [ ] D003 List of monitoring metrics #37  
- [x] D004 pg_stat_statements, tuning opts, pg_stat_kcache #38  
- [ ] D005 track_io_timing, …, auto_explain  #39  
- [ ] D006 Postgres_dba / other toolset - recommend  #40  
- [ ] D007 Postgres-specific tools for troubleshooting  #137  

## E. WAL, Checkpoints

- [ ] E001 WAL/checkpoint settings, IO  #41   
- [ ] E002 Bgwriter, IO  #42  

## F. Autovacuum, Bloat

- [x] F001 < F003 Current autovacuum-related settings  #108, #164    
- [x] F002 < F007 Transaction wraparound check  #16, #171  
- [x] F003 < F006 Autovacuum dead tuples  #164   
- [x] F004 < F001 Heap Bloat estimation #87, #122  
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
- [ ] G006 Default_statistics_target (per table?) #50   

## H. Index Analysis

- [ ] H001 Indexes: invalid #192, #51  
- [x] H002 < H001 Indexes: unused, redundant #51, #180, #170, #168  
- [x] H003 < H002 Missing FK indexes #52, #142, #173  

## J.  Capacity Planning

- [ ] J001 Capacity planning - #54  

## K. SQL query Analysis

- [ ] K001 Globally aggregated query metrics #158, #178, #182, #184  
- [ ] K002 Workload type ("first word" analysis) #159, #178, #179, #182, #184  
- [x] K003 Top-50 queries by total_time  #160, #172, #174, #178, #179, #182, #184, #193

## L. DB Schema Analysis
- [x] L001 (was: H003) Current sizes of DB objects (tables, indexes, mat. views)  #163  
- [ ] L002 (was: H004) Data types being used #53  
- [ ] L003 Integer (int2, int4) out-of-range risks in PKs // calculate capacity remained; optional: predict when capacity will be fully used) https://gitlab.com/postgres-ai-team/postgres-checkup/issues/237

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
