package l003

import (
	"encoding/json"

	checkup ".."
	"../../log"
)

func L003Process(report L003Report) checkup.ReportResult {
	var result checkup.ReportResult
	for host, hostData := range report.Results {
		for table, tableData := range hostData.Data {
			log.Dbg(host, table, tableData)
		}
	}
	return result
}

func L003PreprocessReportData(data map[string]interface{}) {
	var report L003Report
	filePath := data["source_path_full"].(string)
	jsonRaw := checkup.LoadRawJsonReport(filePath)
	if !checkup.CheckUnmarshalResult(json.Unmarshal(jsonRaw, &report)) {
		return
	}
	result := L003Process(report)
	if len(result.Recommendations) > 0 {
	}
	// update data and file
	//checkup.SaveReportResult(data, result)
}
