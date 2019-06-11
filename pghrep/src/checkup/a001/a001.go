package a001

import (
	"encoding/json"

	checkup ".."
)

func A001LoadReportData(filePath string) (bool, A001Report) {
	var report A001Report
	jsonRaw := checkup.LoadRawJsonReport(filePath)

	if !checkup.CheckUnmarshalResult(json.Unmarshal(jsonRaw, &report)) {
		return false, report
	}

	return true, report
}
