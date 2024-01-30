package fyne

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/sirupsen/logrus"
)

func (a *App) clearDisplay() {
	a.parseDisplay.Objects = a.parseDisplay.Objects[:0]
}
func (a *App) refreshDisplay() {
	a.parseDisplay.Refresh()
}
func (a *App) NewParseButton() *widget.Button {
	button := widget.NewButton("解析", func() {
		if a.currentInterface == "" || a.currentStruct == "" {
			return
		}
		missing, wrong, err := a.projectManager.ProjectParser.AnalysisInterface(
			a.currentInterface,
			a.currentStruct)
		if err != nil {
			logrus.WithError(err).Errorf("call AnalysisInterface fail")
			panic(err)
		}

		a.clearDisplay()
		a.parseDisplay.Add(canvas.NewText("没有实现方法:", color.Black))
		for _, fn := range missing {
			a.parseDisplay.Add(canvas.NewText(fn, color.Black))
		}
		a.parseDisplay.Add(canvas.NewText("错误实现的方法:", color.Black))
		for _, fn := range wrong {
			a.parseDisplay.Add(canvas.NewText(fn, color.Black))
		}
		a.refreshDisplay()
	})
	return button
}

func (a *App) NewActionZone() *fyne.Container {
	genButton := a.NewParseButton()
	return container.NewHBox(genButton)
}
