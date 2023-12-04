package main

import (
	"fmt"
	"strings"

	"github.com/aronlt/sharingan/cmd"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func main() {
	parser, err := cmd.NewParser()
	if err != nil {
		panic(err)
	}
	_, structTokens, interfaceTokens := parser.AllTokens()
	interfaceName := getText("接口", interfaceTokens)
	structName := getText("结构体", structTokens)
	missing, wrong, err := parser.AnalysisInterface(interfaceName, structName)
	if err != nil {
		panic(err)
	}
	if len(missing) > 0 {
		fmt.Println("未实现函数:")
		for _, m := range missing {
			fmt.Println(m)
		}
	}
	if len(wrong) > 0 {
		fmt.Println("实现错误函数:")
		for _, m := range wrong {
			fmt.Println(m)
		}
	}
	if len(missing) == 0 && len(wrong) == 0 {
		fmt.Println("结构体实现接口")
	}
}

func getText(hint string, words []string) string {
	app := tview.NewApplication()
	inputField := tview.NewInputField().
		SetLabel(hint + ":").
		SetFieldWidth(30).
		SetDoneFunc(func(key tcell.Key) {
			app.Stop()
		})
	inputField.SetAutocompleteFunc(func(currentText string) (entries []string) {
		if len(currentText) == 0 {
			return
		}
		for _, word := range words {
			if strings.Contains(strings.ToLower(word), strings.ToLower(currentText)) {
				entries = append(entries, word)
			}
		}
		if len(entries) < 1 {
			entries = nil
		}
		return
	})
	inputField.SetAutocompletedFunc(func(text string, index, source int) bool {
		if source != tview.AutocompletedNavigate {
			inputField.SetText(text)
		}
		return source == tview.AutocompletedEnter || source == tview.AutocompletedClick
	})
	if err := app.EnableMouse(true).SetRoot(inputField, true).Run(); err != nil {
		panic(err)
	}

	return inputField.GetText()
}
