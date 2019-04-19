# A. General  / Infrastructural

This group determines the available resources such as hardware characteristics or OS version and settings, as well as current resource usage and PostgreSQL settings. This information helps to identify suboptimal settings and existing or potential issues related to environment where Postgres works. It may be used for creating recommendations for improving database performance. A-group also serves as a basis for some reports in other groups.

### A001 System, CPU, RAM, Virtualization

General information about operational systems where the observed Postgres master and its replicas operate.

Insights:

- Hardware and software differences (OS versions, Linux kernel versions, CPU, Memory). If the observed master and its replicas run on different platforms, it might cause issues with binary replication.
 
- Memory settings tuning. (Examples: is swap enabled? Are huge pages used?) Observing state of memory about memory consumption by database may lead to recommendations of changes to improve system performance.
 
- Information about virtualization type.  

### A002 Postgres Version Information

This report answers the following questions:
- Do all nodes have the same Postgres version?  
- Is the minor version being used up-to-date? Keeping the minor version of the database up-to-date is recommended to decrease chances to encounter with bugs, performance and security issues?
- Is the major version currently supported by the community?  
- Will the major version be supported by the community during the next 12 months?
- If the minor version is not the most recent, are any critical bugfixes released that need to be applied ASAP?


### A003 Postgres Settings

Shows all PostgreSQL settings and their values grouped into categories.
May be used as a quick reference.


### A004 Cluster Information

A quick overview of "what is happening with the observed Postgres nodes?".

The following is included:

- The uptime. Sometimes low uptime may indicate an unplanned, accidental restart of the database.
- General information: how many databases are on one instance, what is their size, replication mode, age of statistics.
- Information about replicas, replication modes, replication delays.
- Ratio of forced checkpoints among all checkpoints registered since statistics reset time. Insights: Frequent checkpoints in most cases create an excessive load on the disk subsystem. Identifying this fact will allow the more optimal disk utilization.
- How big is the observed database (the cluster may have multiple databases)? Insight: if the database is smaller than RAM, there are good chances to avoid intensive disk IO in most operations
- Cache Effectiveness: percentage of buffer pool hits. Insight: if it is not more than 95% on all nodes, it might be a good sign that the buffer pool size needs to be increased.  
- Successful Commits: percentage of successfully committed transactions. Insight: if the value is not more than 99%, it might be a sign of logic issues with application code leading to high rates of ROLLBACK events.
- Temp Files per day: how many temporary files were generated per day in average, since last statistics reset time. Insight: if this value is high (thousands), it is a signal that work_mem should be increased.
- Deadlocks per day. Insight: significant (dozens) daily number of deadlocks is a sign of issues with application logic that needs redesign.

### A005 Extensions

Provides a list of all available and installed (in the current observed database) extensions, with versions. Insight: if there is a newer version of an installed extension, the report will highlight it, meaning that update is needed.

### A006 Postgres Setting Deviations

Helps to check that there are no differences in Postgres configuration on the observed nodes (except `transaction_read_only` and pg_stat_kcache’s `linux_hz`). 

Insights:
- In general, any differences in configuration on master and its replicas might lead to issues in case of failover. An example: the master is tuned, while replicas are not tuned at all or tuned poorly, in the event of failover, a new master cannot operate properly due to poor tuning.

### A007 Altered Settings

There are multiple ways to change database settings globally:
- explicitly, in the configuration file postgresql.conf, and
- implicitly, using 'ALTER SYSTEM' commands.

This report checks if there are settings which were set by implicit (ALTER SYSTEM) way.  

Possible sources of configuration settings (presented in the first column of the report’s table):

* `postgresql.auto.conf`: Changed via 'ALTER SYSTEM' command.
* `%any other file pattern%`: Changed in additional config included to the main one.
* `postgresql.conf`: Non-default values are set in postgresql.conf.
