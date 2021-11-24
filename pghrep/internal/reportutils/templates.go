package reportutils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"text/template"

	"gitlab.com/postgres-ai/postgres-checkup/pghrep/internal/checkup/a002"
	"gitlab.com/postgres-ai/postgres-checkup/pghrep/internal/checkup/a006"
	"gitlab.com/postgres-ai/postgres-checkup/pghrep/internal/checkup/a008"
	"gitlab.com/postgres-ai/postgres-checkup/pghrep/internal/checkup/cfg"
	"gitlab.com/postgres-ai/postgres-checkup/pghrep/internal/checkup/f001"
	"gitlab.com/postgres-ai/postgres-checkup/pghrep/internal/checkup/f002"
	"gitlab.com/postgres-ai/postgres-checkup/pghrep/internal/checkup/f004"
	"gitlab.com/postgres-ai/postgres-checkup/pghrep/internal/checkup/f005"
	"gitlab.com/postgres-ai/postgres-checkup/pghrep/internal/checkup/f008"
	"gitlab.com/postgres-ai/postgres-checkup/pghrep/internal/checkup/g001"
	"gitlab.com/postgres-ai/postgres-checkup/pghrep/internal/checkup/g002"
	"gitlab.com/postgres-ai/postgres-checkup/pghrep/internal/checkup/h001"
	"gitlab.com/postgres-ai/postgres-checkup/pghrep/internal/checkup/h002"
	"gitlab.com/postgres-ai/postgres-checkup/pghrep/internal/checkup/h004"
	"gitlab.com/postgres-ai/postgres-checkup/pghrep/internal/checkup/k000"
	"gitlab.com/postgres-ai/postgres-checkup/pghrep/internal/checkup/l003"
	"gitlab.com/postgres-ai/postgres-checkup/pghrep/internal/helpers"
	"gitlab.com/postgres-ai/postgres-checkup/pghrep/internal/log"
	"gitlab.com/postgres-ai/postgres-checkup/pghrep/internal/orderedmap"
	"gitlab.com/postgres-ai/postgres-checkup/pghrep/internal/pyraconv"
)

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
	if helpers.FileExists(filePath) {
		fileContent, err := ioutil.ReadFile(helpers.GetFilePath(filePath)) // just pass the file name
		if err != nil {
			log.Err("Can't read file:", filePath, err)
			return nil
		}

		return ParseJson(string(fileContent))
	}

	return nil
}

// Load data dependencies
func LoadDependencies(data map[string]interface{}) {
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
	tplFuncMap["AddSession"] = Add
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
	log.Dbg("Data hosts:", hosts)
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
func GenerateMdReports(checkId string, reportData map[string]interface{}, outputDir string) bool {
	category := checkId[0:1]
	reportPrefix := ""

	checkNum, err := strconv.ParseInt(checkId[1:4], 10, 64)
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
		log.Err("Template " + checkId + ".tpl not found")
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
func DetermineMasterReplica(data map[string]interface{}) {
	hostRoles := make(map[string]interface{})
	var sortedReplicas []string
	var keys []int

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
func ReorderHosts(data map[string]interface{}) error {
	hosts := pyraconv.ToInterfaceMap(data["hosts"])
	masterHost := pyraconv.ToString(hosts["master"])
	replicaHosts := pyraconv.ToStringArray(hosts["replicas"])
	var allHosts, hostsWithData []string

	if hosts["master"] != nil {
		allHosts = append(allHosts, masterHost)
	}

	for _, replicaHost := range replicaHosts {
		allHosts = append(allHosts, replicaHost)
	}

	if len(allHosts) == 0 {
		data["reorderedHosts"] = []string{}
		return fmt.Errorf("No hosts")
	}

	// check host data
	for _, host := range allHosts {
		results := pyraconv.ToInterfaceMap(data["results"])
		hostData, ok := results[host]
		if ok && hostData != nil {
			hostsWithData = append(hostsWithData, host)
		}
	}

	if len(hostsWithData) == 0 {
		data["reorderedHosts"] = []string{}
		return fmt.Errorf("No hosts with data")
	}

	master := hostsWithData[0]
	var replicas []string
	replicas = append(replicas, hostsWithData[1:]...)
	reorderedHosts := make(map[string]interface{})
	reorderedHosts["master"] = master
	reorderedHosts["replicas"] = replicas
	data["reorderedHosts"] = reorderedHosts

	return nil
}

func PreprocessReportData(checkId string, config cfg.Config,
	data map[string]interface{}) error {
	switch strings.ToUpper(checkId) {
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
	case "H002":
		h002.H002PreprocessReportData(data)
	case "H004":
		h004.H004PreprocessReportData(data)
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
