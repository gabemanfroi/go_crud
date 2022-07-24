package utils

import (
	"errors"
	"log"
	"os"
)

func GetWorkingDirectory() string {
	directory, err := os.Getwd()
	if err != nil {
		log.Fatal("%s", err.Error())
	}
	return directory
}

func CreateDirectoryIfNotExists(path string) {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			log.Println(err)
		}
	}
}
