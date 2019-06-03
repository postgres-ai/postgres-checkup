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
	if result.P1 ||
		result.P2 ||
		result.P3 ||
		!checkup.InList(result.Conclusions, "The total table (heap) bloat estimate is quite low, just ~5.37% (~23.78 GiB). Hooray! Keep watching it though.") ||
		len(result.Recommendations) != 0 {
		t.Fatal("TestF004Success failed")
	}
	checkup.PrintConclusions(result)
	checkup.PrintRecommendations(result)
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
		BloatRatioAvg:        1.0389965180727816,
	}
	hostResult.Data.HeapBloat = map[string]F004HeapBloat{}
	report.Results = F004ReportHostsResults{"test-host": hostResult}
	result := F004Process(report)
	if !result.P1 ||
		!checkup.InList(result.Conclusions, "[P1] Total table (heap) bloat estimation is 23.78 GiB, it is 25.37% of overall tables size and 25.00% of the DB size. So removing the table bloat can help to reduce the total database size to ~71.32 GiB and to increase the free disk space by 23.78 GiB. Notice that this is only an estimation, sometimes it may be significantly off. Total size of tables is 1.04 times bigger than it could be.") {
		t.Fatal("TestF004TotalExcess failed")
	}
	checkup.PrintConclusions(result)
	checkup.PrintRecommendations(result)
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
	report.Results = F004ReportHostsResults{"test-host": hostResult}
	result := F004Process(report)
	if !result.P2 ||
		!checkup.InListPartial(result.Conclusions, "[P2] There are 6 tables with size > 1 MiB and table bloat estimate >= 40% and < 90%") ||
		!checkup.InListPartial(result.Recommendations, "[P2] To resolve the table bloat issue do both of the following action items") {
		t.Fatal("TestF004SWarnig failed")
	}
	checkup.PrintConclusions(result)
	checkup.PrintRecommendations(result)
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
		BloatRatioAvg:        1.0389965180727816,
	}
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
	if !result.P1 ||
		!checkup.InListPartial(result.Conclusions, "[P1] The following 1 tables have significant size (>1 MiB) and bloat estimate > 90%:") ||
		!checkup.InListPartial(result.Recommendations, "[P1] Reduce and prevent high level of table bloat") {
		t.Fatal("TestF004SWarnig failed")
	}
	checkup.PrintConclusions(result)
	checkup.PrintRecommendations(result)
}
