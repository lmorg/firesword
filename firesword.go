package main

import (
	"flag"
	"fmt"
	"github.com/lmorg/apachelogs"
	"os"
	"runtime"
)

const (
	APP_NAME  = "Firesword"
	VERSION   = "0.8.431 BETA"
	COPYRIGHT = "© 2014-2015 Laurence Morgan"

	FMT_DATE = "02 Jan 2006"
	FMT_TIME = "15:04:05"
)

// flags
var (
	// global
	f_no_smp    bool
	f_no_errors bool

	// Ncurses interface
	f_ncurses bool
	f_sql     string
	f_refresh int64

	// CLI interface
	f_stdout_fmt string
	f_patterns   string

	// Input streams
	f_read_stdin   bool
	f_files_stream FlagStrings
	f_files_static FlagStrings

	// Usage
	f_help1, f_help2, f_help_f, f_help_g, f_version1, f_version2 bool
)

type FlagStrings []string

func (fs *FlagStrings) String() string         { return fmt.Sprint(*fs) }
func (fs *FlagStrings) Set(value string) error { *fs = append(*fs, value); return nil }

func flags() {
	flag.Usage = Usage

	// global
	flag.BoolVar(&f_no_smp, "no-smp", false, "GOMAXPROCS")
	flag.BoolVar(&f_no_errors, "no-errors", false, "surpress errors")

	// Ncurses interface
	flag.BoolVar(&f_ncurses, "n", false, "Ncurses interface")
	flag.Int64Var(&f_refresh, "r", 1, "Ncursers refresh rate")
	flag.StringVar(&f_sql, "sql", "SELECT * FROM default_view", "")

	// CLI interface
	flag.StringVar(&f_stdout_fmt, "fmt", "{ip} {uri} {status} {stitle}", "STDOUT format")
	flag.StringVar(&f_patterns, "grep", "", "filter results")

	// Input streams
	flag.BoolVar(&f_read_stdin, "stdin", false, "")
	flag.Var(&f_files_stream, "f", "tail -f stream")

	// help
	flag.BoolVar(&f_help1, "h", false, "help")
	flag.BoolVar(&f_help2, "?", false, "Same as -h")
	flag.BoolVar(&f_help_f, "hf", false, "format field names")
	flag.BoolVar(&f_help_g, "hg", false, "grep guide")
	flag.BoolVar(&f_version1, "v", false, "version")
	flag.BoolVar(&f_version2, "version", false, "same as -v")

	flag.Parse()
	f_files_static = flag.Args()
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Panic caught:", r)
			os.Exit(2)
		}
	}()

	flags()

	if f_help1 || f_help2 {
		flag.Usage()
		os.Exit(1)
	}

	if f_help_f || f_help_g {
		HelpDetail()
		os.Exit(1)
	}

	if f_version1 || f_version2 {
		fmt.Println(APP_NAME, VERSION, "\n"+COPYRIGHT)
		os.Exit(1)
	}

	if !f_no_smp {
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
