package checkup

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"

	"../orderedmap"
	"../pyraconv"
)

// General for all reports

type ReportHost struct {
	InternalAlias        string `json:"internal_alias"`
	Index                string `json:"index"`
	Role                 string `json:"role"`
	RoleChangeDetectedAt string `json:"role_change_detected_at"`
}

type ReportLastCheck struct {
	Epoch       string `json:"epoch"`
	Timestamptz string `json:"timestamptz"`
	Dir         string `json:"dir"`
}

type ReportHosts map[string]ReportHost

type ReportLastNodes struct {
	Hosts     ReportHosts     `json:"hosts"`
	LastCheck ReportLastCheck `json:"last_check"`
	//	LastCheck
}

func SaveJsonConclusionsRecommendations(data map[string]interface{}, conclusions []string,
	recommendations []string, p1 bool, p2 bool, p3 bool) {
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
		orderedData.Set("p1", p1)
		orderedData.Set("p2", p2)
		orderedData.Set("p3", p3)
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
