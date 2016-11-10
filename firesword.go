package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/lmorg/apachelogs"
	"sync"
)

// App versioning
const (
	AppName   = "Firesword"
	Version   = "3.00.0800"
	Copyright = "© 2014-2016 Laurence Morgan"
)

// Date / time output formatting
const (
	DateFormat     = "02 Jan 2006"
	TimeFormat     = "15:04:05"
	DateTimeFormat = DateFormat + " " + TimeFormat
)

// Command line flags
var (
	// Global
	fNoErrors bool

	// CLI interface
	fStdOutFormat string
	fPatterns     string
	fTrimSlash    bool
	fGroupBy      int

	// Input streams
	fReadStdIn   bool
	fFilesStream FlagStrings
	fFilesStatic FlagStrings

	// Usage
	fHelp1, fHelp2, fHelpF, fHelpG, fVersion1, fVersion2 bool
)

type FlagStrings []string

func (fs *FlagStrings) String() string         { return fmt.Sprint(*fs) }
func (fs *FlagStrings) Set(value string) error { *fs = append(*fs, value); return nil }

func flags() {
	flag.Usage = Usage

	// global
	flag.BoolVar(&fNoErrors, "no-errors", false, "")

	// CLI interface
	flag.StringVar(&fStdOutFormat, "fmt", "{ip} {uri} {status} {stitle}", "")
	flag.StringVar(&fPatterns, "grep", "", "")
	flag.BoolVar(&fTrimSlash, "trim-slash", false, "")
	flag.IntVar(&fGroupBy, "group-by", 0, "")

	// Input streams
	flag.BoolVar(&fReadStdIn, "stdin", false, "")
	flag.Var(&fFilesStream, "f", "")

	// help
	flag.BoolVar(&fHelp1, "h", false, "")
	flag.BoolVar(&fHelp2, "?", false, "")
	flag.BoolVar(&fHelpF, "hf", false, "")
	flag.BoolVar(&fHelpG, "hg", false, "")
	flag.BoolVar(&fVersion1, "v", false, "")
	flag.BoolVar(&fVersion2, "version", false, "")

	flag.Parse()
	fFilesStatic = flag.Args()

	if fHelp1 || fHelp2 {
		flag.Usage()
		os.Exit(1)
	}

	if fHelpF || fHelpG {
		HelpDetail()
		os.Exit(1)
	}

	if fVersion1 || fVersion2 {
		fmt.Println(AppName, Version, "\n"+Copyright)
		os.Exit(1)
	}

	if fPatterns != "" {
		apachelogs.Patterns = PatternDeconstructor(fPatterns)
	}
}

func main() {
	flags()

	var wg sync.WaitGroup

	//GroupsCreate()
	cropUnusedFmt()
	ImportStrLen()

	if fReadStdIn {
		wg.Add(1)
		go ReadStdIn()

	} else if len(fFilesStream) > 0 {
		for i := 0; i < len(fFilesStatic); i++ {
			wg.Add(1)
			go ReadFileStream(fFilesStream[i], &wg)
		}

	} else if len(fFilesStatic) > 0 {
		for i := 0; i < len(fFilesStatic); i++ {
			wg.Add(1)
			go ReadFileStatic(fFilesStatic[i], &wg)
		}

	} else {
		fmt.Println("No input files given. Run with -h for help")
		os.Exit(1)
	}

	wg.Wait()
}
