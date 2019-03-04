package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"../src/orderedmap"
	"../src/pyraconv"
)

var CRITICAL_USAGE int = 90

var Data map[string]interface{}

type prepare string

func saveJsonConclusionsRecommendations(data map[string]interface{}, conclusions []string, recommendations []string) {
	filePath := pyraconv.ToString(data["source_path_full"])
	jsonData, err := ioutil.ReadFile(filePath) // just pass the file name
	if err != nil {
		return
	}
	orderedData := orderedmap.New()
	if err := json.Unmarshal([]byte(jsonData), &orderedData); err != nil {
		return
	} else {
		orderedData.Set("conclusions", conclusions)
		orderedData.Set("recommendations", recommendations)
		resultJson, _ := orderedData.MarshalJSON()
		var out bytes.Buffer
		json.Indent(&out, resultJson, "", "  ")
		jfile, err := os.Create(filePath)
		if err != nil {
			return
		}
		defer jfile.Close()
		out.WriteTo(jfile)
	}
}

func getHostDbpointData(data map[string]interface{}, host string, point string) map[string]interface{} {
	results := pyraconv.ToInterfaceMap(data["results"])
	masterData := pyraconv.ToInterfaceMap(results[host])
	masterData = pyraconv.ToInterfaceMap(masterData["data"])
	hostData := pyraconv.ToInterfaceMap(masterData["db_data"])
	return pyraconv.ToInterfaceMap(hostData[point])
}

func a008ConclusionsRecommendations(data map[string]interface{}) {
	preparedData := a008PrepareData(data)
	var conclusions []string
	var recommendations []string

	for host, hData := range preparedData {
		hostData := pyraconv.ToInterfaceMap(hData)
		for point, pPercent := range hostData {
			pointPercent := int(pyraconv.ToInt64(pPercent))
			if pointPercent > CRITICAL_USAGE {
				// generate recommendation
				pointData := getHostDbpointData(data, host, point)
				conclusion := fmt.Sprintf(":warning: Volume `%s` at host `%s` where placed `%s` is filled more than on %d percent", pointData["mount_point"], host, pointData["path"], CRITICAL_USAGE)
				recommendation := fmt.Sprintf(":warning: Please clean volume `%s` at host `%s` where placed `%s`", pointData["mount_point"], host, pointData["path"])
				conclusions = append(conclusions, conclusion)
				recommendations = append(recommendations, recommendation)
			}
		}
	}

	data["conclusions"] = conclusions
	data["recommendations"] = recommendations
	saveJsonConclusionsRecommendations(data, conclusions, recommendations)
}

func a008PrepareData(data map[string]interface{}) map[string]interface{} {
	prepareData := make(map[string]interface{})

	hosts := pyraconv.ToInterfaceMap(data["hosts"])
	master := pyraconv.ToString(hosts["master"])
	replicas := pyraconv.ToStringArray(hosts["replicas"])

	results := pyraconv.ToInterfaceMap(data["results"])
	masterData := pyraconv.ToInterfaceMap(results[master])
	masterData = pyraconv.ToInterfaceMap(masterData["data"])
	dbMasterData := pyraconv.ToInterfaceMap(masterData["db_data"])
	//master
	preparedMasterData := make(map[string]interface{})
	for point, pointData := range dbMasterData {
		if point == "_keys" {
			continue
		}
		pdata := pyraconv.ToInterfaceMap(pointData)
		usePercent := pyraconv.ToString(pdata["use_percent"])
		usePercent = strings.Replace(usePercent, "%", "", 1)
		percent, _ := strconv.Atoi(usePercent)
		preparedMasterData[point] = percent
	}
	prepareData[master] = preparedMasterData
	// replicas
	for _, replica := range replicas {
		replicaData := pyraconv.ToInterfaceMap(results[master])
		masterData = pyraconv.ToInterfaceMap(replicaData["data"])
		dbReplicaData := pyraconv.ToInterfaceMap(replicaData["db_data"])
		preparedReplicaData := make(map[string]interface{})
		for point, pointData := range dbReplicaData {
			if point == "_keys" {
				continue
			}
			pdata := pyraconv.ToInterfaceMap(pointData)
			usePercent := pyraconv.ToString(pdata["use_percent"])
			usePercent = strings.Replace(usePercent, "%", "", 1)
			percent, _ := strconv.Atoi(usePercent)
			preparedReplicaData[point] = percent
		}
		prepareData[replica] = preparedReplicaData
	}
	return prepareData
}

func (g prepare) Prepare(data map[string]interface{}) map[string]interface{} {
	a008ConclusionsRecommendations(data)
	return data
}

var Preparer prepare
