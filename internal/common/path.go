package common

import (
	"os"
	"path/filepath"
	"sync"
)

const DifficultyFilePath = "data/difficulty.json"
const TagsFilePath = "data/tags.json"
const DataFilePath = "data/data.json"

var resultOnce sync.Once
var fontOnce sync.Once
var resultPath string
var fontPath string

func GetResultPath() string {
	resultOnce.Do(
		func() {
			dirname, err := os.UserHomeDir()
			if err != nil {
				panic(err)
			}
			resultPath = filepath.Join(dirname, "Documents/any-leetcode/result")
		})
	return resultPath
}

func GetFontPath() string {
	fontOnce.Do(
		func() {
			dirname, err := os.UserHomeDir()
			if err != nil {
				panic(err)
			}
			fontPath = filepath.Join(dirname, "Documents/any-leetcode/simkai.ttf")
		})
	return fontPath
}
