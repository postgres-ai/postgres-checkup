package a002

import (
	"strconv"
	"strings"

	"gitlab.com/postgres-ai/postgres-checkup/pghrep/internal/checkup"
)

type A002ReportHostResultData struct {
	Version          string `json:"version"`
	ServerVersionNum string `json:"server_version_num"`
	ServerMajorVer   string `json:"server_major_ver"`
	ServerMinorVer   string `json:"server_minor_ver"`
}

func (v *A002ReportHostResultData) GetMajorVersion() (float64, error) {
	var majorVersion string

	if v.ServerVersionNum[0:1] == "9" {
		majorVersion = v.ServerVersionNum[0:3]
		majorVersion = strings.Replace(majorVersion, "0", ".", 1)
	} else {
		majorVersion = v.ServerVersionNum[0:2]
	}

	return strconv.ParseFloat(majorVersion, 64)
}

func (v *A002ReportHostResultData) GetMinorVersion() (int, error) {
	var minorVersion string

	minorVersion = v.ServerVersionNum[len(v.ServerVersionNum)-2 : len(v.ServerVersionNum)]

	return strconv.Atoi(minorVersion)
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
