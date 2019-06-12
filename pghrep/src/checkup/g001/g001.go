package g001

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	checkup ".."
	"../../fmtutils"
	"../a001"
)

const MIN_PERCENT int64 = 20
const MAX_PERCENT int64 = 80
const MIDDLE_PERCENT int64 = 25
const OOM_PERCENT int64 = 90

const G001_SHARED_BUFFERS_NOT_OPTIMAL string = "G001_SHARED_BUFFERS_NOT_OPTIMAL"
const G001_TUNE_SHARED_BUFFERS string = "G001_TUNE_SHARED_BUFFERS"
const G001_OOM string = "G001_OOM"

func getHostMemTotal(a001 a001.A001Report, host string) (int64, error) {
	var result int64 = -1
	hostData, ok := a001.Results[host]
	if !ok {
		return -1, fmt.Errorf("Value not found")
	}
	val := hostData.Data.Ram.MemTotal
	vals := strings.Split(val, " ")
	if len(vals) > 1 {
		value, err := strconv.Atoi(vals[0])
		unitFactor := fmtutils.GetUnit(vals[1])
		if err != nil {
			return -1, err
		}
		result = int64(value) * unitFactor
	}
	return result, nil
}

func getSettingValue(hostResult G001ReportHostResult, settingName string) (int64, error) {
	setting, ok := hostResult.Data[settingName]

	if !ok {
		return -1, fmt.Errorf("Setting not found")
	}

	intValue, err := strconv.Atoi(setting.Setting)
	unitFactor := fmtutils.GetUnit(setting.Unit)

	if err != nil {
		return -1, err
	}

	if unitFactor == -1 {
		unitFactor = 1
	}

	if intValue == -1 {
		return -1, nil
	} else {
		return int64(intValue) * unitFactor, nil
	}
}

func G001CheckSharedBuffers(report G001Report, a001 a001.A001Report, result checkup.ReportResult) checkup.ReportResult {
	var conclusions []string
	var recommendations []string

	for host, hostData := range report.Results {
		sharedBufferValue, err1 := getSettingValue(hostData, "shared_buffers")
		hostMemTotal, err2 := getHostMemTotal(a001, host)

		if err1 != nil || err2 != nil {
			continue
		}

		minLevel := hostMemTotal / 100 * MIN_PERCENT
		maxLevel := hostMemTotal / 100 * MAX_PERCENT
		recommendedBytes := hostMemTotal / 100 * MIDDLE_PERCENT
		recommendedValue := fmtutils.ByteFormat(float64(recommendedBytes), 2)

		if sharedBufferValue != -1 && sharedBufferValue > maxLevel {
			conclusions = append(conclusions, fmt.Sprintf(MSG_HOST_CONCLUSION_HIGH, host, fmtutils.ByteFormat(float64(hostMemTotal), 2),
				fmtutils.ByteFormat(float64(sharedBufferValue), 2), int(sharedBufferValue/hostMemTotal*100)))
		}

		if (sharedBufferValue != -1) && sharedBufferValue < minLevel {
			conclusions = append(conclusions, fmt.Sprintf(MSG_HOST_CONCLUSION_LOW, host, fmtutils.ByteFormat(float64(hostMemTotal), 2),
				fmtutils.ByteFormat(float64(sharedBufferValue), 2), int(sharedBufferValue/hostMemTotal*100)))
		}

		if (sharedBufferValue != -1) && (sharedBufferValue > maxLevel || sharedBufferValue < minLevel) {
			result.P1 = true
			recommendations = append(recommendations, fmt.Sprintf(MSG_HOST_RECOMMENDATION,
				host, recommendedValue, MIDDLE_PERCENT, fmtutils.ByteFormat(float64(minLevel), 2),
				MIN_PERCENT, fmtutils.ByteFormat(float64(maxLevel), 2), MAX_PERCENT))
		}
	}

	if len(conclusions) > 0 {
		result.AppendConclusion(G001_SHARED_BUFFERS_NOT_OPTIMAL, MSG_SHARED_BUFFERS_NOT_OPTIMAL_CONCLUSION, strings.Join(conclusions, ",  \n"))
	}
	if len(recommendations) > 0 {
		result.AppendRecommendation(G001_SHARED_BUFFERS_NOT_OPTIMAL, MSG_SHARED_BUFFERS_NOT_OPTIMAL_CONCLUSION, strings.Join(recommendations, ",  \n"))
		result.AppendRecommendation(G001_TUNE_SHARED_BUFFERS, MSG_TUNE_SHARED_BUFFERS_RECOMMENDATION)
	}

	return result
}

func G001CheckOOMRisk(report G001Report, a001 a001.A001Report, result checkup.ReportResult) checkup.ReportResult {
	var masterHostName string = checkup.GetMasterHostName(report.LastNodesJson.Hosts)

	for host, hostData := range report.Results {
		if host != masterHostName {
			continue
		}

		var autovacuumWorkMemEffectiveBytes int64
		sharedBufferValueBytes, err1 := getSettingValue(hostData, "shared_buffers")
		autovacuumWorkMemBytes, err2 := getSettingValue(hostData, "autovacuum_work_mem")
		maintenanceWorkMemBytes, err3 := getSettingValue(hostData, "maintenance_work_mem")
		workMemBytes, err4 := getSettingValue(hostData, "work_mem")
		maxConnections, err5 := getSettingValue(hostData, "max_connections")
		autovacuumMaxWorkers, err5 := getSettingValue(hostData, "autovacuum_max_workers")
		hostMemTotal, err6 := getHostMemTotal(a001, host)

		if err1 != nil || err2 != nil || err3 != nil || err4 != nil || err5 != nil || err6 != nil {
			return result
		}

		if autovacuumWorkMemBytes != -1 {
			autovacuumWorkMemEffectiveBytes = autovacuumWorkMemBytes
		} else {
			autovacuumWorkMemEffectiveBytes = maintenanceWorkMemBytes
		}

		usedMem := (sharedBufferValueBytes + (maxConnections * 2 * workMemBytes) +
			(autovacuumWorkMemEffectiveBytes * autovacuumMaxWorkers))
		maxMem := (OOM_PERCENT * hostMemTotal / 100)

		if usedMem >= maxMem {
			result.P1 = true
			result.AppendConclusion(G001_OOM, MSG_OOM_BASE_CONCLUSION, host)
			result.AppendRecommendation(G001_OOM, MSG_OOM_BASE_RECOMMENDATION)
		}
	}

	return result
}

func G001Process(report G001Report, a001 a001.A001Report) checkup.ReportResult {
	var result checkup.ReportResult
	result = G001CheckSharedBuffers(report, a001, result)
	result = G001CheckOOMRisk(report, a001, result)
	return result
}

func G001PreprocessReportData(data map[string]interface{}) {
	var report G001Report
	filePath := data["source_path_full"].(string)
	jsonRaw := checkup.LoadRawJsonReport(filePath)

	if !checkup.CheckUnmarshalResult(json.Unmarshal(jsonRaw, &report)) {
		return
	}

	i := strings.LastIndex(filePath, string(os.PathSeparator))
	path := filePath[:i+1]
	a001FilePath := path + string(os.PathSeparator) + "A001_system_info.json"
	a001, err := a001.A001LoadReportData(a001FilePath)

	if err != nil {
		return
	}

	result := G001Process(report, a001)

	// update data and file
	checkup.SaveReportResult(data, result)
}
