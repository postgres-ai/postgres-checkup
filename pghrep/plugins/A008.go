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
)

type prepare string

const USAGE_CRITICAL int = 90
const USAGE_WARNING int = 70

const MSG_NO_USAGE_RISKS_CONCLUSION string = "No significant risks of out-of-disk-space problem have been detected."
const MSG_USAGE_WARNING_CONCLUSION string = "[P2] Disk %s on %s space usage is %s, it exceeds 70%%. There are some risks of out-of-disk-space problem."
const MSG_USAGE_WARNING_RECCOMENDATION string = "[P2] Add more disk space to %s on %s. It is recommended to keep free disk space less than %d%% " +
	"to reduce risks of out-of-disk-space problem."
const MSG_USAGE_CRITICAL_CONCLUSION string = "Disk %s on %s space usage is %s, it exceeds 90%%. There are significant risks of out-of-disk-space problem. " +
	"In this case, PostgreSQL will stop working and manual fix will be required."
const MSG_USAGE_CRITICAL_RECCOMENDATION string = "[P1] Add more disk space to %s on %s as soon as possible to prevent outage."
const MSG_NETWORK_FS_CONCLUSION string = "[P1] %s on host `%s` %s located on an NFS drive. This might lead to serious issues with Postgres, including downtime and data corruption."
const MSG_NETWORK_FS_RECCOMENDATION string = "[P1] Do not use NFS for Postgres."
const MSG_NOT_RECCOMENDED_FS_CONCLUSION string = "[P3] %s on host `%s` %s located on drive%s where the following filesystems are used: %s%s" +
	". This might mean that Postgres performance and reliability characteristics are worse than it could be in case of use of more popular filesystems (such as ext4)."
const MSG_NOT_RECCOMENDED_FS_RECCOMENDATION string = "[P3] Consider using ext4 for all Postgres directories."
const MSG_NO_RECCOMENDATION string = "No recommendations."

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
	conclusions []string, recommendations []string) (bool, bool, []string, []string) {
	usageWarning := false
	usageCritical := false
	usePercent := strings.Replace(fsItemData.UsePercent, "%", "", 1)
	percent, _ := strconv.Atoi(usePercent)
	if percent < USAGE_WARNING {
		// nothing to do
	}
	if percent >= USAGE_WARNING && percent < USAGE_CRITICAL {
		conclusions = append(conclusions, fmt.Sprintf(MSG_USAGE_WARNING_CONCLUSION, fsItemData.MountPoint,
			host, fsItemData.Used))
		recommendations = append(recommendations, fmt.Sprintf(MSG_USAGE_WARNING_RECCOMENDATION,
			fsItemData.MountPoint, host, USAGE_WARNING))
		usageWarning = true
	}
	if percent >= USAGE_CRITICAL {
		conclusions = append(conclusions, fmt.Sprintf(MSG_USAGE_CRITICAL_CONCLUSION, fsItemData.MountPoint,
			host, fsItemData.Used))
		recommendations = append(recommendations, fmt.Sprintf(MSG_USAGE_CRITICAL_RECCOMENDATION,
			fsItemData.MountPoint, host))
		usageCritical = true
	}
	return usageWarning, usageCritical, conclusions, recommendations
}

// Generate conclusions and recommendatons
func A008Process(report A008Report) checkup.ReportOutcome {
	var result checkup.ReportOutcome
	usageWarning := false
	usageCritical := false
	var nfsConclusions []string
	var notRecConclusions []string
	var conclusions []string
	var recommendations []string
	for host, hostResult := range report.Results {
		var networkFsItems []FsItem
		var notRecommendedFsItems []FsItem
		for _, fsItemData := range hostResult.Data.DbData {
			if isRecommendedFs(strings.ToLower(fsItemData.Fstype)) != true {
				notRecommendedFsItems = append(notRecommendedFsItems, fsItemData)
			}
			if strings.ToLower(fsItemData.Fstype[0:3]) == "nfs" {
				networkFsItems = append(networkFsItems, fsItemData)
			}
			var problem, critical bool
			problem, critical, conclusions, recommendations = checkFsItemUsage(host, fsItemData, conclusions, recommendations)
			usageWarning = usageWarning || problem
			usageCritical = usageCritical || critical
		}
		for _, fsItemData := range hostResult.Data.FsData {
			var problem, critical bool
			problem, critical, conclusions, recommendations = checkFsItemUsage(host, fsItemData, conclusions, recommendations)
			usageWarning = usageWarning || problem
			usageCritical = usageCritical || critical
		}
		if len(networkFsItems) > 0 {
			var networkFsDisks []string
			for _, fsItem := range networkFsItems {
				networkFsDisks = append(networkFsDisks, fsItem.MountPoint)
			}
			result.P1 = true
			areIs := "is"
			if len(networkFsDisks) > 1 {
				areIs = "are"
			}
			nfsConclusions = append(nfsConclusions, fmt.Sprintf(MSG_NETWORK_FS_CONCLUSION,
				strings.Join(networkFsDisks, ", "), host, areIs))
		}
		if len(notRecommendedFsItems) > 0 {
			var notRecommendedFsDisks []string
			var notRecommendedFsDisksFs []string
			for _, fsItem := range notRecommendedFsItems {
				notRecommendedFsDisks = append(notRecommendedFsDisks, fsItem.MountPoint)
				notRecommendedFsDisksFs = append(notRecommendedFsDisksFs, fsItem.Fstype)
			}
			result.P3 = true
			areIs := "is"
			respectively := ""
			s := ""
			if len(notRecommendedFsDisks) > 1 {
				areIs = "are"
				respectively = " respectively"
				s = "s"
			}
			notRecConclusions = append(notRecConclusions, fmt.Sprintf(MSG_NOT_RECCOMENDED_FS_CONCLUSION,
				strings.Join(notRecommendedFsDisks, ", "), host, areIs, s, strings.Join(notRecommendedFsDisksFs, ", "), respectively))
		}

	}
	if usageWarning {
		result.P2 = true
	}
	if usageCritical {
		result.P1 = true
	}
	if !usageWarning && !usageCritical && len(recommendations) == 0 {
		conclusions = append(conclusions, MSG_NO_USAGE_RISKS_CONCLUSION)
	}
	if len(nfsConclusions) > 0 {
		conclusions = append(conclusions, nfsConclusions...)
		recommendations = append(recommendations, MSG_NETWORK_FS_RECCOMENDATION)
	}
	if len(notRecConclusions) > 0 {
		conclusions = append(conclusions, notRecConclusions...)
		recommendations = append(recommendations, MSG_NOT_RECCOMENDED_FS_RECCOMENDATION)
	}
	if len(recommendations) == 0 {
		recommendations = append(recommendations, MSG_NO_RECCOMENDATION)
	}
	result.Conclusions = conclusions
	result.Recommendations = recommendations
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
