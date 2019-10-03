package l003

import (
	"fmt"
	"testing"

	checkup ".."
)

func TestL003Success(t *testing.T) {
	fmt.Println(t.Name())

	var report L003Report
	var hostResult L003ReportHostResult = L003ReportHostResult{
		Data: L003ReportHostResultData{
			Tables: map[string]L003Table{
				"test_schema.orders": L003Table{
					Table:               "test_schema.orders",
					Pk:                  "id",
					Type:                "int4",
					CurrentMaxValue:     80000000,
					CapacityUsedPercent: 3.725,
				},
			},
			MinTableSizeBytes: 0,
		},
	}

	report.Results = L003ReportHostsResults{
		"test-host": hostResult,
	}

	result, err := L003Process(report)
	if err != nil || result.P1 ||
		checkup.ResultInList(result.Conclusions, L003_HIGH_RISKS) ||
		checkup.ResultInList(result.Recommendations, L003_HIGH_RISKS) {
		t.Fatal()
	}

	checkup.PrintResultConclusions(result)
	checkup.PrintResultRecommendations(result)
}

func TestL003P1_1(t *testing.T) {
	fmt.Println(t.Name())

	var report L003Report
	var hostResult L003ReportHostResult = L003ReportHostResult{
		Data: L003ReportHostResultData{
			Tables: map[string]L003Table{
				"test_schema.orders": L003Table{
					Table:               "test_schema.orders",
					Pk:                  "id",
					Type:                "int4",
					CurrentMaxValue:     80000000,
					CapacityUsedPercent: 37.25,
				},
			},
			MinTableSizeBytes: 0,
		},
	}

	report.Results = L003ReportHostsResults{
		"test-host": hostResult,
	}

	result, err := L003Process(report)
	if err != nil || !result.P1 ||
		!checkup.ResultInList(result.Conclusions, L003_HIGH_RISKS) ||
		!checkup.ResultInList(result.Recommendations, L003_HIGH_RISKS) {
		t.Fatal()
	}

	checkup.PrintResultConclusions(result)
	checkup.PrintResultRecommendations(result)
}

func TestL003P1_N(t *testing.T) {
	fmt.Println(t.Name())

	var report L003Report
	var hostResult L003ReportHostResult = L003ReportHostResult{
		Data: L003ReportHostResultData{
			Tables: map[string]L003Table{
				"test_schema.orders_A": L003Table{
					Table:               "test_schema.orders_A",
					Pk:                  "id",
					Type:                "int4",
					CurrentMaxValue:     300000000,
					CapacityUsedPercent: 13.97,
				},
				"test_schema.ordersB": L003Table{
					Table:               "test_schema.orders_B",
					Pk:                  "id",
					Type:                "int4",
					CurrentMaxValue:     80000000,
					CapacityUsedPercent: 37.25,
				},
			},
			MinTableSizeBytes: 0,
		},
	}

	report.Results = L003ReportHostsResults{
		"test-host": hostResult,
	}

	result, err := L003Process(report)
	if err != nil || !result.P1 ||
		!checkup.ResultInList(result.Conclusions, L003_HIGH_RISKS) ||
		!checkup.ResultInList(result.Recommendations, L003_HIGH_RISKS) {
		t.Fatal()
	}

	checkup.PrintResultConclusions(result)
	checkup.PrintResultRecommendations(result)
}

func TestL003P1_6(t *testing.T) {
	fmt.Println(t.Name())

	var report L003Report
	var hostResult L003ReportHostResult = L003ReportHostResult{
		Data: L003ReportHostResultData{
			Tables: map[string]L003Table{
				"test_schema.orders": L003Table{
					Table:               "test_schema.orders",
					Pk:                  "id",
					Type:                "int4",
					CurrentMaxValue:     800000000,
					CapacityUsedPercent: 37.25,
				},
				"test_schema.orders_A": L003Table{
					Table:               "test_schema.orders_A",
					Pk:                  "id",
					Type:                "int4",
					CurrentMaxValue:     300000000,
					CapacityUsedPercent: 33.97,
				},
				"test_schema.orders_B": L003Table{
					Table:               "test_schema.orders_B",
					Pk:                  "id",
					Type:                "int4",
					CurrentMaxValue:     300000000,
					CapacityUsedPercent: 27.97,
				},
				"test_schema.orders_C": L003Table{
					Table:               "test_schema.orders_C",
					Pk:                  "id",
					Type:                "int4",
					CurrentMaxValue:     300000000,
					CapacityUsedPercent: 25.97,
				},
				"test_schema.orders_D": L003Table{
					Table:               "test_schema.orders_D",
					Pk:                  "id",
					Type:                "int4",
					CurrentMaxValue:     300000000,
					CapacityUsedPercent: 23.97,
				},
				"test_schema.orders_E": L003Table{
					Table:               "test_schema.orders_E",
					Pk:                  "id",
					Type:                "int4",
					CurrentMaxValue:     300000000,
					CapacityUsedPercent: 13.97,
				},
			},
			MinTableSizeBytes: 0,
		},
	}

	report.Results = L003ReportHostsResults{
		"test-host": hostResult,
	}

	result, err := L003Process(report)
	if err != nil || !result.P1 ||
		!checkup.ResultInList(result.Conclusions, L003_HIGH_RISKS) ||
		!checkup.ResultInList(result.Recommendations, L003_HIGH_RISKS) {
		t.Fatal()
	}

	checkup.PrintResultConclusions(result)
	checkup.PrintResultRecommendations(result)
}
