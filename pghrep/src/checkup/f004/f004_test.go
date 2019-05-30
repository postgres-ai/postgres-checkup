package f004

import (
	"fmt"
	"testing"

	checkup ".."
)

func TestF004Success(t *testing.T) {
	fmt.Println(t.Name())
	var report F004Report
	var hostResult F004ReportHostResult
	hostResult.Data.HeapBloatTotal = F004HeapBloatTotal{
		Count:                104,
		ExtraSizeBytesSum:    25526329344,
		RealSizeBytesSum:     25526329344,
		BloatSizeBytesSum:    25526329344,
		LiveDataSizeBytesSum: 457681690624,
		BloatRatioPercentAvg: 5.367978121989509,
		BloatRatioAvg:        1.0389965180727816,
	}
	hostResult.Data.DatabaseSizeBytes = 25526329344 * 4
	hostResult.Data.HeapBloat = map[string]F004HeapBloat{
		"table_1": F004HeapBloat{
			Num:               1,
			IsNa:              "",
			TableName:         "table_1",
			RealSize:          "",
			ExtraSizeBytes:    0,
			ExtraRatioPercent: 0.0,
			Extra:             "",
			BloatSizeBytes:    0,
			BloatRatioPercent: 0.0,
			BloatEstimate:     "",
			RealSizeBytes:     0,
			LiveDataSize:      "",
			LiveDataSizeBytes: 0,
			LastVaccuum:       "",
			Fillfactor:        100.0,
			OverridedSettings: false,
			BloatRatio:        1.0,
		},
		"table_2": F004HeapBloat{
			Num:               2,
			IsNa:              "",
			TableName:         "table_2",
			RealSize:          "",
			ExtraSizeBytes:    0,
			ExtraRatioPercent: 0.0,
			Extra:             "",
			BloatSizeBytes:    0,
			BloatRatioPercent: 0.0,
			BloatEstimate:     "",
			RealSizeBytes:     0,
			LiveDataSize:      "",
			LiveDataSizeBytes: 0,
			LastVaccuum:       "",
			Fillfactor:        100.0,
			OverridedSettings: false,
			BloatRatio:        1.0,
		},
	}

	report.Results = F004ReportHostsResults{"test-host": hostResult}
	result := F004Process(report)
	if result.P1 || result.P2 || result.P3 {
		t.Fatal("TestF004Success failed")
	}
	checkup.PrintConclusions(result)
	checkup.PrintReccomendations(result)
}

func TestF004TotalExcess(t *testing.T) {
	fmt.Println(t.Name())
	var report F004Report
	var hostResult F004ReportHostResult
	hostResult.Data.HeapBloatTotal = F004HeapBloatTotal{
		Count:                104,
		ExtraSizeBytesSum:    25526329344,
		RealSizeBytesSum:     25526329344,
		BloatSizeBytesSum:    25526329344,
		LiveDataSizeBytesSum: 457681690624,
		BloatRatioPercentAvg: 25.367978121989509,
		BloatRatioAvg:        1.0389965180727816,
	}
	hostResult.Data.DatabaseSizeBytes = 25526329344 * 4
	hostResult.Data.HeapBloat = map[string]F004HeapBloat{}
	report.Results = F004ReportHostsResults{"test-host": hostResult}
	result := F004Process(report)
	if !result.P1 {
		t.Fatal("TestF004TotalExcess failed")
	}
	checkup.PrintConclusions(result)
	checkup.PrintReccomendations(result)
}

