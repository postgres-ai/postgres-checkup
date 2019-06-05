package l001

import checkup ".."

type L001Table struct {
	Num       int    `json:"num"`
	Table     string `json:"Table"`
	Rows      string `json:"Rows"`
	TotalSize string `json:"Total Size"`
	TableSize string `json:"Table Size"`
	ToastSize string `json:"TOAST Size"`
}

type L001ReportHostResult struct {
	Data      map[string]L001Table    `json:"data"`
	NodesJson checkup.ReportLastNodes `json:"nodes.json"`
}

type L001ReportHostsResults map[string]L001ReportHostResult

type L001Report struct {
	Project       string                  `json:"project"`
	Name          string                  `json:"name"`
	CheckId       string                  `json:"checkId"`
	Timestamptz   string                  `json:"timestamptz"`
	Database      string                  `json:"database"`
	Dependencies  map[string]interface{}  `json:"dependencies"`
	LastNodesJson checkup.ReportLastNodes `json:"last_nodes_json"`
	Results       L001ReportHostsResults  `json:"results"`
}
