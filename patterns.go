package main

import (
	"fmt"
	"github.com/lmorg/apachelogs"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func PatternDeconstructor(cli string) (p []apachelogs.Pattern) {
	rx_op, _ := regexp.Compile(`^([a-z]+)(>|<|=\+|!\+|=~|!~|!=|<>|==|=|~<|~>|\{\}|/|\*)(.*)`)
	for _, s := range strings.Split(cli, ";") {
		pat := rx_op.FindStringSubmatch(s)
		if len(pat) < 4 {
			continue
		}

		var (
			f  apachelogs.FieldID
			op apachelogs.OperatorID
		)

		switch pat[1] {
		case FIELD_IP:
			f = apachelogs.FIELD_IP
		case FIELD_METHOD:
			f = apachelogs.FIELD_METHOD
		case FIELD_PROC:
			f = apachelogs.FIELD_PROC_TIME
		case FIELD_PROTO:
			f = apachelogs.FIELD_PROTOCOL
		case FIELD_QS:
			f = apachelogs.FIELD_QUERY_STRING
		case FIELD_REF:
			f = apachelogs.FIELD_REFERRER
		case FIELD_SIZE:
			f = apachelogs.FIELD_SIZE
		case FIELD_STATUS:
			f = apachelogs.FIELD_STATUS
		case FIELD_TIME:
			f = apachelogs.FIELD_TIME
		case FIELD_DATE:
			f = apachelogs.FIELD_DATE
		case FIELD_DATETIME:
			f = apachelogs.FIELD_DATE_TIME
		case FIELD_URI:
			f = apachelogs.FIELD_URI
		case FIELD_UA:
			f = apachelogs.FIELD_USER_AGENT
		case FIELD_UID:
			f = apachelogs.FIELD_USER_ID
		case FIELD_EPOCH, FIELD_UNIX:
			f = apachelogs.FIELD_DATE_TIME
			t, err := strconv.ParseInt(pat[3], 10, 64)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			pat[3] = time.Unix(t, 0).Format("01-02-2006 15:04")
		default:
			fmt.Printf("Invalid field: %s\n", pat[1])
			os.Exit(1)

		}

		switch pat[2] {
		case "<":
			op = apachelogs.OP_LESS_THAN
		case "=+":
			op = apachelogs.OP_CONTAINS
		case "!+":
			op = apachelogs.OP_NOT_CONTAIN
		case "==", "=":
			op = apachelogs.OP_EQUAL_TO
		case "!=", "<>":
			op = apachelogs.OP_NOT_EQUAL
		case ">":
			op = apachelogs.OP_GREATER_THAN
		case "=~":
			op = apachelogs.OP_REGEX_EQ
		case "!~":
			op = apachelogs.OP_REGEX_NE
		case "~<":
			op = apachelogs.OP_ROUND_DOWN
		case "~>":
			op = apachelogs.OP_ROUND_UP
		case "{}":
			op = apachelogs.OP_REGEX_SUB
		case "/":
			op = apachelogs.OP_DIVIDE
		case "*":
			op = apachelogs.OP_MULTIPLY
		default:
			fmt.Printf("Invalid operator: %s\n", pat[2])
			os.Exit(1)
		}

		new_pat, err := apachelogs.NewPattern(f, op, pat[3])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		p = append(p, new_pat)

	}

	return
}
