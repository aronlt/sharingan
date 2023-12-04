package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/aronlt/sharingan/fileparser"
	"github.com/aronlt/toolkit/terror"
)

const (
	maxDeep = 100
)

// 递归搜索，找出全部go文件, 并且忽略掉proto生成的文件
func visit(path string, deep int) ([]string, error) {
	if deep > maxDeep {
		return nil, fmt.Errorf("can't visit too deep dirs")
	}
	files := make([]string, 0)
	rd, err := os.ReadDir(path)
	if err != nil {
		return nil, terror.Wrapf(err, "can't read dir:%s, error:%s", path, err.Error())
	}
	for _, fi := range rd {
		subPath := filepath.Join(path, fi.Name())
		if fi.IsDir() {
			if fi.Name() != "vendor" {
				if subFiles, err := visit(subPath, deep+1); err == nil {
					files = append(files, subFiles...)
				}
			}
		} else {
			if strings.HasSuffix(subPath, ".go") && !strings.HasSuffix(subPath, ".pb.go") && !strings.HasSuffix(subPath, "_test.go") {
				files = append(files, subPath)
			}
		}
	}
	return files, nil
}

func NewParser() (fileparser.Parser, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, terror.Wrap(err, "call os.Getwd fail")
	}

	files, err := visit(cwd, 0)
	if err != nil {
		return nil, terror.Wrap(err, "call visit fail")
	}

	nodeManager := fileparser.NewParser(cwd)
	for _, file := range files {
		err = nodeManager.Inspect(file)
		if err != nil {
			return nil, terror.Wrapf(err, "can't inspect file:%s, error:%s", file, err.Error())
		}
	}
	nodeManager.ParseStructFunctions()

	return nodeManager, nil
}
