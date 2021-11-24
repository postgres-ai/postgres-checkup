package h002

import (
	"encoding/json"

	"gitlab.com/postgres-ai/postgres-checkup/pghrep/internal/checkup"
)

const H002_UNUSED_INDEXES_FOUND_P2 string = "H002_UNUSED_INDEXES_FOUND_P2"
const H002_UNUSED_INDEXES_FOUND_P3 string = "H002_UNUSED_INDEXES_FOUND_P3"
const H002_UNUSED_INDEXES_FOUND string = "H002_UNUSED_INDEXES_FOUND"
const H002_UNUSED_INDEXES_FOUND_DO string = "H002_UNUSED_INDEXES_FOUND_DO"
const H002_UNUSED_INDEXES_FOUND_UNDO string = "H002_UNUSED_INDEXES_FOUND_UNDO"

const UNUSED_INDEXES_CRITICL_SIZE_PERCENT float64 = 5.0

func H002Process(report H002Report) (checkup.ReportResult, error) {
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

		if len(hostData.Data.NeverUsedIndexes) > 0 && len(hostData.Data.Do) > 0 &&
			len(hostData.Data.UnDo) > 0 {
			if (float64(hostData.Data.NeverUsedIndexesTotal.IndexSizeBytesSum) /
				float64(hostData.Data.DatabaseStat.DatabaseSizeBytes) * 100) > UNUSED_INDEXES_CRITICL_SIZE_PERCENT {
				result.P2 = true
				result.AppendConclusion(H002_UNUSED_INDEXES_FOUND_P2, MSG_UNUSED_INDEXES_FOUND_P2_CONCLUSION,
					len(hostData.Data.NeverUsedIndexes), UNUSED_INDEXES_CRITICL_SIZE_PERCENT)
			} else {
				result.P3 = true
				result.AppendConclusion(H002_UNUSED_INDEXES_FOUND_P3, MSG_UNUSED_INDEXES_FOUND_P3_CONCLUSION,
					len(hostData.Data.NeverUsedIndexes))
			}

			var p = "[P3] "
			if result.P2 {
				p = "[P2] "
			}
			result.AppendRecommendation(H002_UNUSED_INDEXES_FOUND, p+MSG_UNUSED_INDEXES_FOUND_R1)
			result.AppendRecommendation(H002_UNUSED_INDEXES_FOUND, MSG_UNUSED_INDEXES_FOUND_R2)
			result.AppendRecommendation(H002_UNUSED_INDEXES_FOUND, MSG_UNUSED_INDEXES_FOUND_R3)

			var doCode = "```  \n"
			for _, doIndex := range hostData.Data.Do {
				doCode = doCode + doIndex + "  \n"
			}
			doCode = doCode + "```"
			result.AppendRecommendation(H002_UNUSED_INDEXES_FOUND_DO, MSG_UNUSED_INDEXES_FOUND_DO, doCode)

			var undoCode = "```  \n"
			for _, undoIndex := range hostData.Data.UnDo {
				undoCode = undoCode + undoIndex + "  \n"
			}
			undoCode = undoCode + "```"
			result.AppendRecommendation(H002_UNUSED_INDEXES_FOUND_UNDO, MSG_UNUSED_INDEXES_FOUND_UNDO, undoCode)
		}

	}

	return result, nil
}

func H002PreprocessReportData(data map[string]interface{}) {
	var report H002Report
	filePath := data["source_path_full"].(string)
	jsonRaw := checkup.LoadRawJsonReport(filePath)

	if !checkup.CheckUnmarshalResult(json.Unmarshal(jsonRaw, &report)) {
		return
	}

	result, err := H002Process(report)

	if err == nil {
		checkup.SaveReportResult(data, result)
	}
}
