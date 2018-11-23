/* 2018 © Dmitry Udalov dmius@postgres.ai
 2018 © Postgres.ai

Perform a generation *md reports on base of results health checks
Usage: 
pghrep --checkid=XXX --checkdata=file:///path_to_check_results.json --outdir=/home/results
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
)

const (
    DEBUG  = true
)
    
type CheckData struct {
    checkId string
    dependencies interface{}
    hosts []string 
    checkData interface{}
}

type ReportData struct {
    Current string
    Recommended string
    Conclusion string
    Filename string
}


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
    var filePath string
    if strings.HasPrefix(strings.ToLower(name), "file:///") {
        filePath = strings.Replace(name, "file://", "", 1)
        return filePath
    }
    if strings.HasPrefix(strings.ToLower(name), "file://") {
        dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
        if err != nil {
            dbg("Can't determine current path")
        }
        givenPath := strings.Replace(name, "file://", "", 1)
        if strings.HasPrefix(strings.ToLower(givenPath), "/") {
            filePath = dir + givenPath
        } else {
            filePath = dir + "/" + givenPath
        }
        return filePath
    }
    return name
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
    var data interface{}
    if err := json.Unmarshal([]byte(jsonData), &data); err != nil {
        dbg("Can't parse json data:", err)
        return nil
    } else {
        mapData := data.(map[string]interface{})
        return mapData
    }
}

func LoadJsonFile(filePath string) map[string]interface{} {
    if (strings.HasPrefix(strings.ToLower(filePath), "file://") && FileExists(filePath)) {
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
        dbg("Can't load templates")
        return nil
    }

    return templates
}

func generateMdReport(reportData ReportData, outputDir string){
    var outputFileName string
    if strings.HasSuffix(strings.ToLower(outputDir), "/") {
        outputFileName = outputDir + reportData.Filename
    } else {
        outputFileName = outputDir + "/" + reportData.Filename
    }
    _, err := filepath.Abs(filepath.Dir(os.Args[0]))
    f, err := os.OpenFile(outputFileName, os.O_CREATE | os.O_RDWR, 0777)
    if err != nil {
        dbg("Can't create report file", err)
        return
    }
    defer f.Close()   
    f.Truncate(0)
    
    templates := loadTemplates()
    if templates == nil {
        log.Fatal("Can't load template")
    }
    reporTpl := templates.Lookup("report.tpl")
    reporTpl.ExecuteTemplate(f, "report.tpl", reportData)
}

func main() {
    // get input data checkId, checkData
    var checkId string
    var checkData string
    var mapCheckData map[string]interface{}
    checkIdPtr := flag.String("checkid", "", "an check id")
    checkDataPtr := flag.String("checkdata", "", "an report data in json format")
    outDirPtr := flag.String("outdir", "", "an directore where report need save")
    flag.Parse()

    if len(*checkIdPtr) > 0 {
        checkId = *checkIdPtr
    }
    checkData=*checkDataPtr

    if (strings.HasPrefix(strings.ToLower(checkData), "file://") && FileExists(checkData)) {
        mapCheckData = LoadJsonFile(checkData)
        if mapCheckData == nil {
            log.Fatal("ERROR: File given by --checkdata content wrong json data.")
            return
        }
    } else {
        log.Println("ERROR: File given by --checkdata not found, will used as json data", checkData)
        mapCheckData = ParseJson(checkData)
    }

    if mapCheckData != nil {
        checkId = pyraconv.ToString(mapCheckData["checkId"])
    } else {
        log.Fatal("ERROR: Content given by --checkdata is wrong json content.")
    }
    
    if len(checkId) == 0 {
        log.Fatal("ERROR: --checkid value incorrect")
        return
    }

    if len(checkId) == 0 {
        log.Fatal("ERROR: --checkid value incorrect")
        return
    }
    
    checkId = strings.ToLower(checkId)
    loadDependencies(mapCheckData)

    l, err := newLoader()
    if err != nil {
        fmt.Fprintf(os.Stderr, "%v", err)
    }
    defer l.destroy()
    
    objectPath, err := l.get(checkId);
    if err != nil {
        log.Fatal("Cannot find and load plugin.", err)
    }
    result, err := l.call(objectPath, mapCheckData)
    if err != nil {
        fmt.Fprintf(os.Stderr, "%v", err)
    }

    reportData := ReportData{}
    bodyBytes, _ := json.Marshal(result)
    json.Unmarshal(bodyBytes, &reportData)
    if len(result) > 0 && len(reportData.Filename) > 0 {
        var outputDir string
        if len(*outDirPtr) == 0 {
            outputDir = "./"
        } else {
            outputDir = *outDirPtr
        }
        generateMdReport(reportData, outputDir)
    } else {
        log.Fatal("Cannot generate report. Data file or plugin is wrong.")
    }
}
