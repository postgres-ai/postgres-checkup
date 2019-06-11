package f001

import (
	"encoding/json"

	checkup ".."
	"../../log"
)

// Case codes
const F001_AUTOVACUUM_NOT_TUNED string = "F001_AUTOVACUUM_NOT_TUNED"

func F001Process(report F001Report) checkup.ReportResult {
	var result checkup.ReportResult
	for host, hostData := range report.Results {
		for setting, settingData := range hostData.Data.Settings.GlobalSettings {
			log.Dbg(host, setting, settingData)
		}
		for setting, settingData := range hostData.Data.Settings.TableSettings {
			log.Dbg(host, setting, settingData)
		}
	}
	return result
}

func F001PreprocessReportData(data map[string]interface{}) {
	var report F001Report
	filePath := data["source_path_full"].(string)
	jsonRaw := checkup.LoadRawJsonReport(filePath)
	if !checkup.CheckUnmarshalResult(json.Unmarshal(jsonRaw, &report)) {
		return
	}
	result := F001Process(report)
	if len(result.Recommendations) > 0 {
	}
	// update data and file
	//checkup.SaveReportResult(data, result)
}
