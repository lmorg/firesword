package main

import (
	"fmt"
	"github.com/lmorg/apachelogs"
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
		case CLI_STR_IP:
			f = apachelogs.FIELD_IP
		case CLI_STR_METHOD:
			f = apachelogs.FIELD_METHOD
		case CLI_STR_PROC:
			f = apachelogs.FIELD_PROC_TIME
		case CLI_STR_PROTO:
			f = apachelogs.FIELD_PROTOCOL
		case CLI_STR_QS:
			f = apachelogs.FIELD_QUERY_STRING
		case CLI_STR_REF:
			f = apachelogs.FIELD_REFERRER
		case CLI_STR_SIZE:
			f = apachelogs.FIELD_SIZE
		case CLI_STR_STATUS:
			f = apachelogs.FIELD_STATUS
		case CLI_STR_TIME:
			f = apachelogs.FIELD_TIME
		case CLI_STR_DATE:
			f = apachelogs.FIELD_DATE
		case CLI_STR_DATETIME:
			f = apachelogs.FIELD_DATE_TIME
		case CLI_STR_URI:
			f = apachelogs.FIELD_URI
		case CLI_STR_UA:
			f = apachelogs.FIELD_USER_AGENT
		case CLI_STR_UID:
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
