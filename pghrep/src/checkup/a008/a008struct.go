package a008

import checkup ".."

type FsItem struct {
	Fstype     string `json:"fstype"`
	Size       string `json:"size"`
	Avail      string `json:"avail"`
	Used       string `json:"used"`
	UsePercent string `json:"use_percent"`
	MountPoint string `json:"mount_point"`
	Path       string `json:"path"`
	Device     string `json:"device"`
}

type A008ReportHostResultData struct {
	DbData map[string]FsItem `json:"data"`
	FsData map[string]FsItem `json:"fs_data"`
}

type A008ReportHostResult struct {
	Data      A008ReportHostResultData `json:"data"`
	NodesJson checkup.ReportLastNodes  `json:"nodes.json"`
}

type A008ReportHostsResults map[string]A008ReportHostResult

type A008Report struct {
	Project       string                  `json:"project"`
	Name          string                  `json:"name"`
	CheckId       string                  `json:"checkId"`
	Timestamptz   string                  `json:"timestamptz"`
	Database      string                  `json:"database"`
	Dependencies  map[string]interface{}  `json:"dependencies"`
	LastNodesJson checkup.ReportLastNodes `json:"last_nodes_json"`
	Results       A008ReportHostsResults  `json:"results"`
}
