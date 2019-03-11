package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"../src/orderedmap"
	"../src/pyraconv"
)

var CRITICAL_BLOAT float64 = 90

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
		resultJson, jerr := orderedData.MarshalJSON()
		if jerr == nil {
			var out bytes.Buffer
			json.Indent(&out, resultJson, "", "  ")
			jfile, err := os.Create(filePath)
			if err != nil {
				return
			}
			defer jfile.Close()
			out.WriteTo(jfile)
		} else {
			return
		}
	}
}

func f005ConclusionsRecommendations(data map[string]interface{}) {
	var conclusions []string
	var recommendations []string

	hosts := pyraconv.ToInterfaceMap(data["hosts"])
	master := pyraconv.ToString(hosts["master"])
	results := pyraconv.ToInterfaceMap(data["results"])
	masterData := pyraconv.ToInterfaceMap(results[master])
	masterData = pyraconv.ToInterfaceMap(masterData["data"])
	indexesData := pyraconv.ToInterfaceMap(masterData["index_bloat"])
	//master
	preparedData := make(map[string]interface{})
	for index, indexData := range indexesData {
		if index == "_keys" {
			continue
		}
		idata := pyraconv.ToInterfaceMap(indexData)
		bloatPercent := pyraconv.ToString(idata["Bloat ratio"])
		percent, _ := strconv.ParseFloat(bloatPercent, 32)
		if percent > CRITICAL_BLOAT {
			conclusion := fmt.Sprintf(":warning: Index `%s` bloated more than `%.0f` percent", idata["Index (Table)"], CRITICAL_BLOAT)
			recommendation := fmt.Sprintf(":warning: Please check index `%s`", idata["Index (Table)"])
			conclusions = append(conclusions, conclusion)
			recommendations = append(recommendations, recommendation)
		}
		preparedData[index] = percent
	}

	data["conclusions"] = conclusions
	data["recommendations"] = recommendations
	saveJsonConclusionsRecommendations(data, conclusions, recommendations)
}

func (g prepare) Prepare(data map[string]interface{}) map[string]interface{} {
	f005ConclusionsRecommendations(data)
	return data
}

var Preparer prepare
