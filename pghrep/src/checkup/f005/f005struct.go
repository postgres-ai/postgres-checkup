package f005

import (
	checkup ".."
)

type F005IndexBloat struct {
	Num               int     `json:"num"`
	IsNa              string  `json:"is_na"`
	IndexName         string  `json:"index_name"`
	SchemaName        string  `json:"schema_name"`
	TableName         string  `json:"table_name"`
	IndexTableName    string  `json:"index_table_name"`
	RealSizeBytes     int64   `json:"real_size_bytes"`
	Size              string  `json:"size"`
	ExtraRatioPercent float32 `json:"extra_ratio_percent"`
	ExtraSizeBytes    int64   `json:"extra_size_bytes"`
	BloatSizeBytes    int64   `json:"bloat_size_bytes"`
	BloatRatioPercent float32 `json:"bloat_ratio_percent"`
	BloatRatioFactor  float32 `json:"bloat_ratio_factor"`
	LiveDataSizeBytes int64   `json:"live_data_size_bytes"`
	LastVaccuum       string  `json:"last_vaccuum"`
	Fillfactor        float32 `json:"fillfactor"`
	OverriddenSettings bool    `json:"overridden_settings"`
	TableSizeBytes    int64   `json:"table_size_bytes"`
}

// Current database tables list
type F005IndexBloatTotal struct {
	Count                int     `json:"count"`
	ExtraSizeBytesSum    int64   `json:"extra_size_bytes_sum"`
	RealSizeBytesSum     int64   `json:"real_size_bytes_sum"`
	BloatSizeBytesSum    int64   `json:"bloat_size_bytes_sum"`
	BloatRatioFactorAvg  float32 `json:"bloat_ratio_factor_avg"`
	BloatRatioPercentAvg float32 `json:"bloat_ratio_percent_avg"`
	TableSizeBytesSum    float32 `json:"table_size_bytes_sum"`
	LiveDataSizeBytesSum int64   `json:"live_data_size_bytes_sum"`
}

type F005ReportHostResultData struct {
	IndexBloat             map[string]F005IndexBloat `json:"Index_bloat"`
	IndexBloatTotal        F005IndexBloatTotal       `json:"Index_bloat_total"`
	OverriddenSettingsCount int                       `json:"overridden_settings_count"`
	DatabaseSizeBytes      int64                     `json:"database_size_bytes"`
}

type F005ReportHostResult struct {
	Data      F005ReportHostResultData `json:"data"`
	NodesJson checkup.ReportLastNodes  `json:"nodes.json"`
}

type F005ReportHostsResults map[string]F005ReportHostResult

type F005Report struct {
	Project       string                  `json:"project"`
	Name          string                  `json:"name"`
	CheckId       string                  `json:"checkId"`
	Timestamptz   string                  `json:"timestamptz"`
	Database      string                  `json:"database"`
	Dependencies  map[string]interface{}  `json:"dependencies"`
	LastNodesJson checkup.ReportLastNodes `json:"last_nodes_json"`
	Results       F005ReportHostsResults  `json:"results"`
}
