package manager

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/aronlt/sharingan/internal/fileparser"
	"github.com/aronlt/toolkit/terror"
	"github.com/sirupsen/logrus"
)

const (
	maxDeep = 100
)

type ProjectManager struct {
	ProjectDir         string
	ProjectParser      fileparser.Parser
	AllFunctionTokens  []string
	AllStructTokens    []string
	AllInterfaceTokens []string
}

func NewProjectManager() *ProjectManager {
	return &ProjectManager{}
}

// 递归搜索，找出全部go文件, 并且忽略掉proto生成的文件
func (p *ProjectManager) visit(path string, deep int) ([]string, error) {
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
				if subFiles, err := p.visit(subPath, deep+1); err == nil {
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

func (p *ProjectManager) ResetProjectDir(dir string) error {
	p.ProjectDir = dir

	files, err := p.visit(p.ProjectDir, 0)
	if err != nil {
		logrus.WithError(err).Errorf("call visit fail")
		return err
	}

	nodeManager := fileparser.NewParser(p.ProjectDir)
	for _, file := range files {
		err = nodeManager.Inspect(file)
		if err != nil {
			logrus.WithError(err).Errorf("call Inspect file:%s fail", file)
			return err
		}
	}
	nodeManager.ParseStructFunctions()

	p.ProjectParser = nodeManager
	p.AllFunctionTokens, p.AllStructTokens, p.AllInterfaceTokens = p.ProjectParser.AllTokens()
	logrus.Infof("call ResetProjectDir success, dir:%s", dir)
	return nil
}
