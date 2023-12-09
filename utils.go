package main

import (
	"bufio"
	"os"
	"strings"
)

func CleanSlice(s []string) []string {
	unique := make(map[string]bool, len(s))
	var us []string
	for _, elem := range s {
		if len(elem) != 0 {
			if !unique[elem] {
				us = append(us, elem)
				unique[elem] = true
			}
		}
	}

	return us
}

func FileLines(str string) ([]string, error) {
	txtlines := []string{}
	file, err := os.Open(str)

	if err != nil {
		return txtlines, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		txtlines = append(txtlines, scanner.Text())
	}

	return txtlines, nil
}

func Slice2Lowercase(s []string) []string {
	arr := make([]string, len(s))
	for _, elem := range s {
		arr = append(arr, strings.ToLower(elem))
	}
	return arr
}
