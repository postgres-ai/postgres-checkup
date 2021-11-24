package f001

import (
	"gitlab.com/postgres-ai/postgres-checkup/pghrep/internal/checkup"
)

type F001GlobalSetting struct {
	Name           string `json:"name"`
	Setting        string `json:"setting"`
	Unit           string `json:"unit"`
	Сategory       string `json:"category"`
	ShortDesc      string `json:"short_desc"`
	ExtraDesc      string `json:"extra_desc"`
	Context        string `json:"context"`
	Vartype        string `json:"vartype"`
	Source         string `json:"source"`
	MinVal         string `json:"min_val"`
	MaxVal         string `json:"max_val"`
	Numvals        string `json:"enumvals"`
	BootVal        string `json:"boot_val"`
	ResetVal       string `json:"reset_val"`
	Sourcefile     string `json:"sourcefile"`
	Sourceline     int    `json:"sourceline"`
	PendingRestart bool   `json:"pending_restart"`
}

type F001TableSetting struct {
	Namespace  string   `json:"namespace"`
	Relname    string   `json:"relname"`
	Reloptions []string `json:"reloptions"`
}

type F001Settings struct {
	GlobalSettings map[string]F001GlobalSetting `json:"global_settings"`
	TableSettings  map[string]F001TableSetting  `json:"table_settings"`
}

type F001ReportHostResultData struct {
	Settings F001Settings `json:"settings"`
}

type F001ReportHostResult struct {
	Data      F001ReportHostResultData `json:"data"`
	NodesJson checkup.ReportLastNodes  `json:"nodes.json"`
}

type F001ReportHostsResults map[string]F001ReportHostResult

type F001Report struct {
	Project       string                  `json:"project"`
	Name          string                  `json:"name"`
	CheckId       string                  `json:"checkId"`
	Timestamptz   string                  `json:"timestamptz"`
	Database      string                  `json:"database"`
	Dependencies  map[string]interface{}  `json:"dependencies"`
	LastNodesJson checkup.ReportLastNodes `json:"last_nodes_json"`
	Results       F001ReportHostsResults  `json:"results"`
}
