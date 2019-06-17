package k000

import (
	"fmt"
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

func TestK000TotalExcessNodes(t *testing.T) {
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

	report.Results = K000ReportHostsResults{"test-host-1": hostResult, "test-host-2": hostResult}

	result, err := K000Process(report)

	if err != nil || !result.P1 ||
		!checkup.ResultInList(result.Conclusions, K000_EXCESS_QUERY_TOTAL_TIME) ||
		!checkup.ResultInList(result.Recommendations, K000_EXCESS_QUERY_TOTAL_TIME) {
		t.Fatal()
	}

	checkup.PrintResultConclusions(result)
	checkup.PrintResultRecommendations(result)
}
