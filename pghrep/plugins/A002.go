package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"../src/checkup"
	"../src/pyraconv"
	"golang.org/x/net/html"
)

// go get golang.org/x/net/html

const VERSION_SOURCE_URL string = "https://git.postgresql.org/gitweb/?p=postgresql.git;a=tags"

type SupportedVersion struct {
	FirstRelease  string
	FinalRelease  string
	MinorVersions []int
}

var SUPPORTED_VERSIONS map[string]SupportedVersion = map[string]SupportedVersion{
	"11": SupportedVersion{
		FirstRelease:  "2018-10-18",
		FinalRelease:  "2023-11-09",
		MinorVersions: []int{3},
	},
	"10": SupportedVersion{
		FirstRelease:  "2017-10-05",
		FinalRelease:  "2022-11-10",
		MinorVersions: []int{8},
	},
	"9.6": SupportedVersion{
		FirstRelease:  "2016-09-29",
		FinalRelease:  "2021-11-11",
		MinorVersions: []int{13},
	},
	"9.5": SupportedVersion{
		FirstRelease:  "2016-01-07",
		FinalRelease:  "2021-02-11",
		MinorVersions: []int{17},
	},
	"9.4": SupportedVersion{
		FirstRelease:  "2014-12-18",
		FinalRelease:  "2020-02-13",
		MinorVersions: []int{22},
	},
}

type prepare string

type A002ReportHostResultData struct {
	Version          string `json:"version"`
	ServerVersionNum string `json:"server_version_num"`
	ServerMajorVer   string `json:"server_major_ver"`
	ServerMinorVer   string `json:"server_minor_ver"`
}

type A002ReportHostResult struct {
	Data      A002ReportHostResultData `json:"data"`
	NodesJson checkup.ReportLastNodes  `json:"nodes.json"`
}

type A002ReportHostsResults map[string]A002ReportHostResult

type A002Report struct {
	Project       string                  `json:"project"`
	Name          string                  `json:"name"`
	CheckId       string                  `json:"checkId"`
	Timestamptz   string                  `json:"timestamptz"`
	Database      string                  `json:"database"`
	Dependencies  map[string]interface{}  `json:"dependencies"`
	LastNodesJson checkup.ReportLastNodes `json:"last_nodes_json"`
	Results       A002ReportHostsResults  `json:"results"`
}

func A002PrepareVersionInfo() {
	url := VERSION_SOURCE_URL
	fmt.Printf("HTML code of %s ...\n", url)
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	htmlCode, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	domDocTest := html.NewTokenizer(strings.NewReader(string(htmlCode)))
	for tokenType := domDocTest.Next(); tokenType != html.ErrorToken; {
		if tokenType != html.TextToken {
			tokenType = domDocTest.Next()
			continue
		}
		rel := strings.TrimSpace(html.UnescapeString(string(domDocTest.Text())))
		if len(rel) > 3 && rel[0:3] == "REL" {
			if strings.Contains(rel, "BETA") || strings.Contains(rel, "RC") ||
				strings.Contains(rel, "ALPHA") {
				continue
			}
			ver := strings.Split(rel, "_")
			if len(ver) > 2 {
				majorVersion := strings.Replace(ver[0], "REL", "", 1)
				if majorVersion != "" {
					majorVersion = majorVersion + "."
				}
				majorVersion = majorVersion + ver[1]
				minorVersion := ver[2]
				ver, ok := SUPPORTED_VERSIONS[majorVersion]
				if ok {
					mVer, _ := strconv.Atoi(minorVersion)
					ver.MinorVersions = append(ver.MinorVersions, mVer)
					SUPPORTED_VERSIONS[majorVersion] = ver
				}
			}
		}
		tokenType = domDocTest.Next()
	}
}

