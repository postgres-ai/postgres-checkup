/*
2019 © Anatoly Stansler anatoly@postgres.ai
2019 © Dmitry Udalov dmius@postgres.ai
2019 © Postgres.ai
*/

package cfg

import (
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"../../log"

	"golang.org/x/net/html"
)

type Version struct {
	FirstRelease  string
	FinalRelease  string
	MinorVersions []int
}

type Config struct {
	Versions map[string]Version
}

const POSTGRES_RELEASES_URL string = "https://git.postgresql.org/gitweb/?p=postgresql.git;a=tags"
const RELEASE_CODE = "REL"

// TODO(anatoly): Fill up 12 version on release or load this information automatically.
var versionsDefault map[string]Version = map[string]Version{
	"11": Version{
		FirstRelease:  "2018-10-18",
		FinalRelease:  "2023-11-09",
		MinorVersions: []int{0, 1, 2, 3, 4},
	},
	"10": Version{
		FirstRelease:  "2017-10-05",
		FinalRelease:  "2022-11-10",
		MinorVersions: []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
	},
	"9.6": Version{
		FirstRelease:  "2016-09-29",
		FinalRelease:  "2021-11-11",
		MinorVersions: []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14},
	},
	"9.5": Version{
		FirstRelease:  "2016-01-07",
		FinalRelease:  "2021-02-11",
		MinorVersions: []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18},
	},
	"9.4": Version{
		FirstRelease: "2014-12-18",
		FinalRelease: "2020-02-13",
		MinorVersions: []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19,
			20, 21, 22, 23},
	},
	"9.3": Version{
		FirstRelease: "2013-09-09",
		FinalRelease: "2018-11-08",
		MinorVersions: []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19,
			20, 21, 22, 23, 24, 25},
	},
	"9.2": Version{
		FirstRelease: "2012-09-10",
		FinalRelease: "2017-11-09",
		MinorVersions: []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19,
			20, 21, 22, 23, 24},
	},
	"9.1": Version{
		FirstRelease: "2011-09-12",
		FinalRelease: "2016-10-27",
		MinorVersions: []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19,
			20, 21, 22, 23, 24},
	},
	"9.0": Version{
		FirstRelease: "2010-09-20",
		FinalRelease: "2015-10-08",
		MinorVersions: []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19,
			20, 21, 22, 23},
	},
	"8.4": Version{
		FirstRelease: "2009-07-01",
		FinalRelease: "2014-07-24",
		MinorVersions: []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19,
			20, 21, 22},
	},
	"8.3": Version{
		FirstRelease: "2008-02-04",
		FinalRelease: "2013-02-07",
		MinorVersions: []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19,
			20, 21, 22, 23},
	},
	"8.2": Version{
		FirstRelease: "2006-12-05",
		FinalRelease: "2011-12-05",
		MinorVersions: []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19,
			20, 21, 22, 23},
	},
	"8.1": Version{
		FirstRelease: "2005-11-08",
		FinalRelease: "2010-11-08",
		MinorVersions: []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19,
			20, 21, 22, 23},
	},
	"8.0": Version{
		FirstRelease: "2005-01-19",
		FinalRelease: "2010-10-01",
		MinorVersions: []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19,
			20, 21, 22, 23, 24, 25, 26},
	},
	"7.4": Version{
		FirstRelease: "2003-11-17",
		FinalRelease: "2010-10-01",
		MinorVersions: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20,
			21, 22, 23, 24, 25, 26, 27, 28, 29, 30},
	},
	"7.3": Version{
		FirstRelease: "2002-11-27",
		FinalRelease: "2007-11-27",
		MinorVersions: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20,
			21},
	},
	"7.2": Version{
		FirstRelease:  "2002-02-04",
		FinalRelease:  "2007-02-04",
		MinorVersions: []int{1, 2, 3, 4, 5, 6, 7, 8},
	},
	"7.1": Version{
		FirstRelease:  "2001-03-13",
		FinalRelease:  "2006-03-13",
		MinorVersions: []int{1, 2, 3},
	},
	"7.0": Version{
		FirstRelease:  "2000-05-08",
		FinalRelease:  "2005-05-08",
		MinorVersions: []int{2, 3},
	},
	"6.5": Version{
		FirstRelease:  "1999-06-09",
		FinalRelease:  "2004-06-09",
		MinorVersions: []int{1, 2, 3},
	},
	"6.4": Version{
		FirstRelease:  "1998-10-30",
		FinalRelease:  "2003-10-30",
		MinorVersions: []int{2},
	},
	"6.3": Version{
		FirstRelease:  "1998-03-01",
		FinalRelease:  "2003-03-01",
		MinorVersions: []int{2},
	},
}

func NewConfig() Config {
	config := Config{
		Versions: versionsDefault,
	}
	return config
}

// Fetch actual versions from Postgres repository releases.
func (c *Config) LoadVersions() error {
	releases, err := loadPostgresReleases()
	if err != nil {
		return err
	}

	err = fillVersions(c.Versions, releases)
	if err != nil {
		return err
	}

	return nil
}

func loadPostgresReleases() ([]string, error) {
	releases := []string{}

	log.Dbg("Loading Postgres releases...")
	resp, err := http.Get(POSTGRES_RELEASES_URL)
	if err != nil {
		return []string{}, err
	}
	defer resp.Body.Close()

	htmlData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []string{}, err
	}

	tokenizer := html.NewTokenizer(strings.NewReader(string(htmlData)))
	for {
		tokenType := tokenizer.Next()

		if tokenType == html.ErrorToken {
			if tokenizer.Err() == io.EOF {
				break
			}

			return []string{}, tokenizer.Err()
		}

		if tokenType != html.TextToken {
			continue
		}

		text := strings.TrimSpace(html.UnescapeString(string(tokenizer.Text())))

		if !strings.HasPrefix(text, RELEASE_CODE) {
			continue
		}

		releases = append(releases, text)
	}

	return releases, nil
}

// TODO(anatoly): Test.
func fillVersions(versions map[string]Version, releases []string) error {
	minorVersionsMap := make(map[string][]int)
	for majorVersion, _ := range versions {
		minorVersionsMap[majorVersion] = []int{}
	}

	// Samples: REL9_6_14, REL_10_9, REL_11_4, REL_11_BETA4.
	for _, release := range releases {
		release = strings.Replace(release, RELEASE_CODE, "", 1)

		// We need only released versions.
		if strings.Contains(release, "BETA") || strings.Contains(release, "RC") ||
			strings.Contains(release, "ALPHA") {
			continue
		}

		ver := strings.Split(release, "_")
		if len(ver) < 3 {
			continue
		}

		majorVersion := ver[0]
		if majorVersion != "" {
			majorVersion = majorVersion + "."
		}
		majorVersion = majorVersion + ver[1]
		minorVersion := ver[2]

		minorVersions, ok := minorVersionsMap[majorVersion]
		if !ok {
			continue
		}

		minorVersionInt, err := strconv.Atoi(minorVersion)
		if err != nil {
			return err
		}

		minorVersionsMap[majorVersion] = append(minorVersions, minorVersionInt)
	}

	for majorVersion, version := range versions {
		version.MinorVersions = minorVersionsMap[majorVersion]
		versions[majorVersion] = version
	}

	return nil
}
