package f005

import (
	"fmt"
	"testing"

	checkup ".."
)

func TestF005Success(t *testing.T) {
	fmt.Println(t.Name())
	var report F005Report
	var hostResult F005ReportHostResult
	hostResult.Data.DatabaseSizeBytes = 102105317376
	hostResult.Data.IndexBloatTotal = F005IndexBloatTotal{
		Count:                104,
		ExtraSizeBytesSum:    hostResult.Data.DatabaseSizeBytes / 4,
		RealSizeBytesSum:     hostResult.Data.DatabaseSizeBytes / 4,
		BloatSizeBytesSum:    hostResult.Data.DatabaseSizeBytes / 4,
		LiveDataSizeBytesSum: 457681690624,
		BloatRatioPercentAvg: 5.367978121989509,
		BloatRatioAvg:        1.0389965180727816,
	}
	hostResult.Data.IndexBloat = map[string]F005IndexBloat{
		"index_1": F005IndexBloat{
			Num:               1,
			IsNa:              "",
			IndexName:         "index_1",
			TableName:         "table_1",
			ExtraSizeBytes:    0,
			ExtraRatioPercent: 0.0,
			Extra:             "",
			BloatSizeBytes:    0,
			BloatRatioPercent: 0.0,
			RealSizeBytes:     0,
			LiveDataSize:      "",
			LiveDataSizeBytes: 0,
			LastVaccuum:       "",
			Fillfactor:        100.0,
			OverridedSettings: false,
			BloatRatio:        1.0,
		},
		"index_2": F005IndexBloat{
			Num:               2,
			IsNa:              "",
			IndexName:         "index_2",
			TableName:         "table_2",
			ExtraSizeBytes:    0,
			ExtraRatioPercent: 0.0,
			Extra:             "",
			BloatSizeBytes:    0,
			BloatRatioPercent: 0.0,
			RealSizeBytes:     0,
			LiveDataSize:      "",
			LiveDataSizeBytes: 0,
			LastVaccuum:       "",
			Fillfactor:        100.0,
			OverridedSettings: false,
			BloatRatio:        1.0,
		},
	}
	report.Results = F005ReportHostsResults{"test-host": hostResult}
	result := F005Process(report)
	if result.P1 || result.P2 || result.P3 &&
		checkup.InList(result.Conclusions, "The total index bloat estimate is quite low, just ~5.37% (~23.78 GiB). Hooray! Keep watching it though.") {
		t.Fatal("TestF005Success failed")
	}
	checkup.PrintConclusions(result)
	checkup.PrintRecommendations(result)
}

func TestF005TotalExcess(t *testing.T) {
	fmt.Println(t.Name())
	var report F005Report
	var hostResult F005ReportHostResult
	hostResult.Data.DatabaseSizeBytes = 102105317376
	hostResult.Data.IndexBloatTotal = F005IndexBloatTotal{
		Count:                104,
		ExtraSizeBytesSum:    hostResult.Data.DatabaseSizeBytes / 4,
		RealSizeBytesSum:     hostResult.Data.DatabaseSizeBytes / 4,
		BloatSizeBytesSum:    hostResult.Data.DatabaseSizeBytes / 4,
		LiveDataSizeBytesSum: hostResult.Data.DatabaseSizeBytes / 4,
		BloatRatioPercentAvg: 25.367978121989509,
		BloatRatioAvg:        1.0389965180727816,
		TableSizeBytesSum:    25526329344 * 5,
	}
	hostResult.Data.IndexBloat = map[string]F005IndexBloat{}
	report.Results = F005ReportHostsResults{"test-host": hostResult}
	result := F005Process(report)
	if !result.P1 ||
		!checkup.InList(result.Conclusions, "[P1] Total index bloat estimation is 23.78 GiB, it is 25.37% of overall indexes size and 25.00% of the DB size. So removing the index bloat can help to reduce the total database size to ~71.32 GiB and to increase the free disk space by 23.78 GiB. Notice that this is only an estimation, sometimes it may be significantly off. Total size of indexes is 1.04 times bigger than it could be.") ||
		!checkup.InListPartial(result.Recommendations, "[P1] Reduce and prevent high level of index bloat:") {
		t.Fatal("TestF005TotalExcess failed")
	}
	checkup.PrintConclusions(result)
	checkup.PrintRecommendations(result)
}

