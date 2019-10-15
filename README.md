### Demo: [an example of postgres-checkup report](https://gitlab.com/postgres-ai/postgres-checkup-tests/tree/master/1.2.2) (based on CI, multi node).

***Disclaimer: Conclusions, Recommendations – work in progress.**
To treat the data correctly, you need deep Postgres knowledge. Each report
consists of 3 sections: Observations, Conclusions, and Recommendations.
Observations are filled automatically. As for Conclusions and Recommendations
sections, as of June 2019, only several reports have autogeneration for them.*


# About

Postgres Checkup ([postgres-checkup](https://gitlab.com/postgres-ai-team/postgres-checkup)) is a new kind of diagnostics tool for a deep analysis of a Postgres database health. It detects current and potential issues with database performance, scalability and security. It also produces recommendations on how to resolve or prevent them.

A monitoring system will only show current, urgent problems. And postgres-checkup will show sneaking up, deeper problems, that may hit you in the future. It helps to solve many known database administration problems and common pitfalls. It aims to detect issues at a very early stage and to suggest the best ways to prevent them. 
We recommend to run these on a regular basis — weekly, monthly, and quarterly. And also to run these right before and after applying any major change to a database server. Whether it’s a schema or configuration parameter or cluster settings change.


Why do you need postgres-checkup and why it's safe and easy to use:

- *It is unobtrusive*: its impact on the observing system is
close to zero. It does not use any heavy queries, keeping resource usage
very low, and avoiding having the [“observer effect”](https://en.wikipedia.org/wiki/Observer_effect_(information_technology)).
postgres-checkup reports were successfully tested on real-world databases
containing 500,000+ tables and 1,000,000+ indexes.

- *Zero install* (on observed machines): it is able to analyze any Linux
machine (including virtual machines), as well as cloud Postgres instances
(such as Amazon RDs or Google Cloud SQL), not requiring any additional setup
or any changes. It does, hovewer, require a privileged access that a DBA usually
has anyway.

- *Complex analysis*: unlike most monitoring tools, which provide just raw data,
postgres-checkup combines data from various parts of the system (e.g.,
internal Postgres stats are combined with knowledge about system resources
in autovacuum setting and behavior analysis) joining the data into well-formatted
reports aimed to solve particular DBA problems. Also, it analyzes the master
database server together with all its replicas, which is neccessary in such
cases as index analysis or search for settings deviations.

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
1. "Conclusions": what we conclude from the Observations, stated in plain English
in the form that is convenient for engineers who are not DBA experts.
1. "Recommendations": action items, what to do to fix the discovered issues.

Both "Conclusions" and "Recommendations" are to be consumed by engineers who
will make decisions what, how and when to optimize.

# Installation and Usage

## Requirements

For the operator machine (from where the tool will be executed), the following
OS are supported:

* Linux (modern RHEL/CentOS or Debian/Ubuntu; others should work as well, but
are not yet tested);
* MacOS.

There are known cases when postgres-checkup was successfully used on Windows,
althought with some limitations.

The following programs must be installed on the operator machine:

* bash
* psql
* coreutils
* jq >= 1.5
* golang >= 1.8 (no binaries are shipped at the moment)
* awk
* sed
* pandoc *
* wkhtmltopdf >= 0.12.4 *

pandoc and wkhtmltopdf are optional, they are neededed for generating HTML and 
PDF versions of report (options `--html`, `--pdf`).

Nothing special has to be installed on the observed machines. However, they must
run Linux (again: modern RHEL/CentOS or Debian/Ubuntu; others should work as
well, but are not yet tested).

:warning: Only Postgres version 9.6 and higher are currently officially supported.

## How to Install

#### 1. Install required programs

Ubuntu/Debian:
```bash
sudo apt-get update -y
sudo apt-get install -y git postgresql coreutils jq golang

# Optional (to generate PDF/HTML reports)
sudo apt-get install -y pandoc
wget https://github.com/wkhtmltopdf/wkhtmltopdf/releases/download/0.12.4/wkhtmltox-0.12.4_linux-generic-amd64.tar.xz
tar xvf wkhtmltox-0.12.4_linux-generic-amd64.tar.xz
sudo mv wkhtmltox/bin/wkhtmlto* /usr/local/bin
sudo apt-get install -y openssl libssl-dev libxrender-dev libx11-dev libxext-dev libfontconfig1-dev libfreetype6-dev fontconfig
```

CentOS/RHEL:
```bash
sudo yum install -y git postgresql coreutils jq golang

# Optional (to generate PDF/HTML reports)
sudo yum install -y pandoc
wget https://github.com/wkhtmltopdf/wkhtmltopdf/releases/download/0.12.4/wkhtmltox-0.12.4_linux-generic-amd64.tar.xz
tar xvf wkhtmltox-0.12.4_linux-generic-amd64.tar.xz
sudo mv wkhtmltox/bin/wkhtmlto* /usr/local/bin
sudo yum install -y libpng libjpeg openssl icu libX11 libXext libXrender xorg-x11-fonts-Type1 xorg-x11-fonts-75dpi
```

MacOS (assuming that [Homebrew](https://brew.sh/) is installed):
```bash
brew install postgresql coreutils jq golang git

# Optional (to generate PDF/HTML reports)
brew install pandoc Caskroom/cask/wkhtmltopdf
```

#### 2. Clone this repo

```bash
git clone https://gitlab.com/postgres-ai/postgres-checkup.git
# Use --branch to use specific release version. For example, to use version 1.1:
#   git clone --branch 1.1 https://gitlab.com/postgres-ai/postgres-checkup.git
cd postgres-checkup
```

## Example of Use

Let's make a report for a project named `prod1`. Assume that we have two servers,
`db1.vpn.local` and `db1.vpn.local`.

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

Also, you can define a specific way to connect: SSH or `psql`:

`--ssh-hostname db2.vpn.local` - SSH will be used for the connection. SSH port can be defined as well
with option `--ssh-port`.

`--pg-hostname db2.vpn.local` - `psql` will be used for the connection. The port where PostgreSQL
accepts connections can be defined with the option `--pg-port`

In case when `--pg-port` or `--ssh-port` are not defined but `--port` is defined, value of `--port` option
will be used instead of `--pg-port` or `--ssh-port` depending on the current connection type.

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

You can collect and process data separately by specifying working mode name in CLI option `--mode %mode%` or using it as a "command" (`checkup %mode%`).  
Available working modes:  
    `collect` - collect data;
    `process` - generate MD (and, optionally, HTML, PDF) reports with conclusions and recommendations;
    `upload` - upload generated reports to Postgres.ai platform;
    `run` - collect and process data at once. This is the default mode, it is used when no other mode is specified. Note, that upload is not included.

## Docker 🐳

It's possible to use the `postgres-checkup` from a docker container.
The container will run, execute all checks and stop itself.
The check result can be found inside the `artifacts` folder in current directory (pwd).

### Usage with `docker run`

There is an option to run postgres-checkup in a Docker container:

```bash
docker run --rm \
  --name postgres-checkup \
  -e PGPASSWORD="postgres" \
  -v `pwd`/artifacts:/artifacts \
  registry.gitlab.com/postgres-ai/postgres-checkup:latest \
    ./checkup \
      -h hostname \
      -p 5432 \
      --username postgres \
      --dbname postgres \
      --project c \
      -e "$(date +'%Y%m%d')001"
```

In this case some checks (those requiring SSH connection) will be skipped.

If you want to have all supported checks, you have to use SSH access to the
target machine with Postgres database.

If SSH connection to the Postgres server is available, it is possible to pass
SSH keys to the docker container, so postgres-checkup will switch to working via
remote SSH calls, generating all reports (this approach is known to have issues
on Windows, but should work well on Linux and MacOS machines):

```bash
docker run --rm \
  --name postgres-checkup \
  -v "$(pwd)/artifacts:/artifacts" \
  -v "$(echo ~)/.ssh/id_rsa:/root/.ssh/id_rsa:ro" \
  registry.gitlab.com/postgres-ai/postgres-checkup:latest \
  ./checkup \
    -h sshusername@hostname \
    --username my_postgres_user \
    --dbname my_postgres_database \
    --project docker_test_with_ssh \
    -e "$(date +'%Y%m%d')001"
```

If you try to check the local instance of postgres on your host from a container,
you cannot use `localhost` in `-h` parameter. You have to use a bridge between
host OS and Docker Engine. By default, host IP is `172.17.0.1` in `docker0`
network, but it vary depending on configuration. More information [here](https://nickjanetakis.com/blog/docker-tip-65-get-your-docker-hosts-ip-address-from-in-a-container).

## Credits

Some reports are based on or inspired by useful queries created and improved by
various developers, including but not limited to:
 * Jehan-Guillaume (ioguix) de Rorthais https://github.com/ioguix/pgsql-bloat-estimation
 * Alexey Lesovsky, Alexey Ermakov, Maxim Boguk, Ilya Kosmodemiansky et al. from Data Egret (aka PostgreSQL-Consulting) https://github.com/dataegret/pg-utils
 * Josh Berkus, Quinn Weaver et al. from PostgreSQL Experts, Inc. https://github.com/pgexperts/pgx_scripts

Docker support implemented by [Ivan Muratov](https://gitlab.com/binakot).

# The Full List of Reports

## А. General  / Infrastructural

- [x] A001 System information #6 , #56 , #57 , #86
- [x] A002 Version information #68, #21, #86
- [x] A003 Postgres settings  #15, #167, #86
- [x] A004 Cluster information  #7, #58, #59, #86, #162
- [x] A005 Extensions #8, #60, #61, #86, #167
- [x] A006 Postgres setting deviations #9, #62, #63, #86
- [x] A007 Altered settings #18, #86
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
- [ ] C007 Replication slots. Lags. Standby feedbacks

## D. Monitoring / Troubleshooting

- [ ] D001 Logging (syslog?), log_*** #35
- [x] D002 Useful Linux tools  #36
- [ ] D003 List of monitoring metrics #37
- [x] D004 pg_stat_statements and pg_stat_kcache settings #38
- [ ] D005 track_io_timing, …, auto_explain  #39
- [ ] D006 Recommended DBA toolsets: postgres_dba, pgCenter, pgHeroother  #40
- [ ] D007 Postgres-specific tools for troubleshooting  #137

## E. WAL, Checkpoints

- [ ] E001 WAL/checkpoint settings, IO  #41
- [ ] E002 Checkpoints, bgwriter, IO  #42

## F. Autovacuum, Bloat

- [x] F001 < F003 Autovacuum: current settings  #108, #164
- [x] F002 < F007 Autovacuum: transaction ID wraparound check  #16, #171
- [x] F003 < F006 Autovacuum: dead tuples  #164
- [x] F004 < F001 Autovacuum: heap bloat (estimated) #87, #122
- [x] F005 < F002 Autovacuum: index bloat (estimated) #88
- [ ] F006 < F004 Precise heap bloat analysis
- [ ] F007 < F005 Precise index bloat analysis
- [x] F008 < F008 Autovacuum: resource usage #44

## G. Performance / Connections / Memory-related Settings

- [x] G001 Memory-related settings #45, #190
- [x] G002 Connections and current activity #46
- [x] G003 Timeouts, locks, deadlocks #47
- [ ] G004 Query planner (diff) #48
- [ ] G005 I/O settings #49
- [ ] G006 Default_statistics_target (plus per table?) #50

## H. Index Analysis

- [x] H001 Invalid indexes #192, #51
- [x] H002 Unused and redundant indexes #51, #180, #170, #168, #322
- [x] H003 Non-indexed foreign keys #52, #142, #173

## J.  Capacity Planning

- [ ] J001 Capacity planning - #54

## K. SQL query Analysis

- [x] K001 Globally aggregated query metrics #158, #178, #182, #184
- [x] K002 Workload Type ("The First Word" Analysis) #159, #178, #179, #182, #184
- [x] K003 Top-50 queries by total_time #160, #172, #174, #178, #179, #182, #184, #193

## L. DB Schema Analysis
- [x] L001 (was: H003) Table sizes #163
- [ ] L002 (was: H004) Data types being used #53
- [x] L003 Integer (int2, int4) out-of-range risks in PKs // calculate capacity remained; optional: predict when capacity will be fully used) https://gitlab.com/postgres-ai-team/postgres-checkup/issues/237
- [ ] L004 Tables without PK/UK
