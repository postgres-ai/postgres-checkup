package f008

import (
	"fmt"
	"testing"

	"gitlab.com/postgres-ai/postgres-checkup/pghrep/internal/checkup"
	"gitlab.com/postgres-ai/postgres-checkup/pghrep/internal/checkup/a001"
	"gitlab.com/postgres-ai/postgres-checkup/pghrep/internal/checkup/a002"
)

var TestLastNodesJson checkup.ReportLastNodes = checkup.ReportLastNodes{
	Hosts: checkup.ReportHosts{
		"test-host": {
			Role: "master",
		},
	},
}
var HostData map[string]F008Setting = map[string]F008Setting{
	"autovacuum_max_workers": F008Setting{
		Name:    "autovacuum_max_workers",
		Setting: "3",
	},
}

func TestF008MaxWorkersLow(t *testing.T) {
	fmt.Println(t.Name())
	// G001
	var report F008Report
	var hostResult F008ReportHostResult

	hostResult.Data = HostData
	report.Results = F008ReportHostsResults{"test-host": hostResult}
	report.LastNodesJson = TestLastNodesJson

	// A001
	var a001Report a001.A001Report
	var a001HostResult a001.A001ReportHostResult
	a001HostResult.Data = a001.A001ReportHostResultData{
		Cpu: a001.A001ReportCpu{
			CpuCount: "18",
		},
	}
	a001Report.Results = a001.A001ReportHostsResults{"test-host": a001HostResult}

	// A002
	var a002Report a002.A002Report
	var a002HostResult a002.A002ReportHostResult
	a002HostResult.Data = a002.A002ReportHostResultData{
		Version:          "PostgreSQL 9.6.11 on x86_64-pc-linux-gnu (Ubuntu 9.6.11-1.pgdg16.04+1), compiled by gcc (Ubuntu 5.4.0-6ubuntu1~16.04.10) 5.4.0 20160609, 64-bit",
		ServerVersionNum: "90611",
		ServerMajorVer:   "9.6",
		ServerMinorVer:   "11",
	}
	a002Report.Results = a002.A002ReportHostsResults{"test-host": a002HostResult}

	result, err := F008Process(report, a001Report, a002Report)

	if err != nil && !result.P1 ||
		!checkup.ResultInList(result.Conclusions, F008_MAX_WORKERS_LOW) ||
		!checkup.ResultInList(result.Recommendations, F008_MAX_WORKERS_LOW) {
		t.Fatal()
	}

	checkup.PrintResultConclusions(result)
	checkup.PrintResultRecommendations(result)
}

func TestF008CostNotTunedMinus1(t *testing.T) {
	fmt.Println(t.Name())
	// G001
	var report F008Report
	var hostResult F008ReportHostResult

	hostResult.Data = HostData
	hostResult.Data["autovacuum_vacuum_cost_limit"] = F008Setting{
		Name:    "autovacuum_vacuum_cost_limit",
		Setting: "-1",
	}
	hostResult.Data["autovacuum_vacuum_cost_delay"] = F008Setting{
		Name:    "autovacuum_vacuum_cost_delay",
		Setting: "20",
	}
	report.Results = F008ReportHostsResults{"test-host": hostResult}
	report.LastNodesJson = TestLastNodesJson

	// A001
	var a001Report a001.A001Report
	var a001HostResult a001.A001ReportHostResult
	a001HostResult.Data = a001.A001ReportHostResultData{
		Cpu: a001.A001ReportCpu{
			CpuCount: "12",
		},
	}
	a001Report.Results = a001.A001ReportHostsResults{"test-host": a001HostResult}

	// A002
	var a002Report a002.A002Report
	var a002HostResult a002.A002ReportHostResult
	a002HostResult.Data = a002.A002ReportHostResultData{
		Version:          "PostgreSQL 10.0.11 on x86_64-pc-linux-gnu (Ubuntu 9.6.11-1.pgdg16.04+1), compiled by gcc (Ubuntu 5.4.0-6ubuntu1~16.04.10) 5.4.0 20160609, 64-bit",
		ServerVersionNum: "10011",
		ServerMajorVer:   "10",
		ServerMinorVer:   "11",
	}
	a002Report.Results = a002.A002ReportHostsResults{"test-host": a002HostResult}

	result, err := F008Process(report, a001Report, a002Report)

	if err != nil && !result.P1 ||
		!checkup.ResultInList(result.Conclusions, F008_DELAY_NOT_TUNED) ||
		!checkup.ResultInList(result.Recommendations, F008_DELAY_NOT_TUNED) ||
		!checkup.ResultInList(result.Recommendations, F008_DELAY_TUNE) {
		t.Fatal()
	}

	checkup.PrintResultConclusions(result)
	checkup.PrintResultRecommendations(result)
}

func TestF008CostNotTuned200(t *testing.T) {
	fmt.Println(t.Name())
	// G001
	var report F008Report
	var hostResult F008ReportHostResult

	hostResult.Data = HostData
	hostResult.Data["autovacuum_vacuum_cost_limit"] = F008Setting{
		Name:    "autovacuum_vacuum_cost_limit",
		Setting: "200",
	}
	hostResult.Data["autovacuum_vacuum_cost_delay"] = F008Setting{
		Name:    "autovacuum_vacuum_cost_delay",
		Setting: "20",
	}
	report.Results = F008ReportHostsResults{"test-host": hostResult}
	report.LastNodesJson = TestLastNodesJson

	// A001
	var a001Report a001.A001Report
	var a001HostResult a001.A001ReportHostResult
	a001HostResult.Data = a001.A001ReportHostResultData{
		Cpu: a001.A001ReportCpu{
			CpuCount: "12",
		},
	}
	a001Report.Results = a001.A001ReportHostsResults{"test-host": a001HostResult}

	// A002
	var a002Report a002.A002Report
	var a002HostResult a002.A002ReportHostResult
	a002HostResult.Data = a002.A002ReportHostResultData{
		Version:          "PostgreSQL 10.0.11 on x86_64-pc-linux-gnu (Ubuntu 9.6.11-1.pgdg16.04+1), compiled by gcc (Ubuntu 5.4.0-6ubuntu1~16.04.10) 5.4.0 20160609, 64-bit",
		ServerVersionNum: "10011",
		ServerMajorVer:   "10",
		ServerMinorVer:   "11",
	}
	a002Report.Results = a002.A002ReportHostsResults{"test-host": a002HostResult}

	result, err := F008Process(report, a001Report, a002Report)

	if err != nil && !result.P1 ||
		!checkup.ResultInList(result.Conclusions, F008_DELAY_NOT_TUNED) ||
		!checkup.ResultInList(result.Recommendations, F008_DELAY_NOT_TUNED) ||
		!checkup.ResultInList(result.Recommendations, F008_DELAY_TUNE) {
		t.Fatal()
	}

	checkup.PrintResultConclusions(result)
	checkup.PrintResultRecommendations(result)
}