func TestF005Warnig(t *testing.T) {
	fmt.Println(t.Name())
	var report F005Report
	var hostResult F005ReportHostResult
	hostResult.Data.DatabaseSizeBytes = 102105317376
	hostResult.Data.IndexBloatTotal = F005IndexBloatTotal{
		Count:                104,
		ExtraSizeBytesSum:    hostResult.Data.DatabaseSizeBytes / 4,
		RealSizeBytesSum:     hostResult.Data.DatabaseSizeBytes / 4,
		BloatSizeBytesSum:    hostResult.Data.DatabaseSizeBytes / 4,
		LiveDataSizeBytesSum: 457681690624,
		BloatRatioPercentAvg: 5.367978121989509,
		BloatRatioAvg:        1.0389965180727816,
	}
	hostResult.Data.IndexBloat = map[string]F005IndexBloat{
		"index_1": F005IndexBloat{
			Num:               1,
			IsNa:              "",
			IndexName:         "index_1",
			TableName:         "table_1",
			ExtraSizeBytes:    340475904,
			ExtraRatioPercent: 86.9861866889912,
			Extra:             "",
			BloatSizeBytes:    340475904,
			BloatRatioPercent: 86.9861866889912,
			RealSizeBytes:     391413760,
			LiveDataSize:      "",
			LiveDataSizeBytes: 50937856,
			LastVaccuum:       "",
			Fillfactor:        100.0,
			OverridedSettings: false,
			BloatRatio:        7.684142811193309,
		},
		"index_2": F005IndexBloat{
			Num:               2,
			IsNa:              "",
			IndexName:         "index_2",
			TableName:         "table_2",
			ExtraSizeBytes:    3915776,
			ExtraRatioPercent: 59.3788819875776,
			Extra:             "",
			BloatSizeBytes:    3915776,
			BloatRatioPercent: 59.3788819875776,
			RealSizeBytes:     6594560,
			LiveDataSize:      "",
			LiveDataSizeBytes: 2678784,
			LastVaccuum:       "",
			Fillfactor:        100.0,
			OverridedSettings: false,
			BloatRatio:        2.46177370030581,
		},
	}

	report.Results = F005ReportHostsResults{"test-host": hostResult}
	result := F005Process(report)
	if !result.P2 ||
		!checkup.InListPartial(result.Conclusions, "[P2] There are 2 indexes with size > 1 MiB and index bloat estimate >= 40% and < 90%") ||
		!checkup.InListPartial(result.Recommendations, "[P2] Reduce and prevent high level of index bloat") {
		t.Fatal("TestF005SWarnig failed")
	}
	checkup.PrintConclusions(result)
	checkup.PrintRecommendations(result)
}

func TestF005Critical(t *testing.T) {
	fmt.Println(t.Name())
	var report F005Report
	var hostResult F005ReportHostResult
	hostResult.Data.DatabaseSizeBytes = 102105317376
	hostResult.Data.IndexBloatTotal = F005IndexBloatTotal{
		Count:                104,
		ExtraSizeBytesSum:    hostResult.Data.DatabaseSizeBytes / 4,
		RealSizeBytesSum:     hostResult.Data.DatabaseSizeBytes / 4,
		BloatSizeBytesSum:    hostResult.Data.DatabaseSizeBytes / 4,
		LiveDataSizeBytesSum: 457681690624,
		BloatRatioPercentAvg: 5.367978121989509,
		BloatRatioAvg:        1.0389965180727816,
	}
	hostResult.Data.IndexBloat = map[string]F005IndexBloat{
		"index_1": F005IndexBloat{
			Num:               1,
			IsNa:              "",
			IndexName:         "index_1",
			TableName:         "table_1",
			ExtraSizeBytes:    340475904,
			ExtraRatioPercent: 96.9861866889912,
			Extra:             "",
			BloatSizeBytes:    340475904,
			BloatRatioPercent: 96.9861866889912,
			RealSizeBytes:     391413760,
			LiveDataSize:      "",
			LiveDataSizeBytes: 50937856,
			LastVaccuum:       "",
			Fillfactor:        100.0,
			OverridedSettings: false,
			BloatRatio:        7.684142811193309,
		},
		"index_2": F005IndexBloat{
			Num:               2,
			IsNa:              "",
			IndexName:         "index_2",
			TableName:         "table_2",
			ExtraSizeBytes:    3915776,
			ExtraRatioPercent: 59.3788819875776,
			Extra:             "",
			BloatSizeBytes:    3915776,
			BloatRatioPercent: 59.3788819875776,
			RealSizeBytes:     6594560,
			LiveDataSize:      "",
			LiveDataSizeBytes: 2678784,
			LastVaccuum:       "",
			Fillfactor:        100.0,
			OverridedSettings: false,
			BloatRatio:        2.46177370030581,
		},
	}

	report.Results = F005ReportHostsResults{"test-host": hostResult}
	result := F005Process(report)
	if !result.P1 ||
		!checkup.InListPartial(result.Conclusions, "[P1] The following 1 indexes have significant size (>1 MiB) and bloat estimate > 90%") ||
		!checkup.InListPartial(result.Recommendations, "[P1] Reduce and prevent high level of index bloat:") {
		t.Fatal("TestF005SWarnig failed")
	}
	checkup.PrintConclusions(result)
	checkup.PrintRecommendations(result)
}
