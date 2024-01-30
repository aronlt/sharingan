package fyne

import (
	"sort"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/aronlt/toolkit/ds"
	"github.com/sirupsen/logrus"
)

// NewProjectZone 选择项目信息
func (a *App) NewProjectZone() *fyne.Container {
	getProjects := func(filter string) []string {
		allProjects, err := a.projectDirManager.ListAllProject()
		if err != nil {
			logrus.WithError(err).Errorf("call ListAllProject fail")
			panic(err)
		}
		sort.Strings(allProjects)
		allProjects = ds.SliceGetFilter(allProjects, func(i int) bool {
			if a.projectFilter.Text != "" {
				return ds.StrHasContainInsensitive(allProjects[i], a.projectFilter.Text)
			}
			return true
		})
		return allProjects
	}

	newSelectContainer := func(projects []string) *widget.Select {
		selectContainer := widget.NewSelect(projects, func(project string) {
			if project == a.currentProject {
				return
			}
			a.currentInterface = ""
			a.currentStruct = ""
			a.currentProject = ""
			err := a.projectManager.ResetProjectDir(project)
			if err != nil {
				logrus.WithError(err).Errorf("call ResetProjectDir fail, project:%s", project)
				panic(err)
			}
			a.resetStructZone()
			a.resetInterfaceZone()
		})
		return selectContainer
	}

	a.projectFilter = widget.NewEntry()

	projectZone := container.NewVBox(
		widget.NewLabel("选择项目:"),
		a.projectFilter,
		newSelectContainer(getProjects("")))

	a.projectFilter.OnSubmitted = func(s string) {
		projects := getProjects(s)
		index := len(projectZone.Objects) - 1
		projectZone.Objects[index] = newSelectContainer(projects)
	}
	return projectZone
}

func (a *App) resetStructZone() {
	index := len(a.structZone.Objects) - 1
	a.structZone.Objects[index] =
		widget.NewSelect(a.projectManager.AllStructTokens, func(structNode string) {
			a.currentStruct = structNode
		})
}

func (a *App) resetInterfaceZone() {
	index := len(a.interfaceZone.Objects) - 1
	a.interfaceZone.Objects[index] = widget.NewSelect(a.projectManager.AllInterfaceTokens, func(interfaceNode string) {
		a.currentInterface = interfaceNode
	})
}

func (a *App) NewStructZone() *fyne.Container {
	a.structFilter = widget.NewEntry()

	newSelector := func(filter string) *widget.Select {
		tokens := ds.SliceGetFilter(a.projectManager.AllStructTokens, func(i int) bool {
			if filter != "" {
				return ds.StrHasContainInsensitive(a.projectManager.AllStructTokens[i], filter)
			}
			return true
		})
		return widget.NewSelect(tokens, func(structNode string) {
			a.currentStruct = structNode
		})
	}
	a.structZone = container.NewVBox(
		widget.NewLabel("选择结构体:"),
		a.structFilter,
		newSelector(""))
	a.structFilter.OnChanged = func(s string) {
		index := len(a.structZone.Objects) - 1
		a.structZone.Objects[index] = newSelector(s)
	}
	return a.structZone
}

func (a *App) NewInterfaceZone() *fyne.Container {
	a.interfaceFilter = widget.NewEntry()
	newSelector := func(filter string) *widget.Select {
		tokens := ds.SliceGetFilter(a.projectManager.AllInterfaceTokens, func(i int) bool {
			if filter != "" {
				return ds.StrHasContainInsensitive(a.projectManager.AllInterfaceTokens[i], filter)
			}
			return true
		})
		return widget.NewSelect(tokens, func(interfaceNode string) {
			a.currentInterface = interfaceNode
		})
	}
	a.interfaceZone = container.NewVBox(
		widget.NewLabel("选择接口:"),
		a.interfaceFilter,
		newSelector(""))
	a.interfaceFilter.OnChanged = func(s string) {
		index := len(a.interfaceZone.Objects) - 1
		a.interfaceZone.Objects[index] = newSelector(s)
	}

	return a.interfaceZone
}

func (a *App) NewEditZone() *fyne.Container {
	zone := container.NewVBox(
		a.NewProjectZone(),
		a.NewStructZone(),
		a.NewInterfaceZone(),
		a.NewActionZone(),
	)
	return zone
}
