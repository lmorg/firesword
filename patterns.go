package main

import (
	"github.com/lmorg/apachelogs"
	"fmt"
	"os"
	"regexp"
	"strings"
)

func PatternDeconstructor(cli string) (p []apachelogs.Pattern) {
	rx_op, _ := regexp.Compile(`^([a-z]+)(>|<|=\+|!\+|=~|!~|!=|<>|==|=)(.*)`)
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
		default:
			fmt.Printf("Invalid operator: %s\n", pat[2])
			os.Exit(1)
		}

		p = append(p, apachelogs.NewPattern(f, op, pat[3]))

	}

	return
}
