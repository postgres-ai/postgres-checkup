/*
Postgres Healt Reporter

2018-2019 © Dmitry Udalov dmius@postgres.ai
2018-2019 © Postgres.ai

Perform a generation of Markdown report based on JSON results of postgres-checkup
Usage:
pghrep --checkdata=file:///path_to_check_results.json --outdir=/home/results
*/
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"text/template"

	"./checkup/a002"
	"./checkup/a006"
	"./checkup/a008"
	"./checkup/cfg"
	"./checkup/f001"
	"./checkup/f002"
	"./checkup/f004"
	"./checkup/f005"
	"./checkup/f008"
	"./checkup/g001"
	"./checkup/g002"
	"./checkup/h001"
    "./checkup/k000"
    "./checkup/l003"

	"./log"
	"./orderedmap"
	"./pyraconv"
)

var DEBUG bool = false

// Prepropess file paths
// Allow absulute and relative (of pwd) paths with or wothout file:// prefix
// Return absoulute path of file
func GetFilePath(name string) string {
	filePath := name
	// remove file:// prefix
	if strings.HasPrefix(strings.ToLower(filePath), "file://") {
		filePath = strings.Replace(filePath, "file://", "", 1)
	}
	if strings.HasPrefix(strings.ToLower(filePath), "/") {
		// absoulute path will use as is
		return filePath
	} else {
		// for relative path will combine with current path
		curDir, err := os.Getwd()
		if err != nil {
			log.Dbg("Can't determine current path")
		}
		if strings.HasSuffix(strings.ToLower(curDir), "/") {
			filePath = curDir + filePath
		} else {
			filePath = curDir + "/" + filePath
		}
		return filePath
	}
}

// Check file exists
// Allow absulute and relative (of pwd) paths with or wothout file:// prefix
// Return boolean value
func FileExists(name string) bool {
	filePath := GetFilePath(name)
	if _, err := os.Stat(filePath); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

// Parse json data from string to map
// Return map[string]interface{}
func ParseJson(jsonData string) map[string]interface{} {
	orderedData := orderedmap.New()
	if err := json.Unmarshal([]byte(jsonData), &orderedData); err != nil {
		log.Err("Can't parse json data:", err)
		return nil
	} else {
		dt := orderedData.ToInterfaceArray()
		return dt
	}
}

// Load json data from file by path
// Return map[string]interface{}
func LoadJsonFile(filePath string) map[string]interface{} {
	if FileExists(filePath) {
		fileContent, err := ioutil.ReadFile(GetFilePath(filePath)) // just pass the file name
		if err != nil {
			log.Err("Can't read file: ", filePath, err)
			return nil
		}
		return ParseJson(string(fileContent))
	}
	return nil
}

// Load data dependencies
func loadDependencies(data map[string]interface{}) {
	dep := data["dependencies"]
	dependencies := dep.(map[string]interface{})
	for key, value := range dependencies {
		depData := LoadJsonFile(pyraconv.ToString(value))
		dependencies[key] = depData
	}
}

// Load report templates from files
func loadTemplates() *template.Template {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Dbg("Can't determine current path")
	}

	var templates *template.Template
	var allFiles []string
	files, err := ioutil.ReadDir(dir + "/../templates")
	if err != nil {
		fmt.Println(err)
	}
	for _, file := range files {
		fileName := file.Name()
		if strings.HasSuffix(fileName, ".tpl") {
			allFiles = append(allFiles, dir+"/../templates/"+fileName)
		}
	}

	tplFuncMap := make(template.FuncMap)
	tplFuncMap["Split"] = Split
	tplFuncMap["Trim"] = Trim
	tplFuncMap["Replace"] = Replace
	tplFuncMap["Code"] = Code
	tplFuncMap["Nobr"] = Nobr
	tplFuncMap["Br"] = Br
	tplFuncMap["ByteFormat"] = ByteFormat
	tplFuncMap["UnitValue"] = UnitValue
	tplFuncMap["RawIntUnitValue"] = RawIntUnitValue
	tplFuncMap["RoundUp"] = Round
	tplFuncMap["LimitStr"] = LimitStr
	tplFuncMap["Add"] = Add
	tplFuncMap["Sub"] = Sub
	tplFuncMap["Mul"] = Mul
	tplFuncMap["Div"] = Div
	tplFuncMap["NumFormat"] = NumFormat
	tplFuncMap["MsFormat"] = MsFormat
	tplFuncMap["DtFormat"] = DtFormat
	tplFuncMap["RawIntFormat"] = RawIntFormat
	tplFuncMap["RawFloatFormat"] = RawFloatFormat
	tplFuncMap["Int"] = Int
	tplFuncMap["EscapeQuery"] = EscapeQuery
	tplFuncMap["WordWrap"] = WordWrap

	templates, err = template.New("").Funcs(tplFuncMap).ParseFiles(allFiles...)
	if err != nil {
		log.Fatal("Can't load templates", err)
		return nil
	}

	return templates
}

