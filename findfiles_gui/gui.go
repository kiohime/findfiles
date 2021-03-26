package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/widget"
)

var app_window fyne.Window

// /////////////////////////////////////////////////////////////////////////
type modifiedEntry struct {
	widget.Entry
	input *Input_widget
}

// newEscapeEntry - rewriting basic entry widget
func newModifiedEntry() *modifiedEntry {
	entry := &modifiedEntry{}
	entry.ExtendBaseWidget(entry)
	return entry
}

func newEntry() *modifiedEntry {
	return newModifiedEntry()
}

// onEsc - clears entry
func (e *modifiedEntry) onEsc() {
	// fmt.Println(e.Entry.Text)
	e.Entry.SetText("")
}

func (e *modifiedEntry) onEnter(aset *Settings, adata *Data) {
	assert(e)
	assert(e.input)
	// findBtn(input_widget.Entry, screen_widget)
	fmt.Printf("on onEnter AppMode is %v\n", aset.AppMode)
	mainDecider(e.input.Entry, screen_widget, aset, adata)
}

func (e *modifiedEntry) setInput(input *Input_widget) {
	e.input = input
}

// /////////////////////////////////////////////////////////////
type modifiedSelect struct {
	widget.Select
	settings *Input_widget
}

func newModifiedSelect() *modifiedSelect {
	selEntry := &modifiedSelect{}
	selEntry.ExtendBaseWidget(selEntry)
	return selEntry
}

func (s *modifiedSelect) setSettings(settings *Input_widget) {
	assert(s)
	s.settings = settings
}

func newSelect() *modifiedSelect {
	s := newModifiedSelect()

	return s
}

// //////////////////////////////////////////////////////////////////

func newSettings(m *widget.RadioGroup) *modifiedSelect {
	check := newSelect()
	t := []string{"каталоги", "файлы", "И ТО, И ДРУГОЕ"}
	check.Options = t
	check.OnChanged = settingsChanged
	// check.Disable()
	return check
}

func settingsChanged(c string) {
	switch c {
	case "каталоги":
		AppSet.ScanMode = 0
		AppSet.TargetFileName = AppSet.BaseNameDirs
	case "файлы":
		AppSet.ScanMode = 1
		AppSet.TargetFileName = AppSet.BaseNameFiles
	case "И ТО, И ДРУГОЕ":
		AppSet.ScanMode = 2
	}
	fmt.Printf("on settingChanged ScanMode is %v\n", AppSet.ScanMode)
	fmt.Printf("on settingChanged TargetFileName is %v\n", AppSet.TargetFileName)

}

// ////////////////////////////////////////////////////////////////////////
type Input_widget struct {
	Mode   *widget.RadioGroup
	Entry  *modifiedEntry
	Form   *widget.Form
	Option *modifiedSelect
}

func newInputWidget(m *widget.RadioGroup, e *modifiedEntry, f *widget.Form, s *modifiedSelect) *Input_widget {
	assert(m, e, f, s)
	return &Input_widget{Mode: m, Entry: e, Form: f, Option: s}
}

// newModeWidget
func newModeWidget(aset *Settings) *widget.RadioGroup {
	s := widget.NewRadioGroup([]string{"Поиск", "Сканирование"}, func(s string) {
		switch s {
		case "Поиск":
			aset.AppMode = 0
		case "Сканирование":
			aset.AppMode = 1
		}
		fmt.Printf("set AppMode %v\n", aset.AppMode)
	})
	s.SetSelected("Поиск")
	s.Refresh()
	return s
}

// //////////////////////////////////////////////////////////////////

func newForm(i *modifiedEntry, s *widget.Label, aset *Settings, adata *Data) *widget.Form {
	assert(i, s)
	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Data", Widget: i, HintText: "input data"},
		},
		OnCancel: nil,
		OnSubmit: func() { i.onEnter(aset, adata) },

		// OnSubmit: findBtn(i, s),
	}
	i.Validator = validation.NewRegexp(`.+`, "input smthing")
	return form
}

func makeContainerTree(i *Input_widget) *fyne.Container {
	screen_container := container.NewWithoutLayout(screen_widget)
	screen_scroll_container := container.NewScroll(screen_container)
	mode_settings_container := container.NewVBox(
		i.Mode,
		widget.NewSeparator(),
		i.Option,
		// widget.NewSeparator(),
	)
	form_container := container.NewBorder(
		nil, nil, nil, nil, i.Form,
	)
	all_container := container.NewBorder(
		// top
		container.NewVBox(
			container.NewBorder(
				nil, nil, mode_settings_container, nil, form_container,
			),
			widget.NewSeparator(),
			progressBar,
		),
		nil,
		nil,
		nil,
		// other
		screen_scroll_container,
	)
	return all_container
}

func gui(aset *Settings, adata *Data) {
	the_app := app.New()
	app_window = the_app.NewWindow("Notepad")
	app_window.Resize(fyne.NewSize(1000, 500))

	progressBar = widget.NewProgressBarInfinite()
	progressBar.Hide()

	screen_widget = widget.NewLabel("")
	mode_widget := newModeWidget(aset)
	entry_widget := newEntry()
	form_widget := newForm(entry_widget, screen_widget, aset, adata)
	settings_widget := newSettings(mode_widget)
	settings_widget.Selected = "каталоги"
	settings_widget.OnChanged("каталоги")

	input_widget := newInputWidget(mode_widget, entry_widget, form_widget, settings_widget)
	entry_widget.setInput(input_widget)
	settings_widget.setSettings(input_widget)

	tree_container := makeContainerTree(input_widget)

	app_window.SetContent(tree_container)
	app_window.ShowAndRun()
}