func A002CheckAllVersionsIsSame(report A002Report,
	result checkup.ReportOutcome) checkup.ReportOutcome {
	var version string
	var hosts []string
	var vers []string
	diff := false
	for host, hostData := range report.Results {
		if version == "" {
			version = hostData.Data.ServerMajorVer + "." + hostData.Data.ServerMinorVer
		}
		if version != (hostData.Data.ServerMajorVer + "." + hostData.Data.ServerMinorVer) {
			diff = true
		}
		hosts = append(hosts, host)
		vers = append(vers, hostData.Data.ServerMajorVer+"."+hostData.Data.ServerMinorVer)
	}
	respectively := ""
	s := ""
	if len(hosts) > 1 {
		respectively = " respectively"
		s = "s"
	}
	if diff && len(hosts) > 1 {
		result.Conclusions = append(result.Conclusions, fmt.Sprintf("[P2] Not all nodes have the same "+
			"Postgres version. Node"+s+" %s uses Postgres %s"+respectively+".",
			strings.Join(hosts, ", "), strings.Join(vers, ", ")))
		result.Recommendations = append(result.Recommendations, fmt.Sprintf("[P2] Please upgrade "+
			"Postgres so its versions on all nodes match."))
		result.P2 = true
	} else {
		result.Conclusions = append(result.Conclusions, fmt.Sprintf("All nodes have the same Postgres "+
			"version (%s).", version))
	}
	return result
}

func A002CheckMajorVersions(report A002Report, result checkup.ReportOutcome) checkup.ReportOutcome {
	var processed map[string]bool = map[string]bool{}
	for host, hostData := range report.Results {
		if _, vok := processed[hostData.Data.ServerMajorVer]; vok {
			// version already checked
			continue
		}
		ver, ok := SUPPORTED_VERSIONS[hostData.Data.ServerMajorVer]
		if !ok {
			result.Conclusions = append(result.Conclusions, fmt.Sprintf("[P1] Unknown PostgreSQL version %s on %s.",
				hostData.Data.Version, host))
			result.Recommendations = append(result.Recommendations, fmt.Sprintf("[P1] Check PostgreSQL version on %s.", host))
			result.P1 = true
			continue
		}
		from, _ := time.Parse("2006-01-02", ver.FirstRelease)
		to, _ := time.Parse("2006-01-02", ver.FinalRelease)
		yearBeforeFinal := to.AddDate(-1, 0, 0)
		today := time.Now()
		if today.After(to) {
			// already not supported versions
			result.Conclusions = append(result.Conclusions, fmt.Sprintf("[P1] Postgres major version being used is %s and it is "+
				"NOT supported by Postgres community and PGDG (supported ended %s). This is a major issue. New bugs and security "+
				"issues will not be fixed by community and PGDG. You are on your own! Read more: "+
				"[Versioning Policy](https://www.postgresql.org/support/versioning/).",
				hostData.Data.ServerMajorVer, ver.FinalRelease))
			result.Recommendations = append(result.Recommendations, fmt.Sprintf("[P1] Please upgrade Postgres version %s to one of the "+
				"versions supported by the community and PGDG. To minimize downtime, consider using pg_upgrade or one "+
				"of solutions for logical replication.", hostData.Data.ServerMajorVer))
			result.P1 = true
		}
		if today.After(yearBeforeFinal) && today.Before(to) {
			result.Conclusions = append(result.Conclusions, fmt.Sprintf("[P2] Postgres community and PGDG will stop supporting version %s"+
				" within the next 12 months (end of life is scheduled %s). After that, you will be on your own!",
				hostData.Data.ServerMajorVer, ver.FinalRelease))
		}
		if today.After(from) && today.After(to) {
			// ok
			result.Conclusions = append(result.Conclusions, fmt.Sprintf("Postgres major version being used is %s and it is "+
				"currently supported by Postgres community and PGDG (end of life is scheduled %s). It means that in case "+
				"of bugs and security issues, updates (new minor versions) with fixes will be released and available for use."+
				" Read more: [Versioning Policy](https://www.postgresql.org/support/versioning/).",
				hostData.Data.ServerMajorVer, ver.FinalRelease))
		}
		processed[hostData.Data.ServerMajorVer] = true
	}
	return result
}

