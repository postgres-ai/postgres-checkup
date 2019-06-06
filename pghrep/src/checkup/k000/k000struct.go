package k000

import checkup ".."

type K000Table struct {
	Table           string `json:"Table"`
	Pk              string `json:"PK"`
	Type            string `json:"Type"`
	CurrentMaxValue int64  `json:"Current max value"`
}

/*type K00Query struct {
	RowNum                    int64   `json:"rownum"`
	DiffCalls                 int64   `json:"diff_calls"`
	PerSecCalls               float64 `json:"per_sec_calls"`
	PerCallCalls              int64   `json:"per_call_calls"`
	RatioCalls                float64 `json:"ratio_calls"`
	DiffTotalTime             float64 `json:"diff_total_time"`
	PerSecTotalTime           float64 `json:"per_sec_total_time"`
	PerCallTotalTime          float64 `json:"per_call_total_time"`
	RatioTotalTime            float64 `json:"ratio_total_time"`
	DiffRows                  int64   `json:"diff_rows"`
	PerSecRows                float64 `json:"per_sec_rows"`
	PerCallRows               int64   `json:"per_call_rows"`
	RatioRows                 float64 `json:"ratio_rows"`
	DiffSharedBlksHit         int64   `json:"diff_shared_blks_hit"`
	PerSecSharedBlksHit       float64 `json:"per_sec_shared_blks_hit"`
	PerCallSharedBlksHit      float64 `json:"per_call_shared_blks_hit"`
	RatioSharedBlksHit        float64 `json:"ratio_shared_blks_hit"`
	DiffSharedBlksRead        int64   `json:"diff_shared_blks_read"`
	PerSecSharedBlksRead      float64 `json:"per_sec_shared_blks_read"`
	PerCallSharedBlksRead     float64 `json:"per_call_shared_blks_read",`
	RatioSharedBlksRead       float64 `json:"ratio_shared_blks_read"`
	DiffSharedBlksDirtied     int64   `json:"diff_shared_blks_dirtied"`
	PerSecSharedBlksDirtied   float64 `json:"per_sec_shared_blks_dirtied"`
	PerCallSharedBlksDirtied  float64 `json:"per_call_shared_blks_dirtied"`
	RatioSharedBlksDirtied    float64 `json:"ratio_shared_blks_dirtied"`
	DiffSharedBlksWritten     int64   `json:"diff_shared_blks_written"`
	PerSecSharedBlksWritten   int64   `json:"per_sec_shared_blks_written"`
	PerCallSharedBlksWritten  int64   `json:"per_call_shared_blks_written"`
	RatioSharedBlksWritten    float64 `json:"ratio_shared_blks_written"`
	DiffLocalBlksHit          int64   `json:"diff_local_blks_hit"`
	PerSecLocalBlksHit        int64   `json:"per_sec_local_blks_hit"`
	PerCallLocalBlksHit       int64   `json:"per_call_local_blks_hit"`
	RatioLocalBlksHit         float64 `json:"ratio_local_blks_hit"`
	DiffLocalBlksRead         int64   `json:"diff_local_blks_read"`
	PerSecLocalBlksRead       int64   `json:"per_sec_local_blks_read"`
	PerCallLocalBlksRead      int64   `json:"per_call_local_blks_read"`
	RatioLocalBlksRead        float64 `json:"ratio_local_blks_read"`
	DiffLocalBlksDirtied      int64   `json:"diff_local_blks_dirtied"`
	PerSecLocalBlksDirtied    int64   `json:"per_sec_local_blks_dirtied"`
	PerCallLocalBlksDirtied   int64   `json:"per_call_local_blks_dirtied"`
	RatioLocalBlksDirtied     float64 `json:"ratio_local_blks_dirtied"`
	DiffLocalBlksWritten      int64   `json:"diff_local_blks_written"`
	PerSecLocalBlksWritten    int64   `json:"per_sec_local_blks_written"`
	PerCallLocalBlksWritten   int64   `json:"per_call_local_blks_written"`
	RatioLocalBlksWritten     float64 `json:"ratio_local_blks_written"`
	DiffTempBlksRead          int64   `json:"diff_temp_blks_read"`
	PerSecTempBlksRead        int64   `json:"per_sec_temp_blks_read"`
	PerCallTempBlksRead       int64   `json:"per_call_temp_blks_read"`
	RatioTempBlksRead         float64 `json:"ratio_temp_blks_read"`
	DiffTempBlksWritten       int64   `json:"diff_temp_blks_written"`
	PerSecTempBlksWritten     int64   `json:"per_sec_temp_blks_written"`
	PerCallTempBlksWritten    int64   `json:"per_call_temp_blks_written"`
	RatioTempBlksWritten      float64 `json:"ratio_temp_blks_written"`
	DiffBlkReadTime           float64 `json:"diff_blk_read_time"`
	PerSecBlkReadTime         float64 `json:"per_sec_blk_read_time"`
	PerCallBlkReadTime        float64 `json:"per_call_blk_read_time"`
	RatioBlkReadTime          float64 `json:"ratio_blk_read_time"`
	DiffBlkWriteTime          int64   `json:"diff_blk_write_time"`
	PerSecBlkWriteTime        int64   `json:"per_sec_blk_write_time"`
	PerCallBlkWriteTime       int64   `json:"per_call_blk_write_time"`
	RatioBlkWriteTime         float64 `json:"ratio_blk_write_time"`
	DiffKcacheReads           int64   `json:"diff_kcache_reads"`
	PerSecKcacheReads         int64   `json:"per_sec_kcache_reads"`
	PerCallKcacheReads        int64   `json:"per_call_kcache_reads"`
	RatioKcacheReads          float64 `json:"ratio_kcache_reads"`
	DiffKcacheWrites          int64   `json:"diff_kcache_writes"`
	PerSecKcacheWrites        int64   `json:"per_sec_kcache_writes"`
	PerCallKcacheWrites       int64   `json:"per_call_kcache_writes"`
	RatioKcacheWrites         float64 `json:"ratio_kcache_writes"`
	DiffKcacheUserTimeMs      int64   `json:"diff_kcache_user_time_ms"`
	PerSecKcacheUserTimeMs    int64   `json:"per_sec_kcache_user_time_ms"`
	PerCallKcacheUserTimeMs   int64   `json:"per_call_kcache_user_time_ms"`
	RatioKcacheUserTimeMs     float64 `json:"ratio_kcache_user_time_ms"`
	DiffKcacheSystemTimeMs    int64   `json:"diff_kcache_system_time_ms"`
	PerSecKcacheSystemTimeMs  int64   `json:"per_sec_kcache_system_time_ms"`
	PerCallKcacheSystemTimeMs int64   `json:"per_call_kcache_system_time_ms"`
	RatioKcacheSystemTimeMs   float64 `json:"ratio_kcache_system_time_ms"`
}*/

