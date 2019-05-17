package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"../src/checkup"
	"../src/pyraconv"
)

type SupportedVersion struct {
	FirstRelease     string
	FinalRelease     string
	LastMinorVersion int
}

var SUPPORTED_VERSIONS map[string]SupportedVersion = map[string]SupportedVersion{
	"11": SupportedVersion{
		FirstRelease:     "2018-10-18",
		FinalRelease:     "2023-11-09",
		LastMinorVersion: 3,
	},
	"10": SupportedVersion{
		FirstRelease:     "2017-10-05",
		FinalRelease:     "2022-11-10",
		LastMinorVersion: 8,
	},
	"9.6": SupportedVersion{
		FirstRelease:     "2016-09-29",
		FinalRelease:     "2021-11-11",
		LastMinorVersion: 13,
	},
	"9.5": SupportedVersion{
		FirstRelease:     "2016-01-07",
		FinalRelease:     "2021-02-11",
		LastMinorVersion: 17,
	},
	"9.4": SupportedVersion{
		FirstRelease:     "2014-12-18",
		FinalRelease:     "2020-02-13",
		LastMinorVersion: 22,
	},
}

type prepare string

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

// Generate conclusions and recommendatons
func A002(data map[string]interface{}) {
	var conclusions []string
	var recommendations []string
	p1 := false
	p2 := false
	p3 := false
	filePath := pyraconv.ToString(data["source_path_full"])
	jsonRaw := checkup.LoadRawJsonReport(filePath)
	var report A002Report
	err := json.Unmarshal(jsonRaw, &report)
	if err != nil {
		return
	}
	for host, hostData := range report.Results {
		ver, ok := SUPPORTED_VERSIONS[hostData.Data.ServerMajorVer]
		if !ok {
			conclusions = append(conclusions, fmt.Sprintf("Unknown PostgreSQL version on %s.", host))
			recommendations = append(recommendations, fmt.Sprintf("Check PostgreSQL version on %s.", host))
		}
		from, _ := time.Parse("2006-01-02", ver.FirstRelease)
		to, _ := time.Parse("2006-01-02", ver.FinalRelease)
		minorVersion, _ := strconv.ParseInt(hostData.Data.ServerMinorVer, 10, 0)
		if ver.LastMinorVersion > int(minorVersion) {
			conclusions = append(conclusions, fmt.Sprintf("Minor PostgreSQL version on host `%s` is not last.", host))
			recommendations = append(recommendations, fmt.Sprintf("Update minor PostgreSQL version on host `%s`.", host))
			p2 = true
		}
		today := time.Now()
		if today.After(to) || today.Before(from) {
			conclusions = append(conclusions, fmt.Sprintf("Major PostgreSQL version on host `%s` is supported by community now.", host))
			recommendations = append(recommendations, fmt.Sprintf("Update major PostgreSQL version on host `%s`.", host))
			p1 = true
		}
	}
	// update data
	data["conclusions"] = conclusions
	data["recommendations"] = recommendations
	data["p1"] = p1
	data["p2"] = p2
	data["p3"] = p3
	// update file
	checkup.SaveJsonConclusionsRecommendations(data, conclusions, recommendations, p1, p2, p3)
}

// Plugin entry point
func (g prepare) Prepare(data map[string]interface{}) map[string]interface{} {
	A002(data)
	return data
}

var Preparer prepare
