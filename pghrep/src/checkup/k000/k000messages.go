package k000

const MSG_NODE string = "node: %s"
const MSG_NODES string = "nodes: %s"
const MSG_EXCESS_QUERY_TOTAL_TIME_CONCLUSION string = "[P1] For some query groups, total_time > %.0f%% of overall timing. It was observed on the following %s. " +
	"Such a high percentage means that those queries are \"major contributors\" to resource consumption on those nodes. In other words, " +
	"if a query group has `total_time` which is %.0f%% of overall timing, it means that during the observation period, %.0f%% of time CPUs were " +
	"working on that node processing queries from this group.  \n"

const MSG_EXCESS_QUERY_TOTAL_TIME_RECOMMENDATION string = "[P1] For some query groups, total_time > %.0f%% of overall timing. To reduce `total_time` " +
	"for particular query group consider the following tactics:  \n" +
	"    - perform query micro-optimization (take particular query examples related to the group, use `EXPLAIN` and `EXPLAIN (BUFFERS, " +
	"ANALYZE)` to optimize it, also consider using [Joe bot](https://gitlab.com/postgres-ai/joe) and special DB instances) simplifying " +
	"the process of for query optimization, \n" +
	"    - if the frequency of execution is high (check the `calls / second` metric), consider reducing this frequency, changing the " +
	"application code and/or, if it is applicable, applying query results caching.  \n"
