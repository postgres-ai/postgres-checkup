package f004

import (
	checkup ".."
)

type F004HeapBloat struct {
	Num               int     `json:"num"`
	IsNa              string  `json:"is_na"`
	TableName         string  `json:"table_name"`
	RealSize          string  `json:"real_size"`
	ExtraSizeBytes    int64   `json:"extra_size_bytes"`
	ExtraRatioPercent float32 `json:"extra_ratio_percent"`
	Extra             string  `json:"extra"`
	BloatSizeBytes    int64   `json:"bloat_size_bytes"`
	BloatRatioPercent float32 `json:"bloat_ratio_percent"`
	BloatEstimate     string  `json:"bloat_estimate"`
	RealSizeBytes     int64   `json:"real_size_bytes"`
	LiveDataSize      string  `json:"live_data_size"`
	LiveDataSizeBytes int64   `json:"live_data_size_bytes"`
	LastVaccuum       string  `json:"last_vaccuum"`
	Fillfactor        float32 `json:"fillfactor"`
	OverridedSettings bool    `json:"overrided_settings"`
	BloatRatio        float32 `json:"bloat_ratio"`
}

// Current database tables list
type F004HeapBloatTotal struct {
	Count                int     `json:"count"`
	ExtraSizeBytesSum    int64   `json:"extra_size_bytes_sum"`
	RealSizeBytesSum     int64   `json:"real_size_bytes_sum"`
	BloatSizeBytesSum    int64   `json:"bloat_size_bytes_sum"`
	LiveDataSizeBytesSum int64   `json:"live_data_size_bytes_sum"`
	BloatRatioPercentAvg float32 `json:"bloat_ratio_percent_avg"`
	BloatRatioAvg        float32 `json:"bloat_ratio_avg"`
}

type F004ReportHostResultData struct {
	HeapBloat              map[string]F004HeapBloat `json:"heap_bloat"`
	HeapBloatTotal         F004HeapBloatTotal       `json:"heap_bloat_total"`
	OverridedSettingsCount int                      `json:"overrided_settings_count"`
}

type F004ReportHostResult struct {
	Data      F004ReportHostResultData `json:"data"`
	NodesJson checkup.ReportLastNodes  `json:"nodes.json"`
}

type F004ReportHostsResults map[string]F004ReportHostResult

type F004Report struct {
	Project       string                  `json:"project"`
	Name          string                  `json:"name"`
	CheckId       string                  `json:"checkId"`
	Timestamptz   string                  `json:"timestamptz"`
	Database      string                  `json:"database"`
	Dependencies  map[string]interface{}  `json:"dependencies"`
	LastNodesJson checkup.ReportLastNodes `json:"last_nodes_json"`
	Results       F004ReportHostsResults  `json:"results"`
}
