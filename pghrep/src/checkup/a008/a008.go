package a008

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	checkup ".."
	"../../log"
	"../../pyraconv"
	"github.com/dustin/go-humanize/english"
)

const USAGE_CRITICAL int = 90
const USAGE_WARNING int = 70

var FS_RECOMMENDED []string = []string{
	"ext4",
	"xfs",
	"tmpfs",
}

func isRecommendedFs(fs string) bool {
	for _, validFs := range FS_RECOMMENDED {
		if validFs == fs {
			return true
		}
	}
	return false
}

func checkFsItemUsage(host string, fsItemData FsItem,
	result checkup.ReportOutcome) (bool, bool, checkup.ReportOutcome) {
	usageWarning := false
	usageCritical := false
	usePercent := strings.Replace(fsItemData.UsePercent, "%", "", 1)
	percent, _ := strconv.Atoi(usePercent)
	if percent >= USAGE_WARNING && percent < USAGE_CRITICAL {
		result.AppendConclusion(MSG_USAGE_WARNING_CONCLUSION, fsItemData.MountPoint, host, fsItemData.Used)
		result.AppendRecommendation(MSG_USAGE_WARNING_RECOMMENDATION, fsItemData.MountPoint, host, USAGE_WARNING)
		usageWarning = true
		result.P2 = true
	}
	if percent >= USAGE_CRITICAL {
		result.AppendConclusion(MSG_USAGE_CRITICAL_CONCLUSION, fsItemData.MountPoint, host, fsItemData.Used)
		result.AppendRecommendation(MSG_USAGE_CRITICAL_RECOMMENDATION, fsItemData.MountPoint, host)
		usageCritical = true
		result.P1 = true
	}
	return usageWarning, usageCritical, result
}

// Generate conclusions and recommendatons
func A008Process(report A008Report) checkup.ReportOutcome {
	var result checkup.ReportOutcome
	usageWarning := false
	usageCritical := false
	var nfsConclusions []string
	var notRecConclusions []string
	for host, hostResult := range report.Results {
		var notRecommendedFsDisks []string
		var notRecommendedFsDisksFs []string
		var networkFsDisks []string
		for _, fsItemData := range hostResult.Data.DbData {
			if isRecommendedFs(strings.ToLower(fsItemData.Fstype)) != true {
				notRecommendedFsDisks = append(notRecommendedFsDisks, fsItemData.MountPoint)
				notRecommendedFsDisksFs = append(notRecommendedFsDisksFs, fsItemData.Fstype)
			}
			if strings.ToLower(fsItemData.Fstype[0:3]) == "nfs" {
				networkFsDisks = append(networkFsDisks, fsItemData.MountPoint)
			}
			var warning, critical bool
			warning, critical, result = checkFsItemUsage(host, fsItemData, result)
			usageWarning = usageWarning || warning
			usageCritical = usageCritical || critical
		}
		for _, fsItemData := range hostResult.Data.FsData {
			var warning, critical bool
			warning, critical, result = checkFsItemUsage(host, fsItemData, result)
			usageWarning = usageWarning || warning
			usageCritical = usageCritical || critical
		}
		if len(networkFsDisks) > 0 {
			result.P1 = true
			nfsConclusions = append(nfsConclusions, fmt.Sprintf(english.PluralWord(len(networkFsDisks),
				MSG_NETWORK_FS_CONCLUSION_1, MSG_NETWORK_FS_CONCLUSION_N),
				strings.Join(networkFsDisks, ", "), host))
		}
		if len(notRecommendedFsDisks) > 0 {
			result.P3 = true
			notRecConclusions = append(notRecConclusions, fmt.Sprintf(english.PluralWord(len(notRecommendedFsDisks),
				MSG_NOT_RECOMMENDED_FS_CONCLUSION_1, MSG_NOT_RECOMMENDED_FS_CONCLUSION_N),
				strings.Join(notRecommendedFsDisks, ", "), host, strings.Join(notRecommendedFsDisksFs, ", ")))
		}
	}
	if !usageWarning && !usageCritical && len(result.Recommendations) == 0 {
		result.AppendConclusion(MSG_NO_USAGE_RISKS_CONCLUSION)
	}
	if len(nfsConclusions) > 0 {
		result.Conclusions = append(result.Conclusions, nfsConclusions...)
		result.AppendRecommendation(MSG_NETWORK_FS_RECOMMENDATION)
	}
	if len(notRecConclusions) > 0 {
		result.Conclusions = append(result.Conclusions, notRecConclusions...)
		result.AppendRecommendation(MSG_NOT_RECOMMENDED_FS_RECOMMENDATION)
	}
	if len(result.Recommendations) == 0 {
		result.AppendRecommendation(MSG_NO_FS_RECOMMENDATION)
	}
	return result
}

func A008PreprocessReportData(data map[string]interface{}) {
	filePath := pyraconv.ToString(data["source_path_full"])
	jsonRaw := checkup.LoadRawJsonReport(filePath)
	var report A008Report
	err := json.Unmarshal(jsonRaw, &report)
	if err != nil {
		log.Err("Can't load json report to process")
		return
	}
	result := A008Process(report)
	// update data and file
	checkup.SaveConclusionsRecommendations(data, result)
}
