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
			f  apachelogs.AccessFieldId
			op apachelogs.OperatorID
		)

		switch pat[1] {
		case accFieldIp:
			f = apachelogs.AccFieldIp
		case accFieldMethod:
			f = apachelogs.AccFieldMethod
		case accFieldProcTime:
			f = apachelogs.AccFieldProcTime
		case accFieldProtocol:
			f = apachelogs.AccFieldProtocol
		case accFieldQueryString:
			f = apachelogs.AccFieldQueryString
		case accFieldReferrer:
			f = apachelogs.AccFieldReferrer
		case accFieldSize:
			f = apachelogs.AccFieldSize
		case accFieldStatus:
			f = apachelogs.AccFieldStatus
		case accFieldTime:
			f = apachelogs.AccFieldTime
		case accFieldDate:
			f = apachelogs.AccFieldDate
		case accFieldDateTime:
			f = apachelogs.AccFieldDateTime
		case accFieldUri:
			f = apachelogs.AccFieldUri
		case accFieldUserAgent:
			f = apachelogs.AccFieldUserAgent
		case accFieldUserId:
			f = apachelogs.AccFieldUserId
		case accFieldEpoch, accFieldUnix:
			f = apachelogs.AccFieldDateTime
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
			op = apachelogs.OpLessThan
		case "=+":
			op = apachelogs.OpContains
		case "!+":
			op = apachelogs.OpDoesNotContain
		case "==", "=":
			op = apachelogs.OpEqualTo
		case "!=", "<>":
			op = apachelogs.OpNotEqual
		case ">":
			op = apachelogs.OpGreaterThan
		case "=~":
			op = apachelogs.OpRegexEqual
		case "!~":
			op = apachelogs.OpRegexNotEqual
		case "~<":
			op = apachelogs.OpRoundDown
		case "~>":
			op = apachelogs.OpRoundUp
		case "{}":
			op = apachelogs.OpRegexSubstitute
		case "/":
			op = apachelogs.OpDivide
		case "*":
			op = apachelogs.OpMultiply
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
