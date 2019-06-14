package f001

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	checkup ".."
)

const AUTOVACUUM_VACUUM_SCALE_FACTOR_DEFAULT float64 = 0.2
const AUTOVACUUM_VACUUM_THRESHOLD_DEFAULT float64 = 50
const AUTOVACUUM_ANALYZE_SCALE_FACTOR_DEFAULT float64 = 0.1
const AUTOVACUUM_ANALYZE_THRESHOLD_DEFAULT float64 = 50
const AUTOVACUUM_VACUUM_COST_DELAY_DEFAULT float64 = 20
const AUTOVACUUM_VACUUM_COST_LIMIT_DEFAULT float64 = -1

// Case codes
const F001_AUTOVACUUM_NOT_TUNED string = "F001_AUTOVACUUM_NOT_TUNED"
const F001_AUTOVACUUM_TUNE_RECOMMENDATION string = "F001_AUTOVACUUM_TUNE_RECOMMENDATION"

func F001Process(report F001Report) (checkup.ReportResult, error) {
	var result checkup.ReportResult
	var masterHostName string = checkup.GetMasterHostName(report.LastNodesJson.Hosts)

	for host, hostData := range report.Results {
		if host != masterHostName {
			continue
		}

		autovacuumVacuumScaleFactor, found1 := hostData.Data.Settings.GlobalSettings["autovacuum_vacuum_scale_factor"]
		autovacuumVacuumThreshold, found2 := hostData.Data.Settings.GlobalSettings["autovacuum_vacuum_threshold"]
		autovacuumAnalyzeScaleFactor, found3 := hostData.Data.Settings.GlobalSettings["autovacuum_analyze_scale_factor"]
		autovacuumAnalyzeThreshold, found4 := hostData.Data.Settings.GlobalSettings["autovacuum_analyze_threshold"]
		autovacuumVacuumCostDelay, found5 := hostData.Data.Settings.GlobalSettings["autovacuum_vacuum_cost_delay"]
		autovacuumVacuumCostLimit, found6 := hostData.Data.Settings.GlobalSettings["autovacuum_vacuum_cost_limit"]

		if !found1 || !found2 || !found3 || !found4 || !found5 || !found6 {
			return result, fmt.Errorf("Data loading error")
		}

		autovacuumVacuumScaleFactorValue, err1 := strconv.ParseFloat(autovacuumVacuumScaleFactor.Setting, 64)
		autovacuumVacuumThresholdValue, err2 := strconv.ParseFloat(autovacuumVacuumThreshold.Setting, 64)
		autovacuumAnalyzeScaleFactorValue, err3 := strconv.ParseFloat(autovacuumAnalyzeScaleFactor.Setting, 64)
		autovacuumAnalyzeThresholdValue, err4 := strconv.ParseFloat(autovacuumAnalyzeThreshold.Setting, 64)
		autovacuumVacuumCostDelayValue, err5 := strconv.ParseFloat(autovacuumVacuumCostDelay.Setting, 64)
		autovacuumVacuumCostLimitValue, err6 := strconv.ParseFloat(autovacuumVacuumCostLimit.Setting, 64)

		if err1 != nil || err2 != nil || err3 != nil || err4 != nil || err5 != nil || err6 != nil {
			return result, fmt.Errorf("Data loading error")
		}

		var defaultValues []string

		if autovacuumVacuumScaleFactorValue == AUTOVACUUM_VACUUM_SCALE_FACTOR_DEFAULT &&
			autovacuumVacuumThresholdValue == AUTOVACUUM_VACUUM_THRESHOLD_DEFAULT {
			defaultValues = append(defaultValues, fmt.Sprintf(
				"    - `autovacuum_vacuum_scale_factor` = %.1f and `autovacuum_vacuum_threshold` = %.0f  \n",
				AUTOVACUUM_VACUUM_SCALE_FACTOR_DEFAULT, AUTOVACUUM_VACUUM_THRESHOLD_DEFAULT))
		}

		if autovacuumAnalyzeScaleFactorValue == AUTOVACUUM_ANALYZE_SCALE_FACTOR_DEFAULT &&
			autovacuumAnalyzeThresholdValue == AUTOVACUUM_ANALYZE_THRESHOLD_DEFAULT {
			defaultValues = append(defaultValues, fmt.Sprintf(
				"    - `autovacuum_analyze_scale_factor` = %.1f and `autovacuum_analyze_threshold` = %.f  \n",
				AUTOVACUUM_ANALYZE_SCALE_FACTOR_DEFAULT, AUTOVACUUM_ANALYZE_THRESHOLD_DEFAULT))
		}

		if autovacuumVacuumCostDelayValue == AUTOVACUUM_VACUUM_COST_DELAY_DEFAULT &&
			autovacuumVacuumCostLimitValue == AUTOVACUUM_VACUUM_COST_LIMIT_DEFAULT {
			defaultValues = append(defaultValues, fmt.Sprintf(
				"    - `autovacuum_vacuum_cost_delay` = %.0f and `autovacuum_vacuum_cost_limit` = %.0f  \n",
				AUTOVACUUM_VACUUM_COST_DELAY_DEFAULT, AUTOVACUUM_VACUUM_COST_LIMIT_DEFAULT))
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

	return result, nil
}

func F001PreprocessReportData(data map[string]interface{}) {
	var report F001Report
	filePath := data["source_path_full"].(string)
	jsonRaw := checkup.LoadRawJsonReport(filePath)

	if !checkup.CheckUnmarshalResult(json.Unmarshal(jsonRaw, &report)) {
		return
	}

	result, err := F001Process(report)

	if err == nil {
		checkup.SaveReportResult(data, result)
	}
}
