package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"../src/checkup"
	"../src/pyraconv"
	"github.com/dustin/go-humanize/english"
)

type prepare string

const USAGE_CRITICAL int = 90
const USAGE_WARNING int = 70

const MSG_NO_USAGE_RISKS_CONCLUSION string = "No significant risks of out-of-disk-space problem have been detected."
const MSG_USAGE_WARNING_CONCLUSION string = "[P2] Disk %s on %s space usage is %s, it exceeds 70%%. There are some risks of out-of-disk-space problem."
const MSG_USAGE_WARNING_RECOMMENDATION string = "[P2] Add more disk space to %s on %s. It is recommended to keep free disk space less than %d%% " +
	"to reduce risks of out-of-disk-space problem."
const MSG_USAGE_CRITICAL_CONCLUSION string = "Disk %s on %s space usage is %s, it exceeds 90%%. There are significant risks of out-of-disk-space problem. " +
	"In this case, PostgreSQL will stop working and manual fix will be required."
const MSG_USAGE_CRITICAL_RECOMMENDATION string = "[P1] Add more disk space to %s on %s as soon as possible to prevent outage."
const MSG_NETWORK_FS_CONCLUSION_1 string = "[P1] %s on host `%s` is located on an NFS drive. This might lead to serious issues with Postgres, including downtime and data corruption."
const MSG_NETWORK_FS_CONCLUSION_N string = "[P1] %s on host `%s` are located on an NFS drive. This might lead to serious issues with Postgres, including downtime and data corruption."
const MSG_NETWORK_FS_RECOMMENDATION string = "[P1] Do not use NFS for Postgres."
const MSG_NOT_RECOMMENDED_FS_CONCLUSION_1 string = "[P3] %s on host `%s` is located on drive where the following filesystems are used: %s. This might mean that Postgres performance and reliability characteristics are worse than it could be in case of use of more popular filesystems (such as ext4)."
const MSG_NOT_RECOMMENDED_FS_CONCLUSION_N string = "[P3] %s on host `%s` are located on drives where the following filesystems are used: %s respectively. This might mean that Postgres performance and reliability characteristics are worse than it could be in case of use of more popular filesystems (such as ext4)."
const MSG_NOT_RECOMMENDED_FS_RECOMMENDATION string = "[P3] Consider using ext4 for all Postgres directories."
const MSG_NO_FS_RECOMMENDATION string = "No recommendations."

var FS_RECOMMENDED []string = []string{
	"ext4",
	"xfs",
	"tmpfs",
}

type FsItem struct {
	Fstype     string `json:"fstype"`
	Size       string `json:"size"`
	Avail      string `json:"avail"`
	Used       string `json:"used"`
	UsePercent string `json:"use_percent"`
	MountPoint string `json:"mount_point"`
	Path       string `json:"path"`
	Device     string `json:"device"`
}

type A008ReportHostResultData struct {
	DbData map[string]FsItem `json:"data"`
	FsData map[string]FsItem `json:"fs_data"`
}

type A008ReportHostResult struct {
	Data      A008ReportHostResultData `json:"data"`
	NodesJson checkup.ReportLastNodes  `json:"nodes.json"`
}

type A008ReportHostsResults map[string]A008ReportHostResult

type A008Report struct {
	Project       string                  `json:"project"`
	Name          string                  `json:"name"`
	CheckId       string                  `json:"checkId"`
	Timestamptz   string                  `json:"timestamptz"`
	Database      string                  `json:"database"`
	Dependencies  map[string]interface{}  `json:"dependencies"`
	LastNodesJson checkup.ReportLastNodes `json:"last_nodes_json"`
	Results       A008ReportHostsResults  `json:"results"`
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
	if percent < USAGE_WARNING {
		// nothing to do
	}
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

func A008(data map[string]interface{}) {
	filePath := pyraconv.ToString(data["source_path_full"])
	jsonRaw := checkup.LoadRawJsonReport(filePath)
	var report A008Report
	err := json.Unmarshal(jsonRaw, &report)
	if err != nil {
		log.New(os.Stderr, "", 0).Println("Can't load json report to process")
		return
	}
	result := A008Process(report)
	// update data and file
	checkup.SaveConclusionsRecommendations(data, result)
}

// Plugin entry point
func (g prepare) Prepare(data map[string]interface{}) map[string]interface{} {
	A008(data)
	return data
}

var Preparer prepare
