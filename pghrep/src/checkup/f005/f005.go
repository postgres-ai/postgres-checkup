package f005

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	checkup ".."
	"../../fmtutils"
	"../f004"
	"github.com/dustin/go-humanize/english"
)

const F005_TOTAL_BLOAT_EXCESS string = "F005_TOTAL_BLOAT_EXCESS"
const F005_TOTAL_BLOAT_LOW string = "F005_TOTAL_BLOAT_LOW"
const F005_BLOAT_CRITICAL string = "F005_BLOAT_CRITICAL"
const F005_BLOAT_WARNING string = "F005_BLOAT_WARNING"
const F005_BLOAT_CRITICAL_INFO string = "F005_BLOAT_CRITICAL_INFO"
const F005_BLOAT_INFO string = "F005_BLOAT_INFO"
const F005_GENERAL_INFO string = "F005_GENERAL_INFO"
const F005_BLOAT_EXCESS_INFO string = "F005_BLOAT_EXCESS_INFO"
const F005_BLOATED_INDEXES string = "F005_BLOATED_INDEXES"
const F005_BLOATED_TABLE_INDEXES string = "F005_BLOATED_TABLE_INDEXES"

const WARNING_BLOAT_RATIO float32 = 40.0
const CRITICAL_BLOAT_RATIO float32 = 90.0
const CRITICAL_TOTAL_BLOAT_RATIO float32 = 20.0
const MIN_INDEX_SIZE_TO_ANALYZE int64 = 1024 * 1024

func appendIndex(list []string, indexBloatData F005IndexBloat) []string {
	var indexName = indexBloatData.IndexName
	if indexBloatData.SchemaName != "" {
		indexName = indexBloatData.SchemaName + "." + indexName
	}
	return append(list, fmt.Sprintf(INDEX_DETAILS, indexName,
		fmtutils.ByteFormat(float64(indexBloatData.RealSizeBytes), 2),
		indexBloatData.BloatRatioFactor, fmtutils.ByteFormat(float64(indexBloatData.ExtraSizeBytes), 2),
		indexBloatData.BloatRatioPercent))
}

