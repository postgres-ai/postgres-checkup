package f005

const MSG_NO_RECOMMENDATIONS string = "All good, no recommendations here."
const MSG_TOTAL_BLOAT_EXCESS_CONCLUSION string = "[P1] Total index bloat estimation is %s, it is %.2f%% of overall DB size. " +
	"So removing the index bloat can help to reduce the total database size to ~%s and to increase the free disk space by %s. " +
	"Notice that this is only an estimation, sometimes it may be significantly off. Total size of indexes is %.2f times bigger than it could be. " +
	"The most bloated indexes are:  \n %s  \n"
const MSG_TOTAL_BLOAT_LOW_CONCLUSION string = "The total index bloat estimate is quite low, just ~%.2f%% (~%s). Hooray! Keep watching it though."
const MSG_BLOAT_CRITICAL_RECOMMENDATION string = "[P1] Reduce and prevent high level of index bloat:  \n" +
	"    - tune autovacuum: consider more aggressive autovacuum settings (See F001)  \n" +
	"    - reduce index bloat using one of the approaches mentioned below.  \n"
const MSG_BLOAT_WARNING_RECOMMENDATION string = "[P2] Reduce and prevent high level of index bloat:  \n" +
	"    - tune autovacuum: consider more aggressive autovacuum settings (See F001)  \n" +
	"    - reduce index bloat using one of the approaches mentioned below.  \n"
const MSG_BLOAT_GENERAL_RECOMMENDATION_1 string = "If you want to get exact bloat numbers, clone the database, get index sizes, then apply " +
	"`VACUUM FULL` and get new index sizes. This will give the most reliable numbers.  \n"
const MSG_BLOAT_GENERAL_RECOMMENDATION_2 string = "To reduce the index bloat, consider using one of the following:  \n" +
	"    - [`VACUUM FULL`](https://www.postgresql.org/docs/OUR_MAJOR_VERSION/sql-vacuum.html) (:warning:  requires downtime / maintenance window),  \n" +
	"    - [`REINDEX`](https://www.postgresql.org/docs/OUR_MAJOR_VERSION/sql-reindex.html) (`REINDEX INDEX`, `REINDEX TABLE`; :warning:  requires downtime / maintenance window),  \n" +
	"    - `REINDEX CONCURRENTLY` << ONLY IF MAJOR VERSION IS >= 12  \n" +
	"    - recreating indexes online using `CREATE INDEX CONCURRENTLY`, `DROP INDEX CONCURRENTLY` and renaming (not trivial for indexes supporting PK, FK),  << ONLY IF MAJOR VERSION IS < 12  \n" +
	"    - one of the tools reducing the bloat online, without interrupting the operations:  \n" +
	"        - [pg_repack](https://github.com/reorg/pg_repack),  \n" +
	"        - [pg_squeeze](https://github.com/reorg/pg_repack),  \n" +
	"        - [pgcompacttable](https://github.com/dataegret/pgcompacttable).  \n"
const MSG_BLOAT_PX_RECOMMENDATION string = "Read more on this topic:  \n" +
	"    - [Index maintenance](https://wiki.postgresql.org/wiki/Index_Maintenance) (PostgreSQL wiki)  \n" +
	"    - [Btree bloat query](http://blog.ioguix.net/postgresql/2014/11/03/Btree-bloat-query-part-4.html) (2014, ioguix)  \n" +
	"    - [PostgreSQL Index bloat under a microscope](http://pgeoghegan.blogspot.com/2017/07/postgresql-index-bloat-microscope.html) (2017, Peter Geoghegan)  \n" +
	"    - [PostgreSQL Bloat: origins, monitoring and managing](https://www.compose.com/articles/postgresql-bloat-origins-monitoring-and-managing/) (2016, Compose)  \n" +
	"    - [Dealing with significant Postgres database bloat — what are your options?](Dealing with significant Postgres database bloat — what are your options?) (2018, Compass)  \n" +
	"    - [Postgres database bloat analysis](https://about.gitlab.com/handbook/engineering/infrastructure/blueprint/201901-postgres-bloat/) (2019, GitLab)  \n"
const MSG_BLOAT_WARNING_CONCLUSION string = "[P2] There are some indexes with size > 1 MiB and index bloat estimate >= %.0f%% and < %.0f%%:  \n%s  \n"
const MSG_BLOAT_CRITICAL_CONCLUSION string = "[P1] The following indexes have significant size (>1 MiB) and bloat estimate > %.0f%%:  \n%s  \n"
const INDEX_DETAILS string = "    - `%s`: size %s, can be reduced %.2f times, by ~%s (~%.0f%%)  \n"
