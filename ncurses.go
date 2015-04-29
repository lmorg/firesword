package main

import (
	"github.com/lmorg/apachelogs"
	"github.com/lmorg/firesword/sqlite"
	"fmt"
	ui "github.com/gizak/termui"
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
	errors       *ui.List
	list         *ui.List
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
	// start ncurses
	err := ui.Init()
	ui.UseTheme("helloworld")
	if err != nil {
		panic(err)
	}

	// format display
	errors = ui.NewList()
	errors.Height = ERROR_HEIGHT
	errors.Y = ui.TermHeight() - ERROR_HEIGHT
	errors.HasBorder = true
	errors.Border.Label = "Errors:"
	errors.Items = []string{
		fmt.Sprintf("%s %s", APP_NAME, VERSION),
		COPYRIGHT,
	}

	list = ui.NewList()
	list.Height = ui.TermHeight()
	list.HasBorder = true
	list.Border.Label = "Results:"

	ui.Body.AddRows(
		ui.NewRow(
			ui.NewCol(12, 0, list),
		),
		ui.NewRow(
			ui.NewCol(12, 0, errors),
		),
	)

	// calculate layout
	ui.Body.Align()
}

func nRender() {
	ui.Body.Width = ui.TermWidth()
	ui.Body.Align()
	list.Height = ui.TermHeight() - ERROR_HEIGHT
	ui.Render(ui.Body)
}

func nQuit() {
	ui.Close()
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

	// catch panics
	defer nQuit()

	if f_file_stream != "" {
		go ReadFileStream(f_file_stream)

	} else if len(f_files_static) > 0 {
		for i := 0; i < len(f_files_static); i++ {
			wg.Add(1)
			go cliWrapper_ReadFileStatic(f_files_static[i], &wg)
		}

	} else {
		nQuit()
		fmt.Println("No input files given. Run with -h for help")
		os.Exit(1)
	}

	evt := ui.EventCh()
	for {
		select {
		case e := <-evt:
			if e.Type == ui.EventKey && e.Ch == 'q' {
				return
			}
			if e.Type == ui.EventKey && e.Ch == 'i' {
				//errors_backup := errors
				ui.Close()
				if s := Readline(); s != "" {
					f_sql = s
				}
				nInit()
				//errors = errors_backup
				evt = ui.EventCh()
				nRender()
			}
		default:
			list.Items = make([]string, 1)
			rows, err := sqlite.Query(f_sql)
			if err != nil {
				nAddError(err.Error())

			} else {

				cols, _ := rows.Columns()
				for i := 0; i < len(cols); i++ {
					list.Items[0] += nPrintfModifier(cols[i], cols[i]) + " "

				}

				record := 0
				for rows.Next() {
					record++
					if record+ERROR_HEIGHT > ui.TermHeight() {
						break
					}

					err := sqlite.RowsFetch(rows, len(cols))
					if err != nil {
						nAddError(err.Error())

					} else {
						var s string
						for i := 0; i < len(cols); i++ {
							s += nPrintfModifier(cols[i], sqlite.GetField(i)) + " "
						}
						list.Items = append(list.Items, s)
					}
				}

			}

			// render
			nRender()
			time.Sleep(time.Second * time.Duration(f_refresh))
		}
	}
}

func nInsert(access *apachelogs.AccessLog, filename *string) {
	var err error

	err = sqlite.InsertAccess(*access, *filename)
	if err == nil {
		return
	}

	// on error, try 5 more times....
	go func() {
		for i := 0; i < 5; i++ {
			// stagger retries
			time.Sleep(time.Second / 3)

			err = sqlite.InsertAccess(*access, *filename)
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

	//L:
	for {
		//switch result := readline.ReadLine(&prompt); true {
		switch result := readline.ReadLine(&prompt); true {
		case result == nil:
			//println()
			//break L //exit loop with EOF(^D)
			return
		case *result != "": //ignore blank lines
			//println(*result)
			s += *result
			readline.AddHistory(*result) //allow user to recall this line
		}
	}

	return
}

func nPrintfModifier(heading, value string) string {
	var i int

	switch heading {
	case CLI_STR_IP:
		i = cli_sl_ip
	case CLI_STR_METHOD:
		i = cli_sl_method
	case CLI_STR_PROC:
		i = cli_sl_proc
	case CLI_STR_PROTO:
		i = cli_sl_proto
	case CLI_STR_QS:
		i = cli_sl_qs
	case CLI_STR_REF:
		i = cli_sl_ref
	case CLI_STR_SIZE:
		i = cli_sl_size
	case CLI_STR_STATUS:
		//i = cli_sl_status
		if len(value) == 3 {
			return "  " + value + " "
		} else {
			return value
		}
	case CLI_STR_STITLE:
		i = cli_sl_stitle
	case CLI_STR_SDESC:
		i = cli_sl_sdesc
	case CLI_STR_TIME:
		i = cli_sl_time
	case CLI_STR_DATE:
		i = cli_sl_date
	case CLI_STR_DATETIME:
		i = cli_sl_datetime
	case CLI_STR_EPOCH:
		i = cli_sl_epoch
	case CLI_STR_URI:
		i = cli_sl_uri
	case CLI_STR_UA:
		i = cli_sl_ua
	case CLI_STR_UID:
		i = cli_sl_uid
	case CLI_STR_FILE:
		i = cli_sl_file
	case "id":
		i = -7
	case "sql":
		i = -200
	case "#":
		i = 7
	default:
		i = -20
	}

	/*if i < 0 {
		if len(value) > i*-1 {
			return val
		}
	}*/

	return fmt.Sprintf("%"+strconv.Itoa(i)+"s", Trim(value, i))
}
