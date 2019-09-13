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
)

func UploadReport(apiUrl string, token string, project string, path string) error {
	// enumerate files
	var files []string
	var reportPath = path
	var err error

	if !strings.HasSuffix(reportPath, string(os.PathSeparator)) {
		reportPath = reportPath + string(os.PathSeparator)
	}

	epoch, dir, err := GetReportLastCheckData(reportPath)
	if err != nil {
		return fmt.Errorf("Cannot read report information (epoch, dir) from nodes.json. %s", err)
	}

	if dir == "" {
		return fmt.Errorf("Empty report path read from nodes.json")
	}

	files = append(files, reportPath+"nodes.json")
	files, err = ScanPath(reportPath+"json_reports"+string(os.PathSeparator)+dir, files)
	if err != nil {
		return err
	}
	files, err = ScanPath(reportPath+"md_reports"+string(os.PathSeparator)+dir, files)
	if err != nil {
		return err
	}

	if len(files) == 0 {
		return fmt.Errorf("There are no files to upload (report directories: "+
			"'%s/json_files/' and '%s/md_reports')", reportPath, reportPath)
	}

	reportId, cerr := CreateReport(apiUrl, token, project, epoch)
	if cerr != nil {
		return cerr
	}

	log.Dbg("Report created: ", reportId)

	processed := 0
	for _, f := range files {
		log.Dbg(fmt.Sprintf("File '%s' found in project artifcats directory.", f))
		fileType := strings.ToLower(strings.Replace(filepath.Ext(f), ".", "", -1))
		if fileType != "json" && fileType != "sql" && fileType != "md" && fileType != "html" {
			log.Dbg(fmt.Sprintf("File: '%s' skipped (Not json, sql, md or html).", f))
			continue
		}
		log.Dbg(fmt.Sprintf("Try upload file: '%s'.", f))
		uerr := UploadReportFile(apiUrl, token, reportId, f)
		if uerr == nil {
			processed++
			log.Dbg("File: " + f + " uploaded without errors.")
		} else {
			log.Err(fmt.Sprintf("Cannot upload file '%s'. %s", f, uerr))
		}
	}

	log.Msg("Uploaded", processed, "files from", "'"+path+"'.")

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

func GetReportLastCheckData(path string) (string, string, error) {
	nodesJsonPath := path + "nodes.json"
	if _, err := os.Stat(nodesJsonPath); err != nil {
		if os.IsNotExist(err) {
			return "", "", fmt.Errorf("nodes.json is not found")
		}
	}

	jsonRaw := checkup.LoadRawJsonReport(nodesJsonPath)
	var nodesJsonData checkup.ReportLastNodes

	if !checkup.CheckUnmarshalResult(json.Unmarshal(jsonRaw, &nodesJsonData)) {
		return "", "", fmt.Errorf("Unable to load nodes.json data.")
	}

	return nodesJsonData.LastCheck.Epoch, nodesJsonData.LastCheck.Dir, nil
}

func CreateReport(apiUrl string, token string, project string, epoch string) (int64, error) {
	requestData := map[string]interface{}{
		"access_token": token,
		"project":      project,
		"epoch":        epoch,
	}

	response, rerr := MakeRequest(apiUrl, "/rpc/checkup_report_create", requestData)
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

func MakeRequest(apiUrl string, endpoint string, requestData map[string]interface{}) (map[string]interface{}, error) {
	bytesRepresentation, merr := json.Marshal(requestData)
	if merr != nil {
		return nil, merr
	}

	resp, err := http.Post(apiUrl+endpoint, "application/json",
		bytes.NewBuffer(bytesRepresentation))
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}

	json.NewDecoder(resp.Body).Decode(&result)

	return result, nil
}

func UploadReportFile(apiUrl string, token string, reportId int64, path string) error {
	fileType := strings.ToLower(strings.Replace(filepath.Ext(path), ".", "", -1))
	fileName := filepath.Base(path)
	checkId := ""

	if string(fileName[4:5]) == "_" {
		checkId = string(fileName[0:4])
	}

	data, rerr := ioutil.ReadFile(path)
	if rerr != nil {
		return fmt.Errorf("Cannot read file.")
	}

	strData := string(data)

	requestData := map[string]interface{}{
		"access_token":      token,
		"checkup_report_id": reportId,
		"check_id":          checkId,
		"filename":          fileName,
		"data":              strData,
		"type":              fileType,
	}

	response, uerr := MakeRequest(apiUrl, "/rpc/checkup_report_file_post", requestData)
	if uerr != nil {
		return fmt.Errorf("%s", uerr)
	}

	if msg, mok := response["message"]; mok {
		return fmt.Errorf("%s", msg)
	}

	if _, rok := response["report_chunck_id"]; !rok {
		log.Dbg("Response for uploading file '%s': ", response)
		return fmt.Errorf("Response for uploading file '%s' do not content file id.", fileName)
	}

	return nil
}
