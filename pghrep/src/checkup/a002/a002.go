package a002

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	checkup ".."
	"../../log"
	"../config"

	"github.com/dustin/go-humanize/english"
	"golang.org/x/net/html"
)

// Case codes
const A002_NOT_ALL_VERSIONS_SAME string = "A002_NOT_ALL_VERSIONS_SAME"
const A002_ALL_VERSIONS_SAME string = "A002_ALL_VERSIONS_SAME"
const A002_NOT_SUPPORTED_VERSION string = "A002_NOT_SUPPORTED_VERSION"
const A002_LAST_YEAR_SUPPORTED_VERSION string = "A002_LAST_YEAR_SUPPORTED_VERSION"
const A002_SUPPORTED_VERSION string = "A002_SUPPORTED_VERSION"
const A002_NOT_LAST_MAJOR_VERSION string = "A002_NOT_LAST_MAJOR_VERSION"
const A002_UNKNOWN_VERSION string = "A002_UNKNOWN_VERSION"
const A002_LAST_MINOR_VERSION string = "A002_LAST_MINOR_VERSION"
const A002_NOT_LAST_MINOR_VERSION string = "A002_NOT_LAST_MINOR_VERSION"
const A002_GENERAL_INFO_OFFICIAL string = "A002_GENERAL_OFFICIAL"
const A002_GENERAL_INFO_FULL string = "A002_GENERAL_FULL"

func getMajorMinorVersion(serverVersion string) (string, string) {
	var minorVersion string
	var majorVersion string
	minorVersion = serverVersion[len(serverVersion)-2 : len(serverVersion)]
	i, _ := strconv.Atoi(minorVersion)
	minorVersion = strconv.Itoa(i)
	if serverVersion[0:1] == "9" {
		majorVersion = serverVersion[0:3]
		majorVersion = strings.Replace(majorVersion, "0", ".", 1)
	} else {
		majorVersion = serverVersion[0:2]
	}
	return majorVersion, minorVersion
}

func A002CheckAllVersionsIsSame(report A002Report, result checkup.ReportResult) checkup.ReportResult {
	var version string
	var hosts []string
	var vers []string
	diff := false
	for host, hostData := range report.Results {
		majorVersion, minorVersion := getMajorMinorVersion(hostData.Data.ServerVersionNum)
		if version == "" {
			version = majorVersion + "." + minorVersion
		}
		if version != (majorVersion + "." + minorVersion) {
			diff = true
		}
		hosts = append(hosts, host)
		vers = append(vers, majorVersion+"."+minorVersion)
	}
	if diff && len(hosts) > 1 {
		result.AppendConclusion(A002_NOT_ALL_VERSIONS_SAME, english.PluralWord(len(hosts),
			MSG_NOT_ALL_VERSIONS_SAME_CONCLUSION_1, MSG_NOT_ALL_VERSIONS_SAME_CONCLUSION_N),
			strings.Join(hosts, "`, `"), strings.Join(checkup.GetUniques(vers), "`, `"))
		result.AppendRecommendation(A002_NOT_ALL_VERSIONS_SAME, MSG_NOT_ALL_VERSIONS_SAME_RECOMMENDATION)
		result.P2 = true
	} else {
		if len(hosts) > 0 {
			result.AppendConclusion(A002_ALL_VERSIONS_SAME, MSG_ALL_VERSIONS_SAME_CONCLUSION, version)
		}
	}
	return result
}

func A002CheckMajorVersions(report A002Report, config config.Config,
	result checkup.ReportResult) checkup.ReportResult {
	var processed map[string]bool = map[string]bool{}
	for host, hostData := range report.Results {
		majorVersion, _ := getMajorMinorVersion(hostData.Data.ServerVersionNum)
		mjVersion := hostData.Data.ServerVersionNum[0 : len(hostData.Data.ServerVersionNum)-2]
		iMajorVersion, _ := strconv.Atoi(mjVersion)
		if _, vok := processed[majorVersion]; vok {
			// version already checked
			continue
		}
		ver, ok := config.Versions[majorVersion]
		if !ok {
			result.AppendConclusion(A002_UNKNOWN_VERSION, MSG_UNKNOWN_VERSION_CONCLUSION, hostData.Data.Version, host)
			result.AppendRecommendation(A002_UNKNOWN_VERSION, MSG_UNKNOWN_VERSION_RECOMMENDATION, host)
			result.P1 = true
			continue
		}
		from, _ := time.Parse("2006-01-02", ver.FirstRelease)
		to, _ := time.Parse("2006-01-02", ver.FinalRelease)
		yearBeforeFinal := to.AddDate(-1, 0, 0)
		today := time.Now()
		if today.After(to) {
			// already not supported versions
			result.AppendConclusion(MSG_NOT_SUPPORTED_VERSION_CONCLUSION, majorVersion, ver.FinalRelease)
			result.AppendRecommendation(MSG_NOT_SUPPORTED_VERSION_RECOMMENDATION, majorVersion)
			result.P1 = true
		}
		if today.After(yearBeforeFinal) && today.Before(to) {
			// supported last year
			result.AppendConclusion(MSG_LAST_YEAR_SUPPORTED_VERSION_CONCLUSION, majorVersion, ver.FinalRelease)
			result.P2 = true
		}
		if today.After(from) && today.After(to) {
			// ok
			result.AppendConclusion(MSG_SUPPORTED_VERSION_CONCLUSION, majorVersion, ver.FinalRelease)
		}
		// TODO(anatoly)
		/*if MAJOR_VERSIONS[len(MAJOR_VERSIONS)-1] > iMajorVersion {
			result.AppendRecommendation(A002_NOT_LAST_MAJOR_VERSION, MSG_NOT_LAST_MAJOR_VERSION_CONCLUSION, float32(MAJOR_VERSIONS[len(MAJOR_VERSIONS)-1])/100.0)
			result.AppendRecommendation(A002_GENERAL_INFO_FULL, MSG_GENERAL_RECOMMENDATION_1+MSG_GENERAL_RECOMMENDATION_2)
			result.P3 = true
		}*/
		processed[majorVersion] = true
	}
	return result
}

