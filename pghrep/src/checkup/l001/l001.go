package l001

import (
	"encoding/json"

	checkup ".."
	"../../log"
)

func L001Process(report L001Report) checkup.ReportResult {
	var result checkup.ReportResult
	for host, hostData := range report.Results {
		for table, tableData := range hostData.Data {
			log.Dbg(host, table, tableData)
		}
	}
	return result
}

func L001PreprocessReportData(data map[string]interface{}) {
	var report L001Report
	filePath := data["source_path_full"].(string)
	jsonRaw := checkup.LoadRawJsonReport(filePath)
	if !checkup.CheckUnmarshalResult(json.Unmarshal(jsonRaw, &report)) {
		return
	}
	result := L001Process(report)
	if len(result.Recommendations) > 0 {
	}
	// update data and file
	//checkup.SaveReportResult(data, result)
}
