package g001

const MSG_SHARED_BUFFERS_NOT_OPTIMAL_CONCLUSION string = "[P1] Buffer pool size (`shared_buffers`) is far from optimal:\n%s.\n"
const MSG_SHARED_BUFFERS_NOT_OPTIMAL_RECOMMENDATION string = "[P1] Buffer pool size (`shared_buffers`)  value is far from optimal. " +
	"Consider conducting experiments in a special environment (a clone of production DB, replayed or simulated production workload) " +
	"to find optimal `shared_buffers` values for each server. Recommended values of `shared_buffers`:\n%s.\n"
const MSG_HOST_CONCLUSION_HIGH string = "    - server `%s` has %s of RAM, while `shared_buffers` is set to %s, or %.2f%% of RAM – it is too high, " +
	"so memory might not be enough, and [OOM killer](https://en.wikipedia.org/wiki/Out_of_memory) might kill Postgres processes if swapping is disabled."
const MSG_HOST_CONCLUSION_LOW string = "    - server `%s` has %s of RAM, while `shared_buffers` is set to %s, or %.2f%% of RAM – it is too low, " +
	"so it is very likely that Postgres performance is now sub-optimal."
const MSG_HOST_RECOMMENDATION string = "    - server `%s`: %s (%d%%) or a value between %s (%d%%) and %s (%d%%)"
const MSG_TUNE_SHARED_BUFFERS_RECOMMENDATION string = "Useful links related to buffer pool tuning:\n" +
	"    - [PostgreSQL documentation. 19.4. Resource Consumption](https://www.postgresql.org/docs/current/runtime-config-resource.html)\n" +
	"    - [Tuning Your PostgreSQL Server](https://wiki.postgresql.org/wiki/Tuning_Your_PostgreSQL_Server#shared_buffers) (PostgreSQL Wiki)\n" +
	"    - [annotated.conf](https://github.com/jberkus/annotated.conf) (Josh Berkus, 2018)\n"

const MSG_OOM_BASE_CONCLUSION string = "[P1] Potentially high risks of OOM. Memory-related settings on `%s` server look risky: there are potentially " +
	"high risks to have [OOO (out of memory)](https://en.wikipedia.org/wiki/Out_of_memory).\n"
const MSG_OOM_BASE_RECOMMENDATION string = "[P1] Potentially high risks of OOM. Reconsider memory-related settings to minimize risks of OOM.\n"
const MSG_TUNE_MEMORY_RECOMMENDATION string = "Useful links related to memory-related settings:\n" +
	"    - [PostgreSQL documentation. 19.4. Resource Consumption](https://www.postgresql.org/docs/current/runtime-config-resource.html)\n"

const MSG_OOM_SWAP_ENABLED string = "Since swapping is enabled (see A001), it might lead to significant performance degradation.\n"
const MSG_OOM_SWAP_DISABLED string = "Since swapping is disabled (see A001), it might lead to Postgres crashes due to OOM killer's activity.\n"
const MSG_OOM_SHARED_BUFFERS string = "`shared_buffers` is set to %s, which is %.2f%%%% of RAM, making this setting a major contributor to overall memory consumption.  \n"
const MSG_OOM_WORK_MEM_CONNECTIONS string = "`work_mem` is set to %s, and each DB session may use up to this value of memory multiple times " +
	"(for example, if multiple ordering operations with massive data sets are needed), so in case Postgres backends is maxed out " +
	"(`max_connections` value is %d), all backends might consume, say,  `max_connections * 2 * work_mem = %s`, which is %.2f%%%% of RAM. " +
	"It makes `work_mem/max_connections` pair a major contributor to the overall memory consumption.\n"
const MSG_OOM_AUTIVACUUM_WORKMEM_BEGIN string = "`autovacuum_work_mem` is %s"
const MSG_OOM_AUTIVACUUM_WORKMEM_NOTSET string = "(it's set to `-1` so the actual value is inherited from `maintenance_work_mem`)"
const MSG_OOM_AUTIVACUUM_WORKMEM_END string = "and maximum %d autovacuum workers may work simultaneously, so together they may consume up to %s, " +
	"or %.2f%%%% of RAM. It makes `autovacuum_work_mem/autovacuum_max_workers` pair a major contributor to the overall memory consumption.\n"
const MSG_OOM_BASE_RECOMMENDATION_DETAIL string = "First of all, pay attention to the following settings:  \n"
