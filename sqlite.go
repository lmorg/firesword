// Remove the following commented line to compile with ncurses support:
// +build ignore
//
// Ncurses mode also requires sqlite and readline - so compiling with ncurses
// breaks cross-compiling portability for the sake of extra features.

package main

import (
	"fmt"
	"os"
	"os/user"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/lmorg/apachelogs"
	"github.com/lmorg/plasmasword/sqlite"
)

func __nQuit() {
	sqlite.Close()
}

func __nInterface() {
	sqlite.New()
}

func __nInsert(access *apachelogs.AccessLog) {
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