func TestF004Warnig(t *testing.T) {
	fmt.Println(t.Name())
	var report F004Report
	var hostResult F004ReportHostResult
	hostResult.Data.HeapBloatTotal = F004HeapBloatTotal{
		Count:                104,
		ExtraSizeBytesSum:    25526329344,
		RealSizeBytesSum:     25526329344,
		BloatSizeBytesSum:    25526329344,
		LiveDataSizeBytesSum: 457681690624,
		BloatRatioPercentAvg: 5.367978121989509,
		BloatRatioAvg:        1.0389965180727816,
	}

	hostResult.Data.HeapBloat = map[string]F004HeapBloat{
		"table_1": F004HeapBloat{
			Num:               1,
			IsNa:              "",
			TableName:         "table_1",
			RealSize:          "",
			ExtraSizeBytes:    340475904,
			ExtraRatioPercent: 86.9861866889912,
			Extra:             "",
			BloatSizeBytes:    340475904,
			BloatRatioPercent: 86.9861866889912,
			BloatEstimate:     "",
			RealSizeBytes:     391413760,
			LiveDataSize:      "",
			LiveDataSizeBytes: 50937856,
			LastVaccuum:       "",
			Fillfactor:        100.0,
			OverridedSettings: false,
			BloatRatio:        7.684142811193309,
		},
		"table_2": F004HeapBloat{
			Num:               2,
			IsNa:              "",
			TableName:         "table_2",
			RealSize:          "",
			ExtraSizeBytes:    3915776,
			ExtraRatioPercent: 59.3788819875776,
			Extra:             "",
			BloatSizeBytes:    3915776,
			BloatRatioPercent: 59.3788819875776,
			BloatEstimate:     "",
			RealSizeBytes:     6594560,
			LiveDataSize:      "",
			LiveDataSizeBytes: 2678784,
			LastVaccuum:       "",
			Fillfactor:        100.0,
			OverridedSettings: false,
			BloatRatio:        2.46177370030581,
		},
		"table_3": F004HeapBloat{
			Num:               2,
			IsNa:              "",
			TableName:         "table_3",
			RealSize:          "",
			ExtraSizeBytes:    3915776,
			ExtraRatioPercent: 59.3788819875776,
			Extra:             "",
			BloatSizeBytes:    3915776,
			BloatRatioPercent: 59.3788819875776,
			BloatEstimate:     "",
			RealSizeBytes:     6594560,
			LiveDataSize:      "",
			LiveDataSizeBytes: 2678784,
			LastVaccuum:       "",
			Fillfactor:        100.0,
			OverridedSettings: false,
			BloatRatio:        2.46177370030581,
		},
		"table_4": F004HeapBloat{
			Num:               2,
			IsNa:              "",
			TableName:         "table_4",
			RealSize:          "",
			ExtraSizeBytes:    3915776,
			ExtraRatioPercent: 59.3788819875776,
			Extra:             "",
			BloatSizeBytes:    3915776,
			BloatRatioPercent: 59.3788819875776,
			BloatEstimate:     "",
			RealSizeBytes:     6594560,
			LiveDataSize:      "",
			LiveDataSizeBytes: 2678784,
			LastVaccuum:       "",
			Fillfactor:        100.0,
			OverridedSettings: false,
			BloatRatio:        2.46177370030581,
		},
		"table_5": F004HeapBloat{
			Num:               2,
			IsNa:              "",
			TableName:         "table_5",
			RealSize:          "",
			ExtraSizeBytes:    3915776,
			ExtraRatioPercent: 59.3788819875776,
			Extra:             "",
			BloatSizeBytes:    3915776,
			BloatRatioPercent: 59.3788819875776,
			BloatEstimate:     "",
			RealSizeBytes:     6594560,
			LiveDataSize:      "",
			LiveDataSizeBytes: 2678784,
			LastVaccuum:       "",
			Fillfactor:        100.0,
			OverridedSettings: false,
			BloatRatio:        2.46177370030581,
		},
		"table_6": F004HeapBloat{
			Num:               2,
			IsNa:              "",
			TableName:         "table_6",
			RealSize:          "",
			ExtraSizeBytes:    3915776,
			ExtraRatioPercent: 59.3788819875776,
			Extra:             "",
			BloatSizeBytes:    3915776,
			BloatRatioPercent: 59.3788819875776,
			BloatEstimate:     "",
			RealSizeBytes:     6594560,
			LiveDataSize:      "",
			LiveDataSizeBytes: 2678784,
			LastVaccuum:       "",
			Fillfactor:        100.0,
			OverridedSettings: false,
			BloatRatio:        2.46177370030581,
		},
	}
	hostResult.Data.DatabaseSizeBytes = 25526329344 * 4
	report.Results = F004ReportHostsResults{"test-host": hostResult}
	result := F004Process(report)
	if !result.P2 {
		t.Fatal("TestF004SWarnig failed")
	}
	checkup.PrintConclusions(result)
	checkup.PrintReccomendations(result)
}

func TestF004Critical(t *testing.T) {
	fmt.Println(t.Name())
	var report F004Report
	var hostResult F004ReportHostResult
	hostResult.Data.HeapBloatTotal = F004HeapBloatTotal{
		Count:                104,
		ExtraSizeBytesSum:    25526329344,
		RealSizeBytesSum:     25526329344,
		BloatSizeBytesSum:    25526329344,
		LiveDataSizeBytesSum: 457681690624,
		BloatRatioPercentAvg: 5.367978121989509,
		BloatRatioAvg:        1.0389965180727816,
	}
	hostResult.Data.DatabaseSizeBytes = 25526329344 * 4
	hostResult.Data.HeapBloat = map[string]F004HeapBloat{
		"table_1": F004HeapBloat{
			Num:               1,
			IsNa:              "",
			TableName:         "table_1",
			RealSize:          "",
			ExtraSizeBytes:    340475904,
			ExtraRatioPercent: 96.9861866889912,
			Extra:             "",
			BloatSizeBytes:    340475904,
			BloatRatioPercent: 96.9861866889912,
			BloatEstimate:     "",
			RealSizeBytes:     391413760,
			LiveDataSize:      "",
			LiveDataSizeBytes: 50937856,
			LastVaccuum:       "",
			Fillfactor:        100.0,
			OverridedSettings: false,
			BloatRatio:        7.684142811193309,
		},
		"table_2": F004HeapBloat{
			Num:               2,
			IsNa:              "",
			TableName:         "table_2",
			RealSize:          "",
			ExtraSizeBytes:    3915776,
			ExtraRatioPercent: 59.3788819875776,
			Extra:             "",
			BloatSizeBytes:    3915776,
			BloatRatioPercent: 59.3788819875776,
			BloatEstimate:     "",
			RealSizeBytes:     6594560,
			LiveDataSize:      "",
			LiveDataSizeBytes: 2678784,
			LastVaccuum:       "",
			Fillfactor:        100.0,
			OverridedSettings: false,
			BloatRatio:        2.46177370030581,
		},
	}

	report.Results = F004ReportHostsResults{"test-host": hostResult}
	result := F004Process(report)
	if !result.P1 {
		t.Fatal("TestF004SWarnig failed")
	}
	checkup.PrintConclusions(result)
	checkup.PrintReccomendations(result)
}
