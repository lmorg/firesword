package main

import (
	"fmt"
	"github.com/gizak/termui"
	"github.com/lmorg/apachelogs"
	"github.com/lmorg/firesword/sqlite"
	"github.com/shavac/readline"
	"os"
	"os/user"
	"strconv"
	"strings"
	"sync"
	"time"
)

const ERROR_HEIGHT = 5

var (
	errors       *termui.List
	list         *termui.List
	bar_chart    *termui.BarChart
	line_chart   *termui.LineChart
	UseBar       bool
	UseLine      bool
	history_file string
)

func init() {
	usr, _ := user.Current()
	history_file = usr.HomeDir + "/." + strings.ToLower(APP_NAME) + "_history"
}

func nAddError(msg string) {
	if len(errors.Items) >= ERROR_HEIGHT-2 {
		errors.Items = append(
			errors.Items[len(errors.Items)-(ERROR_HEIGHT-2-1):],
			"["+time.Now().Format(FMT_DATE+" "+FMT_TIME)+"] "+msg,
		)

	} else {
		errors.Items = append(
			errors.Items,
			"["+time.Now().Format(FMT_DATE+" "+FMT_TIME)+"] "+msg,
		)
	}
}

func nInit() {
	defer func() {
		if r := recover(); r != nil {
			termui.Close()
			fmt.Println("Panic caught in nInit:", r)
			os.Exit(2)
		}
	}()

	// start ncurses
	err := termui.Init()
	if err != nil {
		//panic(err)
		fmt.Println("Error:", err)
		os.Exit(2)
	}
	termui.UseTheme("helloworld")

	// format display.
	// dynamic sizing go in nRender()
	errors = termui.NewList()
	errors.Height = ERROR_HEIGHT
	errors.HasBorder = true
	errors.Border.Label = "Errors:"
	errors.Items = []string{
		fmt.Sprintf("%s %s", APP_NAME, VERSION),
		COPYRIGHT,
	}

	list = termui.NewList()
	list.HasBorder = true
	list.Border.Label = "Results:"

	// graphs:
	bar_chart = termui.NewBarChart()
	bar_chart.Border.Label = "Graph:"
	bar_chart.BarWidth = 6
	line_chart = termui.NewLineChart()
	line_chart.HasBorder = false
	line_chart.X = 2
}

func nRender() {
	defer func() {
		if r := recover(); r != nil {
			nAddError(fmt.Sprint("Panic caught in nRender:", r))
		}
	}()

	// dynamic sizing:
	errors.Y = termui.TermHeight() - ERROR_HEIGHT
	errors.Width = termui.TermWidth()

	if UseBar {
		bar_chart.Width = termui.TermWidth()
		bar_chart.Height = termui.TermHeight() - ERROR_HEIGHT
		termui.Render(bar_chart, errors)

	} else if UseLine {
		list.Width = termui.TermWidth()
		list.Height = (termui.TermHeight() / 2) - ERROR_HEIGHT
		line_chart.Height = list.Height - 2
		line_chart.Width = termui.TermWidth() - 3
		line_chart.Y = termui.TermHeight() - line_chart.Height - ERROR_HEIGHT - 1
		termui.Render(line_chart, list, errors)

	} else {
		list.Width = termui.TermWidth()
		list.Height = termui.TermHeight() - ERROR_HEIGHT
		termui.Render(list, errors)
	}
}

func nQuit() {
	termui.Close()
	sqlite.Close()
	if err := readline.WriteHistoryFile(history_file); err != nil {
		fmt.Printf("Cannot write history file (%s): %s\n", history_file, err)
	}
}

