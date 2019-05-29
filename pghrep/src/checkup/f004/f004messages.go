package f004

const MSG_NO_RECOMMENDATIONS string = "All good, no recommendations here.  \n"
const MSG_TOTAL_BLOAT_EXCESS_CONCLUSION string = "[P1] Total table (heap) bloat estimation is %s, it is %.2f%% of overall DB size. " +
	"So removing the table bloat can help to reduce the total database size to ~%s and to increase the free disk space by %s. " +
	"Notice that this is only an estimation, sometimes it may be significantly off. Total size of tables is %.2f times bigger than it could be. " +
	"The most bloated tables are:  \n %s  \n"
const MSG_TOTAL_BLOAT_LOW_CONCLUSION string = "The total table (heap) bloat estimate is quite low, just ~%.2f%% (~%s). Hooray! Keep watching it though.  \n"
const MSG_BLOAT_CRITICAL_RECOMMENDATION string = "[P1] Reduce and prevent high level of table bloat:  \n" +
	"    - tune autovacuum: consider more aggressive autovacuum settings (See F001)  \n" +
	"    - reduce table bloat using one of the approaches mentioned below.  \n"
const MSG_BLOAT_WARNING_RECOMMENDATION string = "[P2] To resolve the table bloat issue do both of the following action items:  \n" +
	"    - to prevent high level of bloat in the future, tune autovacuum: consider more aggressive autovacuum settings (see F001);  \n" +
	"    - get rid of current table bloat using one of the approaches mentioned below.  \n"
const MSG_BLOAT_GENERAL_RECOMMENDATION_1 string = "If you want to get exact bloat numbers, clone the database, get table sizes, then apply " +
	"`VACUUM FULL` and get new table sizes. This will give the most reliable numbers.  \n"
const MSG_BLOAT_GENERAL_RECOMMENDATION_2 string = "To reduce the table bloat, consider using one of the following:  \n" +
	"    - [`VACUUM FULL`](https://www.postgresql.org/docs/OUR_MAJOR_VERSION/sql-vacuum.html) (:warning:  requires downtime / maintenance window),  \n" +
	"    - one of the tools reducing the bloat online, without interrupting the operations:  \n" +
	"        - [pg_repack](https://github.com/reorg/pg_repack),  \n" +
	"        - [pg_squeeze](https://github.com/reorg/pg_repack),  \n" +
	"        - [pgcompacttable](https://github.com/dataegret/pgcompacttable).  \n"
const MSG_BLOAT_PX_RECOMMENDATION string = "Read more on this topic:  \n" +
	"    - [Bloat estimation for tables](http://blog.ioguix.net/postgresql/2014/09/10/Bloat-estimation-for-tables.html) (2014, ioguix)  \n" +
	"    - [Show database bloat](https://wiki.postgresql.org/wiki/Show_database_bloat) (PostgreSQL wiki)\n" +
	"    - [PostgreSQL Bloat: origins, monitoring and managing](https://www.compose.com/articles/postgresql-bloat-origins-monitoring-and-managing/) (2016, Compose)  \n" +
	"    - [Dealing with significant Postgres database bloat — what are your options?](https://medium.com/compass-true-north/dealing-with-significant-postgres-database-bloat-what-are-your-options-a6c1814a03a5) (2018, Compass)  \n" +
	"    - [Postgres database bloat analysis](https://about.gitlab.com/handbook/engineering/infrastructure/blueprint/201901-postgres-bloat/) (2019, GitLab)  \n"
const MSG_BLOAT_WARNING_CONCLUSION string = "[P2] There are some tables with size > 1 MiB and table bloat estimate >= %.0f%% and < %.0f%%:  \n%s  \n"
const MSG_BLOAT_CRITICAL_CONCLUSION string = "[P1] The following tables have significant size (>1 MiB) and bloat estimate > %.0f%%:  \n%s  \n"
const TABLE_DETAILS string = "    - `%s`: size %s, can be reduced %.2f times, by ~%s (~%.0f%%)  \n"
