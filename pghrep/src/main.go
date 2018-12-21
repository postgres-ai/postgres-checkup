/*
Postgres Healt Reporter

2018 © Dmitry Udalov dmius@postgres.ai
2018 © Postgres.ai

Perform a generation *md reports on base of results health checks
Usage: 
pghrep --checkdata=file:///path_to_check_results.json --outdir=/home/results
*/
package main

import (
    "fmt"
    "os"
    "flag"
    "strings"
    "encoding/json"
    "io/ioutil"
    "path"
    "path/filepath"
    "./pyraconv"
    "log"
    "text/template"
    "sort"
    "strconv"
    "./orderedmap"
    "./fmtutils"
)

var DEBUG bool = false

// Output debug message
func Dbg(v ...interface{}) {
    if DEBUG {
        message := ""
        for _, value := range v {
            message = message + " " + pyraconv.ToString(value)
        }
        log.Println(">>> DEBUG:", message)
    }
}

// Output debug message
func Err(v ...interface{}) {
    message := ""
    for _, value := range v {
        message = message + " " + pyraconv.ToString(value)
    }
    log.Println(">>> ERROR:", message)
}


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
            Dbg("Can't determine current path")
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
        Err("Can't parse json data:", err)
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
            Err("Can't read file: ", filePath, err)
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
        Dbg("Can't determine current path")
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
            allFiles = append(allFiles, dir + "/../templates/" + fileName)
        }
    }

    tplFuncMap :=  make(template.FuncMap)
    tplFuncMap["Split"] = Split
    tplFuncMap["Trim"] = Trim
    tplFuncMap["Code"] = Code
    tplFuncMap["Nobr"] = Nobr
    tplFuncMap["Br"] = Br
    tplFuncMap["ByteFormat"] = fmtutils.ByteFormat
    tplFuncMap["UnitValue"] = UnitValue
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
    Dbg("Data hosts: ", hosts)
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

// Generate md report (file) on base of reportData and save them to file in outputDir
func generateMdReport(checkId string, reportFilename string, reportData map[string]interface{}, outputDir string) bool{
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
    f, err := os.OpenFile(outputFileName, os.O_CREATE | os.O_RDWR, 0777)
    if err != nil {
        Err("Can't create report file", err)
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
        Err("Template " + checkId + ".tpl not found.")
        getRawData(data)
        reportFileName = "raw.tpl"
        reporTpl = templates.Lookup(reportFileName)
    }
    err = reporTpl.ExecuteTemplate(f, reportFileName, data)
    if err != nil {
        Err("Template execute error is", err)
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
    hosts := pyraconv.ToInterfaceMap(nodes_json["hosts"]);
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

func main() {
    // get input data checkId, checkData
    var checkId string
    var checkData string
    var resultData map[string]interface{}
    checkDataPtr := flag.String("checkdata", "", "an filepath to json report")
    outDirPtr := flag.String("outdir", "", "an directory where report need save")
    debugPtr := flag.Int("debug", 0, "enable debug mode (must be defined 1 or 0 (default))")
    flag.Parse()
    checkData=*checkDataPtr

    if *debugPtr == 1  {
        DEBUG = true
    }

    reportFilename := ""
    if FileExists(checkData) {
        _, file := path.Split(checkData)
        fmt.Println(file)
        reportFilename = strings.Replace(file, ".json", ".md", -1)

        resultData = LoadJsonFile(checkData)
        if resultData == nil {
            log.Fatal("ERROR: File given by --checkdata content wrong json data.")
            return
        }
    } else {
        log.Println("ERROR: File given by --checkdata not found")
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

    l, err := newLoader()
    if err != nil {
        fmt.Fprintf(os.Stderr, "%v", err)
    }
    defer l.destroy()
    
    var reportData map[string]interface{}
    objectPath, err := l.get(checkId);
    if err != nil {
        Dbg("Cannot find and load plugin.", err)
        reportData = resultData
    } else {
        result, err := l.call(objectPath, resultData)
        if err != nil {
            fmt.Fprintf(os.Stderr, "%v", err)
        }
        bodyBytes, _ := json.Marshal(result)
        json.Unmarshal(bodyBytes, &reportData)
    }

    var outputDir string
    if len(*outDirPtr) == 0 {
        outputDir = "./"
    } else {
        outputDir = *outDirPtr
    }

    reportDone := generateMdReport(checkId, reportFilename, reportData, outputDir)
    if ! reportDone  {
        log.Fatal("Cannot generate report. Data file or template is wrong.")
    }
}