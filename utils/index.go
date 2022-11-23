package utils

import "os"

var ERR_LOG_PATH = "launcher_err.log"

func IsFileExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func WriteErrLog(err string) {
	os.WriteFile(ERR_LOG_PATH, []byte(err), 0744)
}
