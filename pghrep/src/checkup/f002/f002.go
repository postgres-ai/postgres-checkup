package f002

import (
	"encoding/json"
	"strings"

	checkup ".."
)

const F002_RISKS_ARE_HIGH = "F002_RISKS_ARE_HIGH"

const CRITICAL_CAPACITY_USAGE float32 = 50.0

// Generate conclusions and recommendatons
func F002Process(report F002Report) checkup.ReportResult {
	var result checkup.ReportResult
	var databases, tables []string
	for _, hostData := range report.Results {
		for db, dbData := range hostData.Data.Databases {
			if dbData.CapacityUsed > CRITICAL_CAPACITY_USAGE {
				databases = append(databases, "    - database `"+db+"`  \n")
				result.P1 = true
			}
		}
		for table, tableData := range hostData.Data.Tables {
			if tableData.CapacityUsed > CRITICAL_CAPACITY_USAGE {
				tables = append(tables, "    - table `"+table+"`  \n")
				result.P1 = true
			}
		}
	}
	items := strings.Join(databases, "") + strings.Join(tables, "")
	if len(databases) > 0 || len(tables) > 0 {
		result.AppendConclusion(F002_RISKS_ARE_HIGH, MSG_RISKS_ARE_HIGH_CONCLUSION_1, items)
		result.AppendConclusion(F002_RISKS_ARE_HIGH, MSG_RISKS_ARE_HIGH_CONCLUSION_2)
		result.AppendRecommendation(F002_RISKS_ARE_HIGH, MSG_RISKS_ARE_HIGH_RECOMMENDATION)
	}
	return result
}

func F002PreprocessReportData(data map[string]interface{}) {
	var report F002Report
	filePath := data["source_path_full"].(string)
	jsonRaw := checkup.LoadRawJsonReport(filePath)
	if !checkup.CheckUnmarshalResult(json.Unmarshal(jsonRaw, &report)) {
		return
	}
	result := F002Process(report)
	checkup.SaveReportResult(data, result)
}
