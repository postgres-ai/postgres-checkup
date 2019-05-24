package a002

import checkup ".."

type A002ReportHostResultData struct {
	Version          string `json:"version"`
	ServerVersionNum string `json:"server_version_num"`
	ServerMajorVer   string `json:"server_major_ver"`
	ServerMinorVer   string `json:"server_minor_ver"`
}

type A002ReportHostResult struct {
	Data      A002ReportHostResultData `json:"data"`
	NodesJson checkup.ReportLastNodes  `json:"nodes.json"`
}

type A002ReportHostsResults map[string]A002ReportHostResult

type A002Report struct {
	Project       string                  `json:"project"`
	Name          string                  `json:"name"`
	CheckId       string                  `json:"checkId"`
	Timestamptz   string                  `json:"timestamptz"`
	Database      string                  `json:"database"`
	Dependencies  map[string]interface{}  `json:"dependencies"`
	LastNodesJson checkup.ReportLastNodes `json:"last_nodes_json"`
	Results       A002ReportHostsResults  `json:"results"`
}
