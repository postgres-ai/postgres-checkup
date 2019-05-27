package f002

const MSG_RISKS_ARE_HIGH_CONCLUSION string = "[P1] Risks of transaction ID wraparound are high for:  \n" +
	"%s  \n" +
	"Approaching 100%% leads to downtime: the system will shut down and refuse to start any new transactions.  \n" +
	"More on this topic:  \n" +
	"- [PostgreSQL Documentation. 24.1.5. Preventing Transaction ID Wraparound Failures](https://www.postgresql.org/docs/current/routine-vacuuming.html#VACUUM-FOR-WRAPAROUND)  \n" +
	"- [The Internals of PostgreSQL. Chapter 5. Concurrency Control. 5.10.1. FREEZE Processing](http://www.interdb.jp/pg/pgsql05.html#_5.10.1.)  \n" +
	"- [Transaction ID Wraparound in Postgres](https://blog.sentry.io/2015/07/23/transaction-id-wraparound-in-postgres) (2015, Sentry blog)  \n" +
	"- [Autovacuum wraparound protection in PostgreSQL](https://www.cybertec-postgresql.com/en/autovacuum-wraparound-protection-in-postgresql/) (2017, Cybertec blog)  \n" +
	"- [What We Learned from the Recent Mandrill Outage](https://mailchimp.com/what-we-learned-from-the-recent-mandrill-outage/) (2019, Mailchimp blog)  \n" +
	"- [Managing Transaction ID Exhaustion (Wraparound) in PostgreSQL](https://info.crunchydata.com/blog/managing-transaction-id-wraparound-in-postgresql) (2019, Crunchy Data blog)  \n"

const MSG_RISKS_ARE_HIGH_RECOMMENDATION string = "[P1] To minimize risks of transaction ID wraparound do the following:  \n" +
	"1. Run `VACUUM FREEZE` for mentioned tables.  \n" +
	"2. Perform autovacuum tuning to ensure that autovacuum has enough resources and runs often enough to minimize risks of transaction ID wraparound. Read articles provided in the \"Conclusions\" section for more details. "

const MSG_NO_RECOMMENDATIONS string = "No recommendations."
