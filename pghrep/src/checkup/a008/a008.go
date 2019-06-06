package a008

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	checkup ".."
	"github.com/dustin/go-humanize/english"
)

// Case codes
const A008_SPACE_USAGE_WARNING string = "A008_SPACE_USAGE_WARNING"
const A008_SPACE_USAGE_CRITICAL string = "A008_SPACE_USAGE_CRITICAL"
const A008_SPACE_USAGE_NORMAL string = "A008_SPACE_USAGE_NORMAL"
const A008_NFS_DISK string = "A008_NFS_DISK"
const A008_NOT_RECOMMENDED_FS string = "A008_NOT_RECOMMENDED_FS"

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
	result checkup.ReportResult) (bool, bool, checkup.ReportResult) {
	usageWarning := false
	usageCritical := false
	usePercent := strings.Replace(fsItemData.UsePercent, "%", "", 1)
	percent, _ := strconv.Atoi(usePercent)
	if percent >= USAGE_WARNING && percent < USAGE_CRITICAL {
		result.AppendConclusion(A008_SPACE_USAGE_WARNING, MSG_USAGE_WARNING_CONCLUSION, fsItemData.MountPoint, host, fsItemData.Used)
		result.AppendRecommendation(A008_SPACE_USAGE_WARNING, MSG_USAGE_WARNING_RECOMMENDATION, fsItemData.MountPoint, host, 100-USAGE_WARNING)
		usageWarning = true
		result.P2 = true
	}
	if percent >= USAGE_CRITICAL {
		result.AppendConclusion(A008_SPACE_USAGE_CRITICAL, MSG_USAGE_CRITICAL_CONCLUSION, fsItemData.MountPoint, host, fsItemData.Used)
		result.AppendRecommendation(A008_SPACE_USAGE_CRITICAL, MSG_USAGE_CRITICAL_RECOMMENDATION, fsItemData.MountPoint, host)
		usageCritical = true
		result.P1 = true
	}
	return usageWarning, usageCritical, result
}

// Generate conclusions and recommendatons
func A008Process(report A008Report) checkup.ReportResult {
	var result checkup.ReportResult
	usageWarning := false
	usageCritical := false
	var nfsConclusions []checkup.ReportResultItem
	var notRecConclusions []checkup.ReportResultItem
	for host, hostResult := range report.Results {
		var notRecommendedFsDisks []string
		var notRecommendedFsDisksFs []string
		var networkFsDisks []string
		for _, fsItemData := range hostResult.Data.DbData {
			if isRecommendedFs(strings.ToLower(fsItemData.Fstype)) != true {
				notRecommendedFsDisks = append(notRecommendedFsDisks, "`"+fsItemData.MountPoint+"`")
				notRecommendedFsDisksFs = append(notRecommendedFsDisksFs, "`"+fsItemData.Fstype+"`")
			}
			if strings.ToLower(fsItemData.Fstype[0:3]) == "nfs" {
				networkFsDisks = append(networkFsDisks, "`"+fsItemData.MountPoint+"`")
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
			nfsConclusions = append(nfsConclusions, checkup.ReportResultItem{
				Id: A008_NFS_DISK,
				Message: fmt.Sprintf(english.PluralWord(len(networkFsDisks),
					MSG_NETWORK_FS_CONCLUSION_1, MSG_NETWORK_FS_CONCLUSION_N),
					english.WordSeries(networkFsDisks, "and"), host),
			})
		}
		if len(notRecommendedFsDisks) > 0 {
			result.P3 = true
			notRecConclusions = append(notRecConclusions, checkup.ReportResultItem{
				Id: A008_NOT_RECOMMENDED_FS,
				Message: fmt.Sprintf(english.PluralWord(len(notRecommendedFsDisks),
					MSG_NOT_RECOMMENDED_FS_CONCLUSION_1, MSG_NOT_RECOMMENDED_FS_CONCLUSION_N),
					english.WordSeries(notRecommendedFsDisks, "and"), host, english.WordSeries(notRecommendedFsDisksFs, "and"))})
		}
	}
	if !usageWarning && !usageCritical && len(result.Recommendations) == 0 {
		result.AppendConclusion(A008_SPACE_USAGE_NORMAL, MSG_NO_USAGE_RISKS_CONCLUSION)
	}
	if len(nfsConclusions) > 0 {
		result.Conclusions = append(result.Conclusions, nfsConclusions...)
		result.AppendRecommendation(A008_NFS_DISK, MSG_NETWORK_FS_RECOMMENDATION)
	}
	if len(notRecConclusions) > 0 {
		result.Conclusions = append(result.Conclusions, notRecConclusions...)
		result.AppendRecommendation(A008_NOT_RECOMMENDED_FS, MSG_NOT_RECOMMENDED_FS_RECOMMENDATION)
	}
	return result
}

func A008PreprocessReportData(data map[string]interface{}) {
	var report A008Report
	filePath := data["source_path_full"].(string)
	jsonRaw := checkup.LoadRawJsonReport(filePath)
	if !checkup.CheckUnmarshalResult(json.Unmarshal(jsonRaw, &report)) {
		return
	}
	result := A008Process(report)
	checkup.SaveReportResult(data, result)
}
