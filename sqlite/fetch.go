package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

var (
	Field []interface{}
)

// this is a bit shit, but I'm not sure there's a fast alternative without Go supporting generics :(
func RowsFetch(rows *sql.Rows, i int) (err error) {
	Field = make([]interface{}, i)
	switch i {
	case 1:
		err = rows.Scan(&Field[0])
	case 2:
		err = rows.Scan(&Field[0], &Field[1])
	case 3:
		err = rows.Scan(&Field[0], &Field[1], &Field[2])
	case 4:
		err = rows.Scan(&Field[0], &Field[1], &Field[2], &Field[3])
	case 5:
		err = rows.Scan(&Field[0], &Field[1], &Field[2], &Field[3], &Field[4])
	case 6:
		err = rows.Scan(&Field[0], &Field[1], &Field[2], &Field[3], &Field[4], &Field[5])
	case 7:
		err = rows.Scan(&Field[0], &Field[1], &Field[2], &Field[3], &Field[4], &Field[5], &Field[6])
	case 8:
		err = rows.Scan(&Field[0], &Field[1], &Field[2], &Field[3], &Field[4], &Field[5], &Field[6], &Field[7])
	case 9:
		err = rows.Scan(&Field[0], &Field[1], &Field[2], &Field[3], &Field[4], &Field[5], &Field[6], &Field[7], &Field[8])
	case 10:
		err = rows.Scan(&Field[0], &Field[1], &Field[2], &Field[3], &Field[4], &Field[5], &Field[6], &Field[7], &Field[8], &Field[9])
	case 11:
		err = rows.Scan(&Field[0], &Field[1], &Field[2], &Field[3], &Field[4], &Field[5], &Field[6], &Field[7], &Field[8], &Field[9], &Field[10])
	case 12:
		err = rows.Scan(&Field[0], &Field[1], &Field[2], &Field[3], &Field[4], &Field[5], &Field[6], &Field[7], &Field[8], &Field[9], &Field[10], &Field[11])
	case 13:
		err = rows.Scan(&Field[0], &Field[1], &Field[2], &Field[3], &Field[4], &Field[5], &Field[6], &Field[7], &Field[8], &Field[9], &Field[10], &Field[11], &Field[12])
	case 14:
		err = rows.Scan(&Field[0], &Field[1], &Field[2], &Field[3], &Field[4], &Field[5], &Field[6], &Field[7], &Field[8], &Field[9], &Field[10], &Field[11], &Field[12], &Field[13])
	case 15:
		err = rows.Scan(&Field[0], &Field[1], &Field[2], &Field[3], &Field[4], &Field[5], &Field[6], &Field[7], &Field[8], &Field[9], &Field[10], &Field[11], &Field[12], &Field[13], &Field[14])
	case 16:
		err = rows.Scan(&Field[0], &Field[1], &Field[2], &Field[3], &Field[4], &Field[5], &Field[6], &Field[7], &Field[8], &Field[9], &Field[10], &Field[11], &Field[12], &Field[13], &Field[14], &Field[15])
	case 17:
		err = rows.Scan(&Field[0], &Field[1], &Field[2], &Field[3], &Field[4], &Field[5], &Field[6], &Field[7], &Field[8], &Field[9], &Field[10], &Field[11], &Field[12], &Field[13], &Field[14], &Field[15], &Field[16])
	case 18:
		err = rows.Scan(&Field[0], &Field[1], &Field[2], &Field[3], &Field[4], &Field[5], &Field[6], &Field[7], &Field[8], &Field[9], &Field[10], &Field[11], &Field[12], &Field[13], &Field[14], &Field[15], &Field[16], &Field[17])
	case 19:
		err = rows.Scan(&Field[0], &Field[1], &Field[2], &Field[3], &Field[4], &Field[5], &Field[6], &Field[7], &Field[8], &Field[9], &Field[10], &Field[11], &Field[12], &Field[13], &Field[14], &Field[15], &Field[16], &Field[17], &Field[18])
	case 20:
		err = rows.Scan(&Field[0], &Field[1], &Field[2], &Field[3], &Field[4], &Field[5], &Field[6], &Field[7], &Field[8], &Field[9], &Field[10], &Field[11], &Field[12], &Field[13], &Field[14], &Field[15], &Field[16], &Field[17], &Field[18], &Field[19])
	default:
		err = errors.New("Must return 1 to 20 fields")
	}

	return
}

func GetField(column int) string {
	switch t := Field[column].(type) {
	default:
		return fmt.Sprintf("{%T}", t)
	case int, int64:
		return fmt.Sprintf("%d", t)
	case string, []uint8:
		return fmt.Sprintf("%s", t)
	case time.Time:
		return t.Format("02 Jan 2006 15:04:05")
	}
	//TODO: probably should do proper conversion here instead of lazy Sprintf
}
