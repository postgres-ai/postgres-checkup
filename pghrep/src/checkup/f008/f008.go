package f008

import (
	"encoding/json"

	checkup ".."
	"../../log"
)

func F008Process(report F008Report) checkup.ReportResult {
	var result checkup.ReportResult
	for host, hostData := range report.Results {
		for setting, settingData := range hostData.Data {
			log.Dbg(host, setting, settingData)
		}
	}
	return result
}

func F008PreprocessReportData(data map[string]interface{}) {
	var report F008Report
	filePath := data["source_path_full"].(string)
	jsonRaw := checkup.LoadRawJsonReport(filePath)
	if !checkup.CheckUnmarshalResult(json.Unmarshal(jsonRaw, &report)) {
		return
	}
	result := F008Process(report)
	if len(result.Recommendations) > 0 {
	}
	// update data and file
	//checkup.SaveReportResult(data, result)
}
