package g002

import (
	"encoding/json"
	"fmt"

	"github.com/dustin/go-humanize/english"

	"gitlab.com/postgres-ai/postgres-checkup/pghrep/internal/checkup"
)

const G002_TX_AGE_MORE_1H string = "G002_TX_AGE_MORE_1H"

func G002Process(report G002Report) (checkup.ReportResult, error) {
	var result checkup.ReportResult
	var hosts []string

	for host, hostData := range report.Results {
		for _, settingData := range hostData.Data {
			if settingData.CurrentState == "idle in transaction" && settingData.TxMore1h > 0 {
				result.P1 = true
				hosts = append(hosts, "`"+host+"`")
			}

			if settingData.CurrentState == "active" && settingData.TxMore1h > 0 {
				hosts = append(hosts, "`"+host+"`")
				result.P1 = true
			}
		}
	}

	if result.P1 && len(hosts) > 0 {
		result.AppendConclusion(G002_TX_AGE_MORE_1H, MSG_TX_AGE_MORE_1H_CONCLUSION,
			fmt.Sprintf(english.PluralWord(len(hosts), MSG_NODE, MSG_NODES), english.WordSeries(hosts, "and")))
		result.AppendRecommendation(G002_TX_AGE_MORE_1H, MSG_TX_AGE_MORE_1H_RECOMMENDATION)
	}

	return result, nil
}

func G002PreprocessReportData(data map[string]interface{}) {
	var report G002Report
	filePath := data["source_path_full"].(string)
	jsonRaw := checkup.LoadRawJsonReport(filePath)

	if !checkup.CheckUnmarshalResult(json.Unmarshal(jsonRaw, &report)) {
		return
	}

	result, err := G002Process(report)

	if err == nil {
		checkup.SaveReportResult(data, result)
	}
}
