package f002

import (
	"encoding/json"
	"fmt"

	checkup ".."
	"../../pyraconv"
)

const CRITICAL_CAPACITY_USAGE float32 = 50.0

// Instace databases list
type F002PerIntace struct {
	Num          int     `json:"num"`
	DatabaseName string  `json:"datname"`
	Age          int     `json:"age"`
	CapacityUsed float32 `json:"capacity_used"`
	Datfrozenxid string  `json:"datfrozenxid"`
	Warning      int     `json:"warning"`
}

// Current database tables list
type F002PerDatabase struct {
	Num               int     `json:"num"`
	Relation          string  `json:"relation"`
	Age               int     `json:"age"`
	CapacityUsed      float32 `json:"capacity_used"`
	RelRelfrozenxid   string  `json:"rel_relfrozenxid"`
	ToastRelfrozenxid string  `json:"toast_relfrozenxid"`
	Warning           int     `json:"warning"`
	OverridedSettings bool    `json:"overrided_settings"`
}

type F002ReportHostResultData struct {
	PerInstance map[string]F002PerIntace   `json:"per_instance"`
	PerDatabase map[string]F002PerDatabase `json:"per_database"`
}

type F002ReportHostResult struct {
	Data      F002ReportHostResultData `json:"data"`
	NodesJson checkup.ReportLastNodes  `json:"nodes.json"`
}

type F002ReportHostsResults map[string]F002ReportHostResult

type F002Report struct {
	Project       string                  `json:"project"`
	Name          string                  `json:"name"`
	CheckId       string                  `json:"checkId"`
	Timestamptz   string                  `json:"timestamptz"`
	Database      string                  `json:"database"`
	Dependencies  map[string]interface{}  `json:"dependencies"`
	LastNodesJson checkup.ReportLastNodes `json:"last_nodes_json"`
	Results       F002ReportHostsResults  `json:"results"`
}

// Generate conclusions and recommendatons
func F002(data map[string]interface{}) {
	var conclusions []string
	var recommendations []string
	p1 := false
	p2 := false
	p3 := false
	filePath := pyraconv.ToString(data["source_path_full"])
	jsonRaw := checkup.LoadRawJsonReport(filePath)
	var report F002Report
	err := json.Unmarshal(jsonRaw, &report)
	if err != nil {
		return
	}
	for host, hostData := range report.Results {
		for table, tableData := range hostData.Data.PerDatabase {
			if tableData.CapacityUsed > CRITICAL_CAPACITY_USAGE {
				conclusions = append(conclusions, fmt.Sprintf("Transaction wraparound risks are large on host `%s` for table `%s`.", host, table))
				recommendations = append(recommendations, fmt.Sprintf("Check transaction wraparound risks on host `%s` for table `%s`.", host, table))
				p1 = true
			}
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
func F002PreprocessReportData(data map[string]interface{}) {
	F002(data)
}
