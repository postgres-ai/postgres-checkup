package a002

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	checkup ".."
	"../../log"
	"../../pyraconv"
	"github.com/dustin/go-humanize/english"
	"golang.org/x/net/html"
)

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

func A002PrepareVersionInfo() {
	url := VERSION_SOURCE_URL
	log.Dbg("HTML code of %s ...\n", url)
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
	curVersions = removeArrayDoubles(curVersions)
	if len(curVersions) > 0 {
		result.AppendConclusion(english.PluralWord(len(curVersions),
			MSG_NOT_LAST_MINOR_VERSION_CONCLUSION_1, MSG_NOT_LAST_MINOR_VERSION_CONCLUSION_N),
			strings.Join(curVersions, ", "), updateVersions[0])
		result.AppendRecommendation(MSG_NOT_LAST_MINOR_VERSION_RECOMMENDATION, updateVersions[0])
		result.P2 = true
	}
	return result
}

func removeArrayDoubles(array []string) []string {
	items := map[string]int{}
	var result []string
	for _, data := range array {
		items[data] = 1
	}
	for key, _ := range items {
		result = append(result, key)
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

func A002PreprocessReportData(data map[string]interface{}) {
	filePath := pyraconv.ToString(data["source_path_full"])
	jsonRaw := checkup.LoadRawJsonReport(filePath)
	var report A002Report
	err := json.Unmarshal(jsonRaw, &report)
	if err != nil {
		log.Err("Can't load json report to process")
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
