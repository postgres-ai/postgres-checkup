package h004

import (
	"encoding/json"

	"gitlab.com/postgres-ai/postgres-checkup/pghrep/internal/checkup"
)

const H004_REDUNDANT_INDEXES_FOUND_P2 string = "H004_REDUNDANT_INDEXES_FOUND_P2"
const H004_REDUNDANT_INDEXES_FOUND_P3 string = "H004_REDUNDANT_INDEXES_FOUND_P3"
const H004_REDUNDANT_INDEXES_FOUND string = "H004_REDUNDANT_INDEXES_FOUND"
const H004_REDUNDANT_INDEXES_FOUND_DO string = "H004_REDUNDANT_INDEXES_FOUND_DO"
const H004_REDUNDANT_INDEXES_FOUND_UNDO string = "H004_REDUNDANT_INDEXES_FOUND_UNDO"

const REDUNDANT_INDEXES_CRITICL_SIZE_PERCENT float64 = 5.0

func H004Process(report H004Report) (checkup.ReportResult, error) {
	var result checkup.ReportResult
	var masterHost = ""

	for host, hostData := range report.LastNodesJson.Hosts {
		if hostData.Role == "master" {
			masterHost = host
			break
		}
	}

	for host, hostData := range report.Results {
		if host != masterHost {
			continue
		}

		if len(hostData.Data.RedundantIndexes) > 0 && len(hostData.Data.Do) > 0 &&
			len(hostData.Data.UnDo) > 0 {
			if (float64(hostData.Data.RedundantIndexesTotal.IndexSizeBytesSum) /
				float64(hostData.Data.DatabaseStat.DatabaseSizeBytes) * 100) > REDUNDANT_INDEXES_CRITICL_SIZE_PERCENT {
				result.P2 = true
				result.AppendConclusion(H004_REDUNDANT_INDEXES_FOUND_P2, MSG_REDUNDANT_INDEXES_FOUND_P2_CONCLUSION,
					len(hostData.Data.RedundantIndexes), REDUNDANT_INDEXES_CRITICL_SIZE_PERCENT)
			} else {
				result.P3 = true
				result.AppendConclusion(H004_REDUNDANT_INDEXES_FOUND_P3, MSG_REDUNDANT_INDEXES_FOUND_P3_CONCLUSION,
					len(hostData.Data.RedundantIndexes))
			}

			var p = "[P3] "
			if result.P2 {
				p = "[P2] "
			}
			result.AppendRecommendation(H004_REDUNDANT_INDEXES_FOUND, p+MSG_REDUNDANT_INDEXES_FOUND_R1)
			result.AppendRecommendation(H004_REDUNDANT_INDEXES_FOUND, MSG_REDUNDANT_INDEXES_FOUND_R2)
			result.AppendRecommendation(H004_REDUNDANT_INDEXES_FOUND, MSG_REDUNDANT_INDEXES_FOUND_R3)

			var doCode = "```  \n"
			for _, doIndex := range hostData.Data.Do {
				doCode = doCode + doIndex + "  \n"
			}
			doCode = doCode + "```"
			result.AppendRecommendation(H004_REDUNDANT_INDEXES_FOUND_DO, MSG_REDUNDANT_INDEXES_FOUND_DO, doCode)

			var undoCode = "```  \n"
			for _, undoIndex := range hostData.Data.UnDo {
				undoCode = undoCode + undoIndex + "  \n"
			}
			undoCode = undoCode + "```"
			result.AppendRecommendation(H004_REDUNDANT_INDEXES_FOUND_UNDO, MSG_REDUNDANT_INDEXES_FOUND_UNDO, undoCode)
		}

	}

	return result, nil
}

func H004PreprocessReportData(data map[string]interface{}) {
	var report H004Report
	filePath := data["source_path_full"].(string)
	jsonRaw := checkup.LoadRawJsonReport(filePath)

	if !checkup.CheckUnmarshalResult(json.Unmarshal(jsonRaw, &report)) {
		return
	}

	result, err := H004Process(report)

	if err == nil {
		checkup.SaveReportResult(data, result)
	}
}
