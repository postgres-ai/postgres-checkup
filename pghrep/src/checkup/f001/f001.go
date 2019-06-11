package f001

import (
	"encoding/json"
	"strconv"
	"strings"

	checkup ".."
)

// Case codes
const F001_AUTOVACUUM_NOT_TUNED string = "F001_AUTOVACUUM_NOT_TUNED"
const F001_AUTOVACUUM_TUNE_RECOMMENDATION string = "F001_AUTOVACUUM_TUNE_RECOMMENDATION"

func F001Process(report F001Report) checkup.ReportResult {
	var result checkup.ReportResult
	var masterHostName string = checkup.GetMasterHostName(report.LastNodesJson.Hosts)

	for host, hostData := range report.Results {
		if host != masterHostName {
			continue
		}

		autovacuumVacuumScaleFactor := hostData.Data.Settings.GlobalSettings["autovacuum_vacuum_scale_factor"]
		autovacuumVacuumThreshold := hostData.Data.Settings.GlobalSettings["autovacuum_vacuum_threshold"]
		autovacuumAnalyzeScaleFactor := hostData.Data.Settings.GlobalSettings["autovacuum_analyze_scale_factor"]
		autovacuumAnalyzeThreshold := hostData.Data.Settings.GlobalSettings["autovacuum_analyze_threshold"]
		autovacuumVacuumCostDelay := hostData.Data.Settings.GlobalSettings["autovacuum_vacuum_cost_delay"]
		autovacuumVacuumCostLimit := hostData.Data.Settings.GlobalSettings["autovacuum_vacuum_cost_limit"]

		autovacuumVacuumScaleFactorValue, _ := strconv.ParseFloat(autovacuumVacuumScaleFactor.Setting, 64)
		autovacuumVacuumThresholdValue, _ := strconv.ParseFloat(autovacuumVacuumThreshold.Setting, 64)
		autovacuumAnalyzeScaleFactorValue, _ := strconv.ParseFloat(autovacuumAnalyzeScaleFactor.Setting, 64)
		autovacuumAnalyzeThresholdValue, _ := strconv.ParseFloat(autovacuumAnalyzeThreshold.Setting, 64)
		autovacuumVacuumCostDelayValue, _ := strconv.ParseFloat(autovacuumVacuumCostDelay.Setting, 64)
		autovacuumVacuumCostLimitValue, _ := strconv.ParseFloat(autovacuumVacuumCostLimit.Setting, 64)

		var defaultValues []string

		if autovacuumVacuumScaleFactorValue == 0.2 && autovacuumVacuumThresholdValue == 50 {
			defaultValues = append(defaultValues, "    - `autovacuum_vacuum_scale_factor` = 0.2 and `autovacuum_vacuum_threshold` = 50  ")

		}

		if autovacuumAnalyzeScaleFactorValue == 0.1 && autovacuumAnalyzeThresholdValue == 50 {
			defaultValues = append(defaultValues, "    - `autovacuum_analyze_scale_factor` = 0.1 and `autovacuum_analyze_threshold` = 50  ")
		}

		if autovacuumVacuumCostDelayValue == 20 && autovacuumVacuumCostLimitValue == -1 {
			defaultValues = append(defaultValues, "    - `autovacuum_vacuum_cost_delay` = 20 and `autovacuum_vacuum_cost_limit` = -1  ")
		}

		if len(defaultValues) > 0 {
			result.P1 = true
			result.AppendConclusion(F001_AUTOVACUUM_NOT_TUNED, MSG_AUTOVACUUM_NOT_TUNED_CONCLUSION, strings.Join(checkup.LimitList(defaultValues), ""))
			result.AppendRecommendation(F001_AUTOVACUUM_NOT_TUNED, MSG_AUTOVACUUM_NOT_TUNED_RECOMMENDATION)
		}

		break
	}

	if result.P1 || result.P2 {
		result.AppendRecommendation(F001_AUTOVACUUM_TUNE_RECOMMENDATION, MSG_AUTOVACUUM_TUNE_RECOMMENDATION)
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
	checkup.SaveReportResult(data, result)
}
