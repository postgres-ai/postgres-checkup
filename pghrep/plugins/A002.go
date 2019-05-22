package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"../src/checkup"
	"../src/pyraconv"
	"github.com/dustin/go-humanize/english"
	"golang.org/x/net/html"
)

type prepare string

const VERSION_SOURCE_URL string = "https://git.postgresql.org/gitweb/?p=postgresql.git;a=tags"

const MSG_WRONG_VERSION_CONCLUSION string = "[P1] Unknown PostgreSQL version %s on %s."
const MSG_WRONG_VERSION_RECOMMENDATION string = "[P1] Check PostgreSQL version on %s."
const MSG_NOT_SUPPORTED_VERSION_CONCLUSION string = "[P1] Postgres major version being used is %s and it is " +
	"NOT supported by Postgres community and PGDG (supported ended %s). This is a major issue. New bugs and security " +
	"issues will not be fixed by community and PGDG. You are on your own! Read more: " +
	"[Versioning Policy](https://www.postgresql.org/support/versioning/)."
const MSG_NOT_SUPPORTED_VERSION_RECOMMENDATION string = "[P1] Please upgrade Postgres version %s to one of the " +
	"versions supported by the community and PGDG. To minimize downtime, consider using pg_upgrade or one " +
	"of solutions for logical replication."
const MSG_LAST_YEAR_SUPPORTED_VERSION_CONCLUSION string = "[P2] Postgres community and PGDG will stop supporting version %s" +
	" within the next 12 months (end of life is scheduled %s). After that, you will be on your own!"
const MSG_SUPPORTED_VERSION_CONCLUSION string = "Postgres major version being used is %s and it is " +
	"currently supported by Postgres community and PGDG (end of life is scheduled %s). It means that in case " +
	"of bugs and security issues, updates (new minor versions) with fixes will be released and available for use." +
	" Read more: [Versioning Policy](https://www.postgresql.org/support/versioning/)."
const MSG_LAST_MINOR_VERSION_CONCLUSION string = "%s is the most up-to-date Postgres minor version in the branch %s."
const MSG_NOT_LAST_MINOR_VERSION_CONCLUSION_1 string = "[P2] Node %s use not the most up-to-date installed minor version%s %s, the newest minor version is %s)."
const MSG_NOT_LAST_MINOR_VERSION_CONCLUSION_N string = "[P2] Nodes %s use not the most up-to-date installed minor versions %s respectively, the newest minor version is %s)."

const MSG_NOT_ALL_VERSIONS_SAME_CONCLUSION_1 string = "[P2] Not all nodes have the same Postgres version. Node %s uses Postgres %s."
const MSG_NOT_ALL_VERSIONS_SAME_CONCLUSION_N string = "[P2] Not all nodes have the same Postgres version. Nodes %s uses Postgres %s respectively."
const MSG_NOT_ALL_VERSIONS_SAME_RECOMMENDATION string = "[P2] Please upgrade Postgres so its versions on all nodes match."
const MSG_ALL_VERSIONS_SAME_CONCLUSION string = "All nodes have the same Postgres version (%s)."

