package f002

import checkup ".."

// Instace databases list
type F002Database struct {
	Num          int     `json:"num"`
	DatabaseName string  `json:"datname"`
	Age          int     `json:"age"`
	CapacityUsed float32 `json:"capacity_used"`
	Datfrozenxid string  `json:"datfrozenxid"`
	Warning      int     `json:"warning"`
}

// Current database tables list
type F002Table struct {
	Num               int     `json:"num"`
	Relation          string  `json:"relation"`
	Age               int     `json:"age"`
	CapacityUsed      float32 `json:"capacity_used"`
	RelRelfrozenxid   string  `json:"rel_relfrozenxid"`
	ToastRelfrozenxid string  `json:"toast_relfrozenxid"`
	Warning           int     `json:"warning"`
	OverriddenSettings bool    `json:"overridden_settings"`
}

type F002ReportHostResultData struct {
	Databases map[string]F002Database `json:"per_instance"`
	Tables    map[string]F002Table    `json:"per_database"`
}

type F002ReportHostResult struct {
	Data      F002ReportHostResultData `json:"data"`
	NodesJson checkup.ReportLastNodes  `json:"nodes.json"`
}

type F002ReportHostsResults map[string]F002ReportHostResult

type F002Report struct {
	Project       string                  `json:"project"`
	Name          string                  `json:"name"`
	CheckId       string                  `json:"checkId"`
	Timestamptz   string                  `json:"timestamptz"`
	Database      string                  `json:"database"`
	Dependencies  map[string]interface{}  `json:"dependencies"`
	LastNodesJson checkup.ReportLastNodes `json:"last_nodes_json"`
	Results       F002ReportHostsResults  `json:"results"`
}