type K00Query struct {
	RowNum                    int64   `json:"rownum"`
	DiffCalls                 float64 `json:"diff_calls"`
	PerSecCalls               float64 `json:"per_sec_calls"`
	PerCallCalls              float64 `json:"per_call_calls"`
	RatioCalls                float64 `json:"ratio_calls"`
	DiffTotalTime             float64 `json:"diff_total_time"`
	PerSecTotalTime           float64 `json:"per_sec_total_time"`
	PerCallTotalTime          float64 `json:"per_call_total_time"`
	RatioTotalTime            float64 `json:"ratio_total_time"`
	DiffRows                  float64 `json:"diff_rows"`
	PerSecRows                float64 `json:"per_sec_rows"`
	PerCallRows               float64 `json:"per_call_rows"`
	RatioRows                 float64 `json:"ratio_rows"`
	DiffSharedBlksHit         float64 `json:"diff_shared_blks_hit"`
	PerSecSharedBlksHit       float64 `json:"per_sec_shared_blks_hit"`
	PerCallSharedBlksHit      float64 `json:"per_call_shared_blks_hit"`
	RatioSharedBlksHit        float64 `json:"ratio_shared_blks_hit"`
	DiffSharedBlksRead        float64 `json:"diff_shared_blks_read"`
	PerSecSharedBlksRead      float64 `json:"per_sec_shared_blks_read"`
	PerCallSharedBlksRead     float64 `json:"per_call_shared_blks_read"`
	RatioSharedBlksRead       float64 `json:"ratio_shared_blks_read"`
	DiffSharedBlksDirtied     float64 `json:"diff_shared_blks_dirtied"`
	PerSecSharedBlksDirtied   float64 `json:"per_sec_shared_blks_dirtied"`
	PerCallSharedBlksDirtied  float64 `json:"per_call_shared_blks_dirtied"`
	RatioSharedBlksDirtied    float64 `json:"ratio_shared_blks_dirtied"`
	DiffSharedBlksWritten     float64 `json:"diff_shared_blks_written"`
	PerSecSharedBlksWritten   float64 `json:"per_sec_shared_blks_written"`
	PerCallSharedBlksWritten  float64 `json:"per_call_shared_blks_written"`
	RatioSharedBlksWritten    float64 `json:"ratio_shared_blks_written"`
	DiffLocalBlksHit          float64 `json:"diff_local_blks_hit"`
	PerSecLocalBlksHit        float64 `json:"per_sec_local_blks_hit"`
	PerCallLocalBlksHit       float64 `json:"per_call_local_blks_hit"`
	RatioLocalBlksHit         float64 `json:"ratio_local_blks_hit"`
	DiffLocalBlksRead         float64 `json:"diff_local_blks_read"`
	PerSecLocalBlksRead       float64 `json:"per_sec_local_blks_read"`
	PerCallLocalBlksRead      float64 `json:"per_call_local_blks_read"`
	RatioLocalBlksRead        float64 `json:"ratio_local_blks_read"`
	DiffLocalBlksDirtied      float64 `json:"diff_local_blks_dirtied"`
	PerSecLocalBlksDirtied    float64 `json:"per_sec_local_blks_dirtied"`
	PerCallLocalBlksDirtied   float64 `json:"per_call_local_blks_dirtied"`
	RatioLocalBlksDirtied     float64 `json:"ratio_local_blks_dirtied"`
	DiffLocalBlksWritten      float64 `json:"diff_local_blks_written"`
	PerSecLocalBlksWritten    float64 `json:"per_sec_local_blks_written"`
	PerCallLocalBlksWritten   float64 `json:"per_call_local_blks_written"`
	RatioLocalBlksWritten     float64 `json:"ratio_local_blks_written"`
	DiffTempBlksRead          float64 `json:"diff_temp_blks_read"`
	PerSecTempBlksRead        float64 `json:"per_sec_temp_blks_read"`
	PerCallTempBlksRead       float64 `json:"per_call_temp_blks_read"`
	RatioTempBlksRead         float64 `json:"ratio_temp_blks_read"`
	DiffTempBlksWritten       float64 `json:"diff_temp_blks_written"`
	PerSecTempBlksWritten     float64 `json:"per_sec_temp_blks_written"`
	PerCallTempBlksWritten    float64 `json:"per_call_temp_blks_written"`
	RatioTempBlksWritten      float64 `json:"ratio_temp_blks_written"`
	DiffBlkReadTime           float64 `json:"diff_blk_read_time"`
	PerSecBlkReadTime         float64 `json:"per_sec_blk_read_time"`
	PerCallBlkReadTime        float64 `json:"per_call_blk_read_time"`
	RatioBlkReadTime          float64 `json:"ratio_blk_read_time"`
	DiffBlkWriteTime          float64 `json:"diff_blk_write_time"`
	PerSecBlkWriteTime        float64 `json:"per_sec_blk_write_time"`
	PerCallBlkWriteTime       float64 `json:"per_call_blk_write_time"`
	RatioBlkWriteTime         float64 `json:"ratio_blk_write_time"`
	DiffKcacheReads           float64 `json:"diff_kcache_reads"`
	PerSecKcacheReads         float64 `json:"per_sec_kcache_reads"`
	PerCallKcacheReads        float64 `json:"per_call_kcache_reads"`
	RatioKcacheReads          float64 `json:"ratio_kcache_reads"`
	DiffKcacheWrites          float64 `json:"diff_kcache_writes"`
	PerSecKcacheWrites        float64 `json:"per_sec_kcache_writes"`
	PerCallKcacheWrites       float64 `json:"per_call_kcache_writes"`
	RatioKcacheWrites         float64 `json:"ratio_kcache_writes"`
	DiffKcacheUserTimeMs      float64 `json:"diff_kcache_user_time_ms"`
	PerSecKcacheUserTimeMs    float64 `json:"per_sec_kcache_user_time_ms"`
	PerCallKcacheUserTimeMs   float64 `json:"per_call_kcache_user_time_ms"`
	RatioKcacheUserTimeMs     float64 `json:"ratio_kcache_user_time_ms"`
	DiffKcacheSystemTimeMs    float64 `json:"diff_kcache_system_time_ms"`
	PerSecKcacheSystemTimeMs  float64 `json:"per_sec_kcache_system_time_ms"`
	PerCallKcacheSystemTimeMs float64 `json:"per_call_kcache_system_time_ms"`
	RatioKcacheSystemTimeMs   float64 `json:"ratio_kcache_system_time_ms"`
	Md5                       string  `json:"md5"`
	Queryid                   string  `json:"queryid"`
	Query                     string  `json:"query"`
	Link                      string  `json:"link"`
	ReadableQueryid           string  `json:"readable_queryid"`
}

