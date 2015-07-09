package main

import (
	"bufio"
	"fmt"
	"github.com/ActiveState/tail"
	"github.com/lmorg/apachelogs"
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
		s := scanner.Text()

		access, err, matched := apachelogs.ParseAccessLine(s)

		if err == nil && matched {
			access.FileName = "<STDIN>"
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
