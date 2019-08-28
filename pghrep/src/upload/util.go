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
	var err error
	files, err = ScanPath(path, files)
	if err != nil {
		return err
	}

	if len(files) == 0 {
		return fmt.Errorf("There are no files to upload (artifacts directory: '%s')", path)
	}

	reportId, cerr := CreateReport(apiUrl, token, project, path)
	if cerr != nil {
		return cerr
	}

	processed := 0
	for _, f := range files {
		fileType := strings.ToLower(strings.Replace(filepath.Ext(f), ".", "", -1))
		if fileType != "json" && fileType != "sql" && fileType != "md" && fileType != "html" {
			continue
		}

		uerr := UploadReportFile(apiUrl, token, reportId, f)
		if uerr == nil {
			processed++
		} else {
			log.Err(fmt.Sprintf("Cannot upload file %s. %s", f, uerr))
		}
	}

  log.Msg("Uploaded", processed, "files from", "'" + path + "'.")

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
			return "", fmt.Errorf("nodes.json is not found")
		}
	}

	jsonRaw := checkup.LoadRawJsonReport(nodesJsonPath)
	var nodesJsonData checkup.ReportLastNodes

	if !checkup.CheckUnmarshalResult(json.Unmarshal(jsonRaw, &nodesJsonData)) {
		return "", fmt.Errorf("Unable to load nodes.json data.")
	}

	return nodesJsonData.LastCheck.Epoch, nil
}

func CreateReport(apiUrl string, token string, project string, path string) (int64, error) {
	epoch, err := GetReportEpoch(path)

	if err != nil {
		return -1, err
	}

	requestData := map[string]interface{}{
		"access_token": token,
		"project":      project,
		"epoch":        epoch,
	}

	response, rerr := MakeRequest(apiUrl, "/rpc/post_checkup_report", requestData)
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

	response, uerr := MakeRequest(apiUrl, "/rpc/post_checkup_report_chunk", requestData)
	if uerr != nil {
		return fmt.Errorf("%s", uerr)
	}
	if msg, mok := response["message"]; mok {
		return fmt.Errorf("%s", msg)
	}

	return nil
}
