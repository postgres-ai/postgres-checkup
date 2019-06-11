package a001

import (
	"encoding/json"

	checkup ".."
)

func A001LoadReportData(filePath string) (A001Report, bool) {
	var report A001Report
	jsonRaw := checkup.LoadRawJsonReport(filePath)

	if !checkup.CheckUnmarshalResult(json.Unmarshal(jsonRaw, &report)) {
		return report, false
	}

	return report, true
}
