package f002

import (
	"fmt"
	"testing"

	checkup ".."
)

func printConclusions(result checkup.ReportOutcome) {
	for _, conclusion := range result.Conclusions {
		fmt.Println("C:  ", conclusion)
	}
}

func printReccomendations(result checkup.ReportOutcome) {
	for _, recommendation := range result.Recommendations {
		fmt.Println("R:  ", recommendation)
	}
}

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
	if result.P1 || result.P2 || result.P3 {
		t.Fatal("TestF002Success failed")
	}
	printConclusions(result)
	printReccomendations(result)
}

func TestF002ChecDatabases(t *testing.T) {
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
	if !result.P1 {
		t.Fatal("TestF002Sucess failed")
	}
	printConclusions(result)
	printReccomendations(result)
}

func TestF002ChecTables(t *testing.T) {
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
			Num:               1,
			Relation:          "table_1",
			Age:               0,
			CapacityUsed:      43,
			RelRelfrozenxid:   "",
			ToastRelfrozenxid: "",
			Warning:           0,
			OverridedSettings: false,
		},
		"table_2": F002Table{
			Num:               1,
			Relation:          "table_2",
			Age:               0,
			CapacityUsed:      76,
			RelRelfrozenxid:   "",
			ToastRelfrozenxid: "",
			Warning:           0,
			OverridedSettings: false,
		},
	}

	report.Results = F002ReportHostsResults{"test-host": hostResult}
	result := F002Process(report)
	if !result.P1 {
		t.Fatal("TestF002Sucess failed")
	}
	printConclusions(result)
	printReccomendations(result)
}

func TestF002ChecDatabaseTables(t *testing.T) {
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
			Num:               1,
			Relation:          "table_1",
			Age:               0,
			CapacityUsed:      43,
			RelRelfrozenxid:   "",
			ToastRelfrozenxid: "",
			Warning:           0,
			OverridedSettings: false,
		},
		"table_2": F002Table{
			Num:               1,
			Relation:          "table_2",
			Age:               0,
			CapacityUsed:      76,
			RelRelfrozenxid:   "",
			ToastRelfrozenxid: "",
			Warning:           0,
			OverridedSettings: false,
		},
	}

	report.Results = F002ReportHostsResults{"test-host": hostResult}
	result := F002Process(report)
	if !result.P1 {
		t.Fatal("TestF002Sucess failed")
	}
	printConclusions(result)
	printReccomendations(result)
}
