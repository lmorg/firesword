package main

import (
	"bufio"
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
	apachelogs.ReadAccessLog(filename, stdout_handler, errOut)
}

func ReadSTDIN() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		s := scanner.Text()

		access, err, matched := apachelogs.ParseAccessLine(s)

		if err == nil && matched {
			access.FileName = "<STDIN>"
			stdout_handler(access)

		} else if err != nil && !f_no_errors {
			errOut(err)
			os.Exit(2)
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
	if f_no_errors {
		return
	}

	stderr_handler(err.Error())
}
