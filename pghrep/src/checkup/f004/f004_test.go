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
		BloatRatioFactorAvg:  1.0389965180727816,
	}
	hostResult.Data.DatabaseSizeBytes = 25526329344 * 4
	hostResult.Data.HeapBloat = map[string]F004HeapBloat{
		"table_1": F004HeapBloat{
			Num:               1,
			IsNa:              "",
			TableName:         "table_1",
			ExtraSizeBytes:    0,
			ExtraRatioPercent: 0.0,
			BloatSizeBytes:    0,
			BloatRatioPercent: 0.0,
			RealSizeBytes:     0,
			LiveDataSizeBytes: 0,
			LastVaccuum:       "",
			Fillfactor:        100.0,
			OverridedSettings: false,
			BloatRatioFactor:  1.0,
		},
		"table_2": F004HeapBloat{
			Num:               2,
			IsNa:              "",
			TableName:         "table_2",
			ExtraSizeBytes:    0,
			ExtraRatioPercent: 0.0,
			BloatSizeBytes:    0,
			BloatRatioPercent: 0.0,
			RealSizeBytes:     0,
			LiveDataSizeBytes: 0,
			LastVaccuum:       "",
			Fillfactor:        100.0,
			OverridedSettings: false,
			BloatRatioFactor:  1.0,
		},
	}

	report.Results = F004ReportHostsResults{"test-host": hostResult}
	result := F004Process(report)
	if result.P1 ||
		result.P2 ||
		result.P3 ||
		!checkup.ResultInList(result.Conclusions, F004_TOTAL_BLOAT_LOW) ||
		len(result.Recommendations) != 0 {
		t.Fatal("TestF004Success failed")
	}
	checkup.PrintResultConclusions(result)
	checkup.PrintResultRecommendations(result)
}

func TestF004TotalExcess(t *testing.T) {
	fmt.Println(t.Name())
	var report F004Report
	var hostResult F004ReportHostResult
	hostResult.Data.DatabaseSizeBytes = 102105317376
	hostResult.Data.HeapBloatTotal = F004HeapBloatTotal{
		Count:                104,
		ExtraSizeBytesSum:    25526329344,
		RealSizeBytesSum:     25526329344,
		BloatSizeBytesSum:    25526329344,
		LiveDataSizeBytesSum: 457681690624,
		BloatRatioPercentAvg: 25.367978121989509,
		BloatRatioFactorAvg:  1.0389965180727816,
	}
	hostResult.Data.HeapBloat = map[string]F004HeapBloat{}
	report.Results = F004ReportHostsResults{"test-host": hostResult}
	result := F004Process(report)
	if !result.P1 ||
		!checkup.ResultInList(result.Conclusions, F004_TOTAL_BLOAT_EXCESS) {
		t.Fatal("TestF004TotalExcess failed")
	}
	checkup.PrintResultConclusions(result)
	checkup.PrintResultRecommendations(result)
}

