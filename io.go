package main

import (
	"bufio"
	"os"
	"sync"

	"github.com/ActiveState/tail"
	"github.com/lmorg/apachelogs"
)

func ReadStdIn() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		s := scanner.Text()

		access, err, matched := apachelogs.ParseAccessLine(s)

		if err == nil && matched {
			access.FileName = "<STDIN>"
			PrintAccessLogs(access)

		} else if err != nil {
			errOut(err)
		}
	}
}

func ReadFileStream(filename string, wg *sync.WaitGroup) {
	defer wg.Done()

	t, err := tail.TailFile(filename, tail.Config{Follow: true})
	if err != nil {
		// TODO: this is shit!!
		panic(err)
	}
	for line := range t.Lines {
		if line.Err != nil {
			errOut(line.Err)
		}
		apachelogs.ParseAccessLine(line.Text)
	}
}

func ReadFileStatic(filename string, wg *sync.WaitGroup) {
	defer wg.Done()

	apachelogs.ReadAccessLog(filename, PrintAccessLogs, errOut)
}

func errOut(err error) {
	if fNoErrors {
		return
	}

	PrintStdError(err.Error())
}
