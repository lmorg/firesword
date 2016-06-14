package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"

	"github.com/lmorg/apachelogs"
)

// App versioning
const (
	APP_NAME  = "Plasmasword"
	VERSION   = "2.00.600 BETA"
	COPYRIGHT = "© 2014-2016 Laurence Morgan"
)

// Date / time output formatting
const (
	FMT_DATE     = "02 Jan 2006"
	FMT_TIME     = "15:04:05"
	FMT_DATETIME = FMT_DATE + " " + FMT_TIME
)

// Command line flags
var (
	// Global
	f_no_smp    bool
	f_no_errors bool

	// CLI interface
	f_stdout_fmt string
	f_patterns   string
	f_trim_slash bool

	// Input streams
	f_read_stdin   bool
	f_files_stream FlagStrings
	f_files_static FlagStrings

	// Output streams
	f_write_sqlite string

	// Usage
	f_help1, f_help2, f_help_f, f_help_g, f_version1, f_version2 bool

	// Lazy fix to check if compiled with ncurses.
	// Ncurses can be enabled or disabled via '// +build ignore' (without quotes)
	// at the top of ncurses.go
	//
	// Ncurses mode also requires sqlite and readline - so compiling with ncurses
	// breaks cross-compiling portability for the sake of extra features.
	//ncurses_compiled bool
)

type FlagStrings []string

func (fs *FlagStrings) String() string         { return fmt.Sprint(*fs) }
func (fs *FlagStrings) Set(value string) error { *fs = append(*fs, value); return nil }

func flags() {
	flag.Usage = Usage

	// global
	flag.BoolVar(&f_no_errors, "no-errors", false, "")

	// CLI interface
	flag.StringVar(&f_stdout_fmt, "fmt", "{ip} {uri} {status} {stitle}", "")
	flag.StringVar(&f_patterns, "grep", "", "")
	flag.BoolVar(&f_trim_slash, "trim-slash", false, "")

	// Input streams
	flag.BoolVar(&f_read_stdin, "stdin", false, "")
	flag.Var(&f_files_stream, "f", "")

	// Output streams
	flag.StringVar(&f_write_sqlite, "sqlout", "", "")

	// help
	flag.BoolVar(&f_help1, "h", false, "")
	flag.BoolVar(&f_help2, "?", false, "")
	flag.BoolVar(&f_help_f, "hf", false, "")
	flag.BoolVar(&f_help_g, "hg", false, "")
	flag.BoolVar(&f_version1, "v", false, "")
	flag.BoolVar(&f_version2, "version", false, "")

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

	cliInterface()
}
