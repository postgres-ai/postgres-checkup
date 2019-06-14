package g002

import (
	"fmt"
	"testing"

	checkup ".."
)

func TestG002Success(t *testing.T) {
	fmt.Println(t.Name())
	// G001
	var report G002Report
	var hostResult G002ReportHostResult
	hostResult.Data = map[string]G002Connection{
		"1": G002Connection{
			CurrentState: "idle",
			TxMore1h:     0,
		},
		"2": G002Connection{
			CurrentState: "idle in transaction",
			TxMore1h:     0,
		},
		"3": G002Connection{
			CurrentState: "active",
			TxMore1h:     0,
		},
	}
	report.Results = G002ReportHostsResults{"test-host": hostResult}

	result, err := G002Process(report)

	if err != nil || result.P1 {
		t.Fatal()
	}

	checkup.PrintResultConclusions(result)
	checkup.PrintResultRecommendations(result)
}

func TestG002IdleInTransaction(t *testing.T) {
	fmt.Println(t.Name())
	// G001
	var report G002Report
	var hostResult G002ReportHostResult
	hostResult.Data = map[string]G002Connection{
		"1": G002Connection{
			CurrentState: "idle",
			TxMore1h:     0,
		},
		"2": G002Connection{
			CurrentState: "idle in transaction",
			TxMore1h:     4,
		},
	}
	report.Results = G002ReportHostsResults{"test-host": hostResult}

	result, err := G002Process(report)

	if err != nil || !result.P1 || !checkup.ResultInList(result.Conclusions, G002_TX_AGE_MORE_1H) {
		t.Fatal()
	}

	checkup.PrintResultConclusions(result)
	checkup.PrintResultRecommendations(result)
}

func TestG002ActiveTransaction(t *testing.T) {
	fmt.Println(t.Name())
	// G001
	var report G002Report
	var hostResult G002ReportHostResult
	hostResult.Data = map[string]G002Connection{
		"1": G002Connection{
			CurrentState: "idle",
			TxMore1h:     0,
		},
		"2": G002Connection{
			CurrentState: "active",
			TxMore1h:     2,
		},
	}
	report.Results = G002ReportHostsResults{"test-host": hostResult}

	result, err := G002Process(report)

	if err != nil || !result.P1 || !checkup.ResultInList(result.Conclusions, G002_TX_AGE_MORE_1H) {
		t.Fatal()
	}

	checkup.PrintResultConclusions(result)
	checkup.PrintResultRecommendations(result)
}

func TestG002ActiveTransactionNodes(t *testing.T) {
	fmt.Println(t.Name())
	// G001
	var report G002Report
	var hostResult G002ReportHostResult
	hostResult.Data = map[string]G002Connection{
		"1": G002Connection{
			CurrentState: "idle",
			TxMore1h:     0,
		},
		"2": G002Connection{
			CurrentState: "active",
			TxMore1h:     2,
		},
	}
	report.Results = G002ReportHostsResults{"test-host-1": hostResult, "test-host-2": hostResult}

	result, err := G002Process(report)

	if err != nil || !result.P1 || !checkup.ResultInList(result.Conclusions, G002_TX_AGE_MORE_1H) {
		t.Fatal()
	}

	checkup.PrintResultConclusions(result)
	checkup.PrintResultRecommendations(result)
}
