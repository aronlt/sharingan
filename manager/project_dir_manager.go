package manager

import (
	"github.com/aronlt/toolkit/tio"
	"github.com/sirupsen/logrus"
)

var GoHomeDir string

type ProjectDirManager struct {
}

func NewProjectDirManager() *ProjectDirManager {
	return &ProjectDirManager{}
}

func (p *ProjectDirManager) ListAllProject() ([]string, error) {
	allDirs, _, err := tio.ReadDir(GoHomeDir)
	if err != nil {
		logrus.WithError(err).Errorf("call tio.ReadDir fail, dir:%s", GoHomeDir)
		return nil, err
	}
	logrus.Infof("call ReadDir success")
	return allDirs, nil
}
