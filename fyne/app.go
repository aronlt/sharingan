package fyne

import (
	"embed"
	"fmt"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/aronlt/sharingan/internal/common"
	"github.com/aronlt/sharingan/manager"
	"github.com/aronlt/toolkit/tio"
	"github.com/golang/freetype/truetype"
	"github.com/pkg/errors"
)

//go:embed font
var f embed.FS

type App struct {
	projectDirManager *manager.ProjectDirManager
	projectManager    *manager.ProjectManager
	currentProject    string
	projectFilter     *widget.Entry
	currentStruct     string
	structFilter      *widget.Entry
	currentInterface  string
	interfaceFilter   *widget.Entry
	structZone        *fyne.Container
	interfaceZone     *fyne.Container
	goPathZone        *fyne.Container
	parseDisplay      *fyne.Container
	window            fyne.Window
}

func init() {
	fontData, err := f.ReadFile("font/simkai.ttf")
	if err != nil {
		panic(err)
	}
	_, err = truetype.Parse(fontData)
	if err != nil {
		panic(err)
	}
	tio.WriteFile(common.GetFontPath(), fontData, false)

	os.Setenv("FYNE_FONT", common.GetFontPath())
}

func chooseDirectory(app *App) {
	dialog.ShowFolderOpen(func(dir fyne.ListableURI, err error) {
		if err != nil {
			dialog.ShowError(err, app.window)
			panic(err)
		}
		if dir != nil {
			fmt.Println(dir.Path())
			manager.GoHomeDir = dir.Path()
			editZone := app.NewEditZone()
			displayZone := app.NewDisplayZone()
			content := container.NewVBox(
				editZone,
				displayZone,
			)
			app.window.SetContent(content)
		} else {
			err = errors.New("empty dir")
			dialog.ShowError(err, app.window)
			panic(err)
		}
	}, app.window)
}

func (a *App) Init() {
	main := container.NewVBox(
		container.NewVBox(
			widget.NewButton("选择go的项目主目录", func() {
				chooseDirectory(a)
			}),
		),
	)
	a.window.SetContent(main)
	a.clearDisplay()
}

func NewApp() *App {
	window := app.New().NewWindow("go分析器")
	window.Resize(fyne.Size{
		Width:  1000,
		Height: 1000,
	})
	return &App{
		window:            window,
		projectDirManager: manager.NewProjectDirManager(),
		projectManager:    manager.NewProjectManager(),
		currentProject:    "",
		currentStruct:     "",
		currentInterface:  "",
		parseDisplay:      container.NewVBox(),
	}
}

func (a *App) Start() {
	a.window.ShowAndRun()
}
