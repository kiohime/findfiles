package ffgui

import (
	"errors"
	"fmt"
	"os"

	"github.com/kiohime/findfiles/pkg/walk"
)

// ////////////////////////////////////////////////

//rootCheck - Получает переменную пути, проверет - файл или каталог. Если файл, то выдает ошибку
func rootCheck(r string) error {
	fmt.Println("## rootcheck")
	file, fileOpenErr := os.Open(r)
	if fileOpenErr != nil {
		return fileOpenErr
	}
	defer file.Close()
	fi, err := file.Stat()
	switch {
	case !fi.IsDir():
		err = errors.New("root var is file")
	default:
		fmt.Println("root var is directory")
	}
	return err
}

////////////////////////////////////////////////

//writeBaser - пишет базу
func writeBaser(aset *Settings, idata []string) ([]string, error) {
	fmt.Println("# writebaser")
	// проверка переменной пути на то, является ли та настоящим путем, если нет - остановить программу
	// root := ""
	var err error

	for _, root := range idata {
		err = rootCheck(root)
		if err != nil {
			return nil, fmt.Errorf("root checking error : %v\n", err)
		}
	}
	// err = errors.New("TEST_ERROR")

	fmt.Printf("ScanMode is %v", aset.ScanMode)
	// exportFullPath := aset.WorkDir + aset.TargetFileName

	// костыль под задел для обработки нескольких путей
	firstRoot := idata[0]

	// парсит
	result, errs := walk.StartWalk(
		[]string{firstRoot},
		walk.Options{
			WalkScanMode: aset.ScanMode,
		},
	)
	// adata.Cache = res

	if len(errs) > 0 {
		var err string
		// fmt.Printf("error in walking : %q\n", err)
		for i, e := range errs {
			err += fmt.Sprintf("	%v : %v\n", i, e)
		}
	}
	fmt.Println("# END writebaser")
	return result, nil
}

// ///////////////////////////////////////////////
