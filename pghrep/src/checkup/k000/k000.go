package k000

import (
	"encoding/json"
	"fmt"

	checkup ".."
	"github.com/dustin/go-humanize/english"
)

const MAX_QUERY_TOTAL_TIME float64 = 30 // Percents from all workload during the analyzed period

const K000_EXCESS_QUERY_TOTAL_TIME string = "K000_EXCESS_QUERY_TOTAL_TIME"

func K000Process(report K000Report) (checkup.ReportResult, error) {
	var result checkup.ReportResult
	var hosts []string
	for host, hostData := range report.Results {
		for _, queryData := range hostData.Data.Queries {
			if queryData.RatioTotalTime > MAX_QUERY_TOTAL_TIME {
				result.P1 = true
				hosts = append(hosts, "`"+host+"`")
			}
		}
	}

	if result.P1 && len(hosts) > 0 {
		result.AppendConclusion(K000_EXCESS_QUERY_TOTAL_TIME, MSG_EXCESS_QUERY_TOTAL_TIME_CONCLUSION, MAX_QUERY_TOTAL_TIME,
			fmt.Sprintf(english.PluralWord(len(hosts), MSG_NODE, MSG_NODES), english.WordSeries(hosts, "and")),
			MAX_QUERY_TOTAL_TIME, MAX_QUERY_TOTAL_TIME)
		result.AppendRecommendation(K000_EXCESS_QUERY_TOTAL_TIME, MSG_EXCESS_QUERY_TOTAL_TIME_RECOMMENDATION, MAX_QUERY_TOTAL_TIME)
	}

	return result, nil
}

func K000PreprocessReportData(data map[string]interface{}) {
	var report K000Report
	filePath := data["source_path_full"].(string)
	jsonRaw := checkup.LoadRawJsonReport(filePath)

	if !checkup.CheckUnmarshalResult(json.Unmarshal(jsonRaw, &report)) {
		return
	}

	result, err := K000Process(report)

	if err == nil {
		checkup.SaveReportResult(data, result)
	}
}