func A002CheckMinorVersions(report A002Report, result checkup.ReportOutcome) checkup.ReportOutcome {
	var updateHosts []string
	var curVersions []string
	var updateVersions []string
	var processed map[string]bool = map[string]bool{}
	for host, hostData := range report.Results {
		if _, vok := processed[hostData.Data.ServerMinorVer]; vok {
			// version already checked
			continue
		}
		ver, ok := SUPPORTED_VERSIONS[hostData.Data.ServerMajorVer]
		if !ok {
			result.Conclusions = append(result.Conclusions, fmt.Sprintf("[P1] Unknown PostgreSQL version %s on %s.",
				hostData.Data.Version, host))
			result.Recommendations = append(result.Recommendations, fmt.Sprintf("[P1] Check PostgreSQL version on %s.", host))
			result.P1 = true
			continue
		}
		sort.Ints(ver.MinorVersions)
		lastVersion := ver.MinorVersions[len(ver.MinorVersions)-1]
		minorVersion, _ := strconv.Atoi(hostData.Data.ServerMinorVer)
		if minorVersion == lastVersion {
			result.Conclusions = append(result.Conclusions, fmt.Sprintf("%s is the most up-to-date Postgres minor version in the branch %s.",
				hostData.Data.ServerMajorVer+"."+hostData.Data.ServerMinorVer, hostData.Data.ServerMajorVer))
			processed[hostData.Data.ServerMinorVer] = true
		} else {
			updateHosts = append(updateHosts, host)
			curVersions = append(curVersions, hostData.Data.ServerMajorVer+"."+hostData.Data.ServerMinorVer)
			updateVersions = append(updateVersions, hostData.Data.ServerMajorVer+"."+strconv.Itoa(lastVersion))
		}
	}
	s := ""
	respectively := ""
	if len(updateHosts) > 1 {
		s = "s"
		respectively = " respectively"
	}
	if len(updateHosts) > 0 {
		result.Conclusions = append(result.Conclusions, fmt.Sprintf("[P2] Node"+s+" %s use not the most up-to-date installed minor version"+s+
			" %s"+respectively+", the newest minor version"+s+" is %s).",
			strings.Join(updateHosts, ", "), strings.Join(curVersions, ", "), updateVersions[0]))
		result.Recommendations = append(result.Recommendations, fmt.Sprintf("[P2] Please upgrade Postgres to the most recent minor version: %s.",
			updateVersions[0]))
		result.P2 = true
	}
	return result
}

func A002Process(report A002Report) checkup.ReportOutcome {
	var result checkup.ReportOutcome
	A002PrepareVersionInfo()
	result = A002CheckAllVersionsIsSame(report, result)
	result = A002CheckMajorVersions(report, result)
	result = A002CheckMinorVersions(report, result)
	return result
}

func A002(data map[string]interface{}) {
	filePath := pyraconv.ToString(data["source_path_full"])
	jsonRaw := checkup.LoadRawJsonReport(filePath)
	var report A002Report
	err := json.Unmarshal(jsonRaw, &report)
	if err != nil {
		return
	}
	result := A002Process(report)
	if len(result.Recommendations) == 0 {
		result.Recommendations = append(result.Recommendations, "No recommendations.")
	} else {
		result.Recommendations = append(result.Recommendations, "  ")
		result.Recommendations = append(result.Recommendations, "For more information about minor and major upgrades see:")
		result.Recommendations = append(result.Recommendations, " - Official documentation: https://www.postgresql.org/docs") ///XX.YY/upgrading.html
		result.Recommendations = append(result.Recommendations, " - [Major-version upgrading with minimal downtime](https://www.depesz.com/2016/11/08/major-version-upgrading-with-minimal-downtime/) (depesz.com)")
		result.Recommendations = append(result.Recommendations, " - [Upgrading PostgreSQL on AWS RDS with minimum or zero downtime](https://medium.com/preply-engineering/postgres-multimaster-34f2446d5e14)")
		result.Recommendations = append(result.Recommendations, " - [Near-Zero Downtime Automated Upgrades of PostgreSQL Clusters in Cloud](https://www.2ndquadrant.com/en/blog/near-zero-downtime-automated-upgrades-postgresql-clusters-cloud/) (2ndQuadrant.com)")
		result.Recommendations = append(result.Recommendations, " - [Updating a 50 terabyte PostgreSQL database](https://medium.com/adyen/updating-a-50-terabyte-postgresql-database-f64384b799e7)")
	}

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
	A002(data)
	return data
}

var Preparer prepare
