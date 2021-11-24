/*
Postgres Health Reporter

Copyright © Postgres.ai

Perform a generation of Markdown report based on JSON results of postgres-checkup
Usage:
pghrep --checkdata=file:///path_to_check_results.json --outdir=/home/results
*/
package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"gitlab.com/postgres-ai/postgres-checkup/pghrep/internal/checkup/cfg"
	"gitlab.com/postgres-ai/postgres-checkup/pghrep/internal/config"
	"gitlab.com/postgres-ai/postgres-checkup/pghrep/internal/helpers"
	"gitlab.com/postgres-ai/postgres-checkup/pghrep/internal/log"
	"gitlab.com/postgres-ai/postgres-checkup/pghrep/internal/pyraconv"
	"gitlab.com/postgres-ai/postgres-checkup/pghrep/internal/reportutils"
	"gitlab.com/postgres-ai/postgres-checkup/pghrep/internal/upload"
)

func main() {
	// get input data checkId, checkData
	var checkId string
	var checkData string
	var resultData map[string]interface{}

	checkDataPtr := flag.String("checkdata", "", "an filepath to json report")
	outDirPtr := flag.String("outdir", "", "an directory where report need save")
	debugPtr := flag.Int("debug", 0, "enable debug mode (must be defined 1 or 0 (default))")
	modeDataPtr := flag.String("mode", "", "working mode: 'generate' (default), 'upload'")
	tokenDataPtr := flag.String("token", "", "API token to upload reports to remove server")
	projectDataPtr := flag.String("project", "", "target project used during uploading")
	pathDataPtr := flag.String("path", "", "path to artifacts directory used during uploading")
	apiUrlDataPtr := flag.String("apiurl", "", "API URL for reports uploading")

	flag.Parse()
	checkData = *checkDataPtr

	LISTLIMIT, rlerr := strconv.Atoi(os.Getenv("LISTLIMIT"))
	if rlerr != nil {
		LISTLIMIT = 50
	}

	if *debugPtr == 1 {
		log.DEBUG = true
	} else {
		log.DEBUG = false
	}

	switch *modeDataPtr {
	case "upload":
		token := *tokenDataPtr
		project := *projectDataPtr
		path := *pathDataPtr
		apiUrl := *apiUrlDataPtr

		if len(token) == 0 {
			log.Fatal("Token is not defined")
		}
		if len(apiUrl) == 0 {
			log.Fatal("API URL is not defined")
		}
		if len(project) == 0 {
			log.Fatal("Project (for reports uploading) is not defined")
		}
		if len(path) == 0 {
			log.Fatal("Artifacts directory is not defined")
		}

		err := upload.UploadReport(apiUrl, token, project, path)
		if err != nil {
			log.Err(err)
			os.Exit(1)
		}

		return
	case "loadcfg":
		path := *pathDataPtr
		if len(path) == 0 {
			log.Fatal("Config path is not defined")
		}

		conf, err := config.LoadConfig(path)
		if err != nil {
			log.Fatal(fmt.Sprintf("Cannot load config. %s", err))
			os.Exit(1)
		}

		if len(conf) == 0 {
			log.Fatal(fmt.Sprintf("Config '%s' is empty", path))
		}

		config.OutputConfig(conf)

		return
	}

	if helpers.FileExists(checkData) {
		resultData = reportutils.LoadJsonFile(checkData)

		if resultData == nil {
			log.Fatal("File given by --checkdata content wrong json data")
			return
		}

		resultData["source_path_full"] = checkData
		resultData["source_path_parts"] = strings.Split(checkData, string(os.PathSeparator))
	} else {
		log.Err("File given by --checkdata not found")
		return
	}

	if resultData != nil {
		checkId = pyraconv.ToString(resultData["checkId"])
	} else {
		log.Fatal("Content defined by '--checkdata' is invalid JSON")
	}

	checkId = strings.ToUpper(checkId)
	reportutils.LoadDependencies(resultData)
	reportutils.DetermineMasterReplica(resultData)

	err := reportutils.ReorderHosts(resultData)
	if err != nil {
		log.Err("There is no data to generate the report")
	}

	newConfig := cfg.NewConfig()

	err = reportutils.PreprocessReportData(checkId, newConfig, resultData)
	if err != nil {
		log.Fatal(err)
	}

	resultData["LISTLIMIT"] = LISTLIMIT
	var outputDir string

	if len(*outDirPtr) == 0 {
		outputDir = "./"
	} else {
		outputDir = *outDirPtr
	}

	reportDone := reportutils.GenerateMdReports(checkId, resultData, outputDir)
	if !reportDone {
		log.Fatal("Cannot generate report. Data file or template is wrong")
	}
}