func A002CheckMinorVersions(report A002Report, config config.Config,
	result checkup.ReportResult) checkup.ReportResult {
	var updateHosts []string
	var curVersions []string
	var updateVersions []string
	var processed map[string]bool = map[string]bool{}
	for host, hostData := range report.Results {
		majorVersion, minorVersion := getMajorMinorVersion(hostData.Data.ServerVersionNum)
		if _, vok := processed[minorVersion]; vok {
			// version already checked
			continue
		}
		ver, ok := config.Versions[majorVersion]
		if !ok {
			result.AppendConclusion(A002_NOT_SUPPORTED_VERSION, MSG_NOT_SUPPORTED_VERSION_CONCLUSION, majorVersion, ver.FinalRelease)
			result.AppendRecommendation(A002_NOT_SUPPORTED_VERSION, MSG_NOT_SUPPORTED_VERSION_RECOMMENDATION, majorVersion)
			result.P1 = true
			continue
		}
		sort.Ints(ver.MinorVersions)
		lastVersion := ver.MinorVersions[len(ver.MinorVersions)-1]
		intMinorVersion, _ := strconv.Atoi(minorVersion)
		if intMinorVersion >= lastVersion {
			result.AppendConclusion(A002_LAST_MINOR_VERSION, MSG_LAST_MINOR_VERSION_CONCLUSION,
				majorVersion+"."+minorVersion, majorVersion)
			processed[minorVersion] = true
		} else {
			updateHosts = append(updateHosts, host)
			curVersions = append(curVersions, majorVersion+"."+minorVersion)
			updateVersions = append(updateVersions, majorVersion+"."+strconv.Itoa(lastVersion))
		}
	}
	curVersions = checkup.GetUniques(curVersions)
	if len(curVersions) > 0 {
		result.AppendConclusion(A002_NOT_LAST_MINOR_VERSION, english.PluralWord(len(curVersions),
			MSG_NOT_LAST_MINOR_VERSION_CONCLUSION_1, MSG_NOT_LAST_MINOR_VERSION_CONCLUSION_N),
			strings.Join(curVersions, "`, `"), updateVersions[0])
		result.AppendRecommendation(A002_NOT_LAST_MINOR_VERSION, MSG_NOT_LAST_MINOR_VERSION_RECOMMENDATION, updateVersions[0])
		result.P2 = true
	}
	return result
}

func A002PreprocessReportData(data map[string]interface{}, config config.Config) {
	report, err := A002LoadReportData(data["source_path_full"].(string))

	if err != nil {
		return
	}

	result := A002Process(report, config)
	if len(result.Recommendations) > 0 && !checkup.ResultInList(result.Recommendations, A002_GENERAL_INFO_FULL) {
		result.AppendRecommendation(A002_GENERAL_INFO_OFFICIAL, MSG_GENERAL_RECOMMENDATION_1)
	}
	// update data and file
	checkup.SaveReportResult(data, result)
}

func A002LoadReportData(filePath string) (A002Report, error) {
	var report A002Report
	jsonRaw := checkup.LoadRawJsonReport(filePath)

	if !checkup.CheckUnmarshalResult(json.Unmarshal(jsonRaw, &report)) {
		return report, fmt.Errorf("Unable to load A002 report.")
	}

	return report, nil
}

func A002Process(report A002Report, config config.Config) checkup.ReportResult {
	var result checkup.ReportResult
	A002PrepareVersionInfo()
	result = A002CheckAllVersionsIsSame(report, result)
	result = A002CheckMajorVersions(report, config, result)
	result = A002CheckMinorVersions(report, config, result)
	return result
}
