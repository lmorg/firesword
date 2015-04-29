package main

import (
	"apachelogs"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

// command line field names
const (
	CLI_STR_IP       = "ip"
	CLI_STR_METHOD   = "method"
	CLI_STR_PROC     = "proc"
	CLI_STR_PROTO    = "proto"
	CLI_STR_QS       = "qs"
	CLI_STR_REF      = "ref"
	CLI_STR_SIZE     = "size"
	CLI_STR_STATUS   = "status"
	CLI_STR_STITLE   = "stitle"
	CLI_STR_SDESC    = "sdesc"
	CLI_STR_TIME     = "time"
	CLI_STR_DATE     = "date"
	CLI_STR_DATETIME = "datetime"
	CLI_STR_EPOCH    = "epoch"
	CLI_STR_URI      = "uri"
	CLI_STR_UA       = "ua"
	CLI_STR_UID      = "uid"
	CLI_STR_FILE     = "file"
)

// defaults
var (
	cli_sl_ip       int = -15
	cli_sl_method   int = -4
	cli_sl_proc     int = 9
	cli_sl_proto    int = -8
	cli_sl_qs       int = -20
	cli_sl_ref      int = 30
	cli_sl_size     int = 9
	cli_sl_status   int = -3
	cli_sl_stitle   int = -20
	cli_sl_sdesc    int = -40
	cli_sl_time     int = -8
	cli_sl_date     int = -11
	cli_sl_datetime int = -20
	cli_sl_epoch    int = 10
	cli_sl_uri      int = 40
	cli_sl_ua       int = -20
	cli_sl_uid      int = -10
	cli_sl_file     int = -10
)

// CLI main()
func cliInterface() {
	var wg sync.WaitGroup

	ImportStrLen()

	if f_read_stdin {
		wg.Add(1)
		go ReadSTDIN()

	} else if f_file_stream != "" {
		wg.Add(1)
		go cliWrapper_ReadFileStream(f_file_stream, &wg)

	} else if len(f_files_static) > 0 {

		for i := 0; i < len(f_files_static); i++ {
			wg.Add(1)
			go cliWrapper_ReadFileStatic(f_files_static[i], &wg)
		}

	} else {
		fmt.Println("No input files given. Run with -h for help")
		os.Exit(1)
	}

	wg.Wait()
}

// Cycles through the stdout format string looking for string lengths.
// Runs on start of application.
func ImportStrLen() {
	getStrLen(CLI_STR_IP, &cli_sl_ip)
	getStrLen(CLI_STR_METHOD, &cli_sl_method)
	getStrLen(CLI_STR_PROC, &cli_sl_proc)
	getStrLen(CLI_STR_PROTO, &cli_sl_proto)
	getStrLen(CLI_STR_QS, &cli_sl_qs)
	getStrLen(CLI_STR_REF, &cli_sl_ref)
	getStrLen(CLI_STR_SIZE, &cli_sl_size)
	getStrLen(CLI_STR_STATUS, &cli_sl_status)
	getStrLen(CLI_STR_STITLE, &cli_sl_stitle)
	getStrLen(CLI_STR_SDESC, &cli_sl_sdesc)
	getStrLen(CLI_STR_TIME, &cli_sl_time)
	getStrLen(CLI_STR_DATE, &cli_sl_date)
	getStrLen(CLI_STR_EPOCH, &cli_sl_epoch)
	getStrLen(CLI_STR_URI, &cli_sl_uri)
	getStrLen(CLI_STR_UA, &cli_sl_ua)
	getStrLen(CLI_STR_UID, &cli_sl_uid)
	getStrLen(CLI_STR_FILE, &cli_sl_file)
}

// Gets the string length for each item and then reformat stdout format string
// for tighter output loops
func getStrLen(item string, val *int) {
	rx, _ := regexp.Compile(`{` + item + `,(\-?[0-9]+)}`)
	found := rx.FindAllStringSubmatch(f_stdout_fmt, -1)

	if len(found) == 0 || len(found[0]) == 0 {
		return
	}

	*val, _ = strconv.Atoi(found[0][1])
	//if *val == "0" {
	//	*val = ""
	//}

	for i, _ := range found[0] {
		f_stdout_fmt = strings.Replace(
			f_stdout_fmt,
			fmt.Sprintf("{%s,%s}", item, found[0][i]),
			fmt.Sprintf("{%s}", item),
			-1)
	}
}

/*func CompilePrintf(item string, val *int) {
	rx, _ := regexp.Compile(`{` + item + `,(\-?[0-9]+)}`)
	found := rx.FindAllStringSubmatch(f_stdout_fmt, -1)

	if len(found) == 0 || len(found[0]) == 0 {
		return
	}

	*val, _ = strconv.Atoi(found[0][1])

	for i, _ := range found[0] {
		f_stdout_fmt = strings.Replace(
			f_stdout_fmt,
			fmt.Sprintf("{%s,%s}", item, found[0][i]),
			fmt.Sprintf("{%s}", item),
			-1)
	}
}
*/
func cliWrapper_ReadFileStream(filename string, wg *sync.WaitGroup) {
	defer wg.Done()
	ReadFileStream(filename)
}

func cliWrapper_ReadFileStatic(filename string, wg *sync.WaitGroup) {
	defer wg.Done()
	ReadFileStatic(filename)
}

func PrintAccessLogs(access *apachelogs.AccessLog, filename string) {
	s := f_stdout_fmt
	formatSTDOUT(&s, CLI_STR_IP, access.IP, cli_sl_ip)
	formatSTDOUT(&s, CLI_STR_METHOD, access.Method, cli_sl_method)
	formatSTDOUT(&s, CLI_STR_PROC, strconv.Itoa(access.ProcTime), cli_sl_proc)
	formatSTDOUT(&s, CLI_STR_PROTO, access.Protocol, cli_sl_proto)
	formatSTDOUT(&s, CLI_STR_QS, access.QueryString, cli_sl_qs)
	formatSTDOUT(&s, CLI_STR_REF, access.Referrer, cli_sl_ref)
	formatSTDOUT(&s, CLI_STR_SIZE, strconv.Itoa(access.Size), cli_sl_size)
	formatSTDOUT(&s, CLI_STR_STATUS, access.Status.A, cli_sl_status)
	formatSTDOUT(&s, CLI_STR_STITLE, access.Status.Title(), cli_sl_stitle)
	formatSTDOUT(&s, CLI_STR_SDESC, access.Status.Description(), cli_sl_sdesc)
	formatSTDOUT(&s, CLI_STR_TIME, access.DateTime.Format(FMT_TIME), cli_sl_time)
	formatSTDOUT(&s, CLI_STR_DATE, access.DateTime.Format(FMT_DATE), cli_sl_date)
	formatSTDOUT(&s, CLI_STR_EPOCH, strconv.FormatInt(access.DateTime.Unix(), 10), cli_sl_epoch)
	formatSTDOUT(&s, CLI_STR_URI, access.URI, cli_sl_uri)
	formatSTDOUT(&s, CLI_STR_UA, access.UserAgent, cli_sl_ua)
	formatSTDOUT(&s, CLI_STR_UID, access.UserID, cli_sl_uid)
	formatSTDOUT(&s, CLI_STR_FILE, filename, cli_sl_file)
	fmt.Println(s)
}

func formatSTDOUT(source *string, search, value string, length int) {
	slength := strconv.Itoa(length)
	if slength == "0" {
		slength = ""
	}

	value = Trim(fmt.Sprintf("%"+slength+"s", value), length)
	*source = strings.Replace(*source, "{"+search+"}", value, -1)
}
