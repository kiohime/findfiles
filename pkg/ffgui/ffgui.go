package ffgui

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"github.com/kiohime/findfiles/pkg/fzscreen"
)

var (
	progressBar             *widget.ProgressBarInfinite
	entry_widget            *modifiedEntry
	global_screen_container *fyne.Container
)

// //////////////////////////////////////////////////////////////////////

var AppData *Data

type Data struct {
	InputData []string
	PrintData []string
}

// ///////////////СТРУКТУРА НАСТРОЙКИ ПРИЛОЖЕНИЯ/////////////////////////////

// AppSet - глобальная переменная с настройками приложения
var AppSet *Settings

// Settings - глобальная структура с настройками приложения
type Settings struct {
	AppMode        int
	ScanMode       int
	RootDir        string
	WorkDir        string
	BaseNameFiles  string
	BaseNameDirs   string
	TargetFileName string
}

// блок инициализации: установка рабочего пути для файлов базы и поиска
func Initialize(aset *Settings) error {
	fmt.Println("### initialize")

	aset.BaseNameFiles = "filesearch_files.txt"
	aset.BaseNameDirs = "filesearch_directory.txt"

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("UserHomeDir error: %v", err)
	}
	filePath := filepath.Join(homeDir, ".config", "kiohime", "file.txt")
	filePathDir := filepath.Dir(filePath)
	aset.WorkDir = filePathDir + "\\"
	err = os.MkdirAll(filePathDir, 0777)
	// 0666 for files
	if err != nil {
		return fmt.Errorf("Database error: %v", err)
	}
	fmt.Println("### END initialize")
	return nil
}

//bahniFile - создает файл с именем inputName из массива данных inputData, данные добавляются через ньюлайн
func bahniFile(fName string, fData *[]string) error {
	fmt.Println("#### bahnifile")
	fmt.Println(fName)
	// fmt.Println(inputData)

	// err := os.Remove(inputName)
	// создание файла по полному пути, вставка значений из кэша отрисовки с обрезкой лишняка
	file, err := os.OpenFile(fName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		fmt.Println("#### END bahnifile")
		return fmt.Errorf("bahni file 1 : %v", err)
	}
	s := strings.Join(*fData, "\n")
	_, e := file.WriteString(s)

	if e != nil {
		err = e
		fmt.Println("#### END bahnifile")
		return fmt.Errorf("bahni file 2 : %v", err)
	}

	// закрытие файла
	err = file.Close()
	if err != nil {
		fmt.Println("#### END bahnifile")
		return fmt.Errorf("bahni file 3 : %v", err)
	}
	return nil
}

func printFile(pathname string) *[]string {
	p := []string{"1", "2"}
	return &p
}

type fzScreenWidget struct {
	Entry  *widget.Entry
	Screen *fyne.Container
}

func (c *fzScreenWidget) Update(s []string) {
	global_screen_container.Objects = nil
	pattern := c.Entry.Text
	result := fzscreen.Render(s, pattern)
	for _, obj := range result.Objects {
		global_screen_container.Objects = append(global_screen_container.Objects, obj)
	}
}

// mainDecider - executes searching mode, depending on switch
func mainDecider(input *modifiedEntry, scr *fyne.Container, aset *Settings, adata *Data) {
	assert(input, scr)
	fmt.Printf("on decider AppMode is %v\n", aset.AppMode)

	progressBar.Start()
	progressBar.Show()
	adata.InputData = nil
	adata.InputData = append(adata.InputData, input.Text)
	resultData := executer(aset, adata)
	resultScreen := fyne.NewContainer()

	texter := &fzScreenWidget{
		Entry:  &input.Entry,
		Screen: resultScreen,
	}
	input.Entry.OnCursorChanged = func() { texter.Update(resultData) }

	// newScr :=

	// scr.Text = result
	input.Text = ""
	input.Refresh()
	progressBar.Hide()
	progressBar.Stop()
}

