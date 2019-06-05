package g001

import (
	"encoding/json"

	checkup ".."
	"../../log"
)

func G001Process(report G001Report) checkup.ReportResult {
	var result checkup.ReportResult
	for host, hostData := range report.Results {
		for setting, settingData := range hostData.Data {
			log.Dbg(host, setting, settingData)
		}
	}
	return result
}

func G001PreprocessReportData(data map[string]interface{}) {
	var report G001Report
	filePath := data["source_path_full"].(string)
	jsonRaw := checkup.LoadRawJsonReport(filePath)
	if !checkup.CheckUnmarshalResult(json.Unmarshal(jsonRaw, &report)) {
		return
	}
	result := G001Process(report)
	if len(result.Recommendations) > 0 {
	}
	// update data and file
	//checkup.SaveReportResult(data, result)
}