// Prepare raw json data for every host
func getRawData(data map[string]interface{}) {
	// for every host get data
	var rawData []interface{}
	hosts := pyraconv.ToInterfaceMap(data["hosts"])
	log.Dbg("Data hosts: ", hosts)
	results := pyraconv.ToInterfaceMap(data["results"])
	masterName := pyraconv.ToString(hosts["master"])
	masterResults := pyraconv.ToInterfaceMap(results[masterName])
	masterData := pyraconv.ToInterfaceMap(masterResults["data"])
	masterJson, err := json.Marshal(masterData)
	if err == nil {
		masterItem := make(map[string]interface{})
		masterItem["host"] = masterName
		masterItem["data"] = string(masterJson)
		rawData = append(rawData, masterItem)
	}
	replicas := pyraconv.ToStringArray(hosts["replicas"])
	for _, host := range replicas {
		hostResults := pyraconv.ToInterfaceMap(results[host])
		hostData := pyraconv.ToInterfaceMap(hostResults["data"])
		hostJson, err := json.Marshal(hostData)
		if err == nil {
			hostItem := make(map[string]interface{})
			hostItem["host"] = host
			hostItem["data"] = string(hostJson)
			rawData = append(rawData, hostItem)
		}
	}
	data["rawData"] = rawData
}

/*
Generate MD reports by given check Id
CheckId can be either ID of concrete check (e.g. H003) or represent the whole category (e.g. K000)
*/
func generateMdReports(checkId string, reportData map[string]interface{}, outputDir string) bool {
	category := checkId[0:1]
	checkNum, err := strconv.ParseInt(checkId[1:4], 10, 64)

	reportPrefix := ""
	if checkNum != 0 {
		reportPrefix = checkId // specified check given
	} else {
		reportPrefix = category // category given
	}

	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Err(err)
		return false
	}
	files, err := ioutil.ReadDir(dir + "/../templates")
	if err != nil {
		log.Err(err)
		return false
	}
	for _, file := range files {
		fileName := file.Name()
		if strings.HasPrefix(fileName, reportPrefix) && strings.HasSuffix(fileName, ".tpl") {
			curCheckId := fileName[0:4]
			outputFileName := strings.Replace(fileName, ".tpl", ".md", -1)
			reportData["checkId"] = curCheckId
			if !generateMdReport(curCheckId, outputFileName, reportData, outputDir) {
				log.Err("Can't generate report " + outputFileName + " based on " + checkId + " json data")
				return false
			}
		}
	}

	return true
}

// Generate md report (file) on base of reportData and save them to file in outputDir
func generateMdReport(checkId string, reportFilename string, reportData map[string]interface{}, outputDir string) bool {
	var outputFileName string
	if len(reportFilename) > 0 {
		outputFileName = reportFilename
	} else {
		outputFileName = checkId + ".md"
	}
	if strings.HasSuffix(strings.ToLower(outputDir), "/") {
		outputFileName = outputDir + outputFileName
	} else {
		outputFileName = outputDir + "/" + outputFileName
	}
	_, err := filepath.Abs(filepath.Dir(os.Args[0]))
	f, err := os.OpenFile(outputFileName, os.O_CREATE|os.O_RDWR, 0777)
	if err != nil {
		log.Err("Can't create report file", err)
		return false
	}
	defer f.Close()
	f.Truncate(0)

	templates := loadTemplates()
	if templates == nil {
		log.Fatal("Can't load template")
	}
	reportFileName := checkId + ".tpl"
	reporTpl := templates.Lookup(reportFileName)
	data := reportData
	if reporTpl == nil {
		log.Err("Template " + checkId + ".tpl not found.")
		getRawData(data)
		reportFileName = "raw.tpl"
		reporTpl = templates.Lookup(reportFileName)
	}
	err = reporTpl.ExecuteTemplate(f, reportFileName, data)
	if err != nil {
		log.Err("Template execute error is", err)
		defer os.Remove(outputFileName)
		return false
	} else {
		return true
	}
}

