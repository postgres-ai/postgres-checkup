package helpers

import (
	"os"
	"strings"

	"gitlab.com/postgres-ai/postgres-checkup/pghrep/internal/log"
)

// Prepropess file paths
// Allow absulute and relative (of pwd) paths with or wothout file:// prefix
// Return absoulute path of file
func GetFilePath(name string) string {
	filePath := name
	// remove file:// prefix
	if strings.HasPrefix(strings.ToLower(filePath), "file://") {
		filePath = strings.Replace(filePath, "file://", "", 1)
	}

	if strings.HasPrefix(strings.ToLower(filePath), "/") /* for *nix paths */ ||
		(len(filePath) > 1 && filePath[1] == ':') /* for windows path */ {
		// absoulute path will use as is
		return filePath
	} else {
		// for relative path will combine with current path
		curDir, err := os.Getwd()
		if err != nil {
			log.Dbg("Can't determine current path")
		}

		if strings.HasSuffix(strings.ToLower(curDir), "/") {
			filePath = curDir + filePath
		} else {
			filePath = curDir + "/" + filePath
		}

		return filePath
	}
}

// Check file exists
// Allow absulute and relative (of pwd) paths with or wothout file:// prefix
// Return boolean value
func FileExists(name string) bool {
	filePath := GetFilePath(name)
	if _, err := os.Stat(filePath); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}

	return true
}