type K000HostData struct {
	StartTimestamptz       string              `json:"start_timestamptz"`
	EndTimestamptz         string              `json:"end_timestamptz"`
	PeriodSeconds          float64             `json:"period_seconds"`
	PeriodAge              string              `json:"period_age"`
	AbsoluteErrorCalls     float64             `json:"absolute_error_calls"`
	AbsoluteErrorTotalTime float64             `json:"absolute_error_total_time"`
	RelativeErrorCalls     float64             `json:"relative_error_calls"`
	RelativeErrorTotalTime float64             `json:"relative_error_total_time"`
	Queries                map[string]K00Query `json:"queries"`
	Aggregated             map[string]K00Query `json:"aggregated"`
	WorkloadType           map[string]K00Query `json:"workload_type"`
}

type K000ReportHostResult struct {
	Data      K000HostData            `json:"data"`
	NodesJson checkup.ReportLastNodes `json:"nodes.json"`
}

type K000ReportHostsResults map[string]K000ReportHostResult

type K000Report struct {
	Project       string                  `json:"project"`
	Name          string                  `json:"name"`
	CheckId       string                  `json:"checkId"`
	Timestamptz   string                  `json:"timestamptz"`
	Database      string                  `json:"database"`
	Dependencies  map[string]interface{}  `json:"dependencies"`
	LastNodesJson checkup.ReportLastNodes `json:"last_nodes_json"`
	Results       K000ReportHostsResults  `json:"results"`
}
