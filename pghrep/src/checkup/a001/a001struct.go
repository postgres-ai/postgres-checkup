package a001

import checkup ".."

type A001ReportCpu struct {
}

type A001ReportRam struct {
	MemTotal  string `json:"MemTotal"`
	SwapTotal string `json:"SwapTotal"`
}

type A001ReportSystem struct {
}

type A001ReportVirtualization struct {
}

type A001ReportHostResultData struct {
	Cpu            A001ReportCpu            `json:"cpu"`
	Ram            A001ReportRam            `json:"ram"`
	System         A001ReportSystem         `json:"system"`
	Virtualization A001ReportVirtualization `json:"virtualization"`
}

type A001ReportHostResult struct {
	Data      A001ReportHostResultData `json:"data"`
	NodesJson checkup.ReportLastNodes  `json:"nodes.json"`
}

type A001ReportHostsResults map[string]A001ReportHostResult

type A001Report struct {
	Project       string                  `json:"project"`
	Name          string                  `json:"name"`
	CheckId       string                  `json:"checkId"`
	Timestamptz   string                  `json:"timestamptz"`
	Database      string                  `json:"database"`
	Dependencies  map[string]interface{}  `json:"dependencies"`
	LastNodesJson checkup.ReportLastNodes `json:"last_nodes_json"`
	Results       A001ReportHostsResults  `json:"results"`
}