// Sort hosts on master and replicas by role and index.
// Return map {"master":name string, "replicas":[replica1 string, replica2 string]}
func determineMasterReplica(data map[string]interface{}) {
	hostRoles := make(map[string]interface{})
	var sortedReplicas []string
	replicas := make(map[int]string)
	nodes_json := pyraconv.ToInterfaceMap(data["last_nodes_json"])
	hosts := pyraconv.ToInterfaceMap(nodes_json["hosts"])
	hostRoles["master"] = nil
	for host, value := range hosts {
		hostData := pyraconv.ToInterfaceMap(value)
		if hostData["role"] == "master" {
			hostRoles["master"] = host
		} else {
			if host != "_keys" {
				index, _ := strconv.Atoi(pyraconv.ToString(hostData["index"]))
				replicas[index] = host
			}
		}
	}
	var keys []int
	for k := range replicas {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	for _, k := range keys {
		sortedReplicas = append(sortedReplicas, replicas[k])
	}

	hostRoles["replicas"] = sortedReplicas
	data["hosts"] = hostRoles
}

/*
  Replace master on replica#1 if master not defined
*/
func reorderHosts(data map[string]interface{}) {
	hosts := pyraconv.ToInterfaceMap(data["hosts"])
	masterHost := pyraconv.ToString(hosts["master"])
	replicaHosts := pyraconv.ToStringArray(hosts["replicas"])
	var allHosts []string
	if hosts["master"] != nil {
		allHosts = append(allHosts, masterHost)
	}
	for _, replicaHost := range replicaHosts {
		allHosts = append(allHosts, replicaHost)
	}
	if len(allHosts) == 0 {
		return
	}
	// check host data
	var hostsWithData []string
	for _, host := range allHosts {
		results := pyraconv.ToInterfaceMap(data["results"])
		hostData, ok := results[host]
		if ok && hostData != nil {
			hostsWithData = append(hostsWithData, host)
		}
	}
	master := hostsWithData[0]
	var replicas []string
	replicas = append(replicas, hostsWithData[1:]...)
	reorderedHosts := make(map[string]interface{})
	reorderedHosts["master"] = master
	reorderedHosts["replicas"] = replicas
	data["reorderedHosts"] = reorderedHosts
}

func main() {
	// get input data checkId, checkData
	var checkId string
	var checkData string
	var resultData map[string]interface{}
	checkDataPtr := flag.String("checkdata", "", "an filepath to json report")
	outDirPtr := flag.String("outdir", "", "an directory where report need save")
	debugPtr := flag.Int("debug", 0, "enable debug mode (must be defined 1 or 0 (default))")
	flag.Parse()
	checkData = *checkDataPtr
	LISTLIMIT, rlerr := strconv.Atoi(os.Getenv("LISTLIMIT"))
	if rlerr != nil {
		LISTLIMIT = 50
	}

	if *debugPtr == 1 {
		DEBUG = true
	}

	if FileExists(checkData) {
		resultData = LoadJsonFile(checkData)
		if resultData == nil {
			log.Fatal("ERROR: File given by --checkdata content wrong json data.")
			return
		}
		resultData["source_path_full"] = checkData
		resultData["source_path_parts"] = strings.Split(checkData, string(os.PathSeparator))
	} else {
		log.Err("ERROR: File given by --checkdata not found")
		return
	}

	if resultData != nil {
		checkId = pyraconv.ToString(resultData["checkId"])
	} else {
		log.Fatal("ERROR: Content given by --checkdata is wrong json content.")
	}

	checkId = strings.ToUpper(checkId)
	loadDependencies(resultData)
	determineMasterReplica(resultData)
	reorderHosts(resultData)

	config := cfg.NewConfig()

	err := preprocessReportData(checkId, config, resultData)
	if err != nil {
		log.Fatal(err)
	}

	resultData["LISTLIMIT"] = LISTLIMIT
	var outputDir string
	if len(*outDirPtr) == 0 {
		outputDir = "./"
	} else {
		outputDir = *outDirPtr
	}
	reportDone := generateMdReports(checkId, resultData, outputDir)
	if !reportDone {
		log.Fatal("Cannot generate report. Data file or template is wrong.")
	}
}

func preprocessReportData(checkId string, config cfg.Config,
	data map[string]interface{}) error {
	switch checkId {
	case "A002":
		// Try to load actual Postgres versions.
		err := config.LoadVersions()
		if err != nil {
			// TODO(anatoly): Add warnings to the beginning of MD report.
			log.Err("Cannot load latest Postgres versions. Recommendations may be inaccurate.", err)
		}

		a002.A002PreprocessReportData(data, config)
	case "A006":
		a006.A006PreprocessReportData(data)
	case "A008":
		a008.A008PreprocessReportData(data)
	case "H001":
		h001.H001PreprocessReportData(data)
	case "F001":
		f001.F001PreprocessReportData(data)
	case "F002":
		f002.F002PreprocessReportData(data)
	case "F004":
		f004.F004PreprocessReportData(data)
	case "F005":
		f005.F005PreprocessReportData(data)
	case "F008":
		f008.F008PreprocessReportData(data)
	case "G001":
		g001.G001PreprocessReportData(data)
	case "G002":
		g002.G002PreprocessReportData(data)
	case "K000":
		k000.K000PreprocessReportData(data)
	case "L003":
		l003.L003PreprocessReportData(data)
	}

	return nil
}
