package main

import (
	"github.com/lmorg/apachelogs"
	"bufio"
	"fmt"
	"github.com/ActiveState/tail"
	"os"
	"sync"
)

func ReadFileStreamWrapper(filename string, wg *sync.WaitGroup) {
	defer wg.Done()
	ReadFileStream(filename)
}

func ReadFileStaticWrapper(filename string, wg *sync.WaitGroup) {
	defer wg.Done()

	if f_ncurses {
		apachelogs.ReadAccessLog(filename, nInsert, errOut)
	} else {
		apachelogs.ReadAccessLog(filename, PrintAccessLogs, errOut)
	}

}

func ReadSTDIN() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		// TODO: investigate if scanner.Text saves a string(b) later
		b := scanner.Bytes()

		access, err, matched := apachelogs.ParseAccessLine(string(b))

		if err == nil && matched {
			access.FileName = "/dev/stdin"
			if f_ncurses {
				nInsert(access)

			} else {
				PrintAccessLogs(access)
			}
		}
	}
}

func ReadFileStream(filename string) {
	t, err := tail.TailFile(filename, tail.Config{Follow: true})
	if err != nil {
		// TODO: this is shit.
		panic(err)
	}
	for line := range t.Lines {
		apachelogs.ParseAccessLine(line.Text)
	}
}

func errOut(err error) {
	if f_ncurses {
		nAddError(err.Error())
	} else {
		fmt.Println("ERROR:", err)
	}
}
