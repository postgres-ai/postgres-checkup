package g002

const MSG_NODE string = "node: %s"
const MSG_NODES string = "nodes: %s"
const MSG_TX_AGE_MORE_1H_CONCLUSION string = "[P1] There are transactions with transaction age > 1 hour. It was observed on the following " +
	"%s. For OLTP databases, it is important to avoid long-lasting transactions. At the moment of this report generation, " +
	"such transactions were detected. Long-lasting  transactions lead to two big issues in the database, both affecting the system " +
	"performance negatively:  \n" +
	"    - higher risks of having locking issues (unless such transactions are read-only and do not involve explicit locking),  \n" +
	"    - VACUUM cannot process some entries in tables and indexes, and as a result, bloat levels increase more and faster than usual."
const MSG_TX_AGE_MORE_1H_RECOMMENDATION string = "[P1] There are transactions with transaction age > 1 hour. For better understanding, " +
	"refer to monitoring (add transaction-related graphs there if they are missing; it is important to split data by `state` " +
	"in `pg_stat_activity`). Consider the following tactics to avoid long-lasting transactions:  \n" +
	"    - split long transactions to smaller ones – ideally, OLTP workload should not have transactions longer than a few seconds;  \n" +
	"    - if long-lasting transactions often appear in 'idle in transaction' state, this is a sign that delays happen on " +
	"client's (application's) side; it is very important to reduce such delays as possible;  \n" +
	"    - if long-lasting transactions are read-only (for example, dumping tables using pg_dump, exporting data using " +
	"regular `SELECT` or `COPY`, or building some analytical reports), consider offloading this work to a replica; it is important " +
	"that such replica works with `host_standby_feedback = off` and is allowed to lag significantly in applying WALs.  \n"
