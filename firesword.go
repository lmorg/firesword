package main

import (
	"github.com/lmorg/apachelogs"
	"flag"
	"fmt"
	"os"
	"runtime"
)

const (
	APP_NAME  = "Firesword"
	VERSION   = "0.7.250 BETA"
	COPYRIGHT = "© 2014-2015 Laurence Morgan"

	FMT_DATE = "02 Jan 2006"
	FMT_TIME = "15:04:05"
)

// flags
var (
	// global
	f_nosmp bool
	f_debug bool

	// Ncurses interface
	f_ncurses bool
	f_sql     string
	f_refresh int64

	// CLI interface
	f_stdout_fmt string
	f_patterns   string

	// Input streams
	f_read_stdin   bool
	f_file_stream  string
	f_files_static []string

	// Usage
	f_help1, f_help2, f_help_f, f_help_g, f_version1, f_version2 bool
)

func flags() {
	flag.Usage = Usage

	// global
	flag.BoolVar(&f_nosmp, "nosmp", false, "GOMAXPROCS")
	flag.BoolVar(&f_debug, "debug", false, "debug mode")

	// Ncurses interface
	flag.BoolVar(&f_ncurses, "n", false, "Ncurses interface")
	flag.Int64Var(&f_refresh, "r", 1, "Ncursers refresh rate (seconds)")
	flag.StringVar(&f_sql, "sql", "SELECT * FROM default_view", "")

	// CLI interface
	flag.StringVar(&f_stdout_fmt, "f", "{ip} {uri} {status} {stitle}", "STDOUT format")
	flag.StringVar(&f_patterns, "grep", "", "filter results")

	// Input streams
	flag.BoolVar(&f_read_stdin, "stdin", false, "Read from STDIN")
	flag.StringVar(&f_file_stream, "file", "", "Read from file stream (filename required)")

	// help
	flag.BoolVar(&f_help1, "h", false, "Prints this message")
	flag.BoolVar(&f_help2, "?", false, "Same as -h")
	flag.BoolVar(&f_help_f, "hf", false, "Prints format field names")
	flag.BoolVar(&f_help_g, "hg", false, "Prints grep guide")
	flag.BoolVar(&f_version1, "v", false, "Prints version number")
	flag.BoolVar(&f_version2, "version", false, "Prints version number")

	flag.Parse()
	f_files_static = flag.Args()
}

func main() {
	flags()

	if f_debug {
		apachelogs.Debug = true
	}

	if f_help1 || f_help2 {
		flag.Usage()
		os.Exit(1)
	}

	if f_help_f || f_help_g {
		ManPage()
		os.Exit(1)
	}

	if f_version1 || f_version2 {
		fmt.Println(APP_NAME, VERSION, "\n"+COPYRIGHT)
		os.Exit(1)
	}

	//if f_gomaxprocs {
	if !f_nosmp {
		runtime.GOMAXPROCS(runtime.NumCPU())
	}

	if f_patterns != "" {
		apachelogs.Patterns = PatternDeconstructor(f_patterns)
	}

	if f_ncurses {
		nInterface()
	} else {
		cliInterface()
	}
}

func Trim(s string, length int) string {
	switch {
	case length > 0:
		return lTrim(s, length)

	case length < 0:
		return rTrim(s, -length)
	}

	return s
}

func lTrim(s string, length int) string {
	if len(s) <= length {
		return s
	}

	return "…" + s[len(s)-length+1:]
}

func rTrim(s string, length int) string {
	if len(s) <= length {
		return s
	}

	return s[:length-1] + "…"
}

func Round(val, rounder int) int {
	return int(val/rounder) * rounder
}
