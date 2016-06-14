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

	// Usage
	f_help1, f_help2, f_help_f, f_help_g, f_version1, f_version2 bool

	// Output handlers to manage between CLI and ncurses modes
	stdout_handler func(access *apachelogs.AccessLog)
	stderr_handler func(message string)
	main_handler   func()

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
	flag.BoolVar(&f_no_errors, "no-errors", false, "surpress errors")

	// CLI interface
	flag.StringVar(&f_stdout_fmt, "fmt", "{ip} {uri} {status} {stitle}", "STDOUT format")
	flag.StringVar(&f_patterns, "grep", "", "filter results")
	flag.BoolVar(&f_trim_slash, "trim-slash", false, "")

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

	stdout_handler = PrintAccessLogs
	stderr_handler = PrintStdError
	main_handler = cliInterface

	main_handler()
}
