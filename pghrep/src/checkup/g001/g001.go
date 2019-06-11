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

const G001_SHARED_BUFFERS_NOT_OPTIMAL string = "G001_SHARED_BUFFERS_NOT_OPTIMAL"
const G001_TUNE_SHARED_BUFFERS string = "G001_TUNE_SHARED_BUFFERS"

func getHostMemTotal(a001 a001.A001Report, host string) int64 {
	var result int64 = -1
	hostData, ok := a001.Results[host]
	if !ok {
		return -1
	}
	val := hostData.Data.Ram.MemTotal
	vals := strings.Split(val, " ")
	if len(vals) > 1 {
		value, err := strconv.Atoi(vals[0])
		unitFactor := fmtutils.GetUnit(vals[1])
		if err != nil {
			return -1
		}
		result = int64(value) * unitFactor
	}
	return result
}

func G001Process(report G001Report, a001 a001.A001Report) checkup.ReportResult {
	var result checkup.ReportResult
	var conclusions []string
	var recommendations []string
	for host, hostData := range report.Results {
		var sharedBufferValue int64 = -1
		var hostMemTotal int64

		hostMemTotal = getHostMemTotal(a001, host)
		recommendedValue := fmtutils.ByteFormat(float64(hostMemTotal/100*MIDDLE_PERCENT), 2)

		if hostMemTotal == -1 {
			continue
		}

		minLevel := hostMemTotal / 100 * MIN_PERCENT
		maxLevel := hostMemTotal / 100 * MAX_PERCENT
		sharedBuffers, ok := hostData.Data["shared_buffers"]

		if !ok {
			continue
		}

		value, err := strconv.Atoi(sharedBuffers.Setting)
		unitFactor := fmtutils.GetUnit(sharedBuffers.Unit)
		if err != nil {
			continue
		}

		sharedBufferValue = int64(value) * unitFactor

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
	ok, a001 := a001.A001LoadReportData(a001FilePath)

	if !ok {
		return
	}

	result := G001Process(report, a001)

	// update data and file
	checkup.SaveReportResult(data, result)
}
