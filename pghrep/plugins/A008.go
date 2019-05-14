package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"../src/checkup"
	"../src/pyraconv"
)

const CRITICAL_USAGE int = 90
const PROBLEM_USAGE int = 70

type prepare string

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

type ReportHostResultData struct {
	DbData map[string]FsItem `json:"db_data"`
	FsData map[string]FsItem `json:"fs_data"`
}

type ReportHostResult struct {
	Data      ReportHostResultData    `json:"data"`
	NodesJson checkup.ReportLastNodes `json:"nodes.json"`
}

type ReportHostsResults map[string]ReportHostResult

type Report struct {
	Project       string                  `json:"project"`
	Name          string                  `json:"name"`
	CheckId       string                  `json:"checkId"`
	Timestamptz   string                  `json:"timestamptz"`
	Database      string                  `json:"database"`
	Dependencies  map[string]interface{}  `json:"dependencies"`
	LastNodesJson checkup.ReportLastNodes `json:"last_nodes_json"`
	Results       ReportHostsResults      `json:"results"`
}

func readData(filePath string) *Report {
	file, err := os.Open(filePath)
	if err != nil {
		return nil
	}
	defer file.Close()
	jsonRaw, err := ioutil.ReadAll(file)
	if err != nil {
		return nil
	}
	var report Report
	err = json.Unmarshal(jsonRaw, &report)
	if err != nil {
		return nil
	}
	return &report
}

func checkFsItem(fsItemData FsItem) (bool, bool, bool, bool) {
	nfs := false
	less70p := false
	in7090p := false
	more90p := false
	if strings.ToLower(fsItemData.Fstype) == "nfs" {
		nfs = true
	}
	usePercent := strings.Replace(fsItemData.UsePercent, "%", "", 1)
	percent, _ := strconv.Atoi(usePercent)
	if percent < PROBLEM_USAGE {
		less70p = true
	}
	if percent >= PROBLEM_USAGE && percent < CRITICAL_USAGE {
		in7090p = true
	}
	if percent >= CRITICAL_USAGE {
		more90p = true
	}
	return nfs, less70p, in7090p, more90p
}

// check specified disk
func checkDisk(fsItemData FsItem, conclusions []string, recommendations []string) ([]string, []string) {
	if strings.ToLower(fsItemData.Fstype) == "nfs" {
		conclusions = append(conclusions, "NFS found")
		recommendations = append(recommendations, "NFS found")
	}
	usePercent := strings.Replace(fsItemData.UsePercent, "%", "", 1)
	percent, _ := strconv.Atoi(usePercent)
	if percent < PROBLEM_USAGE && fsItemData.Fstype != "tmpfs" {
		conclusions = append(conclusions, "No significant risks of out-of-disk-space problem have been detected.")
		recommendations = append(recommendations, "No recommendations.")
	}
	if percent >= PROBLEM_USAGE && percent < CRITICAL_USAGE && fsItemData.Fstype != "tmpfs" {
		conclusions = append(conclusions, "[P2] Free disk space exceeds 70%. There are some risks of out-of-disk-space problem.")
		recommendations = append(recommendations, "[P2] Add more disk space. It is recommended to keep free disk space less than 70%. to reduce risks of out-of-disk-space problem.")
	}
	if percent >= CRITICAL_USAGE && fsItemData.Fstype != "tmpfs" {
		conclusions = append(conclusions, "[P1] Free disk space exceeds 90%. There are significant risks of out-of-disk-space problem. In this case, PostgreSQL will stop working and manual fix will be required.")
		recommendations = append(recommendations, "[P1] Add more disk space as soon as possible to prevent outage.")
	}
	return conclusions, recommendations
}

// Generate conclusions and recommendatons
func A008(data map[string]interface{}) {
	nfs := false
	less70p := false
	in7090p := false
	more90p := false
	p1 := false
	p2 := false
	var conclusions []string
	var recommendations []string
	filePath := pyraconv.ToString(data["source_path_full"])
	report := readData(filePath)
	if report == nil {
		return
	}
	for _, hostResult := range report.Results {
		for _, fsItemData := range hostResult.Data.DbData {
			n, l, i, m := checkFsItem(fsItemData)
			nfs = nfs || n
			less70p = less70p || l
			in7090p = in7090p || i
			more90p = more90p || m
		}
		for _, fsItemData := range hostResult.Data.FsData {
			n, l, i, m := checkFsItem(fsItemData)
			nfs = nfs || n
			less70p = less70p || l
			in7090p = in7090p || i
			more90p = more90p || m
		}
	}
	if nfs {
		conclusions = append(conclusions, "NFS disk found")
		recommendations = append(recommendations, "NFS disk found")
		p1 = true
	}
	if more90p {
		conclusions = append(conclusions, "[P1] Free disk space exceeds 90%. There are significant risks of out-of-disk-space problem. "+
			"In this case, PostgreSQL will stop working and manual fix will be required.")
		recommendations = append(recommendations, "[P1] Add more disk space as soon as possible to prevent outage.")
		p1 = true
	}
	if in7090p {
		conclusions = append(conclusions, "[P2] Free disk space exceeds 70%. There are some risks of out-of-disk-space problem.")
		recommendations = append(recommendations, "[P2] Add more disk space. It is recommended to keep free disk space less than 70%. "+
			"To reduce risks of out-of-disk-space problem.")
		p2 = true
	}
	if less70p && !in7090p && !more90p {
		conclusions = append(conclusions, "No significant risks of out-of-disk-space problem have been detected.")
		if len(recommendations) == 0 {
			recommendations = append(recommendations, "No recommendations.")
		}
	}
    // update data
    data["conclusions"] = conclusions
	data["recommendations"] = recommendations
	data["p1"] = p1
	data["p2"] = p2
    // update file
	checkup.SaveJsonConclusionsRecommendations(data, conclusions, recommendations, p1, p2)
}

// Plugin entry point
func (g prepare) Prepare(data map[string]interface{}) map[string]interface{} {
	A008(data)
	return data
}

var Preparer prepare
