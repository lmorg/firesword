package main

import (
	"bytes"
	"fmt"
	"github.com/lmorg/apachelogs"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// Command line field names.
// The code duplication is nasty, but realistically it needs to be
// both string and byte slice for performance reasons. Having this
// duplication compiled in saves conversion overhead in tight loops.
var (
	accFieldIp          = "ip"
	accFieldMethod      = "method"
	accFieldProcTime    = "proc"
	accFieldProtocol    = "proto"
	accFieldQueryString = "qs"
	accFieldReferrer    = "ref"
	accFieldSize        = "size"
	accFieldStatus      = "status"
	accFieldStatusTitle = "stitle"
	accFieldStatusDesc  = "sdesc"
	accFieldTime        = "time"
	accFieldDate        = "date"
	accFieldDateTime    = "datetime"
	accFieldEpoch       = "epoch"
	accFieldUnix        = "unix"
	accFieldUri         = "uri"
	accFieldUserAgent   = "ua"
	accFieldUserId      = "uid"
	accFieldFileName    = "file"

	// As odd as seems duplicating field names from the above,
	// this method is quicker at printing to stdout as we can
	// alternate between strings and character arrays depending
	// on the core libs being used.
	bAccFieldIp          = []byte(accFieldIp)
	bAccFieldMethod      = []byte(accFieldMethod)
	bAccFieldProcTime    = []byte(accFieldProcTime)
	bAccFieldProtocol    = []byte(accFieldProtocol)
	bAccFieldQueryString = []byte(accFieldQueryString)
	bAccFieldReferrer    = []byte(accFieldReferrer)
	bAccFieldSize        = []byte(accFieldSize)
	bAccFieldStatus      = []byte(accFieldStatus)
	bAccFieldStatusTitle = []byte(accFieldStatusTitle)
	bAccFieldStatusDesc  = []byte(accFieldStatusDesc)
	bAccFieldTime        = []byte(accFieldTime)
	bAccFieldDate        = []byte(accFieldDate)
	bAccFieldDateTime    = []byte(accFieldDateTime)
	bAccFieldEpoch       = []byte(accFieldEpoch)
	bAccFieldUnix        = []byte(accFieldUnix)
	bAccFieldUri         = []byte(accFieldUri)
	bAccFieldUserAgent   = []byte(accFieldUserAgent)
	bAccFieldUserId      = []byte(accFieldUserId)
	bAccFieldFileName    = []byte(accFieldFileName)
)

// Field lengths
var (
	lenIp          int = -15
	lenMethod      int = -4
	lenProcTime    int = 9
	lenProtocol    int = -8
	lenQueryString int = -20
	lenReferrer    int = 30
	lenSize        int = 9
	lenStatus      int = -3
	lenStatusTitle int = -20
	lenStatusDesc  int = -40
	lenTime        int = -8
	lenDate        int = -11
	lenDateTime    int = -20
	lenEpoch       int = 10
	lenUri         int = 40
	lenUserAgent   int = -20
	lenUserId      int = -10
	lenFileName    int = -20
)

// Weird variable, but this allows printf-style padding using fast
// slicing rather than slow fmt.Printf([]interface{}...)'s
var spaces = "                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                    "

// Used for concatenating byte slices with append()
// (the braces are part of the field matching in the command line string)
var (
	braceOpen  []byte = []byte{'{'}
	braceClose byte   = '}'
)

// Check if fields are being used in --fmt. Do this up front and
// then set the field length to zero to disable if field unused.
// This saves searching through the --fmt string for each value
// against every log line outputted (slow).
func cropUnusedFmt() {
	if !strings.Contains(fStdOutFormat, accFieldIp) {
		lenIp = 0
	}

	if !strings.Contains(fStdOutFormat, accFieldMethod) {
		lenMethod = 0
	}

	if !strings.Contains(fStdOutFormat, accFieldProcTime) {
		lenProcTime = 0
	}

	if !strings.Contains(fStdOutFormat, accFieldProtocol) {
		lenProtocol = 0
	}

	if !strings.Contains(fStdOutFormat, accFieldQueryString) {
		lenQueryString = 0
	}

	if !strings.Contains(fStdOutFormat, accFieldReferrer) {
		lenReferrer = 0
	}

	if !strings.Contains(fStdOutFormat, accFieldSize) {
		lenSize = 0
	}

	if !strings.Contains(fStdOutFormat, accFieldStatus) {
		lenStatus = 0
	}

	if !strings.Contains(fStdOutFormat, accFieldStatusTitle) {
		lenStatusTitle = 0
	}

	if !strings.Contains(fStdOutFormat, accFieldStatusDesc) {
		lenStatusDesc = 0
	}

	if !strings.Contains(fStdOutFormat, accFieldTime) {
		lenTime = 0
	}

	if !strings.Contains(fStdOutFormat, accFieldDate) {
		lenDate = 0
	}

	if !strings.Contains(fStdOutFormat, accFieldDateTime) {
		lenDateTime = 0
	}

	if !strings.Contains(fStdOutFormat, accFieldEpoch) && !strings.Contains(fStdOutFormat, accFieldUnix) {
		lenEpoch = 0
	}

	if !strings.Contains(fStdOutFormat, accFieldUri) {
		lenUri = 0
	}

	if !strings.Contains(fStdOutFormat, accFieldUserAgent) {
		lenUserAgent = 0
	}

	if !strings.Contains(fStdOutFormat, accFieldUserId) {
		lenUserId = 0
	}

	if !strings.Contains(fStdOutFormat, accFieldFileName) {
		lenFileName = 0
	}

}

// Cycles through the stdout format string looking for string lengths.
// Runs on start of application.
func ImportStrLen() {
	getStrLen(accFieldIp, &lenIp)
	getStrLen(accFieldMethod, &lenMethod)
	getStrLen(accFieldProcTime, &lenProcTime)
	getStrLen(accFieldProtocol, &lenProtocol)
	getStrLen(accFieldQueryString, &lenQueryString)
	getStrLen(accFieldReferrer, &lenReferrer)
	getStrLen(accFieldSize, &lenSize)
	getStrLen(accFieldStatus, &lenStatus)
	getStrLen(accFieldStatusTitle, &lenStatusTitle)
	getStrLen(accFieldStatusDesc, &lenStatusDesc)
	getStrLen(accFieldTime, &lenTime)
	getStrLen(accFieldDate, &lenDate)
	getStrLen(accFieldDateTime, &lenDateTime)
	getStrLen(accFieldEpoch, &lenEpoch)
	getStrLen(accFieldUnix, &lenEpoch)
	getStrLen(accFieldUri, &lenUri)
	getStrLen(accFieldUserAgent, &lenUserAgent)
	getStrLen(accFieldUserId, &lenUserId)
	getStrLen(accFieldFileName, &lenFileName)
}

// Gets the string length for each item and then reformat stdout format string
// for tighter output loops
func getStrLen(item string, val *int) {
	rx, _ := regexp.Compile(`{` + item + `,(\-?[0-9]+)}`)
	found := rx.FindAllStringSubmatch(fStdOutFormat, -1)

	if len(found) == 0 || len(found[0]) == 0 {
		return
	}

	*val, _ = strconv.Atoi(found[0][1])

	for i := range found[0] {
		if found[0][i] == "0" {
			fmt.Println("0 (zero) is not allowed as field length.")
			os.Exit(1)
		}

		fStdOutFormat = strings.Replace(
			fStdOutFormat,
			fmt.Sprintf("{%s,%s}", item, found[0][i]),
			fmt.Sprintf("{%s}", item),
			-1)
	}
}

func PrintStdError(message string) {
	b := append([]byte(message), '\n')
	os.Stderr.Write(b)
}

func PrintAccessLogs(access *apachelogs.AccessLine) {
	if fTrimSlash {

		if len(access.URI) > 1 && access.URI[len(access.URI)-1] == '/' {
			access.URI = access.URI[:len(access.URI)-1]

		} else if access.URI == "/" {
			access.URI = "-"
		}
	}

	b := []byte(fStdOutFormat)
	if lenIp != 0 {
		formatStdOut(&b, bAccFieldIp, access.IP, lenIp)
	}

	if lenMethod != 0 {
		formatStdOut(&b, bAccFieldMethod, access.Method, lenMethod)
	}

	if lenProcTime != 0 {
		formatStdOut(&b, bAccFieldProcTime, strconv.Itoa(access.ProcTime), lenProcTime)
	}

	if lenProtocol != 0 {
		formatStdOut(&b, bAccFieldProtocol, access.Protocol, lenProtocol)
	}

	if lenQueryString != 0 {
		formatStdOut(&b, bAccFieldQueryString, access.QueryString, lenQueryString)
	}

	if lenReferrer != 0 {
		formatStdOut(&b, bAccFieldReferrer, access.Referrer, lenReferrer)
	}

	if lenSize != 0 {
		formatStdOut(&b, bAccFieldSize, strconv.Itoa(access.Size), lenSize)
	}

	if lenStatus != 0 {
		formatStdOut(&b, bAccFieldStatus, access.Status.A, lenStatus)
	}

	if lenStatusTitle != 0 {
		formatStdOut(&b, bAccFieldStatusTitle, access.Status.Title(), lenStatusTitle)
	}

	if lenStatusDesc != 0 {
		formatStdOut(&b, bAccFieldStatusDesc, access.Status.Description(), lenStatusDesc)
	}

	if lenTime != 0 {
		formatStdOut(&b, bAccFieldTime, access.DateTime.Format(TimeFormat), lenTime)
	}

	if lenDate != 0 {
		formatStdOut(&b, bAccFieldDate, access.DateTime.Format(DateFormat), lenDate)
	}

	if lenDateTime != 0 {
		formatStdOut(&b, bAccFieldDateTime, access.DateTime.Format(DateTimeFormat), lenDateTime)
	}

	if lenEpoch != 0 {
		s := strconv.FormatInt(access.DateTime.Unix(), 10)
		formatStdOut(&b, bAccFieldEpoch, s, lenEpoch)
		formatStdOut(&b, bAccFieldUnix, s, lenEpoch)
	}

	if lenUri != 0 {
		formatStdOut(&b, bAccFieldUri, access.URI, lenUri)
	}

	if lenUserAgent != 0 {
		formatStdOut(&b, bAccFieldUserAgent, access.UserAgent, lenUserAgent)
	}

	if lenUserId != 0 {
		formatStdOut(&b, bAccFieldUserId, access.UserID, lenUserId)
	}

	if lenFileName != 0 {
		formatStdOut(&b, bAccFieldFileName, access.FileName, lenFileName)
	}

	b = append(b, '\n')
	os.Stdout.Write(b)
}

func formatStdOut(source *[]byte, search []byte, value string, length int) {
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
	search = append(search, braceClose)
	*source = bytes.Replace(*source, append(braceOpen, search...), []byte(value), -1)
}
