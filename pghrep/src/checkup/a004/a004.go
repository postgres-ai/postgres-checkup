package a004

import (
	"encoding/json"

	checkup ".."
	"../../log"
)

func A004Process(report A004Report) checkup.ReportResult {
	var result checkup.ReportResult
	for host, hostData := range report.Results {
		for metric, metricData := range hostData.Data.GeneralInfo {
			log.Dbg(host, metric, metricData)
		}
		for database, size := range hostData.Data.DatabaseSizes {
			log.Dbg(host, database, size)
		}
	}
	return result
}

func A004PreprocessReportData(data map[string]interface{}) {
	var report A004Report
	filePath := data["source_path_full"].(string)
	jsonRaw := checkup.LoadRawJsonReport(filePath)
	if !checkup.CheckUnmarshalResult(json.Unmarshal(jsonRaw, &report)) {
		return
	}
	result := A004Process(report)
	if len(result.Recommendations) > 0 {
	}
	// update data and file
	//checkup.SaveReportResult(data, result)
}
