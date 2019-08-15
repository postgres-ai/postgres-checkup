package upload

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"../checkup"
	"../log"
	"gopkg.in/yaml.v2"
)

type UploadConfig struct {
	ApiUrl string `yaml:"apiUrl"`
}

var uploadConfig = UploadConfig{}

func loadConfig() error {
	name := "upload.yaml"
	configPath := getConfigPath(name)
	if _, err := os.Stat(configPath); err != nil {
		if os.IsNotExist(err) {
			configPath = getConfigPath("default.upload.yaml")
			name = "default.upload.yaml"
		}
	}

	b, err := ioutil.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("Error loading %s config file. Default values are used.", name)
	}

	err = yaml.Unmarshal(b, &uploadConfig)
	if err != nil {
		return fmt.Errorf("Error parsing %s config. Default values are used.", name)
	}

	log.Dbg("Config loaded", uploadConfig)
	return nil
}

func getConfigPath(name string) string {
	bindir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	dir, _ := filepath.Abs(filepath.Dir(bindir))
	path := dir + string(os.PathSeparator) + "config" + string(os.PathSeparator) + name
	return path
}

func UploadReport(token string, project string, path string) error {
	cerr := loadConfig()
	if cerr != nil {
		return fmt.Errorf("Upload configuration not found")
	}
	// enumerate files
	var files []string
	var err error
	files, err = ScanPath(path, files)
	if err != nil {
		return err
	}

	if len(files) == 0 {
		return fmt.Errorf("Files to upload not found")
	}

	// create report
	reportId, cerr := CreateReport(token, project, path)
	if cerr != nil {
		return cerr
	}

	processed := 0
	for _, f := range files {
		uerr := UploadReportFile(token, reportId, f)
		if uerr == nil {
			processed++
		}
	}

	log.Msg("Uploaded", processed, "files from", len(files), "of report.")

	return nil
}

func ScanPath(path string, files []string) ([]string, error) {
	result := files
	dirFiles, err := ioutil.ReadDir(path + string(os.PathSeparator))
	if err != nil {
		return nil, err
	}

	for _, f := range dirFiles {
		if f.IsDir() {
			var sderr error
			result, sderr = ScanPath(path+string(os.PathSeparator)+f.Name(), result)
			if sderr != nil {
				log.Dbg(sderr)
			}
		} else {
			result = append(result, path+string(os.PathSeparator)+f.Name())
		}
	}

	return result, nil
}

func GetReportEpoch(path string) (string, error) {
	nodesJsonPath := path + string(os.PathSeparator) + "nodes.json"
	if _, err := os.Stat(nodesJsonPath); err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("File nodes.json not found")
		}
	}

	jsonRaw := checkup.LoadRawJsonReport(nodesJsonPath)
	var nodesJsonData checkup.ReportLastNodes

	if !checkup.CheckUnmarshalResult(json.Unmarshal(jsonRaw, &nodesJsonData)) {
		return "", fmt.Errorf("Unable to load nodes.json data.")
	}

	return nodesJsonData.LastCheck.Epoch, nil
}

func CreateReport(token string, project string, path string) (int64, error) {
	epoch, err := GetReportEpoch(path)

	if err != nil {
		return -1, err
	}

	requestData := map[string]interface{}{
		"access_token": token,
		"project":      project,
		"epoch":        epoch,
	}

	response, rerr := MakeRequest("post_checkup_report", requestData)
	if rerr != nil {
		return -1, rerr
	}

	var intId int64 = 0
	var iok bool = false
	floatId, fok := response["report_id"].(float64)
	if !fok {
		intId, iok = response["report_id"].(int64)
		if iok {
			return intId, nil
		}
	} else {
		return int64(floatId), nil
	}

	if msg, mok := response["message"]; mok {
		return -1, fmt.Errorf("%s", msg)
	} else {
		return -1, fmt.Errorf("Cannot create report.")
	}
}

func MakeRequest(endpoint string, requestData map[string]interface{}) (map[string]interface{}, error) {
	bytesRepresentation, merr := json.Marshal(requestData)
	if merr != nil {
		return nil, merr
	}

	resp, err := http.Post(uploadConfig.ApiUrl+endpoint, "application/json", bytes.NewBuffer(bytesRepresentation))
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}

	json.NewDecoder(resp.Body).Decode(&result)

	return result, nil
}

func UploadReportFile(token string, reportId int64, path string) error {
	fileType := strings.ToLower(strings.Replace(filepath.Ext(path), ".", "", -1))
	fileName := filepath.Base(path)
	checkId := ""

	if fileType != "json" && fileType != "sql" && fileType != "md" && fileType != "html" {
		return fmt.Errorf("Unsupported file type.")
	}

	if string(fileName[4:5]) == "_" {
		checkId = string(fileName[0:4])
	}

	// read file
	data, rerr := ioutil.ReadFile(path) // just pass the file name
	if rerr != nil {
		return fmt.Errorf("Cannot read file.")
	}

	strData := string(data) // convert content to a 'string'

	requestData := map[string]interface{}{
		"access_token":      token,
		"checkup_report_id": reportId,
		"check_id":          checkId,
		"filename":          fileName,
		"data":              strData,
		"type":              fileType,
	}

	response, uerr := MakeRequest("post_checkup_report_chunk", requestData)
	if uerr != nil {
		if msg, mok := response["hint"]; mok {
			return fmt.Errorf("%s", msg)
		} else {
			return fmt.Errorf("Cannot upload file. %s", uerr)
		}
	}

	return nil
}
