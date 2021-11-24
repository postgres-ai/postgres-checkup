package f004

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/dustin/go-humanize/english"

	"gitlab.com/postgres-ai/postgres-checkup/pghrep/internal/checkup"
	"gitlab.com/postgres-ai/postgres-checkup/pghrep/internal/fmtutils"
)

const F004_TOTAL_BLOAT_EXCESS string = "F004_TOTAL_BLOAT_EXCESS"
const F004_TOTAL_BLOAT_LOW string = "F004_TOTAL_BLOAT_LOW"
const F004_BLOAT_CRITICAL string = "F004_BLOAT_CRITICAL"
const F004_BLOATED_TABLES string = "F004_BLOATED_TABLES"
const F004_BLOAT_WARNING string = "F004_BLOAT_WARNING"
const F004_BLOAT_INFO string = "F004_BLOAT_INFO"
const F004_GENERAL_INFO string = "F004_GENERAL_INFO"

const WARNING_BLOAT_RATIO float32 = 40.0
const CRITICAL_BLOAT_RATIO float32 = 90.0
const CRITICAL_TOTAL_BLOAT_RATIO float32 = 20.0
const MIN_TABLE_SIZE_TO_ANALYZE int64 = 1024 * 1024

func appendTable(list []string, tableBloatData F004HeapBloat) []string {
	return append(list, fmt.Sprintf(TABLE_DETAILS, tableBloatData.TableName,
		fmtutils.ByteFormat(float64(tableBloatData.RealSizeBytes), 2),
		tableBloatData.BloatRatioFactor, fmtutils.ByteFormat(float64(tableBloatData.ExtraSizeBytes), 2),
		tableBloatData.BloatRatioPercent))
}

