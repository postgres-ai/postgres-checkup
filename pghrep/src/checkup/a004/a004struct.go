package a004

import checkup ".."

type A004Metric struct {
	Metric string `json:"metric"`
	Value  string `json:"value"`
}

type A004ReportHostResultData struct {
	GeneralInfo   map[string]A004Metric `json:"general_info"`
	DatabaseSizes map[string]int64      `json:"database_sizes"`
}

type A004ReportHostResult struct {
	Data      A004ReportHostResultData `json:"data"`
	NodesJson checkup.ReportLastNodes  `json:"nodes.json"`
}

type A004ReportHostsResults map[string]A004ReportHostResult

type A004Report struct {
	Project       string                  `json:"project"`
	Name          string                  `json:"name"`
	CheckId       string                  `json:"checkId"`
	Timestamptz   string                  `json:"timestamptz"`
	Database      string                  `json:"database"`
	Dependencies  map[string]interface{}  `json:"dependencies"`
	LastNodesJson checkup.ReportLastNodes `json:"last_nodes_json"`
	Results       A004ReportHostsResults  `json:"results"`
}
