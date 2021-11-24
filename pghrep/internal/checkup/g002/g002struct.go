package g002

import (
	"gitlab.com/postgres-ai/postgres-checkup/pghrep/internal/checkup"
)

type G002Connection struct {
	Num          int    `json:"num"`
	User         string `json:"user"`
	Database     string `json:"database"`
	CurrentState string `json:"current_state"`
	Count        int    `json:"count"`
	StateMore1m  int    `json:"state_changed_more_1m_ago"`
	StateMore1h  int    `json:"state_changed_more_1h_ago"`
	TxMore1m     int    `json:"tx_age_more_1m"`
	TxMore1h     int    `json:"tx_age_more_1h"`
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
