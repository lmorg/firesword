package main

import (
	"fmt"
	"github.com/lmorg/apachelogs"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

// command line field names
const (
	FIELD_IP       = "ip"
	FIELD_METHOD   = "method"
	FIELD_PROC     = "proc"
	FIELD_PROTO    = "proto"
	FIELD_QS       = "qs"
	FIELD_REF      = "ref"
	FIELD_SIZE     = "size"
	FIELD_STATUS   = "status"
	FIELD_STITLE   = "stitle"
	FIELD_SDESC    = "sdesc"
	FIELD_TIME     = "time"
	FIELD_DATE     = "date"
	FIELD_DATETIME = "datetime"
	FIELD_EPOCH    = "epoch"
	FIELD_URI      = "uri"
	FIELD_UA       = "ua"
	FIELD_UID      = "uid"
	FIELD_FILE     = "file"
)

// field lengths
var (
	len_ip       int = -15
	len_method   int = -4
	len_proc     int = 9
	len_proto    int = -8
	len_qs       int = -20
	len_ref      int = 30
	len_size     int = 9
	len_status   int = -3
	len_stitle   int = -20
	len_sdesc    int = -40
	len_time     int = -8
	len_date     int = -11
	len_datetime int = -20
	len_epoch    int = 10
	len_uri      int = 40
	len_ua       int = -20
	len_uid      int = -10
	len_file     int = -10
)

// CLI main()
func cliInterface() {
	var wg sync.WaitGroup

	ImportStrLen()

	if f_read_stdin {
		wg.Add(1)
		go ReadSTDIN()

	} else if len(f_files_stream) > 0 {
		for i := 0; i < len(f_files_static); i++ {
			wg.Add(1)
			go ReadFileStreamWrapper(f_files_stream[i], &wg)
		}

	} else if len(f_files_static) > 0 {
		for i := 0; i < len(f_files_static); i++ {
			wg.Add(1)
			go ReadFileStaticWrapper(f_files_static[i], &wg)
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
	getStrLen(FIELD_IP, &len_ip)
	getStrLen(FIELD_METHOD, &len_method)
	getStrLen(FIELD_PROC, &len_proc)
	getStrLen(FIELD_PROTO, &len_proto)
	getStrLen(FIELD_QS, &len_qs)
	getStrLen(FIELD_REF, &len_ref)
	getStrLen(FIELD_SIZE, &len_size)
	getStrLen(FIELD_STATUS, &len_status)
	getStrLen(FIELD_STITLE, &len_stitle)
	getStrLen(FIELD_SDESC, &len_sdesc)
	getStrLen(FIELD_TIME, &len_time)
	getStrLen(FIELD_DATE, &len_date)
	getStrLen(FIELD_EPOCH, &len_epoch)
	getStrLen(FIELD_URI, &len_uri)
	getStrLen(FIELD_UA, &len_ua)
	getStrLen(FIELD_UID, &len_uid)
	getStrLen(FIELD_FILE, &len_file)
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

	for i, _ := range found[0] {
		f_stdout_fmt = strings.Replace(
			f_stdout_fmt,
			fmt.Sprintf("{%s,%s}", item, found[0][i]),
			fmt.Sprintf("{%s}", item),
			-1)
	}
}

func PrintAccessLogs(access *apachelogs.AccessLog) {
	s := f_stdout_fmt
	formatSTDOUT(&s, FIELD_IP, access.IP, len_ip)
	formatSTDOUT(&s, FIELD_METHOD, access.Method, len_method)
	formatSTDOUT(&s, FIELD_PROC, strconv.Itoa(access.ProcTime), len_proc)
	formatSTDOUT(&s, FIELD_PROTO, access.Protocol, len_proto)
	formatSTDOUT(&s, FIELD_QS, access.QueryString, len_qs)
	formatSTDOUT(&s, FIELD_REF, access.Referrer, len_ref)
	formatSTDOUT(&s, FIELD_SIZE, strconv.Itoa(access.Size), len_size)
	formatSTDOUT(&s, FIELD_STATUS, access.Status.A, len_status)
	formatSTDOUT(&s, FIELD_STITLE, access.Status.Title(), len_stitle)
	formatSTDOUT(&s, FIELD_SDESC, access.Status.Description(), len_sdesc)
	formatSTDOUT(&s, FIELD_TIME, access.DateTime.Format(FMT_TIME), len_time)
	formatSTDOUT(&s, FIELD_DATE, access.DateTime.Format(FMT_DATE), len_date)
	formatSTDOUT(&s, FIELD_EPOCH, strconv.FormatInt(access.DateTime.Unix(), 10), len_epoch)
	formatSTDOUT(&s, FIELD_URI, access.URI, len_uri)
	formatSTDOUT(&s, FIELD_UA, access.UserAgent, len_ua)
	formatSTDOUT(&s, FIELD_UID, access.UserID, len_uid)
	formatSTDOUT(&s, FIELD_FILE, access.FileName, len_file)
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
