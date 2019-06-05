package g002

import checkup ".."

type G002Connection struct {
	Num          int    `json:"num"`
	User         string `json:"User"`
	Db           string `json:"DB"`
	CurrentState string `json:"Current State"`
	Count        int    `json:"Count"`
	StateMore1m  int    `json:"State changed >1m ago"`
	StateMore1h  int    `json:"State changed >1h ago"`
	TxMore1m     int    `json:"Tx age >1m"`
	TmMore1h     int    `json:"Tx age >1h"`
}

type G002ReportHostResult struct {
	Data      map[string]G002Connection `json:"data"`
	NodesJson checkup.ReportLastNodes   `json:"nodes.json"`
}

type G002ReportHostsResults map[string]G002ReportHostResult

type G002Report struct {
	Project       string                  `json:"project"`
	Name          string                  `json:"name"`
	CheckId       string                  `json:"checkId"`
	Timestamptz   string                  `json:"timestamptz"`
	Database      string                  `json:"database"`
	Dependencies  map[string]interface{}  `json:"dependencies"`
	LastNodesJson checkup.ReportLastNodes `json:"last_nodes_json"`
	Results       G002ReportHostsResults  `json:"results"`
}
