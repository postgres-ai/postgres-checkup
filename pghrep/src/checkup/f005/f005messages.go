package f005

const MSG_NO_RECOMMENDATIONS string = "All good 👍"

const MSG_TOTAL_BLOAT_EXCESS_CONCLUSION string = "[P1] Total index (btree only) bloat estimation is ~%s, it is %.2f%% of overall indexes size and %.2f%% of the DB size. " +
	"Removing the index bloat will reduce the total DB size down to ~%s. Free disk space will be increased by ~%s. " +
	"Total size of indexes is %.2f times bigger than it could be. " +
	"Notice that this is only an estimation, sometimes it may be significantly off.\n"

const MSG_TOTAL_BLOAT_LOW_CONCLUSION string = "The estimated index (btree only) bloat in this DB is low, just ~%.2f%% (~%s). No action is needed now. Keep watching it though.\n"
const MSG_BLOAT_CRITICAL_RECOMMENDATION string = "[P1] Reduce and prevent a high level of index bloat:\n" +
	"    - to prevent a high level of bloat in the future, tune autovacuum: consider more aggressive autovacuum settings (see F001);\n" +
	"    - eliminate or reduce the current index bloat using one of the approaches listed below.\n"
const MSG_BLOAT_WARNING_RECOMMENDATION string = "[P2] Consider the following:\n" +
	"    - to prevent a high level of bloat in the future, tune autovacuum: consider more aggressive autovacuum settings (see F001);\n" +
	"    - eliminate or reduce the current index bloat using one of the approaches listed below.\n"
const MSG_BLOAT_GENERAL_RECOMMENDATION_1 string = "If you want to get exact bloat numbers, clone the database, get index sizes, then apply " +
	"database-wide `VACUUM FULL` (it eliminates all the bloat), and gets new table sizes. Then compare old and new numbers.\n"
const MSG_BLOAT_GENERAL_RECOMMENDATION_2 string = "To reduce the index bloat, consider one of the following approaches:\n" +
	"    - [`VACUUM FULL`](https://www.postgresql.org/docs/current/sql-vacuum.html) (:warning:  requires downtime / maintenance window),\n" +
	"    - [`REINDEX`](https://www.postgresql.org/docs/current/sql-reindex.html) (`REINDEX INDEX`, `REINDEX TABLE`; :warning:  requires downtime / maintenance window),\n" +
	"    - recreating indexes online using `CREATE INDEX CONCURRENTLY`, `DROP INDEX CONCURRENTLY` and renaming (not trivial for indexes supporting PK, FK) // `REINDEX CONCURRENTLY` is available in Postgres 12+,\n" +
	"    - one of the tools reducing the bloat online, without interrupting the operations:  \n" +
	"        - [pg_repack](https://github.com/reorg/pg_repack),\n" +
	"        - [pg_squeeze](https://github.com/reorg/pg_repack),\n" +
	"        - [pgcompacttable](https://github.com/dataegret/pgcompacttable).\n"
const MSG_BLOAT_PX_RECOMMENDATION string = "Read more on this topic:\n" +
	"    - [Index maintenance](https://wiki.postgresql.org/wiki/Index_Maintenance) (PostgreSQL wiki)\n" +
	"    - [Btree bloat query](http://blog.ioguix.net/postgresql/2014/11/03/Btree-bloat-query-part-4.html) (2014, ioguix)\n" +
	"    - [PostgreSQL Index bloat under a microscope](http://pgeoghegan.blogspot.com/2017/07/postgresql-index-bloat-microscope.html) (2017, Peter Geoghegan)\n" +
	"    - [PostgreSQL Bloat: origins, monitoring and managing](https://www.compose.com/articles/postgresql-bloat-origins-monitoring-and-managing/) (2016, Compose)  \n" +
	"    - [Dealing with significant Postgres database bloat — what are your options?](Dealing with significant Postgres database bloat — what are your options?) (2018, Compass)\n" +
	"    - [Postgres database bloat analysis](https://about.gitlab.com/handbook/engineering/infrastructure/blueprint/201901-postgres-bloat/) (2019, GitLab)\n"
const MSG_BLOAT_WARNING_CONCLUSION_1 string = "[P2] There is %d index with size > 1 MiB and index bloat estimate >= %.0f%% and < %.0f%%:  \n%s  \n"
const MSG_BLOAT_CRITICAL_CONCLUSION_1 string = "[P1] The following %d index has significant size (>1 MiB) and bloat estimate > %.0f%%:  \n%s  \n"
const MSG_BLOAT_WARNING_CONCLUSION_N string = "[P2] There are %d indexes with size > 1 MiB and index bloat estimate >= %.0f%% and < %.0f%%:  \n%s  \n"
const MSG_BLOAT_CRITICAL_CONCLUSION_N string = "[P1] The following %d indexes have significant size (>1 MiB) and bloat estimate > %.0f%%:  \n%s  \n"

const INDEX_DETAILS string = "    - `%s`: size %s, can be reduced %.2f times, by ~%s (~%.0f%%)\n"
