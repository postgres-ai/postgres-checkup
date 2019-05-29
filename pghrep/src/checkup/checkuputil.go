package checkup

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"../orderedmap"
	"../pyraconv"
)

// General for all reports

const MSG_ALL_GOOD_CONCLUSION string = "Hooray, all good. Keep this up!"
const MSG_NO_RECOMMENDATION string = "No recommendations."
const MSG_ETC_ITEM string = "    - etc."
const RECOMMENDATION_ITEMS_LIMIT int = 5

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

type ReportOutcome struct {
	P1              bool
	P2              bool
	P3              bool
	Conclusions     []string
	Recommendations []string
}

func (r *ReportOutcome) AppendConclusion(conclusion string, a ...interface{}) {
	r.Conclusions = append(r.Conclusions, fmt.Sprintf(conclusion, a...))
}

func (r *ReportOutcome) AppendRecommendation(reccomendation string,
	a ...interface{}) {
	if reccomendation != "" {
		r.Recommendations = append(r.Recommendations, fmt.Sprintf(reccomendation, a...))
	}
}

func LoadRawJsonReport(filePath string) []byte {
	file, err := os.Open(filePath)
	if err != nil {
		return []byte{}
	}
	defer file.Close()
	jsonRaw, err := ioutil.ReadAll(file)
	if err != nil {
		return []byte{}
	}
	return jsonRaw
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
		resultJson, merr := orderedData.MarshalJSON()
		if merr != nil {
			return
		}
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

func SaveConclusionsRecommendations(data map[string]interface{},
	result ReportOutcome) map[string]interface{} {
	if len(result.Conclusions) == 0 {
		result.AppendConclusion(MSG_ALL_GOOD_CONCLUSION)
	}
	if len(result.Recommendations) == 0 {
		result.AppendRecommendation(MSG_NO_RECOMMENDATION)
	}
	data["conclusions"] = result.Conclusions
	data["recommendations"] = result.Recommendations
	data["p1"] = result.P1
	data["p2"] = result.P2
	data["p3"] = result.P3
	SaveJsonConclusionsRecommendations(data, result.Conclusions, result.Recommendations, result.P1, result.P2, result.P3)
	return data
}

func PrintConclusions(result ReportOutcome) {
	for _, conclusion := range result.Conclusions {
		fmt.Println("C:  ", conclusion)
	}
}

func PrintReccomendations(result ReportOutcome) {
	for _, recommendation := range result.Recommendations {
		fmt.Println("R:  ", recommendation)
	}
}

func GetUniques(array []string) []string {
	items := map[string]bool{}
	for _, item := range array {
		items[item] = true
	}

	res := make([]string, len(items))
	i := 0
	for key, _ := range items {
		res[i] = key
		i++
	}
	return res
}

func LimitList(array []string) []string {
	if len(array) <= RECOMMENDATION_ITEMS_LIMIT {
		return array
	} else {
		limitedArray := array[0:RECOMMENDATION_ITEMS_LIMIT]
		limitedArray = append(limitedArray, MSG_ETC_ITEM)
		return limitedArray
	}
}
