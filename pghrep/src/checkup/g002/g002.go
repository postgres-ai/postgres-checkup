package g002

import (
	"encoding/json"

	checkup ".."
	"../../log"
)

func G002Process(report G002Report) checkup.ReportResult {
	var result checkup.ReportResult
	for host, hostData := range report.Results {
		for setting, settingData := range hostData.Data {
			log.Dbg(host, setting, settingData)
		}
	}
	return result
}

func G002PreprocessReportData(data map[string]interface{}) {
	var report G002Report
	filePath := data["source_path_full"].(string)
	jsonRaw := checkup.LoadRawJsonReport(filePath)
	if !checkup.CheckUnmarshalResult(json.Unmarshal(jsonRaw, &report)) {
		return
	}
	result := G002Process(report)
	if len(result.Recommendations) > 0 {
	}
	// update data and file
	//checkup.SaveReportResult(data, result)
}