// Generate conclusions and recommendatons
func F004Process(report F004Report) checkup.ReportResult {
	var result checkup.ReportResult
	// check total values
	var criticalTables []string
	var warningTables []string
	var bloatedTables []string
	totalBloatIsCritical := false
	var totalData F004HeapBloatTotal
	var databaseSize int64
	for _, hostData := range report.Results {
		sortedTables := checkup.GetItemsSortedByNum(hostData.Data.HeapBloat)
		databaseSize = hostData.Data.DatabaseSizeBytes
		totalData = hostData.Data.HeapBloatTotal
		if hostData.Data.HeapBloatTotal.BloatRatioPercentAvg > CRITICAL_TOTAL_BLOAT_RATIO {
			totalBloatIsCritical = true
			result.P1 = true
		}
		for _, table := range sortedTables {
			heapBloatData := hostData.Data.HeapBloat[table]
			if (heapBloatData.RealSizeBytes > MIN_TABLE_SIZE_TO_ANALYZE) && (heapBloatData.BloatRatioPercent >= WARNING_BLOAT_RATIO) &&
				(heapBloatData.BloatRatioPercent < CRITICAL_BLOAT_RATIO) && (heapBloatData.BloatRatioFactor > 0) {
				warningTables = appendTable(warningTables, heapBloatData)
			}
			if (heapBloatData.RealSizeBytes > MIN_TABLE_SIZE_TO_ANALYZE) && (heapBloatData.BloatRatioPercent >= CRITICAL_BLOAT_RATIO) &&
				(heapBloatData.BloatRatioFactor > 0) {
				criticalTables = appendTable(criticalTables, heapBloatData)
			}
			if (heapBloatData.RealSizeBytes > MIN_TABLE_SIZE_TO_ANALYZE) && (heapBloatData.BloatRatioPercent >= WARNING_BLOAT_RATIO) {
				bloatedTables = append(bloatedTables, "`"+heapBloatData.TableName+"`")
			}
		}
	}
	if totalBloatIsCritical {
		result.AppendConclusion(F004_TOTAL_BLOAT_EXCESS, MSG_TOTAL_BLOAT_EXCESS_CONCLUSION,
			fmtutils.ByteFormat(float64(totalData.BloatSizeBytesSum), 2),
			totalData.BloatRatioPercentAvg,
			float64(float64(totalData.BloatSizeBytesSum)/float64(databaseSize)*100),
			fmtutils.ByteFormat(float64(databaseSize-totalData.BloatSizeBytesSum), 2),
			fmtutils.ByteFormat(float64(totalData.BloatSizeBytesSum), 2),
			totalData.BloatRatioFactorAvg)
		result.P1 = true
	} else {
		result.AppendConclusion(F004_TOTAL_BLOAT_LOW, MSG_TOTAL_BLOAT_LOW_CONCLUSION, totalData.BloatRatioPercentAvg,
			fmtutils.ByteFormat(float64(totalData.BloatSizeBytesSum), 2))
	}
	if len(criticalTables) > 0 {
		result.AppendConclusion(F004_BLOAT_CRITICAL, english.PluralWord(len(criticalTables), MSG_BLOAT_CRITICAL_CONCLUSION_1,
			MSG_BLOAT_CRITICAL_CONCLUSION_N), len(criticalTables), CRITICAL_BLOAT_RATIO,
			strings.Join(checkup.LimitList(criticalTables), ""))
		result.AppendRecommendation(F004_BLOAT_CRITICAL, MSG_BLOAT_CRITICAL_RECOMMENDATION)
		result.P1 = true
	}
	if len(warningTables) > 0 {
		result.AppendConclusion(F004_BLOAT_WARNING, english.PluralWord(len(warningTables), MSG_BLOAT_WARNING_CONCLUSION_1,
			MSG_BLOAT_WARNING_CONCLUSION_N), len(warningTables), WARNING_BLOAT_RATIO, CRITICAL_BLOAT_RATIO,
			strings.Join(checkup.LimitList(warningTables), ""))
		if !result.P1 {
			result.AppendRecommendation(F004_BLOAT_WARNING, MSG_BLOAT_WARNING_RECOMMENDATION)
		}
		result.P2 = true
	}
	if len(bloatedTables) > 0 {
		result.AppendRecommendation(F004_BLOATED_TABLES, MSG_BLOAT_WARNING_RECOMMENDATION_TABLES, WARNING_BLOAT_RATIO, strings.Join(bloatedTables, ", "))
	}
	if len(result.Recommendations) > 0 {
		result.AppendRecommendation(F004_GENERAL_INFO, MSG_BLOAT_GENERAL_RECOMMENDATION_1)
		result.AppendRecommendation(F004_GENERAL_INFO, MSG_BLOAT_GENERAL_RECOMMENDATION_2)
	}
	if result.P1 || result.P2 {
		result.AppendRecommendation(F004_BLOAT_INFO, MSG_BLOAT_PX_RECOMMENDATION)
	}
	return result
}

func F004LoadReportData(filePath string) (F004Report, error) {
	var report F004Report
	jsonRaw := checkup.LoadRawJsonReport(filePath)

	if !checkup.CheckUnmarshalResult(json.Unmarshal(jsonRaw, &report)) {
		return report, fmt.Errorf("Unable to load F004 report.")
	}

	return report, nil
}

func F004PreprocessReportData(data map[string]interface{}) {
	var report F004Report
	var err error
	filePath := data["source_path_full"].(string)

	report, err = F004LoadReportData(filePath)
	if err != nil {
		return
	}

	result := F004Process(report)
	checkup.SaveReportResult(data, result)
}

func F004GetBloatedTables(report F004Report) []string {
	var bloatedTables []string
	for _, hostData := range report.Results {
		sortedTables := checkup.GetItemsSortedByNum(hostData.Data.HeapBloat)
		for _, table := range sortedTables {
			heapBloatData := hostData.Data.HeapBloat[table]
			if (heapBloatData.RealSizeBytes > MIN_TABLE_SIZE_TO_ANALYZE) && (heapBloatData.BloatRatioPercent >= WARNING_BLOAT_RATIO) {
				bloatedTables = append(bloatedTables, heapBloatData.TableName)
			}
		}
	}
	return bloatedTables
}
