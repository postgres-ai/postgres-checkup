package f008

import checkup ".."

type F008Setting struct {
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

type F008ReportHostResult struct {
	Data      map[string]F008Setting  `json:"data"`
	NodesJson checkup.ReportLastNodes `json:"nodes.json"`
}

type F008ReportHostsResults map[string]F008ReportHostResult

type F008Report struct {
	Project       string                  `json:"project"`
	Name          string                  `json:"name"`
	CheckId       string                  `json:"checkId"`
	Timestamptz   string                  `json:"timestamptz"`
	Database      string                  `json:"database"`
	Dependencies  map[string]interface{}  `json:"dependencies"`
	LastNodesJson checkup.ReportLastNodes `json:"last_nodes_json"`
	Results       F008ReportHostsResults  `json:"results"`
}
