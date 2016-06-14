package main

import (
	"bytes"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/lmorg/apachelogs"
)

// Command line field names.
// The code duplication is nasty, but realistically it needs to be
// both string and byte slice for performance reasons.
var (
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
	FIELD_UNIX     = "unix"
	FIELD_URI      = "uri"
	FIELD_UA       = "ua"
	FIELD_UID      = "uid"
	FIELD_FILE     = "file"

	// As odd as seems duplicating field names from the above,
	// this method is quicker out printing to STDOUT as we can
	// alternate between strings and character arrays depending
	// on the core libs.
	bFIELD_IP       = []byte(FIELD_IP)
	bFIELD_METHOD   = []byte(FIELD_METHOD)
	bFIELD_PROC     = []byte(FIELD_PROC)
	bFIELD_PROTO    = []byte(FIELD_PROTO)
	bFIELD_QS       = []byte(FIELD_QS)
	bFIELD_REF      = []byte(FIELD_REF)
	bFIELD_SIZE     = []byte(FIELD_SIZE)
	bFIELD_STATUS   = []byte(FIELD_STATUS)
	bFIELD_STITLE   = []byte(FIELD_STITLE)
	bFIELD_SDESC    = []byte(FIELD_SDESC)
	bFIELD_TIME     = []byte(FIELD_TIME)
	bFIELD_DATE     = []byte(FIELD_DATE)
	bFIELD_DATETIME = []byte(FIELD_DATETIME)
	bFIELD_EPOCH    = []byte(FIELD_EPOCH)
	bFIELD_UNIX     = []byte(FIELD_UNIX)
	bFIELD_URI      = []byte(FIELD_URI)
	bFIELD_UA       = []byte(FIELD_UA)
	bFIELD_UID      = []byte(FIELD_UID)
	bFIELD_FILE     = []byte(FIELD_FILE)
)

// Field lengths
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
	len_file     int = -20
)

// Weird variable, but this allows printf-style padding using fast
// slicing rather than slow fmt.Printf([]interface{}...)'s
var spaces = "                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                    "

// Used for concatenating byte slices with append()
// (the braces are part of the field matching in the command line string)
var (
	brace_open  []byte = []byte{'{'}
	brace_close byte   = '}'
)

// CLI main()
func cliInterface() {
	var wg sync.WaitGroup

	cropUnusedFmt()
	ImportStrLen()

	if f_write_sqlite != "" {
		sqlite.New()
	}

	if f_read_stdin {
		wg.Add(1)
		go ReadSTDIN()

	} else if len(f_files_stream) > 0 {
		for i := 0; i < len(f_files_static); i++ {
			wg.Add(1)
			go ReadFileStream(f_files_stream[i], &wg)
		}

	} else if len(f_files_static) > 0 {
		for i := 0; i < len(f_files_static); i++ {
			wg.Add(1)
			go ReadFileStatic(f_files_static[i], &wg)
		}

	} else {
		fmt.Println("No input files given. Run with -h for help")
		os.Exit(1)
	}

	wg.Wait()

	if f_write_sqlite != "" {
		sqlite.Dump(f_write_sqlite)
		sqlite.Close()
	}
}

// Check if fields are being used in --fmt. Do this up front and
// then set the field length to zero to disable if field unused.
// This saves searching through the --fmt string for each value
// against every log line outputted (slow).
func cropUnusedFmt() {
	if !strings.Contains(f_stdout_fmt, FIELD_IP) {
		len_ip = 0
	}

	if !strings.Contains(f_stdout_fmt, FIELD_METHOD) {
		len_method = 0
	}

	if !strings.Contains(f_stdout_fmt, FIELD_PROC) {
		len_proc = 0
	}

	if !strings.Contains(f_stdout_fmt, FIELD_PROTO) {
		len_proto = 0
	}

	if !strings.Contains(f_stdout_fmt, FIELD_QS) {
		len_qs = 0
	}

	if !strings.Contains(f_stdout_fmt, FIELD_REF) {
		len_ref = 0
	}

	if !strings.Contains(f_stdout_fmt, FIELD_SIZE) {
		len_size = 0
	}

	if !strings.Contains(f_stdout_fmt, FIELD_STATUS) {
		len_status = 0
	}

	if !strings.Contains(f_stdout_fmt, FIELD_STITLE) {
		len_stitle = 0
	}

	if !strings.Contains(f_stdout_fmt, FIELD_SDESC) {
		len_sdesc = 0
	}

	if !strings.Contains(f_stdout_fmt, FIELD_TIME) {
		len_time = 0
	}

	if !strings.Contains(f_stdout_fmt, FIELD_DATE) {
		len_date = 0
	}

	if !strings.Contains(f_stdout_fmt, FIELD_DATETIME) {
		len_datetime = 0
	}

	if !strings.Contains(f_stdout_fmt, FIELD_EPOCH) && !strings.Contains(f_stdout_fmt, FIELD_UNIX) {
		len_epoch = 0
	}

	if !strings.Contains(f_stdout_fmt, FIELD_URI) {
		len_uri = 0
	}

	if !strings.Contains(f_stdout_fmt, FIELD_UA) {
		len_ua = 0
	}

	if !strings.Contains(f_stdout_fmt, FIELD_UID) {
		len_uid = 0
	}

	if !strings.Contains(f_stdout_fmt, FIELD_FILE) {
		len_file = 0
	}

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
	getStrLen(FIELD_DATETIME, &len_datetime)
	getStrLen(FIELD_EPOCH, &len_epoch)
	getStrLen(FIELD_UNIX, &len_epoch)
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
		if found[0][i] == "0" {
			fmt.Println("0 (zero) is not allowed as field length.")
			os.Exit(1)
		}

		f_stdout_fmt = strings.Replace(
			f_stdout_fmt,
			fmt.Sprintf("{%s,%s}", item, found[0][i]),
			fmt.Sprintf("{%s}", item),
			-1)
	}
}

