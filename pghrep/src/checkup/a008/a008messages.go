package a008

const MSG_NO_USAGE_RISKS_CONCLUSION string = "No significant risks of out-of-disk-space problem have been detected."
const MSG_USAGE_WARNING_CONCLUSION string = "[P2] Disk `%s` on `%s` space usage is %s, it exceeds 70%%. There are some risks of out-of-disk-space problem."
const MSG_USAGE_WARNING_RECOMMENDATION string = "[P2] Add more disk space to `%s` on `%s`. It is recommended to keep free disk space more than %d%% " +
	"to reduce risks of out-of-disk-space problem."
const MSG_USAGE_CRITICAL_CONCLUSION string = "[P1] Disk `%s` on `%s` space usage is %s, it exceeds 90%%. There are significant risks of out-of-disk-space problem. " +
	"In this case, PostgreSQL will stop working and manual fix will be required."
const MSG_USAGE_CRITICAL_RECOMMENDATION string = "[P1] Add more disk space to `%s` on `%s` as soon as possible to prevent outage."
const MSG_NETWORK_FS_CONCLUSION_1 string = "[P1] `%s` on host `%s` is located on an NFS drive. This might lead to serious issues with Postgres, including downtime and data corruption."
const MSG_NETWORK_FS_CONCLUSION_N string = "[P1] `%s` on host `%s` are located on an NFS drive. This might lead to serious issues with Postgres, including downtime and data corruption."
const MSG_NETWORK_FS_RECOMMENDATION string = "[P1] Do not use NFS for Postgres."
const MSG_NOT_RECOMMENDED_FS_CONCLUSION_1 string = "[P3] `%s` on host `%s` is located on drive where the following filesystems are used: `%s`. This might mean that Postgres performance and reliability characteristics are worse than it could be in case of use of more popular filesystems (such as ext4)."
const MSG_NOT_RECOMMENDED_FS_CONCLUSION_N string = "[P3] `%s` on host `%s` are located on drives where the following filesystems are used: `%s` respectively. This might mean that Postgres performance and reliability characteristics are worse than it could be in case of use of more popular filesystems (such as ext4)."
const MSG_NOT_RECOMMENDED_FS_RECOMMENDATION string = "[P3] Consider using ext4 for all Postgres directories."
const MSG_NO_FS_RECOMMENDATION string = "No recommendations."
