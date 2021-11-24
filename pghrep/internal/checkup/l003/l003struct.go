package l003

import (
	"gitlab.com/postgres-ai/postgres-checkup/pghrep/internal/checkup"
)

type L003Table struct {
	Table               string  `json:"table"`
	Pk                  string  `json:"pk"`
	Type                string  `json:"type"`
	CurrentMaxValue     int64   `json:"current_max_value"`
	CapacityUsedPercent float64 `json:"capacity_used_percent"`
}

type L003ReportHostResultData struct {
	Tables            map[string]L003Table `json:"tables"`
	MinTableSizeBytes int64                `json:"min_table_size_bytes"`
}

type L003ReportHostResult struct {
	Data      L003ReportHostResultData `json:"data"`
	NodesJson checkup.ReportLastNodes  `json:"nodes.json"`
}

type L003ReportHostsResults map[string]L003ReportHostResult

type L003Report struct {
	Project       string                  `json:"project"`
	Name          string                  `json:"name"`
	CheckId       string                  `json:"checkId"`
	Timestamptz   string                  `json:"timestamptz"`
	Database      string                  `json:"database"`
	Dependencies  map[string]interface{}  `json:"dependencies"`
	LastNodesJson checkup.ReportLastNodes `json:"last_nodes_json"`
	Results       L003ReportHostsResults  `json:"results"`
}
