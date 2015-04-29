package main

import (
	"apachelogs"
	"bufio"
	"compress/gzip"
	"fmt"
	"os"
)

func ReadSTDIN() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		b := scanner.Bytes()
		ParseLog("/dev/stdin", &b)
	}
}

func ReadFileStream(filename string) {
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		b := scanner.Bytes()
		ParseLog(filename, &b)
	}
}

func ReadFileStatic(filename string) {
	var (
		reader *bufio.Reader
		err    error
	)

	fi, err := os.Open(filename)
	if err != nil {
		errOut(err)
		return
	}
	defer fi.Close()

	if filename[len(filename)-3:] == ".gz" {
		fz, err := gzip.NewReader(fi)
		if err != nil {
			errOut(err)
			return
		}
		reader = bufio.NewReader(fz)
	} else {
		reader = bufio.NewReader(fi)
	}

	for {
		b, _, err := reader.ReadLine()
		if err != nil {
			if err.Error() != "EOF" {
				errOut(err)
			}
			return
		}
		ParseLog(filename, &b)
	}

	return
}

func ParseLog(filename string, b *[]byte) {
	access, err, matched := apachelogs.ParseAccessLine(b, errOut)

	if err == nil && matched {
		if f_ncurses {
			nInsert(&access, &filename)

		} else {
			PrintAccessLogs(&access, filename)
		}
	}
}

func errOut(err error) {
	if f_ncurses {
		nAddError(err.Error())
	} else {
		fmt.Println("ERROR:", err)
	}
}
