package f001

const MSG_AUTOVACUUM_NOT_TUNED_CONCLUSION string = "[P1] Autovacuum is not well-tuned. The following parameters are default, meaning that autovacuum behavior is far from optimal for an OLTP workload leading to higher levels of bloat in tables and indexes, lagging statistics:  \n%s  \n"
const MSG_AUTOVACUUM_NOT_TUNED_RECOMMENDATION string = "[P1] Autovacuum is not well-tuned. Consider its tuning for your workload. The links below can be helpful.  \n"
const MSG_AUTOVACUUM_TUNE_RECOMMENDATION string = "Useful links related to autovacuum tuning:  \n" +
	"    - [PostgreSQL Documentation. 19.10. Automatic Vacuuming](https://www.postgresql.org/docs/current/runtime-config-autovacuum.html)  \n" +
	"    - [Autovacuum Tuning Basics](https://www.2ndquadrant.com/en/blog/autovacuum-tuning-basics/) (2ndQuadrant, 2017)  \n" +
	"    - [Visualizing & Tuning Postgres Autovacuum](https://pganalyze.com/blog/visualizing-and-tuning-postgres-autovacuum) (pganalyze, 2017)  \n" +
	"    - [A Case Study of Tuning Autovacuum in Amazon RDS for PostgreSQL](https://aws.amazon.com/ru/blogs/database/a-case-study-of-tuning-autovacuum-in-amazon-rds-for-postgresql/) (AWS, 2018)  \n" +
	"    - [Understanding Autovacuum](https://www.youtube.com/watch?v=GqrBp0gyNHs) (video, 55 min, Citus Data, 2016)  \n"