func TestF004Warnig(t *testing.T) {
	fmt.Println(t.Name())
	var report F004Report
	var hostResult F004ReportHostResult
	hostResult.Data.DatabaseSizeBytes = 102105317376
	hostResult.Data.HeapBloatTotal = F004HeapBloatTotal{
		Count:                104,
		ExtraSizeBytesSum:    hostResult.Data.DatabaseSizeBytes / 4,
		RealSizeBytesSum:     hostResult.Data.DatabaseSizeBytes / 4,
		BloatSizeBytesSum:    hostResult.Data.DatabaseSizeBytes / 4,
		LiveDataSizeBytesSum: 457681690624,
		BloatRatioPercentAvg: 5.367978121989509,
		BloatRatioFactorAvg:  1.0389965180727816,
	}
	hostResult.Data.HeapBloat = map[string]F004HeapBloat{
		"table_1": F004HeapBloat{
			Num:               1,
			IsNa:              "",
			TableName:         "table_1",
			ExtraSizeBytes:    340475904,
			ExtraRatioPercent: 86.9861866889912,
			BloatSizeBytes:    340475904,
			BloatRatioPercent: 86.9861866889912,
			RealSizeBytes:     391413760,
			LiveDataSizeBytes: 50937856,
			LastVaccuum:       "",
			Fillfactor:        100.0,
			OverridedSettings: false,
			BloatRatioFactor:  7.684142811193309,
		},
		"table_2": F004HeapBloat{
			Num:               2,
			IsNa:              "",
			TableName:         "table_2",
			ExtraSizeBytes:    3915776,
			ExtraRatioPercent: 59.3788819875776,
			BloatSizeBytes:    3915776,
			BloatRatioPercent: 59.3788819875776,
			RealSizeBytes:     6594560,
			LiveDataSizeBytes: 2678784,
			LastVaccuum:       "",
			Fillfactor:        100.0,
			OverridedSettings: false,
			BloatRatioFactor:  2.46177370030581,
		},
		"table_3": F004HeapBloat{
			Num:               3,
			IsNa:              "",
			TableName:         "table_3",
			ExtraSizeBytes:    3915776,
			ExtraRatioPercent: 59.3788819875776,
			BloatSizeBytes:    3915776,
			BloatRatioPercent: 59.3788819875776,
			RealSizeBytes:     6594560,
			LiveDataSizeBytes: 2678784,
			LastVaccuum:       "",
			Fillfactor:        100.0,
			OverridedSettings: false,
			BloatRatioFactor:  2.46177370030581,
		},
		"table_4": F004HeapBloat{
			Num:               4,
			IsNa:              "",
			TableName:         "table_4",
			ExtraSizeBytes:    3915776,
			ExtraRatioPercent: 59.3788819875776,
			BloatSizeBytes:    3915776,
			BloatRatioPercent: 59.3788819875776,
			RealSizeBytes:     6594560,
			LiveDataSizeBytes: 2678784,
			LastVaccuum:       "",
			Fillfactor:        100.0,
			OverridedSettings: false,
			BloatRatioFactor:  2.46177370030581,
		},
		"table_5": F004HeapBloat{
			Num:               5,
			IsNa:              "",
			TableName:         "table_5",
			ExtraSizeBytes:    3915776,
			ExtraRatioPercent: 59.3788819875776,
			BloatSizeBytes:    3915776,
			BloatRatioPercent: 59.3788819875776,
			RealSizeBytes:     6594560,
			LiveDataSizeBytes: 2678784,
			LastVaccuum:       "",
			Fillfactor:        100.0,
			OverridedSettings: false,
			BloatRatioFactor:  2.46177370030581,
		},
		"table_6": F004HeapBloat{
			Num:               6,
			IsNa:              "",
			TableName:         "table_6",
			ExtraSizeBytes:    3915776,
			ExtraRatioPercent: 59.3788819875776,
			BloatSizeBytes:    3915776,
			BloatRatioPercent: 59.3788819875776,
			RealSizeBytes:     6594560,
			LiveDataSizeBytes: 2678784,
			LastVaccuum:       "",
			Fillfactor:        100.0,
			OverridedSettings: false,
			BloatRatioFactor:  2.46177370030581,
		},
	}
	report.Results = F004ReportHostsResults{"test-host": hostResult}
	result := F004Process(report)
	if !result.P2 ||
		!checkup.ResultInList(result.Conclusions, F004_BLOAT_WARNING) ||
		!checkup.ResultInList(result.Recommendations, F004_BLOAT_WARNING) {
		t.Fatal("TestF004SWarnig failed")
	}
	checkup.PrintResultConclusions(result)
	checkup.PrintResultRecommendations(result)
}

func TestF004Critical(t *testing.T) {
	fmt.Println(t.Name())
	var report F004Report
	var hostResult F004ReportHostResult
	hostResult.Data.DatabaseSizeBytes = 102105317376
	hostResult.Data.HeapBloatTotal = F004HeapBloatTotal{
		Count:                104,
		ExtraSizeBytesSum:    hostResult.Data.DatabaseSizeBytes / 4,
		RealSizeBytesSum:     hostResult.Data.DatabaseSizeBytes / 4,
		BloatSizeBytesSum:    hostResult.Data.DatabaseSizeBytes / 4,
		LiveDataSizeBytesSum: 457681690624,
		BloatRatioPercentAvg: 5.367978121989509,
		BloatRatioFactorAvg:  1.0389965180727816,
	}
	hostResult.Data.HeapBloat = map[string]F004HeapBloat{
		"table_1": F004HeapBloat{
			Num:               1,
			IsNa:              "",
			TableName:         "table_1",
			ExtraSizeBytes:    340475904,
			ExtraRatioPercent: 96.9861866889912,
			BloatSizeBytes:    340475904,
			BloatRatioPercent: 96.9861866889912,
			RealSizeBytes:     391413760,
			LiveDataSizeBytes: 50937856,
			LastVaccuum:       "",
			Fillfactor:        100.0,
			OverridedSettings: false,
			BloatRatioFactor:  7.684142811193309,
		},
		"table_2": F004HeapBloat{
			Num:               2,
			IsNa:              "",
			TableName:         "table_2",
			ExtraSizeBytes:    3915776,
			ExtraRatioPercent: 59.3788819875776,
			BloatSizeBytes:    3915776,
			BloatRatioPercent: 59.3788819875776,
			RealSizeBytes:     6594560,
			LiveDataSizeBytes: 2678784,
			LastVaccuum:       "",
			Fillfactor:        100.0,
			OverridedSettings: false,
			BloatRatioFactor:  2.46177370030581,
		},
	}

	report.Results = F004ReportHostsResults{"test-host": hostResult}
	result := F004Process(report)
	if !result.P1 ||
		!checkup.ResultInList(result.Conclusions, F004_BLOAT_CRITICAL) ||
		!checkup.ResultInList(result.Recommendations, F004_BLOAT_CRITICAL) {
		t.Fatal("TestF004SWarnig failed")
	}
	checkup.PrintResultConclusions(result)
	checkup.PrintResultRecommendations(result)
}