func PrintStdError(message string) {
	b := append([]byte(message), '\n')
	os.Stderr.Write(b)
}

func PrintAccessLogs(access *apachelogs.AccessLog) {
	if f_trim_slash {

		if len(access.URI) > 1 && access.URI[len(access.URI)-1] == '/' {
			access.URI = access.URI[:len(access.URI)-1]

		} else if access.URI == "/" {
			access.URI = "-"
		}
	}

	if f_write_sqlite != "" {
		go func() {
			if err := sqlite.InsertAccess(access); err != nil {
				PrintStdError(err.Error())
			}
		}()
	}

	b := []byte(f_stdout_fmt)
	if len_ip != 0 {
		formatSTDOUTb(&b, bFIELD_IP, access.IP, len_ip)
	}

	if len_method != 0 {
		formatSTDOUTb(&b, bFIELD_METHOD, access.Method, len_method)
	}

	if len_proc != 0 {
		formatSTDOUTb(&b, bFIELD_PROC, strconv.Itoa(access.ProcTime), len_proc)
	}

	if len_proto != 0 {
		formatSTDOUTb(&b, bFIELD_PROTO, access.Protocol, len_proto)
	}

	if len_qs != 0 {
		formatSTDOUTb(&b, bFIELD_QS, access.QueryString, len_qs)
	}

	if len_ref != 0 {
		formatSTDOUTb(&b, bFIELD_REF, access.Referrer, len_ref)
	}

	if len_size != 0 {
		formatSTDOUTb(&b, bFIELD_SIZE, strconv.Itoa(access.Size), len_size)
	}

	if len_status != 0 {
		formatSTDOUTb(&b, bFIELD_STATUS, access.Status.A, len_status)
	}

	if len_stitle != 0 {
		formatSTDOUTb(&b, bFIELD_STITLE, access.Status.Title(), len_stitle)
	}

	if len_sdesc != 0 {
		formatSTDOUTb(&b, bFIELD_SDESC, access.Status.Description(), len_sdesc)
	}

	if len_time != 0 {
		formatSTDOUTb(&b, bFIELD_TIME, access.DateTime.Format(FMT_TIME), len_time)
	}

	if len_date != 0 {
		formatSTDOUTb(&b, bFIELD_DATE, access.DateTime.Format(FMT_DATE), len_date)
	}

	if len_datetime != 0 {
		formatSTDOUTb(&b, bFIELD_DATETIME, access.DateTime.Format(FMT_DATETIME), len_datetime)
	}

	if len_epoch != 0 {
		s := strconv.FormatInt(access.DateTime.Unix(), 10)
		formatSTDOUTb(&b, bFIELD_EPOCH, s, len_epoch)
		formatSTDOUTb(&b, bFIELD_UNIX, s, len_epoch)
	}

	if len_uri != 0 {
		formatSTDOUTb(&b, bFIELD_URI, access.URI, len_uri)
	}

	if len_ua != 0 {
		formatSTDOUTb(&b, bFIELD_UA, access.UserAgent, len_ua)
	}

	if len_uid != 0 {
		formatSTDOUTb(&b, bFIELD_UID, access.UserID, len_uid)
	}

	if len_file != 0 {
		formatSTDOUTb(&b, bFIELD_FILE, access.FileName, len_file)
	}

	b = append(b, '\n')
	os.Stdout.Write(b)
}

func formatSTDOUTb(source *[]byte, search []byte, value string, length int) {
	if length > 0 {
		if len(value) > length {
			value = "…" + value[len(value)-length+1:]
		} else {
			value = spaces[:length-len(value)] + value
		}

	} else {
		l := -length
		if len(value) > l {
			value = value[:l-1] + "…"
		} else {
			value = value + spaces[:l-len(value)]
		}

	}
	search = append(search, brace_close)
	*source = bytes.Replace(*source, append(brace_open, search...), []byte(value), -1)
}
