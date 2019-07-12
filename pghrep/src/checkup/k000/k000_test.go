package k000

import (
	"fmt"
	"strings"
	"testing"

	checkup ".."
)

func TestK000Success(t *testing.T) {
	fmt.Println(t.Name())

	var report K000Report
	var hostResult K000ReportHostResult
	hostResult.Data = K000HostData{
		Queries: map[string]K00Query{
			"1": K00Query{
				RatioTotalTime: 10.0,
			},
		},
	}

	report.Results = K000ReportHostsResults{"test-host": hostResult}

	result, err := K000Process(report)

	if err != nil || result.P1 {
		t.Fatal()
	}

	checkup.PrintResultConclusions(result)
	checkup.PrintResultRecommendations(result)
}

func TestK000TotalExcess(t *testing.T) {
	fmt.Println(t.Name())

	var report K000Report
	var hostResult K000ReportHostResult
	hostResult.Data = K000HostData{
		Queries: map[string]K00Query{
			"1": K00Query{
				RatioTotalTime: 40.0,
			},
		},
	}

	report.Results = K000ReportHostsResults{"test-host": hostResult}

	result, err := K000Process(report)

	if err != nil || !result.P1 ||
		!checkup.ResultInList(result.Conclusions, K000_EXCESS_QUERY_TOTAL_TIME) ||
		!checkup.ResultInList(result.Recommendations, K000_EXCESS_QUERY_TOTAL_TIME) {
		t.Fatal()
	}

	checkup.PrintResultConclusions(result)
	checkup.PrintResultRecommendations(result)
}

func TestK000HostTwice(t *testing.T) {
	fmt.Println(t.Name())

	var report K000Report
	var hostResult K000ReportHostResult
	hostResult.Data = K000HostData{
		Queries: map[string]K00Query{
			"1": K00Query{
				RatioTotalTime: 40.0,
			},
			"2": K00Query{
				RatioTotalTime: 34.0,
			},
		},
	}

	report.Results = K000ReportHostsResults{"test-host-1": hostResult, "test-host-2": hostResult}

	result, err := K000Process(report)

	if err != nil || !result.P1 ||
		!checkup.ResultInList(result.Conclusions, K000_EXCESS_QUERY_TOTAL_TIME) ||
		!checkup.ResultInList(result.Recommendations, K000_EXCESS_QUERY_TOTAL_TIME) {
		t.Fatal()
	}

	conclusion, err2 := checkup.GetResultItem(result.Conclusions, K000_EXCESS_QUERY_TOTAL_TIME)
	if err2 != nil || strings.Contains(conclusion.Message, "`test-host-1`, `test-host-1`") ||
		strings.Contains(conclusion.Message, "`test-host-2`, `test-host-2`") ||
		strings.Contains(conclusion.Message, "`test-host-2` and `test-host-2`") ||
		strings.Contains(conclusion.Message, "`test-host-1` and `test-host-1`") {
		t.Fatal("Same host twice in conlusion")
	}

	checkup.PrintResultConclusions(result)
	checkup.PrintResultRecommendations(result)
}
