package f008

const MSG_MAX_WORKERS_LOW_CONCLUSION string = "[P1] Maximum number of autovacuum workers is too low. System has %d CPUs (see A001), but only 3 autovacuum" +
	" workers are allowed (according to the `autovacuum_max_workers` setting) . This means that with increasing" +
	" workload, if significant part of this workload consists of modifying queries, autovacuum might start" +
	" lagging, maxing out `autovacuum_max_workers` value.  \n"

const MSG_MAX_WORKERS_LOW_RECOMMENDATION string = "[P1] Maximum number of autovacuum workers is too low. Consider raising `autovacuum_max_workers`.  \n" +
	"Consider values from %d to %d, depending on your workload. However, if 100%% or almost 100%% of the " +
	"workload is read-only queries and/or number of tables is less than 10, it might make sense to leave the " +
	"default value of `autovacuum_max_workers`. Refer to K001 to understand the workload and to L001 to see " +
	"the table list with sizes."

const MSG_AUTOVACUUM_COST_DELAY_NOT_TUNED_CONCLUSION string = "[P1] Autovacuum cost delay and limit are not tuned. In Postgres versions prior to version 12, " +
	"the effective values of `autovacuum_vacuum_cost_limit` and `autovacuum_vacuum_cost_delay` are too " +
	"conservative, so autovacuum is throttled. Roughly speaking, the default settings mean, that all " +
	"autovacuum workers (except those which process tables with individual, per-table settings) can read data with combined read throughput only up to " +
	"~8 MiB/s. This is extremely low for modern disk systems, and with growing sizes of tables and indexes it " +
	"might lead to cases when some objects are processed by autovacuum during many hours: it will take ~30 hours " +
	"to read 1 TiB of data if the allowed throughput is just 8 MiB/s. In some cases, it might lead to serious issues " +
	"such as performance degradation due to inability to process big tables in timed fashion and growing bloat, " +
	"and even to such critical issues as transaction ID wraparound. At the same time, if disk system is powerful " +
	"enough these risks can be easily mitigated by reducing throttling for autovacuum. In Postgres 12, it was " +
	"decided to reduce `autovacuum_vacuum_cost_delay` 10 times (from 20 ms to 2 ms)."

const MSG_AUTOVACUUM_COST_DELAY_NOT_TUNED_RECOMMENDATION string = "[P1] Autovacuum cost delay and limit are not tuned. Consider raising `autovacuum_vacuum_cost_limit` or " +
	"reducing `autovacuum_vacuum_cost_delay`: for example, consider reducing `autovacuum_vacuum_cost_delay` to 2 milliseconds, as it was done with default  " +
	"settings in Postgres 12. For more fine-grained tuning, analyze disk capabilities (first of all, random read and random write troughput) and logs with " +
	"`log_autovacuum_min_duration = 0`, and do a few iterations of tuning."

const MSG_AUTOVACUUM_COST_DELAY_TUNE_RECOMMENDATION string = "Useful links related to autovacuum tuning:  \n" +
	"    - [PostgreSQL Documentation. 19.10. Automatic Vacuuming](https://www.postgresql.org/docs/%.1f/runtime-config-autovacuum.html)  \n" +
	"    - [Autovacuum Tuning Basics](https://www.2ndquadrant.com/en/blog/autovacuum-tuning-basics/) (2ndQuadrant, 2017)  \n" +
	"    - [Visualizing & Tuning Postgres Autovacuum](https://pganalyze.com/blog/visualizing-and-tuning-postgres-autovacuum) (pganalyze, 2017)  \n" +
	"    - [A Case Study of Tuning Autovacuum in Amazon RDS for PostgreSQL](https://aws.amazon.com/ru/blogs/database/a-case-study-of-tuning-autovacuum-in-amazon-rds-for-postgresql/) (AWS, 2018)  \n" +
	"    - [Understanding Autovacuum](https://www.youtube.com/watch?v=GqrBp0gyNHs) (video, 55 min, Citus Data, 2016)"
