package f008

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	checkup ".."
	"../a001"
	"../a002"
)

const AUTOVACUUM_MAX_WORKERS_DEFAULT int = 3
const AUTOVACUUM_VACUUM_COST_LIMIT_DEFAULT int = -1
const VACUUM_COST_LIMIT_DEFAULT int = 200
const AUTOVACUUM_VACUUM_COST_DELAY_DEFAULT int = 20

const F008_MAX_WORKERS_LOW string = "F008_MAX_WORKERS_LOW"
const F008_DELAY_NOT_TUNED string = "F008_DELAY_NOT_TUNED"
const F008_DELAY_TUNE string = "F008_DELAY_TUNE"

func getCpuCount(a001 a001.A001Report, host string) (int, error) {
	hostData, ok := a001.Results[host]

	if !ok {
		return -1, fmt.Errorf("Value not found")
	}

	value := hostData.Data.Cpu.CpuCount
	count, err := strconv.Atoi(value)

	if err != nil {
		return -1, err
	}

	return count, nil
}

func getMajorVersion(a002 a002.A002Report, host string) (float64, error) {
	hostData, ok := a002.Results[host]

	if !ok {
		return -1, fmt.Errorf("Value not found")
	}

	return hostData.Data.GetMajorVersion()
}

func F008CheckMaxWorkers(report F008Report, a001Report a001.A001Report,
	result checkup.ReportResult) (checkup.ReportResult, error) {
	for host, hostData := range report.Results {
		var cpuCount, maxWorkersCount int
		var err1, err2 error
		cpuCount, err1 = getCpuCount(a001Report, host)
		autovacuumMaxWorkers, found := hostData.Data["autovacuum_max_workers"]

		if found {
			maxWorkersCount, err2 = strconv.Atoi(autovacuumMaxWorkers.Setting)
		}

		if err1 != nil || err2 != nil {
			return result, fmt.Errorf("Data loading error")
		}

		if cpuCount < 16 {
			continue
		} else {
			if maxWorkersCount == AUTOVACUUM_MAX_WORKERS_DEFAULT {
				// P1
				result.P1 = true
				result.AppendConclusion(F008_MAX_WORKERS_LOW, MSG_MAX_WORKERS_LOW_CONCLUSION, cpuCount)
				result.AppendRecommendation(F008_MAX_WORKERS_LOW, MSG_MAX_WORKERS_LOW_RECOMMENDATION,
					cpuCount/4, cpuCount/2)
			}
		}

	}
	return result, nil
}

func F008CheckAutovacuumCost(report F008Report, a002Report a002.A002Report,
	result checkup.ReportResult) (checkup.ReportResult, error) {
	for host, hostData := range report.Results {
		majorVersion, err1 := getMajorVersion(a002Report, host)
		autovacuumVacuumCostLimit, found1 := hostData.Data["autovacuum_vacuum_cost_limit"]
		//vacuumCostLimit, found2 := hostData.Data["vacuum_cost_limit"]
		autovacuumVacuumCostDelay, found3 := hostData.Data["autovacuum_vacuum_cost_delay"]

		if err1 != nil || !found1 || !found3 {
			return result, fmt.Errorf("Data loading error")
		}

		if majorVersion >= 12 {
			return result, nil
		}

		autovacuumVacuumCostLimitValue, err2 := strconv.Atoi(autovacuumVacuumCostLimit.Setting)
		autovacuumVacuumCostDelayValue, err3 := strconv.Atoi(autovacuumVacuumCostDelay.Setting)

		if err2 != nil || err3 != nil {
			return result, fmt.Errorf("Data loading error")
		}

		if (autovacuumVacuumCostLimitValue == AUTOVACUUM_VACUUM_COST_LIMIT_DEFAULT ||
			autovacuumVacuumCostLimitValue == VACUUM_COST_LIMIT_DEFAULT) &&
			autovacuumVacuumCostDelayValue == AUTOVACUUM_VACUUM_COST_DELAY_DEFAULT {
			result.P1 = true
			result.AppendConclusion(F008_DELAY_NOT_TUNED, MSG_AUTOVACUUM_COST_DELAY_NOT_TUNED_CONCLUSION)
			result.AppendRecommendation(F008_DELAY_NOT_TUNED, MSG_AUTOVACUUM_COST_DELAY_NOT_TUNED_RECOMMENDATION)
			result.AppendRecommendation(F008_DELAY_TUNE, MSG_AUTOVACUUM_COST_DELAY_TUNE_RECOMMENDATION, majorVersion)
		}
	}
	return result, nil
}

func F008Process(report F008Report, a001Report a001.A001Report,
	a002Report a002.A002Report) (checkup.ReportResult, error) {
	var result checkup.ReportResult
	var err1, err2 error

	result, err1 = F008CheckMaxWorkers(report, a001Report, result)
	result, err2 = F008CheckAutovacuumCost(report, a002Report, result)

	if err1 != nil || err2 != nil {
		return result, fmt.Errorf("Errors during analyze data")
	} else {
		return result, nil
	}
}

func F008PreprocessReportData(data map[string]interface{}) {
	var report F008Report
	filePath := data["source_path_full"].(string)
	jsonRaw := checkup.LoadRawJsonReport(filePath)

	if !checkup.CheckUnmarshalResult(json.Unmarshal(jsonRaw, &report)) {
		return
	}

	i := strings.LastIndex(filePath, string(os.PathSeparator))
	path := filePath[:i+1]
	a001FilePath := path + string(os.PathSeparator) + "A001_system_info.json"
	a001, err1 := a001.A001LoadReportData(a001FilePath)

	a002FilePath := path + string(os.PathSeparator) + "A002_pgversion.json"
	a002, err2 := a002.A002LoadReportData(a002FilePath)

	if err1 != nil || err2 != nil {
		return
	}

	result, err := F008Process(report, a001, a002)

	if err == nil || (err != nil && len(result.Recommendations) > 0) {
		// update data and file only if processed successful or some recommendations found
		checkup.SaveReportResult(data, result)
	}
}
