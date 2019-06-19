/*
2019 © Anatoly Stansler anatoly@postgres.ai
2019 © Dmitry Udalov dmius@postgres.ai
2019 © Postgres.ai
*/

package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

type SupportedVersion struct {
	FirstRelease  string
	FinalRelease  string
	MinorVersions []int
}

type Config struct {
	MajorVersions []int
	SupportedVersions map[string]SupportedVersion
}

const VERSIONS_SOURCE_URL string = "https://git.postgresql.org/gitweb/?p=postgresql.git;a=tags"
const VERSIONS_REL_TEXT = "REL"

const SUPPORTED_VERSIONS_DEFAULT map[string]SupportedVersion = map[string]SupportedVersion{
	"11": SupportedVersion{
		FirstRelease:  "2018-10-18",
		FinalRelease:  "2023-11-09",
		MinorVersions: []int,
	},
	"10": SupportedVersion{
		FirstRelease:  "2017-10-05",
		FinalRelease:  "2022-11-10",
		MinorVersions: []int,
	},
	"9.6": SupportedVersion{
		FirstRelease:  "2016-09-29",
		FinalRelease:  "2021-11-11",
		MinorVersions: []int,
	},
	"9.5": SupportedVersion{
		FirstRelease:  "2016-01-07",
		FinalRelease:  "2021-02-11",
		MinorVersions: []int,
	},
	"9.4": SupportedVersion{
		FirstRelease:  "2014-12-18",
		FinalRelease:  "2020-02-13",
		MinorVersions: []int,
	},
}

func (c *Config) Load() (error) {
	versions, majorVersions, err := loadVersions()
	if err != nil {
		return err
	}
	c.SupportedVersions = supportedVersions
	c.MajorVersions = majorVersions
}

func loadVersions() (supportedVersions map[string]SupportedVersion, majorVersions []int, error) {
	supportedVersions := SUPPORTED_VERSIONS_DEFAULT
	majorVersionsMap := make(map[int]bool)

	log.Dbg("Loading versions...")
	resp, err := http.Get(VERSIONS_SOURCE_URL)
	if err != nil {
		return []int, err
	}
	defer resp.Body.Close()

	html, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []int, err
	}

	tokenizer := html.NewTokenizer(strings.NewReader(string(html)))
	for {
		tokenType := tokenizer.Next()

		if tokenType == html.ErrorToken {
			return []int, tokenizer.Err()
		}

		if tokenType != html.TextToken {
			continue
		}

		text := strings.TrimSpace(html.UnescapeString(string(tokenizer.Text())))

		// Samples: REL9_6_14, REL_10_9, REL_11_4, REL_11_BETA4.
		if !strings.HasPrefix(text, VERSIONS_REL_TEXT) {
			continue
		}
		text = strings.Replace(text, VERSIONS_REL_TEXT, "", 1)

		// We need only released versions.
		if strings.Contains(text, "BETA") || strings.Contains(text, "RC") ||
			strings.Contains(text, "ALPHA") {
			continue
		}

		ver := strings.Split(text, "_")
		if len(ver) < 3 {
			continue
		}

		majorVersion := ver[0]
		if majorVersion != "" {
			majorVersion = majorVersion + "."
		}
		majorVersion = majorVersion + ver[1]

		majorVersionNum := strings.Replace(majorVersion, ".", "0", 1)
		if len(majorVersionNum) < 3 {
			majorVersionNum = majorVersionNum + "00"
		}
		majorVersionInt, err := strconv.Atoi(majorVersionNum)
		if err != nil {
			return []int, err
		}
		majorVersionsMap[majorVersionInt] = true

		minorVersion := ver[2]

		version, ok := supportedVersions[majorVersion]
		if ok {
			minorVersionInt, _ := strconv.Atoi(minorVersion)
			ver.MinorVersions = append(ver.MinorVersions, minorVersionInt)
			supportedVersions[majorVersion] = version
		}
	}

	majorVersionsList []int
	for ver, _ := range majorVersionsMap {
		majorVersionsList = append(majorVersionsList, ver)
	}
	sort.Ints(majorVersionsList)

	return supportedVersions, majorVersionsList, nil
}
