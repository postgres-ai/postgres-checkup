package a001

import (
	"encoding/json"
	"fmt"

	checkup ".."
)

func A001LoadReportData(filePath string) (A001Report, error) {
	var report A001Report
	jsonRaw := checkup.LoadRawJsonReport(filePath)

	if !checkup.CheckUnmarshalResult(json.Unmarshal(jsonRaw, &report)) {
		return report, fmt.Errorf("Unable to load A001 report.")
	}

	return report, nil
}
