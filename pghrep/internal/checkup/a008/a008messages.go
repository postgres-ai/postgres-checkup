package a008

const MSG_NO_USAGE_RISKS_CONCLUSION string = "Risks of running out of free disk space in the nearest future are low."
const MSG_USAGE_WARNING_CONCLUSION string = "[P2] `%s` on machine `%s` has %s of space used, it exceeds 70%%. Risks of running out of free disk space in the nearest future are significant. "
const MSG_USAGE_WARNING_RECOMMENDATION string = "[P2] Increase free space on `%s` on machine `%s`. It is recommended to keep more than %d%% disk space free " +
	"to reduce risks of running out of free disk space."
const MSG_USAGE_CRITICAL_CONCLUSION string = "[P1] `%s` on machine `%s` has %s of space used, it exceeds 90%%. Risks of running out of free disk space in the nearest future are high. " +
	"If it happens, PostgreSQL shuts down and a manual fix is required."
const MSG_USAGE_CRITICAL_RECOMMENDATION string = "[P1] Increase free space on `%s` on machine `%s` as soon as possible to prevent service outage."
const MSG_NETWORK_FS_CONCLUSION_1 string = "[P1] %s on host `%s` uses [NFS](https://en.wikipedia.org/wiki/Network_File_System). This might lead to serious issues with Postgres, including downtime and data corruption."
const MSG_NETWORK_FS_CONCLUSION_N string = "[P1] %s on host `%s` use [NFS](https://en.wikipedia.org/wiki/Network_File_System). This might lead to serious issues with Postgres, including downtime and data corruption."
const MSG_NETWORK_FS_RECOMMENDATION string = "[P1] Never use NFS to run Postgres."
const MSG_NOT_RECOMMENDED_FS_CONCLUSION_1 string = "[P3] %s on machine `%s` is located on drive where the following filesystems are used: %s. This might mean that Postgres performance and reliability characteristics are worse than it could be in case of use of more popular filesystems (such as ext4)."
const MSG_NOT_RECOMMENDED_FS_CONCLUSION_N string = "[P3] %s on machine `%s` are located on drives where the following filesystems are used: %s respectively. This might mean that Postgres performance and reliability characteristics are worse than they could be in case of use of more popular filesystems (such as ext4)."
const MSG_NOT_RECOMMENDED_FS_RECOMMENDATION string = "[P3] Consider using [ext4](https://en.wikipedia.org/wiki/Ext4) for all Postgres directories."
