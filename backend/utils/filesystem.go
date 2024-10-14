package utils

import (
	"backend/config"
	"fmt"
	"io/ioutil"
	"path/filepath"
)

// ReadVarFile returns the content of a file in var directory
func ReadVarFile(subpath string) (string, error) {
	iConfig := config.GetConfig()
	tplFilePath := filepath.Join(iConfig.GetVarDirAbsPath(), subpath)
	tplContent, err := ioutil.ReadFile(tplFilePath)

	if err != nil {
		return "", fmt.Errorf("could not read file '%s': %v", tplFilePath, err)
	}

	return string(tplContent), nil
}
