package main

import (
	"bufio"
	"fmt"
	"os"

	"fyne.io/fyne/v2"
)

// пауза
func keyWait() {
	fmt.Printf("Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

// TypedKey - overriding default TypedKey method in fyne.Focusable, adding switch
// ТУТ ПРЕСЕТЫ ДЛЯ КЛАВИАТУРЫ
func (e *modifiedEntry) TypedKey(key *fyne.KeyEvent) {
	switch key.Name {
	case fyne.KeyEscape:
		e.onEsc()
	case fyne.KeyEnter, fyne.KeyReturn:
		e.onEnter(AppSet, AppData)
	default:
		e.Entry.TypedKey(key)
	}
}