const MSG_NOT_LAST_MINOR_VERSION_RECOMMENDATION string = "[P2] Please upgrade Postgres to the most recent minor version: %s."
const MSG_NO_RECOMMENDATION string = "No recommendations."
const MSG_GENERAL_RECOMMENDATION string = "  \n" +
	"For more information about minor and major upgrades see:  \n" +
	" - Official documentation: https://www.postgresql.org/docs  \n" + ///XX.YY/upgrading.html
	" - [Major-version upgrading with minimal downtime](https://www.depesz.com/2016/11/08/major-version-upgrading-with-minimal-downtime/) (depesz.com)  \n" +
	" - [Upgrading PostgreSQL on AWS RDS with minimum or zero downtime](https://medium.com/preply-engineering/postgres-multimaster-34f2446d5e14)  \n" +
	" - [Near-Zero Downtime Automated Upgrades of PostgreSQL Clusters in Cloud](https://www.2ndquadrant.com/en/blog/near-zero-downtime-automated-upgrades-postgresql-clusters-cloud/) (2ndQuadrant.com)  \n" +
	" - [Updating a 50 terabyte PostgreSQL database](https://medium.com/adyen/updating-a-50-terabyte-postgresql-database-f64384b799e7)  \n"

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
	if diff && len(hosts) > 1 {
		result.AppendConclusion(english.PluralWord(len(hosts),
			MSG_NOT_ALL_VERSIONS_SAME_CONCLUSION_1, MSG_NOT_ALL_VERSIONS_SAME_CONCLUSION_N),
			strings.Join(hosts, ", "), strings.Join(vers, ", "))
		result.AppendRecommendation(MSG_NOT_ALL_VERSIONS_SAME_RECOMMENDATION)
		result.P2 = true
	} else {
		result.AppendConclusion(MSG_ALL_VERSIONS_SAME_CONCLUSION, version)
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
			result.AppendConclusion(MSG_WRONG_VERSION_CONCLUSION, hostData.Data.Version, host)
			result.AppendRecommendation(MSG_WRONG_VERSION_RECOMMENDATION, host)
			result.P1 = true
			continue
		}
		from, _ := time.Parse("2006-01-02", ver.FirstRelease)
		to, _ := time.Parse("2006-01-02", ver.FinalRelease)
		yearBeforeFinal := to.AddDate(-1, 0, 0)
		today := time.Now()
		if today.After(to) {
			// already not supported versions
			result.AppendConclusion(MSG_NOT_SUPPORTED_VERSION_CONCLUSION, hostData.Data.ServerMajorVer, ver.FinalRelease)
			result.AppendRecommendation(MSG_NOT_SUPPORTED_VERSION_RECOMMENDATION, hostData.Data.ServerMajorVer)
			result.P1 = true
		}
		if today.After(yearBeforeFinal) && today.Before(to) {
			// supported last year
			result.AppendConclusion(MSG_LAST_YEAR_SUPPORTED_VERSION_CONCLUSION, hostData.Data.ServerMajorVer, ver.FinalRelease)
			result.P2 = true
		}
		if today.After(from) && today.After(to) {
			// ok
			result.AppendConclusion(MSG_SUPPORTED_VERSION_CONCLUSION, hostData.Data.ServerMajorVer, ver.FinalRelease)
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
			result.AppendConclusion(MSG_NOT_SUPPORTED_VERSION_CONCLUSION, hostData.Data.ServerMajorVer, ver.FinalRelease)
			result.AppendRecommendation(MSG_NOT_SUPPORTED_VERSION_RECOMMENDATION, hostData.Data.ServerMajorVer)
			result.P1 = true
			continue
		}
		sort.Ints(ver.MinorVersions)
		lastVersion := ver.MinorVersions[len(ver.MinorVersions)-1]
		minorVersion, _ := strconv.Atoi(hostData.Data.ServerMinorVer)
		if minorVersion >= lastVersion {
			result.AppendConclusion(MSG_LAST_MINOR_VERSION_CONCLUSION,
				hostData.Data.ServerMajorVer+"."+hostData.Data.ServerMinorVer, hostData.Data.ServerMajorVer)
			processed[hostData.Data.ServerMinorVer] = true
		} else {
			updateHosts = append(updateHosts, host)
			curVersions = append(curVersions, hostData.Data.ServerMajorVer+"."+hostData.Data.ServerMinorVer)
			updateVersions = append(updateVersions, hostData.Data.ServerMajorVer+"."+strconv.Itoa(lastVersion))
		}
	}
	if len(updateHosts) > 0 {
		result.AppendConclusion(english.PluralWord(len(updateHosts),
			MSG_NOT_LAST_MINOR_VERSION_CONCLUSION_1, MSG_NOT_LAST_MINOR_VERSION_CONCLUSION_N),
			strings.Join(updateHosts, ", "), strings.Join(curVersions, ", "), updateVersions[0])
		result.AppendRecommendation(MSG_NOT_LAST_MINOR_VERSION_RECOMMENDATION, updateVersions[0])
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
		log.New(os.Stderr, "", 0).Println("Can't load json report to process")
		return
	}
	result := A002Process(report)
	if len(result.Recommendations) == 0 {
		result.AppendRecommendation(MSG_NO_RECOMMENDATION)
	} else {
		result.AppendRecommendation(MSG_GENERAL_RECOMMENDATION)
	}

	// update data and file
	checkup.SaveConclusionsRecommendations(data, result)
}

// Plugin entry point
func (g prepare) Prepare(data map[string]interface{}) map[string]interface{} {
	A002(data)
	return data
}

var Preparer prepare
