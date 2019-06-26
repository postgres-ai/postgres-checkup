package l003

import (
	"encoding/json"
	"fmt"
	"strings"

	checkup ".."
	"github.com/dustin/go-humanize/english"
)

const MAX_RATIO_PERCENT float64 = 10

const L003_HIGH_RISKS string = "L003_HIGH_RISKS"

func L003Process(report L003Report) (checkup.ReportResult, error) {
	var result checkup.ReportResult
	var tables []string
	for _, hostData := range report.Results {
		for _, tableData := range hostData.Data {
			if tableData.CapacityUsedPercent > MAX_RATIO_PERCENT {
				tables = append(tables, fmt.Sprintf(MSG_HIGH_RISKS_TABLE,
					tableData.Table, tableData.CurrentMaxValue, tableData.CapacityUsedPercent, tableData.Type))
			}
		}
	}

	if len(tables) > 0 {
		result.P1 = true
		result.AppendConclusion(L003_HIGH_RISKS,
			english.PluralWord(len(tables), MSG_HIGH_RISKS_CONCLUSION_1, MSG_HIGH_RISKS_CONCLUSION_N),
			strings.Join(tables, ";\n"))
		result.AppendRecommendation(L003_HIGH_RISKS, MSG_HIGH_RISKS_RECOMMENDATION)
	}

	return result, nil
}

func L003PreprocessReportData(data map[string]interface{}) {
	var report L003Report
	filePath := data["source_path_full"].(string)
	jsonRaw := checkup.LoadRawJsonReport(filePath)

	if !checkup.CheckUnmarshalResult(json.Unmarshal(jsonRaw, &report)) {
		return
	}

	result, err := L003Process(report)

	if err == nil {
		checkup.SaveReportResult(data, result)
	}

}
