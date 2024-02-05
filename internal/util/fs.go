package util

import (
	"encoding/json"
	"os"
)

func IsExist(path string) bool {
	_, err := os.Lstat(path)
	if err != nil {
		return !os.IsNotExist(err)
	}
	return true
}

func IsFileExist(path string) bool {
	fileInfo, err := os.Lstat(path)
	if err != nil {
		return !os.IsNotExist(err)
	}
	return fileInfo.Mode().IsRegular()
}

func IsDirExists(path string) bool {
	fileInfo, err := os.Lstat(path)
	if err != nil {
		return !os.IsNotExist(err)
	}
	return fileInfo.IsDir()
}

func CanSymbolLink(path string) bool {
	fileInfo, err := os.Lstat(path)
	if err != nil {
		return os.IsNotExist(err)
	}
	return fileInfo.Mode()&os.ModeSymlink != 0
}

func ReadJsonFromFile(path string, v interface{}) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	return json.NewDecoder(f).Decode(v)
}

func WriteJsonToFile(path string, v interface{}) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	return json.NewEncoder(f).Encode(v)
}
