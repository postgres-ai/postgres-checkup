package k000

const MSG_NODE string = "node: %s"
const MSG_NODES string = "nodes: %s"
const MSG_EXCESS_QUERY_TOTAL_TIME_CONCLUSION string = "[P1] For some query groups, `total_time` > %.2f%% of overall timing, observed on: %s. " +
	"Such a high percentage means that those queries are \"major contributors\" to resource consumption on those nodes. In other words, " +
	"if a query group has `total_time` which is %.2f%% of overall timing, it means that during the observation period, %.2f%% of time CPUs were " +
	"working on that node processing queries from this group.\n"

const MSG_EXCESS_QUERY_TOTAL_TIME_RECOMMENDATION string = "[P1] For some query groups, `total_time` > %.2f%% of overall timing. To reduce `total_time` " +
	"for particular query group consider the following tactics:  \n" +
	"    - perform query micro-optimization (take particular query examples related to the group, use `EXPLAIN` and `EXPLAIN (BUFFERS, " +
	"ANALYZE)` to optimize it; consider using [Joe](https://gitlab.com/postgres-ai/joe) to boost the optimization process);\n" +
	"    - if the frequency of execution is high (check the `calls / second` metric), try to find a way to reduce the frequency, changing the " +
	"application code and/or, if applicable, applying caching.\n"
