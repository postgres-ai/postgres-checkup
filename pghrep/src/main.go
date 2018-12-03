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
    "path/filepath"
    "./pyraconv"
    "log"
    "text/template"
    "sort"
    "strconv"
)

var DEBUG bool = false

func dbg(v ...interface{}) {
    if DEBUG {
        message := ""
        for _, value := range v {
            message = message + " " + pyraconv.ToString(value)
        }
        log.Println(">>> DEBUG:", message)
    }
}

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
            dbg("Can't determine current path")
        }
        if strings.HasSuffix(strings.ToLower(curDir), "/") {
            filePath = curDir + filePath
        } else {
            filePath = curDir + "/" + filePath
        }
        return filePath
    }
}

// Exists reports whether the named file or directory exists.
func FileExists(name string) bool {
    filePath := GetFilePath(name)
    dbg("File path", filePath)
    if _, err := os.Stat(filePath); err != nil {
        if os.IsNotExist(err) {
            return false
        }
    }
    return true
}

func ParseJson(jsonData string) map[string]interface{} {
    var data map[string]interface{}
    if err := json.Unmarshal([]byte(jsonData), &data); err != nil {
        dbg("Can't parse json data:", err)
        return nil
    } else {
        return data
    }
}

func LoadJsonFile(filePath string) map[string]interface{} {
    if FileExists(filePath) {
        fileContent, err := ioutil.ReadFile(GetFilePath(filePath)) // just pass the file name
        if err != nil {
            log.Println("Can't read file: ", filePath, err)
            return nil
        }
        return ParseJson(string(fileContent))
    }
    return nil
}

func loadDependencies(data map[string]interface{}) {
    dep := data["dependencies"]
    dependencies := dep.(map[string]interface{})
    for key, value := range dependencies {
        depData := LoadJsonFile(pyraconv.ToString(value))
        dependencies[key] = depData
    }
}

func loadTemplates() *template.Template {
    dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
    if err != nil {
        dbg("Can't determine current path")
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

    templates, err = template.ParseFiles(allFiles...)
    if err != nil {
        dbg("Can't load templates", err)
        return nil
    }

    return templates
}

func generateMdReport(checkId string, reportData map[string]interface{}, outputDir string) bool{
    var outputFileName string
    if strings.HasSuffix(strings.ToLower(outputDir), "/") {
        outputFileName = outputDir + checkId + ".md"
    } else {
        outputFileName = outputDir + "/" + checkId + ".md"
    }
    _, err := filepath.Abs(filepath.Dir(os.Args[0]))
    f, err := os.OpenFile(outputFileName, os.O_CREATE | os.O_RDWR, 0777)
    if err != nil {
        dbg("Can't create report file", err)
        return false
    }
    defer f.Close()
    f.Truncate(0)

    templates := loadTemplates()
    if templates == nil {
        log.Fatal("Can't load template")
    }
    dbg("Find", checkId + ".tpl")
    reporTpl := templates.Lookup(checkId + ".tpl")
    if reporTpl != nil {
        err = reporTpl.ExecuteTemplate(f, checkId + ".tpl", reportData)
        if err != nil {
            dbg("Template execute error is", err)
            defer os.Remove(outputFileName)
            return false
        } else {
            return true
        }
    } else {
        dbg("Template " + checkId + ".tpl not found.")
        defer os.Remove(outputFileName)
        return false
    }
}

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
            index, _ := strconv.Atoi(pyraconv.ToString(hostData["index"]))
            replicas[index] = host
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

    if FileExists(checkData) {
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
        dbg("Cannot find and load plugin.", err)
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

    reportDone := generateMdReport(checkId, reportData, outputDir)
    if ! reportDone  {
        log.Fatal("Cannot generate report. Data file or template is wrong.")
    }
}