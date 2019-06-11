package g001

const MSG_SHARED_BUFFERS_NOT_OPTIMAL_CONCLUSION string = "[P1] shared_buffers value is far from optimal:  \n%s.  \n"
const MSG_SHARED_BUFFERS_NOT_OPTIMAL_RECOMMENDATION string = "[P1] shared_buffers value is far from optimal. Consider conducting experiments in a special environment (a clone of production DB, replayed or simulated production workload) to find optimal `shared_buffers` values for each server. Recommended values of `shared_buffers`:  \n%s  \n"
const MSG_HOST_CONCLUSION_HIGH string = "    - server `%s` has %s of RAM, while `shared_buffers` is set to %s, or %d%% of RAM – it is too high, so memory might not be enough, and [OOM killer](https://en.wikipedia.org/wiki/Out_of_memory) might kill Postgres processes since swap is disabled"
const MSG_HOST_CONCLUSION_LOW string = "    - server `%s` has %s of RAM, while `shared_buffers` is set to %s, or %d%% of RAM – it is too low, so Postgres performance on this server is sub-optimal"
const MSG_HOST_RECOMMENDATION string = "    - server `%s`: %s (%d%%) or a value between %s (%d%%) and %s (%d%%)"
const MSG_TUNE_SHARED_BUFFERS_RECOMMENDATION string = "Useful links related to `shared_buffers` tuning:  \n" +
	"    - [PostgreSQL documentation. 19.4. Resource Consumption](https://www.postgresql.org/docs/current/runtime-config-resource.html)  \n" +
	"    - [Tuning Your PostgreSQL Server](https://wiki.postgresql.org/wiki/Tuning_Your_PostgreSQL_Server#shared_buffers) (PostgreSQL Wiki)  \n" +
	"    - [annotated.conf](https://github.com/jberkus/annotated.conf) (Josh Berkus, 2018)  \n"
