package f002

import (
	"fmt"
	"testing"

	"gitlab.com/postgres-ai/postgres-checkup/pghrep/internal/checkup"
)

func TestF002Success(t *testing.T) {
	fmt.Println(t.Name())
	var report F002Report
	var hostResult F002ReportHostResult
	hostResult.Data.Databases = map[string]F002Database{
		"database_1": F002Database{
			Num:          1,
			DatabaseName: "database_1",
			Age:          1,
			CapacityUsed: 34,
			Datfrozenxid: "",
			Warning:      0,
		},
		"database_2": F002Database{
			Num:          2,
			DatabaseName: "database_2",
			Age:          1,
			CapacityUsed: 37,
			Datfrozenxid: "",
			Warning:      0,
		},
	}
	report.Results = F002ReportHostsResults{"test-host": hostResult}
	result := F002Process(report)
	if result.P1 ||
		result.P2 ||
		result.P3 ||
		len(result.Conclusions) > 0 ||
		len(result.Recommendations) > 0 {
		t.Fatal("TestF002Success failed")
	}
	checkup.PrintResultConclusions(result)
	checkup.PrintResultRecommendations(result)
}

func TestF002CheckDatabases(t *testing.T) {
	fmt.Println(t.Name())
	var report F002Report
	var hostResult F002ReportHostResult
	hostResult.Data.Databases = map[string]F002Database{
		"database_1": F002Database{
			Num:          1,
			DatabaseName: "database_1",
			Age:          1,
			CapacityUsed: 76,
			Datfrozenxid: "",
			Warning:      0,
		},
		"database_2": F002Database{
			Num:          2,
			DatabaseName: "database_2",
			Age:          1,
			CapacityUsed: 37,
			Datfrozenxid: "",
			Warning:      0,
		},
	}
	report.Results = F002ReportHostsResults{"test-host": hostResult}
	result := F002Process(report)
	if !result.P1 ||
		!checkup.ResultInList(result.Conclusions, F002_RISKS_ARE_HIGH) ||
		!checkup.ResultInList(result.Recommendations, F002_RISKS_ARE_HIGH) {
		t.Fatal("TestF002CheckDatabases failed")
	}
	checkup.PrintResultConclusions(result)
	checkup.PrintResultRecommendations(result)
}

func TestF002CheckTables(t *testing.T) {
	fmt.Println(t.Name())
	var report F002Report
	var hostResult F002ReportHostResult
	hostResult.Data.Databases = map[string]F002Database{
		"database_1": F002Database{
			Num:          1,
			DatabaseName: "database_1",
			Age:          1,
			CapacityUsed: 36,
			Datfrozenxid: "",
			Warning:      0,
		},
		"database_2": F002Database{
			Num:          2,
			DatabaseName: "database_2",
			Age:          1,
			CapacityUsed: 37,
			Datfrozenxid: "",
			Warning:      0,
		},
	}
	hostResult.Data.Tables = map[string]F002Table{
		"table_1": F002Table{
			Num:                1,
			Relation:           "table_1",
			Age:                0,
			CapacityUsed:       43,
			RelRelfrozenxid:    "",
			ToastRelfrozenxid:  "",
			Warning:            0,
			OverriddenSettings: false,
		},
		"table_2": F002Table{
			Num:                1,
			Relation:           "table_2",
			Age:                0,
			CapacityUsed:       76,
			RelRelfrozenxid:    "",
			ToastRelfrozenxid:  "",
			Warning:            0,
			OverriddenSettings: false,
		},
	}
	report.Results = F002ReportHostsResults{"test-host": hostResult}
	result := F002Process(report)
	if !result.P1 ||
		!checkup.ResultInList(result.Conclusions, F002_RISKS_ARE_HIGH) ||
		!checkup.ResultInList(result.Recommendations, F002_RISKS_ARE_HIGH) {
		t.Fatal("TestF002CheckTables failed")
	}
	checkup.PrintResultConclusions(result)
	checkup.PrintResultRecommendations(result)
}

func TestF002CheckDatabaseTables(t *testing.T) {
	fmt.Println(t.Name())
	var report F002Report
	var hostResult F002ReportHostResult
	hostResult.Data.Databases = map[string]F002Database{
		"database_1": F002Database{
			Num:          1,
			DatabaseName: "database_1",
			Age:          1,
			CapacityUsed: 56,
			Datfrozenxid: "",
			Warning:      0,
		},
		"database_2": F002Database{
			Num:          2,
			DatabaseName: "database_2",
			Age:          1,
			CapacityUsed: 37,
			Datfrozenxid: "",
			Warning:      0,
		},
	}
	hostResult.Data.Tables = map[string]F002Table{
		"table_1": F002Table{
			Num:                1,
			Relation:           "table_1",
			Age:                0,
			CapacityUsed:       43,
			RelRelfrozenxid:    "",
			ToastRelfrozenxid:  "",
			Warning:            0,
			OverriddenSettings: false,
		},
		"table_2": F002Table{
			Num:                1,
			Relation:           "table_2",
			Age:                0,
			CapacityUsed:       76,
			RelRelfrozenxid:    "",
			ToastRelfrozenxid:  "",
			Warning:            0,
			OverriddenSettings: false,
		},
	}

	report.Results = F002ReportHostsResults{"test-host": hostResult}
	result := F002Process(report)
	if !result.P1 ||
		!checkup.ResultInList(result.Conclusions, F002_RISKS_ARE_HIGH) ||
		!checkup.ResultInList(result.Recommendations, F002_RISKS_ARE_HIGH) {
		t.Fatal("TestF002CheckDatabaseTables failed")
	}
	checkup.PrintResultConclusions(result)
	checkup.PrintResultRecommendations(result)
}
