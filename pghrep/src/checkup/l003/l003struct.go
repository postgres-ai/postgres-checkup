package l003

import checkup ".."

type L003Table struct {
	Table           string `json:"Table"`
	Pk              string `json:"PK"`
	Type            string `json:"Type"`
	CurrentMaxValue int64  `json:"Current max value"`
	//	CapacityUsedPercent float32 `json:"Capacity used, %%"`
}

type L003ReportHostResult struct {
	Data      map[string]L003Table    `json:"data"`
	NodesJson checkup.ReportLastNodes `json:"nodes.json"`
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
