package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

//readBaser - reading mode : сканирует существующую базу, возвращает данные, соответствующие запросу и ошибку
func readBaser(aset *Settings, idata []string) ([]string, error) {
	result := []string{}
	fmt.Println("### readbaser")

	// проверка на кол-во аргументов на 24 марта 2021 аргумент всегда один,
	argLen := len(idata)
	// fmt.Println("11111  ", argLen)
	fmt.Println()
	if argLen == 0 {
		err := errors.New("no search arguments was inputed")
		return nil, err
	}
	// так что эта часть кода вообщето не нужна?

	base := aset.WorkDir + aset.TargetFileName
	// fmt.Println("base", base)
	f, err := os.Open(base)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	// Splits on newlines by default.
	scanner := bufio.NewScanner(f)
	// countScan := 0
	for scanner.Scan() {
		// fmt.Println("######################################")
		// fmt.Println("counter is ", countScan)
		line := scanner.Text()
		line = strings.TrimSpace(line)
		lineFile := filepath.Base(line)

		for _, a := range idata {
			// fmt.Println("argument is ", a)
			if strings.Contains(lineFile, a) {
				result = append(result, line)
			}
		}
	}

	return result, err
}
