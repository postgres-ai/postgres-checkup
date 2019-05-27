package f004

import (
	"encoding/json"
	"strings"

	"../../log"

	checkup ".."
	"../../fmtutils"
	"../../pyraconv"
)

const WARNING_BLOAT_RATIO float32 = 40.0        //90.0
const CRITICAL_BLOAT_RATIO float32 = 90.0       //90.0
const CRITICAL_TOTAL_BLOAT_RATIO float32 = 20.0 //90.0

// Generate conclusions and recommendatons
func F004Process(report F004Report) checkup.ReportOutcome {
	var result checkup.ReportOutcome
	// check total values
	var top5tables []string
	var criticalTables []string
	var warningTables []string
	totalBloatIsCritical := false
	var totalData F004HeapBloatTotal
	i := 0
	for _, hostData := range report.Results {
		totalData = hostData.Data.HeapBloatTotal
		if hostData.Data.HeapBloatTotal.BloatRatioPercentAvg > CRITICAL_TOTAL_BLOAT_RATIO {
			totalBloatIsCritical = true
			result.P1 = true
		}
		for table, heapBloatData := range hostData.Data.HeapBloat {
			if totalBloatIsCritical && heapBloatData.RealSizeBytes > 1024*1024 && i < 5 {
				top5tables = append(top5tables, "- "+table+": "+fmtutils.ByteFormat(float64(heapBloatData.RealSizeBytes), 2)+"  \n")
				i++
			}
			if (heapBloatData.RealSizeBytes > 1024*1024) && (heapBloatData.BloatRatioPercent >= WARNING_BLOAT_RATIO) &&
				(heapBloatData.BloatRatioPercent < CRITICAL_BLOAT_RATIO) {
				warningTables = append(warningTables, "- "+table+": "+fmtutils.ByteFormat(float64(heapBloatData.RealSizeBytes), 2)+"  \n")
			}
			if heapBloatData.RealSizeBytes > 1024*1024 && heapBloatData.BloatRatioPercent >= CRITICAL_BLOAT_RATIO {
				criticalTables = append(criticalTables, "- "+table+": "+fmtutils.ByteFormat(float64(heapBloatData.RealSizeBytes), 2)+"  \n")
			}

		}
	}
	if totalBloatIsCritical {
		result.AppendConclusion(MSG_TOTAL_BLOAT_EXCESS_CONCLUSION,
			fmtutils.ByteFormat(float64(totalData.BloatSizeBytesSum), 2),
			totalData.BloatRatioPercentAvg,
			fmtutils.ByteFormat(float64(totalData.BloatSizeBytesSum), 2),
			fmtutils.ByteFormat(float64(totalData.BloatSizeBytesSum), 2),
			totalData.BloatRatioAvg,
			strings.Join(top5tables, ""))
		result.P1 = true
	} else {
		result.AppendConclusion(MSG_TOTAL_BLOAT_LOW_CONCLUSION, fmtutils.ByteFormat(float64(totalData.BloatSizeBytesSum), 2))
	}
	if len(warningTables) > 0 {
		result.AppendConclusion(MSG_BLOAT_WARNING_CONCLUSION, WARNING_BLOAT_RATIO, CRITICAL_BLOAT_RATIO, strings.Join(warningTables, ""))
		result.AppendRecommendation(MSG_BLOAT_WARNING_RECOMMENDATION)
		result.P2 = true
	}
	if len(criticalTables) > 0 {
		result.AppendConclusion(MSG_BLOAT_CRITICAL_CONCLUSION, CRITICAL_BLOAT_RATIO, strings.Join(criticalTables, ""))
		result.AppendRecommendation(MSG_BLOAT_CRITICAL_RECOMMENDATION)
		result.P1 = true
	}
	if len(result.Recommendations) > 0 {
		result.AppendRecommendation(MSG_BLOAT_GENERAL_RECOMMENDATION)
	}
	if result.P1 || result.P2 {
		result.AppendRecommendation(MSG_BLOAT_PX_RECOMMENDATION)
	}
	return result
}

func F004PreprocessReportData(data map[string]interface{}) {
	filePath := pyraconv.ToString(data["source_path_full"])
	jsonRaw := checkup.LoadRawJsonReport(filePath)
	var report F004Report
	err := json.Unmarshal(jsonRaw, &report)
	if err != nil {
		log.Err("Cannot load json report to process")
		return
	}
	result := F004Process(report)
	if len(result.Recommendations) == 0 {
		result.AppendRecommendation(MSG_NO_RECOMMENDATIONS)
	}
	checkup.SaveConclusionsRecommendations(data, result)
}