// Generate conclusions and recommendatons
func F005Process(report F005Report, bloatedTables []string) checkup.ReportResult {
	var result checkup.ReportResult
	// check total values
	var top5Indexes []string
	var criticalIndexes []string
	var warningIndexes []string
	var bloatedIndexes []string
	var bloatedTableIndexes []string
	totalBloatIsCritical := false
	var totalData F005IndexBloatTotal
	var databaseSize int64
	i := 0
	for _, hostData := range report.Results {
		sortedIndexes := checkup.GetItemsSortedByNum(hostData.Data.IndexBloat)
		databaseSize = hostData.Data.DatabaseSizeBytes
		totalData = hostData.Data.IndexBloatTotal
		if hostData.Data.IndexBloatTotal.BloatRatioPercentAvg > CRITICAL_TOTAL_BLOAT_RATIO {
			totalBloatIsCritical = true
			result.P1 = true
		}
		for _, index := range sortedIndexes {
			indexBloatData := hostData.Data.IndexBloat[index]
			if totalBloatIsCritical && indexBloatData.RealSizeBytes > MIN_INDEX_SIZE_TO_ANALYZE && i < 5 &&
				(indexBloatData.BloatRatioFactor > 0) {
				top5Indexes = appendIndex(top5Indexes, indexBloatData)
				i++
			}
			if indexBloatData.RealSizeBytes > MIN_INDEX_SIZE_TO_ANALYZE && indexBloatData.BloatRatioPercent >= CRITICAL_BLOAT_RATIO &&
				(indexBloatData.BloatRatioFactor > 0) {
				criticalIndexes = appendIndex(criticalIndexes, indexBloatData)
			}
			if (indexBloatData.RealSizeBytes > MIN_INDEX_SIZE_TO_ANALYZE) && (indexBloatData.BloatRatioPercent >= WARNING_BLOAT_RATIO) &&
				(indexBloatData.BloatRatioPercent < CRITICAL_BLOAT_RATIO) && (indexBloatData.BloatRatioFactor > 0) {
				warningIndexes = appendIndex(warningIndexes, indexBloatData)
			}
			if (indexBloatData.RealSizeBytes > MIN_INDEX_SIZE_TO_ANALYZE) && (indexBloatData.BloatRatioPercent >= WARNING_BLOAT_RATIO) {
				var indexName = indexBloatData.IndexName
				if indexBloatData.SchemaName != "" {
					indexName = indexBloatData.SchemaName + "." + indexName
				}
				if ok, _ := checkup.StringInArray(indexBloatData.TableName, bloatedTables); ok {
					bloatedTableIndexes = append(bloatedTableIndexes, "`"+indexName+"`")
				} else {
					bloatedIndexes = append(bloatedIndexes, "`"+indexName+"`")
				}
			}
		}
	}
	if totalBloatIsCritical {
		result.AppendConclusion(F005_TOTAL_BLOAT_EXCESS, MSG_TOTAL_BLOAT_EXCESS_CONCLUSION,
			fmtutils.ByteFormat(float64(totalData.BloatSizeBytesSum), 2),
			totalData.BloatRatioPercentAvg,
			float64(float64(totalData.BloatSizeBytesSum)/float64(databaseSize)*100),
			fmtutils.ByteFormat(float64(databaseSize-totalData.BloatSizeBytesSum), 2),
			fmtutils.ByteFormat(float64(totalData.BloatSizeBytesSum), 2),
			totalData.BloatRatioFactorAvg)
		result.P1 = true
		result.AppendRecommendation(F005_BLOAT_CRITICAL_INFO, MSG_BLOAT_CRITICAL_RECOMMENDATION)
	} else {
		result.AppendConclusion(F005_TOTAL_BLOAT_LOW, MSG_TOTAL_BLOAT_LOW_CONCLUSION, totalData.BloatRatioPercentAvg,
			fmtutils.ByteFormat(float64(totalData.BloatSizeBytesSum), 2))
	}
	if len(criticalIndexes) > 0 {
		result.AppendConclusion(F005_BLOAT_CRITICAL,
			english.PluralWord(len(criticalIndexes), MSG_BLOAT_CRITICAL_CONCLUSION_1, MSG_BLOAT_CRITICAL_CONCLUSION_N),
			len(criticalIndexes), CRITICAL_BLOAT_RATIO,
			strings.Join(checkup.LimitList(criticalIndexes), ""))

		if !checkup.ResultInList(result.Recommendations, F005_BLOAT_CRITICAL_INFO) {
			result.AppendRecommendation(F005_BLOAT_CRITICAL_INFO, MSG_BLOAT_CRITICAL_RECOMMENDATION)
		}
		result.P1 = true
	}
	if len(warningIndexes) > 0 {
		result.AppendConclusion(F005_BLOAT_WARNING, english.PluralWord(len(warningIndexes), MSG_BLOAT_WARNING_CONCLUSION_1,
			MSG_BLOAT_WARNING_CONCLUSION_N), len(warningIndexes), WARNING_BLOAT_RATIO, CRITICAL_BLOAT_RATIO,
			strings.Join(checkup.LimitList(warningIndexes), ""))
		if !result.P1 {
			result.AppendRecommendation(F005_BLOAT_WARNING, MSG_BLOAT_WARNING_RECOMMENDATION)
		}
		result.P2 = true
	}
	if len(bloatedIndexes) > 0 {
		result.AppendRecommendation(F005_BLOATED_INDEXES, MSG_BLOAT_WARNING_RECOMMENDATION_INDEXES, WARNING_BLOAT_RATIO,
			strings.Join(bloatedIndexes, ", "))
	}
	if len(bloatedTableIndexes) > 0 {
		result.AppendRecommendation(F005_BLOATED_TABLE_INDEXES, MSG_BLOAT_WARNING_RECOMMENDATION_TABLE_INDEXES,
			WARNING_BLOAT_RATIO, strings.Join(bloatedTableIndexes, ", "))
	}
	if len(result.Recommendations) > 0 {
		result.AppendRecommendation(F005_GENERAL_INFO, MSG_BLOAT_GENERAL_RECOMMENDATION_1)
		result.AppendRecommendation(F005_GENERAL_INFO, MSG_BLOAT_GENERAL_RECOMMENDATION_2)
	}
	if result.P1 || result.P2 {
		result.AppendRecommendation(F005_BLOAT_EXCESS_INFO, MSG_BLOAT_PX_RECOMMENDATION)
	}
	return result
}

func F005PreprocessReportData(data map[string]interface{}) {
	var report F005Report
	filePath := data["source_path_full"].(string)
	jsonRaw := checkup.LoadRawJsonReport(filePath)
	if !checkup.CheckUnmarshalResult(json.Unmarshal(jsonRaw, &report)) {
		return
	}

	i := strings.LastIndex(filePath, string(os.PathSeparator))
	path := filePath[:i+1]
	f004FilePath := path + string(os.PathSeparator) + "F004_heap_bloat.json"
	f004Report, err := f004.F004LoadReportData(f004FilePath)
	var bloatedTables []string
	if err == nil {
		bloatedTables = f004.F004GetBloatedTables(f004Report)
	}

	result := F005Process(report, bloatedTables)
	checkup.SaveReportResult(data, result)
}
