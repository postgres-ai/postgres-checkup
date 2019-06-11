package g001

import (
	"fmt"
	"testing"

	checkup ".."
	"../a001"
)

func TestG001Success(t *testing.T) {
	fmt.Println(t.Name())
	// G001
	var report G001Report
	var hostResult G001ReportHostResult

	hostResult.Data = map[string]G001Setting{
		"shared_buffers": G001Setting{
			Name:    "shared_buffers",
			Setting: "4194304",
			Unit:    "8kB",
		},
	}
	report.Results = G001ReportHostsResults{"test-host": hostResult}

	// A001
	var a001Report a001.A001Report
	var a001HostResult a001.A001ReportHostResult
	a001HostResult.Data = a001.A001ReportHostResultData{
		Ram: a001.A001ReportRam{
			MemTotal: "65888240 kB",
		},
	}
	a001Report.Results = a001.A001ReportHostsResults{"test-host": a001HostResult}

	result := G001Process(report, a001Report)

	if result.P1 || result.P2 || result.P3 &&
		checkup.ResultInList(result.Conclusions, G001_SHARED_BUFFERS_NOT_OPTIMAL) {
		t.Fatal("TestG001Success failed")
	}

	checkup.PrintResultConclusions(result)
	checkup.PrintResultRecommendations(result)
}

func TestG001Low(t *testing.T) {
	fmt.Println(t.Name())
	// G001
	var report G001Report
	var hostResult G001ReportHostResult

	hostResult.Data = map[string]G001Setting{
		"shared_buffers": G001Setting{
			Name:    "shared_buffers",
			Setting: "1235404",
			Unit:    "8kB",
		},
	}
	report.Results = G001ReportHostsResults{"test-host": hostResult}

	// A001
	var a001Report a001.A001Report
	var a001HostResult a001.A001ReportHostResult
	a001HostResult.Data = a001.A001ReportHostResultData{
		Ram: a001.A001ReportRam{
			MemTotal: "65888240 kB",
		},
	}
	a001Report.Results = a001.A001ReportHostsResults{"test-host": a001HostResult}

	result := G001Process(report, a001Report)

	if !result.P1 ||
		!checkup.ResultInList(result.Conclusions, G001_SHARED_BUFFERS_NOT_OPTIMAL) ||
		!checkup.ResultInList(result.Recommendations, G001_SHARED_BUFFERS_NOT_OPTIMAL) ||
		!checkup.ResultInList(result.Recommendations, G001_TUNE_SHARED_BUFFERS) {
		t.Fatal("TestG001Low failed")
	}

	checkup.PrintResultConclusions(result)
	checkup.PrintResultRecommendations(result)
}

func TestG001High(t *testing.T) {
	fmt.Println(t.Name())
	// G001
	var report G001Report
	var hostResult G001ReportHostResult

	hostResult.Data = map[string]G001Setting{
		"shared_buffers": G001Setting{
			Name:    "shared_buffers",
			Setting: "6753544",
			Unit:    "8kB",
		},
	}
	report.Results = G001ReportHostsResults{"test-host": hostResult}

	// A001
	var a001Report a001.A001Report
	var a001HostResult a001.A001ReportHostResult
	a001HostResult.Data = a001.A001ReportHostResultData{
		Ram: a001.A001ReportRam{
			MemTotal: "65888240 kB",
		},
	}
	a001Report.Results = a001.A001ReportHostsResults{"test-host": a001HostResult}

	result := G001Process(report, a001Report)

	if !result.P1 ||
		!checkup.ResultInList(result.Conclusions, G001_SHARED_BUFFERS_NOT_OPTIMAL) ||
		!checkup.ResultInList(result.Recommendations, G001_SHARED_BUFFERS_NOT_OPTIMAL) ||
		!checkup.ResultInList(result.Recommendations, G001_TUNE_SHARED_BUFFERS) {
		t.Fatal("TestG001High failed")
	}

	checkup.PrintResultConclusions(result)
	checkup.PrintResultRecommendations(result)
}
