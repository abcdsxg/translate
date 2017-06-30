package main

import (
	"Translate/base"
	"Translate/site"
	"fmt"

	"github.com/andlabs/ui"
)

var inputEntry *ui.Entry
var cibaLabel *ui.Label
var dictLabel *ui.Label
var shanbeiLabel *ui.Label
var youdaoLabel *ui.Label
func main() {
	err := ui.Main(func() {
		mainBox := ui.NewVerticalBox()
		window := ui.NewWindow("划词翻译", 500, 500, false)
		window.SetMargined(true)
		initInputBox(mainBox)
		initResultBox(mainBox)
		initWindow(window, mainBox)

	})
	if err != nil {
		panic(err)
	}
}

func initInputBox(mainBox *ui.Box) {
	inputBox := ui.NewHorizontalBox()
	inputEntry = ui.NewEntry()
	inputBtn := ui.NewButton("查询")
	inputBtn.OnClicked(search)
	inputBox.Append(inputEntry, false)
	inputBox.Append(inputBtn, false)
	inputBox.SetPadded(true)
	mainBox.Append(inputBox, false)
}

func initResultBox(mainBox *ui.Box) {

	dictBox := ui.NewVerticalBox()
	cibaBox := ui.NewVerticalBox()
	shanbeiBox := ui.NewVerticalBox()
	youdaoBox := ui.NewVerticalBox()

	dictTab := ui.NewTab()
	cibaTab := ui.NewTab()
	shanbeiTab := ui.NewTab()
	youdaoTab := ui.NewTab()

	initLabel()

	dictTab.InsertAt("海词", 0, dictLabel)
	cibaTab.InsertAt("词霸", 0, cibaLabel)
	shanbeiTab.InsertAt("扇贝", 0, shanbeiLabel)
	youdaoTab.InsertAt("有道", 0, youdaoLabel)

	dictTab.SetMargined(0,true)
	cibaTab.SetMargined(0,true)
	shanbeiTab.SetMargined(0,true)
	youdaoTab.SetMargined(0,true)

	dictBox.Append(dictTab, false)
	cibaBox.Append(cibaTab, false)
	shanbeiBox.Append(shanbeiTab, false)
	youdaoBox.Append(youdaoTab, false)

	mainBox.Append(dictBox, false)
	mainBox.Append(cibaBox, false)
	mainBox.Append(shanbeiBox, false)
	mainBox.Append(youdaoBox, false)

}
func initLabel() {
	cibaLabel = ui.NewLabel("\n\n\n")
	dictLabel = ui.NewLabel("\n\n\n")
	shanbeiLabel = ui.NewLabel("\n\n\n")
	youdaoLabel = ui.NewLabel("\n\n\n")
}

func initWindow(window *ui.Window, mainBox *ui.Box) {
	window.SetChild(mainBox)
	window.OnClosing(func(*ui.Window) bool {
		ui.Quit()
		return true
	})
	window.Show()
}

func search(*ui.Button) {
	var translateSites = map[*ui.Label]base.Translate{
		dictLabel:    site.Dict{},
		cibaLabel:    site.Ciba{},
		shanbeiLabel: site.Shanbei{},
		youdaoLabel:  site.Youdao{},
	}
	for label, s := range translateSites {
		text := inputEntry.Text()
		fmt.Println(text)
		go getResult(label, s, text)
	}
}
func getResult(label *ui.Label, s base.Translate, text string) {
	_, translates, phonetics, audios := s.GetTranslate(text)
	if s.GetName() == "dict" && audios != nil {
		base.Play("http://audio.dict.cn/" + audios[0])
	}
	ui.QueueMain(func() {
		label.SetText(translates + "\n" + phonetics)
	})
	if s.GetName() == "youdao" && (label.Text() == "\n\n\n" || label.Text() == "  渣渣有道词典又罢工啦！\n\n"){
		ui.QueueMain(func() {
			label.SetText("  渣渣有道词典又罢工啦！\n\n")
		})
	}
}
