package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"../src/checkup"
	"../src/pyraconv"
)

type prepare string

const CRITICAL_USAGE int = 90
const PROBLEM_USAGE int = 70

var VALID_FS []string = []string{
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

func isValidFs(fs string) bool {
	for _, validFs := range VALID_FS {
		if validFs == fs {
			return true
		}
	}
	return false
}

func checkFsItem(host string, fsItemData FsItem,
	conclusions []string, recommendations []string) (bool, bool, bool, bool) {
	nfs := false
	problemUsage := false
	criticalUsage := false
	nonExt4 := false

	if isValidFs(strings.ToLower(fsItemData.Fstype)) != true {
		nonExt4 = true
	}
	if strings.ToLower(fsItemData.Fstype[0:3]) == "nfs" {
		nfs = true
	}
	usePercent := strings.Replace(fsItemData.UsePercent, "%", "", 1)
	percent, _ := strconv.Atoi(usePercent)
	if percent < PROBLEM_USAGE {
		// nothing to do
	}
	if percent >= PROBLEM_USAGE && percent < CRITICAL_USAGE {
		conclusions = append(conclusions, fmt.Sprintf("[P2] Disk %s on %s space usage is %s, it exceeds 70%%. "+
			"There are some risks of out-of-disk-space problem.", fsItemData.MountPoint, host, fsItemData.Used))
		recommendations = append(recommendations, fmt.Sprintf("[P2] Add more disk space to %s on %s. "+
			"It is recommended to keep free disk space less than 70%% "+
			"to reduce risks of out-of-disk-space problem.", fsItemData.MountPoint, host))
		problemUsage = true
	}
	if percent >= CRITICAL_USAGE {
		conclusions = append(conclusions, fmt.Sprintf("Disk %s on %s space usage is %s, it exceeds 90%%. "+
			"There are significant risks of out-of-disk-space problem. "+
			"In this case, PostgreSQL will stop working and manual fix will be required.",
			fsItemData.MountPoint, host, fsItemData.Used))
		recommendations = append(recommendations, fmt.Sprintf("[P1] Add more disk space to %s on %s as "+
			"soon as possible to prevent outage.", fsItemData.MountPoint, host))
		criticalUsage = true
	}
	return problemUsage, criticalUsage, nfs, nonExt4
}

// Generate conclusions and recommendatons
func A008Process(report A008Report) checkup.ReportOutcome {
	var result checkup.ReportOutcome
	problemUsage := false
	criticalUsage := false
	var nfsConclusions []string
	var nExtConclusions []string
	var conclusions []string
	var recommendations []string
	for host, hostResult := range report.Results {
		var nfsItems []FsItem
		var notRecItems []FsItem
		for _, fsItemData := range hostResult.Data.DbData {
			problem, critical, nfs, notrec := checkFsItem(host, fsItemData, conclusions, recommendations)
			problemUsage = problemUsage || problem
			criticalUsage = criticalUsage || critical
			if nfs {
				nfsItems = append(nfsItems, fsItemData)
			}
			if notrec {
				notRecItems = append(notRecItems, fsItemData)
			}
		}
		for _, fsItemData := range hostResult.Data.FsData {
			problem, critical, _, _ := checkFsItem(host, fsItemData, conclusions, recommendations)
			problemUsage = problemUsage || problem
			criticalUsage = criticalUsage || critical
		}
		if len(nfsItems) > 0 {
			var nfsDisks []string
			for _, nfsItem := range nfsItems {
				nfsDisks = append(nfsDisks, nfsItem.MountPoint)
			}
			result.P1 = true
			areIs := "is"
			if len(nfsDisks) > 1 {
				areIs = "are"
			}
			nfsConclusions = append(nfsConclusions, fmt.Sprintf("[P1] %s on host `%s` "+areIs+" located on an NFS drive. "+
				"This might lead to serious issues with Postgres, including downtime and data corruption.",
				strings.Join(nfsDisks, ", "), host))
		}
		if len(notRecItems) > 0 {
			var nExtDisks []string
			var nExtDiskFs []string
			for _, nExtItem := range notRecItems {
				nExtDisks = append(nExtDisks, nExtItem.MountPoint)
				nExtDiskFs = append(nExtDiskFs, nExtItem.Fstype)
			}
			result.P3 = true
			areIs := "is"
			respectively := ""
			s := ""
			if len(nExtDisks) > 1 {
				areIs = "are"
				respectively = " respectively"
				s = "s"
			}
			nExtConclusions = append(nExtConclusions, fmt.Sprintf("[P3] %s on host `%s` "+areIs+
				" located on drive"+s+" where the following filesystems are used: %s"+respectively+
				". This might mean that Postgres performance and reliability characteristics are worse than it "+
				"could be in case of use of more popular filesystems (such as ext4).",
				strings.Join(nExtDisks, ", "), host, strings.Join(nExtDiskFs, ", ")))
		}

	}
	if problemUsage {
		result.P2 = true
	}
	if criticalUsage {
		result.P1 = true
	}
	if !problemUsage && !criticalUsage && len(recommendations) == 0 {
		conclusions = append(conclusions, "No significant risks of out-of-disk-space problem have been detected.")
	}
	if len(nfsConclusions) > 0 {
		conclusions = append(conclusions, nfsConclusions...)
		recommendations = append(recommendations, "[P1] Do not use NFS for Postgres.")
	}
	if len(nExtConclusions) > 0 {
		conclusions = append(conclusions, nExtConclusions...)
		recommendations = append(recommendations, "[P3] Consider using ext4 for all Postgres directories.")
	}
	if len(recommendations) == 0 {
		recommendations = append(recommendations, "No recommendations.")
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
		return
	}
	result := A008Process(report)
	// update data
	data["conclusions"] = result.Conclusions
	data["recommendations"] = result.Recommendations
	data["p1"] = result.P1
	data["p2"] = result.P2
	data["p3"] = result.P3
	// update file
	checkup.SaveJsonConclusionsRecommendations(data, result.Conclusions, result.Recommendations, result.P1, result.P2, result.P3)
}

// Plugin entry point
func (g prepare) Prepare(data map[string]interface{}) map[string]interface{} {
	A008(data)
	return data
}

var Preparer prepare
