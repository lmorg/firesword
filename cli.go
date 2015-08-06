package main

import (
	"bytes"
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

/*
 * Everything from here downwards is an experimental switch from strings and
 * fmt functions which use interface{}..., to []byte, append and os.Stdout.Write.
 *
 * I'm actually seeing 3 seconds shaved off my benchmarks - which I know is
 * meaningless to anyone reading this without sample sizes and machine specs.
 * But it gives you an idea as to why this ugly code exists.
 *
 * One day I might actually rewrite this entire project and it's sister package
 * (apachelogs) to run entirely on []byte. However apathy will probably prevail...
 * and it's not as if this utility is slow anyway.
 *
 * A more pressing requirement would be to making this code a little more readable
 * (eg removing the commented out code) and testing it thoroughly for bugs.
 * However for now I'm going to ship this beta code; you can uncomment the
 * following two functions then delete everything else that proceeds it should
 * you run into any annoying show-stopper bugs. Or pull the code from Github as
 * I may have fixed things by the time you read this :)
 */

/*
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
*/
/*
func formatSTDOUT(source *string, search, value string, length int) {
	slength := strconv.Itoa(length)
	if slength == "0" {
		slength = ""
	}

	value = Trim(fmt.Sprintf("%"+slength+"s", value), length)
	*source = strings.Replace(*source, "{"+search+"}", value, -1)
}
*/

var (
	bFIELD_IP       = []byte("ip")
	bFIELD_METHOD   = []byte("method")
	bFIELD_PROC     = []byte("proc")
	bFIELD_PROTO    = []byte("proto")
	bFIELD_QS       = []byte("qs")
	bFIELD_REF      = []byte("ref")
	bFIELD_SIZE     = []byte("size")
	bFIELD_STATUS   = []byte("status")
	bFIELD_STITLE   = []byte("stitle")
	bFIELD_SDESC    = []byte("sdesc")
	bFIELD_TIME     = []byte("time")
	bFIELD_DATE     = []byte("date")
	bFIELD_DATETIME = []byte("datetime")
	bFIELD_EPOCH    = []byte("epoch")
	bFIELD_URI      = []byte("uri")
	bFIELD_UA       = []byte("ua")
	bFIELD_UID      = []byte("uid")
	bFIELD_FILE     = []byte("file")
)

func PrintStdError(message string) {
	os.Stderr.Write([]byte(message))
}

func PrintAccessLogs(access *apachelogs.AccessLog) {
	b := []byte(f_stdout_fmt)
	formatSTDOUTb(&b, bFIELD_IP, access.IP, len_ip)
	formatSTDOUTb(&b, bFIELD_METHOD, access.Method, len_method)
	formatSTDOUTb(&b, bFIELD_PROC, strconv.Itoa(access.ProcTime), len_proc)
	formatSTDOUTb(&b, bFIELD_PROTO, access.Protocol, len_proto)
	formatSTDOUTb(&b, bFIELD_QS, access.QueryString, len_qs)
	formatSTDOUTb(&b, bFIELD_REF, access.Referrer, len_ref)
	formatSTDOUTb(&b, bFIELD_SIZE, strconv.Itoa(access.Size), len_size)
	formatSTDOUTb(&b, bFIELD_STATUS, access.Status.A, len_status)
	formatSTDOUTb(&b, bFIELD_STITLE, access.Status.Title(), len_stitle)
	formatSTDOUTb(&b, bFIELD_SDESC, access.Status.Description(), len_sdesc)
	formatSTDOUTb(&b, bFIELD_TIME, access.DateTime.Format(FMT_TIME), len_time)
	formatSTDOUTb(&b, bFIELD_DATE, access.DateTime.Format(FMT_DATE), len_date)
	formatSTDOUTb(&b, bFIELD_EPOCH, strconv.FormatInt(access.DateTime.Unix(), 10), len_epoch)
	formatSTDOUTb(&b, bFIELD_URI, access.URI, len_uri)
	formatSTDOUTb(&b, bFIELD_UA, access.UserAgent, len_ua)
	formatSTDOUTb(&b, bFIELD_UID, access.UserID, len_uid)
	formatSTDOUTb(&b, bFIELD_FILE, access.FileName, len_file)
	b = append(b, '\n')
	os.Stdout.Write(b)
}

func formatSTDOUTb(source *[]byte, search []byte, value string, length int) {
	//b = []byte(value)
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
	search = append(search, bc)
	*source = bytes.Replace(*source, append(bo, search...), []byte(value), -1)
}

var (
	bo []byte = []byte{'{'}
	bc byte   = '}'
)

var spaces = "                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                    "
