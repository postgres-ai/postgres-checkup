package k000

import (
	"encoding/json"

	checkup ".."
	"../../log"
)

func K000Process(report K000Report) checkup.ReportResult {
	var result checkup.ReportResult
	for host, hostData := range report.Results {
		for table, tableData := range hostData.Data.Queries {
			log.Dbg(host, table, tableData)
		}
	}
	return result
}

func K000PreprocessReportData(data map[string]interface{}) {
	var report K000Report
	filePath := data["source_path_full"].(string)
	jsonRaw := checkup.LoadRawJsonReport(filePath)
	if !checkup.CheckUnmarshalResult(json.Unmarshal(jsonRaw, &report)) {
		return
	}
	result := K000Process(report)
	if len(result.Recommendations) > 0 {
	}
	// update data and file
	//checkup.SaveReportResult(data, result)
}
