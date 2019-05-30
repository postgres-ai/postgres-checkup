package f005

import (
	"encoding/json"
	"fmt"
	"strings"

	"../../log"

	checkup ".."
	"../../fmtutils"
	"../../pyraconv"
)

const WARNING_BLOAT_RATIO float32 = 40.0
const CRITICAL_BLOAT_RATIO float32 = 90.0
const CRITICAL_TOTAL_BLOAT_RATIO float32 = 20.0
const MIN_TABLE_SIZE_TO_ANALYZE int64 = 1024 * 1024

func appendIndex(list []string, indexBloatData F005IndexBloat) []string {
	return append(list, fmt.Sprintf(INDEX_DETAILS, indexBloatData.IndexName,
		fmtutils.ByteFormat(float64(indexBloatData.RealSizeBytes), 2),
		indexBloatData.BloatRatio, fmtutils.ByteFormat(float64(indexBloatData.ExtraSizeBytes), 2),
		indexBloatData.BloatRatioPercent))
}

// Generate conclusions and recommendatons
func F005Process(report F005Report) checkup.ReportOutcome {
	var result checkup.ReportOutcome
	// check total values
	var top5Indexes []string
	var criticalIndexes []string
	var warningIndexes []string
	totalBloatIsCritical := false
	var totalData F005IndexBloatTotal
	var databaseSize int64
	i := 0
	for _, hostData := range report.Results {
		databaseSize = hostData.Data.DatabaseSizeBytes
		totalData = hostData.Data.IndexBloatTotal
		if hostData.Data.IndexBloatTotal.BloatRatioPercentAvg > CRITICAL_TOTAL_BLOAT_RATIO {
			totalBloatIsCritical = true
			result.P1 = true
		}
		for _, indexBloatData := range hostData.Data.IndexBloat {
			if totalBloatIsCritical && indexBloatData.RealSizeBytes > MIN_TABLE_SIZE_TO_ANALYZE && i < 5 {
				top5Indexes = appendIndex(top5Indexes, indexBloatData)
				i++
			}
			if indexBloatData.RealSizeBytes > MIN_TABLE_SIZE_TO_ANALYZE && indexBloatData.BloatRatioPercent >= CRITICAL_BLOAT_RATIO {
				criticalIndexes = appendIndex(criticalIndexes, indexBloatData)
			}
			if (indexBloatData.RealSizeBytes > MIN_TABLE_SIZE_TO_ANALYZE) && (indexBloatData.BloatRatioPercent >= WARNING_BLOAT_RATIO) &&
				(indexBloatData.BloatRatioPercent < CRITICAL_BLOAT_RATIO) {
				warningIndexes = appendIndex(warningIndexes, indexBloatData)
			}
		}
	}
	if totalBloatIsCritical {
		result.AppendConclusion(MSG_TOTAL_BLOAT_EXCESS_CONCLUSION,
			fmtutils.ByteFormat(float64(totalData.BloatSizeBytesSum), 2),
			totalData.BloatRatioPercentAvg,
			float64(float64(totalData.BloatSizeBytesSum)/float64(databaseSize)*100),
			fmtutils.ByteFormat(float64(databaseSize-totalData.BloatSizeBytesSum), 2),
			fmtutils.ByteFormat(float64(totalData.BloatSizeBytesSum), 2),
			totalData.BloatRatioAvg)
		result.P1 = true
	} else {
		result.AppendConclusion(MSG_TOTAL_BLOAT_LOW_CONCLUSION, totalData.BloatRatioPercentAvg,
			fmtutils.ByteFormat(float64(totalData.BloatSizeBytesSum), 2))
	}
	if len(criticalIndexes) > 0 {
		result.AppendConclusion(MSG_BLOAT_CRITICAL_CONCLUSION, len(criticalIndexes), CRITICAL_BLOAT_RATIO,
			strings.Join(checkup.LimitList(criticalIndexes), ""))
		result.AppendRecommendation(MSG_BLOAT_CRITICAL_RECOMMENDATION)
		result.P1 = true
	}
	if len(warningIndexes) > 0 {
		result.AppendConclusion(MSG_BLOAT_WARNING_CONCLUSION, len(warningIndexes), WARNING_BLOAT_RATIO, CRITICAL_BLOAT_RATIO, strings.Join(checkup.LimitList(warningIndexes), ""))
		if !result.P1 {
			result.AppendRecommendation(MSG_BLOAT_WARNING_RECOMMENDATION)
		}
		result.P2 = true
	}
	if len(result.Recommendations) > 0 {
		result.AppendRecommendation(MSG_BLOAT_GENERAL_RECOMMENDATION_1)
		result.AppendRecommendation(MSG_BLOAT_GENERAL_RECOMMENDATION_2)
	}
	if result.P1 || result.P2 {
		result.AppendRecommendation(MSG_BLOAT_PX_RECOMMENDATION)
	}
	return result
}

func F005PreprocessReportData(data map[string]interface{}) {
	filePath := pyraconv.ToString(data["source_path_full"])
	jsonRaw := checkup.LoadRawJsonReport(filePath)
	var report F005Report
	err := json.Unmarshal(jsonRaw, &report)
	if err != nil {
		log.Err("Cannot load json report to process")
		return
	}
	result := F005Process(report)
	if len(result.Recommendations) == 0 {
		result.AppendRecommendation(MSG_NO_RECOMMENDATIONS)
	}
	checkup.SaveConclusionsRecommendations(data, result)
}
