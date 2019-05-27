package f002

import (
	"encoding/json"
	"strings"

	"../../log"
	"../../pyraconv"

	checkup ".."
)

const CRITICAL_CAPACITY_USAGE float32 = 50.0

// Generate conclusions and recommendatons
func F002Process(report F002Report) checkup.ReportOutcome {
	var result checkup.ReportOutcome
	var databases, tables []string
	for _, hostData := range report.Results {
		for db, dbData := range hostData.Data.Databases {
			if dbData.CapacityUsed > CRITICAL_CAPACITY_USAGE {
				databases = append(databases, "- database `"+db+"`  \n")
				result.P1 = true
			}
		}
		for table, tableData := range hostData.Data.Tables {
			if tableData.CapacityUsed > CRITICAL_CAPACITY_USAGE {
				tables = append(tables, "- table `"+table+"`  \n")
				result.P1 = true
			}
		}
	}
	items := strings.Join(databases, "") + strings.Join(tables, "")
	if len(databases) > 0 || len(tables) > 0 {
		result.AppendConclusion(MSG_RISKS_ARE_HIGH_CONCLUSION, items)
		result.AppendRecommendation(MSG_RISKS_ARE_HIGH_RECOMMENDATION)
	}
	return result
}

func F002PreprocessReportData(data map[string]interface{}) {
	filePath := pyraconv.ToString(data["source_path_full"])
	jsonRaw := checkup.LoadRawJsonReport(filePath)
	var report F002Report
	err := json.Unmarshal(jsonRaw, &report)
	if err != nil {
		log.Err("Cannot load json report to process")
		return
	}
	result := F002Process(report)
	if len(result.Recommendations) == 0 {
		result.AppendRecommendation(MSG_NO_RECOMMENDATIONS)
	}
	checkup.SaveConclusionsRecommendations(data, result)
}