//executer - запускает программу в устновленном режиме
func executer(aset *Settings, adata *Data) []string {
	result := []string{}
	allErrs := []error{}
	fmt.Printf("on executer AppMode is %v\n", aset.AppMode)
	switch aset.AppMode {
	// режим поиска по базе...
	case 0:
		readBase := []string{}
		// allReadBase := []string{}
		readErrors := []error{}
		var rError error
		switch {
		// ...нужны и каталоги, и файлы
		case aset.ScanMode == 2:
			for i := 0; i < 2; i++ {
				aset.ScanMode = i
				switch aset.ScanMode {
				case 0:
					aset.TargetFileName = aset.BaseNameDirs
				case 1:
					aset.TargetFileName = aset.BaseNameFiles
				}
				// запуск поиска, получаем массив данных и ошибку
				// если есть ошибки при поиске, передаем их в массив ошибок поиска по базе
				rb := []string{}
				rb, rError = readBaser(aset, adata.InputData)
				readBase = append(readBase, rb...)

				// fmt.Printf("%v\t%v\n", i, rb)
				if rError != nil {
					readErrors = append(readErrors, rError)
					aset.ScanMode = 2
					return nil
				}
			}
			// fmt.Println("**************", readBase)
			aset.ScanMode = 2

		// ...нужны или каталоги, или файлы отдельно
		case aset.ScanMode == 0 || aset.ScanMode == 1:
			// запуск поиска, получаем массив данных
			// если есть ошибки при поиске, передаем их в массив ошибок поиска по базе
			readBase, rError = readBaser(aset, adata.InputData)
			if rError != nil {
				readErrors = append(readErrors, rError)
				return nil
			}
		}

		// составляем строку из массива данных для вывода в экран гуи
		// result = strings.Join(readBase, "\n")
		result = readBase

	// /////////////////////////////////////
	// режим сканирования пути...
	case 1:
		writeBase := []string{}
		writeErrors := []error{}
		var wError error
		switch {
		// ...нужны и каталоги, и файлы
		case aset.ScanMode == 2:
			for i := 0; i < 2; i++ {
				aset.ScanMode = i
				switch aset.ScanMode {
				case 0:
					aset.TargetFileName = aset.BaseNameDirs
				case 1:
					aset.TargetFileName = aset.BaseNameFiles
				}
				// запуск сканирования, получаем массив данных и ошибку
				// если есть ошибки при сканировании, передаем их в массив ошибок сканирования пути
				writeBase, wError = writeBaser(aset, adata.InputData)
				if wError != nil {
					writeErrors = append(writeErrors, wError)
					aset.ScanMode = 2
					return nil
				}
				// создаем файл результат
				// если возникла ошибка при создании файла результата сканирования, добавляем ошибку в массив ошибок сканирования пути
				write := aset.WorkDir + aset.TargetFileName
				writeFileError := bahniFile(write, &writeBase)
				if writeFileError != nil {
					wfe := fmt.Errorf("error in creating file : %v", writeFileError)
					writeErrors = append(writeErrors, wfe)
				}
			}
			aset.ScanMode = 2

		// ...нужны или каталоги, или файлы отдельно
		case aset.ScanMode == 0 || aset.ScanMode == 1:
			// запуск сканирования, получаем массив данных и ошибку
			// если есть ошибки при сканировании, передаем их в массив ошибок сканирования пути
			writeBase, wError = writeBaser(aset, adata.InputData)
			if wError != nil {
				writeErrors = append(writeErrors, wError)
				return nil
			}
			// создаем файл результат
			// если возникла ошибка при создании файла результата сканирования, добавляем ошибку в массив ошибок сканирования пути
			write := aset.WorkDir + aset.TargetFileName
			writeFileError := bahniFile(write, &writeBase)
			if writeFileError != nil {
				wfe := fmt.Errorf("error in creating file : %v", writeFileError)
				writeErrors = append(writeErrors, wfe)
			}
		}
		// добавляем массив ошибок сканирования пути в единый массив экзекутора
		allErrs = append(allErrs, writeErrors...)

	}
	if len(allErrs) > 0 {
		var err string
		// fmt.Printf("error in walking : %q\n", err)
		for i, e := range allErrs {
			err += fmt.Sprintf("	%v : %v\n", i, e)
		}
		fmt.Printf(err)
	}
	fmt.Println("THE END")
	return result
}

// assert - проверяет на наличие критических ошибок
func assert(args ...interface{}) {
	okGlobal := true
	msg := "assertion failed:\n"
	for i, arg := range args {
		// fmt.Printf("____%v: %v\n", i, arg)
		ok := true
		switch t := arg.(type) {
		case bool:
			ok = t
		case error:
			ok = t == nil
		case nil:
			ok = false
			fmt.Printf("NILNIL\n")
		}
		if reflect.ValueOf(arg).IsNil() {
			ok = false
		}

		if !ok {
			_, file, line, _ := runtime.Caller(1)
			msg += fmt.Sprintf("\t(%v, %T) %v: %v\n", i, arg, file, line)
		}
		okGlobal = okGlobal && ok
	}
	if !okGlobal {
		panic(msg)
	}
}

// func main() {
// 	AppSet = &Settings{}
// 	AppData = &Data{}
// 	err := Initialize(AppSet)
// 	if err != nil {
// 		fmt.Printf("Error in initialisation : %q\n", err)
// 		KeyWait()
// 		os.Exit(1)
// 	}
// 	Gui(AppSet, AppData)
// }