func nInterface() {
	var wg sync.WaitGroup
	sqlite.New()

	// start ncurses
	nInit()

	// load readline history
	if err := readline.ReadHistoryFile(history_file); err != nil {
		nAddError(fmt.Sprintf("Cannot open history file(%s): %s\n", history_file, err))
	}

	// catch returns and panics
	defer func() {
		if r := recover(); r != nil {
			nAddError(fmt.Sprint("Panic caught in nInterface:", r))
		} else {
			nQuit()
		}
	}()

	if len(f_files_stream) > 0 {
		for i := 0; i < len(f_files_stream); i++ {
			wg.Add(1)
			go ReadFileStaticWrapper(f_files_stream[i], &wg)
		}

	}

	if len(f_files_static) > 0 {
		for i := 0; i < len(f_files_static); i++ {
			wg.Add(1)
			go ReadFileStaticWrapper(f_files_static[i], &wg)
		}

	} else if len(f_files_stream) == 0 {
		nQuit()
		fmt.Println("No input files given. Run with -h for help")
		os.Exit(1)
	}

	go nListener()

	for {
		list.Items = make([]string, 1)
		rows, err := sqlite.Query(f_sql)
		if err != nil {
			nAddError(err.Error())

		} else {

			cols, _ := rows.Columns()
			for i := 0; i < len(cols); i++ {
				list.Items[0] += nPrintfModifier(cols[i], cols[i]) + " "
			}

			var (
				record int
				data   []int
				labels []string
				plots  []float64
			)

			for rows.Next() {
				record++
				if (!UseBar && !UseLine && record+ERROR_HEIGHT > termui.TermHeight()) || record >= 200 {
					break
				}

				//if record+ERROR_HEIGHT > termui.TermHeight() {
				//	break
				//}

				err := sqlite.RowsFetch(rows, len(cols))
				if err != nil {
					nAddError(err.Error())

				} else {
					var s string
					for i := 0; i < len(cols); i++ {
						s += nPrintfModifier(cols[i], sqlite.GetField(i)) + " "
					}
					list.Items = append(list.Items, s)

					if UseBar {
						// Slower but robust:
						//atoi, _ := strconv.Atoi(sqlite.GetField(0))
						//data = append(data, atoi)
						// Quicker, but no erorr handling:
						data = append(data, int(sqlite.Field[0].(int64)))
						labels = append(labels, sqlite.GetField(1))
					} else if UseLine {
						// Slower but robust:
						// build line graph
						//atoi, _ := strconv.Atoi(sqlite.GetField(0))
						//plots = append(plots, atoi)
						// Quicker, but no erorr handling:
						plots = append(plots, float64(sqlite.Field[0].(int64)))
					}
				}
			}

			if UseBar {
				bar_chart.Data = data
				bar_chart.DataLabels = labels
			} else if UseLine {
				line_chart.Data = plots
			}
		}

		// render
		nRender()
		time.Sleep(time.Second * time.Duration(f_refresh))
	}
	//}
}

func nListener() {
	evt := termui.EventCh()
	for {
		e := <-evt
		if e.Type == termui.EventKey {
			switch e.Ch {
			case 'q':
				nQuit()
				os.Exit(0)

			case 'i':
				//errors_backup := errors
				termui.Close()
				if s := Readline(); s != "" {
					f_sql = s
				}
				nInit()
				//errors = errors_backup
				evt = termui.EventCh()
				nRender()

			case 'b':
				UseBar = !UseBar
				UseLine = false

			case 'l':
				UseLine = !UseLine
				UseBar = false
			}
		}
	}
}

func nInsert(access *apachelogs.AccessLog) {
	var err error

	err = sqlite.InsertAccess(*access)
	if err == nil {
		return
	}

	// Bit of a kludge, but on error try 5 more times in staggered intervals.
	// This gets around most of the locking issues but isn't fool proof.
	go func() {
		for i := 0; i < 5; i++ {
			// stagger retries
			time.Sleep(time.Second / 3)

			err = sqlite.InsertAccess(*access)
			if err == nil {
				return
			}
		}

		nAddError(err.Error())
	}()

}

func Readline() (s string) {
	defer func() {
		s = strings.Replace(s, "\n", " ", -1)
		s = strings.TrimSpace(s)
		fmt.Print("\n")
	}()

	prompt := "SQL> "
	//loop until ReadLine returns nil (signalling EOF)

	for {
		switch result := readline.ReadLine(&prompt); true {
		case result == nil: // ^d quit
			return
		case *result != "": //ignore blank lines
			s += *result
			readline.AddHistory(*result) //allow user to recall this line
		}
	}

	return
}

func nPrintfModifier(heading, value string) string {
	var i int

	switch heading {
	case FIELD_IP:
		i = len_ip
	case FIELD_METHOD:
		i = len_method
	case FIELD_PROC:
		i = len_proc
	case FIELD_PROTO:
		i = len_proto
	case FIELD_QS:
		i = len_qs
	case FIELD_REF:
		i = len_ref
	case FIELD_SIZE:
		i = len_size
	case FIELD_STATUS:
		if len(value) == 3 {
			return "  " + value + " "
		} else {
			return value
		}
	case FIELD_STITLE:
		i = len_stitle
	case FIELD_SDESC:
		i = len_sdesc
	case FIELD_TIME:
		i = len_time
	case FIELD_DATE:
		i = len_date
	case FIELD_DATETIME:
		i = len_datetime
	case FIELD_EPOCH:
		i = len_epoch
	case FIELD_URI:
		i = len_uri
	case FIELD_UA:
		i = len_ua
	case FIELD_UID:
		i = len_uid
	case FIELD_FILE:
		i = len_file
	case "id":
		i = -7
	case "sql":
		i = -200
	case "#":
		i = 7
	default:
		i = -20
	}

	return fmt.Sprintf("%"+strconv.Itoa(i)+"s", Trim(value, i))
}
