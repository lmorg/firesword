package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/lmorg/apachelogs"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func New() {
	// empty string == in memory
	//if filename == "" {
	//filename := ":memory:"
	//filename := "logs.db"
	//}

	var err error
	//db, err = sql.Open("sqlite3", fmt.Sprintf("file:%s?cache=shared&mode=rwc", filename))
	db, err = sql.Open("sqlite3", "file:memdb1?mode=memory&cache=shared")
	//init_db, err = sql.Open("sqlite3", "file:test.db?cache=shared")
	if err != nil {
		log.Fatalln("could not open database:", err)
	}

	_, err = db.Exec(_SQL_CREATE_TABLE)
	if err != nil {
		log.Fatalln("could not create table:", err)
	}

	/*
		// views
		view := func(sql string) {
			_, err = db.Exec(sql)
			if err != nil {
				log.Fatalln("could not create view:", err)
			}
		}

		view(_VIEW_ALL)
		view(_VIEW_LATEST_NON_200)
		view(_VIEW_LATEST_PROC)
		view(_VIEW_LATEST_304)
		view(_VIEW_COUNT_STATUS)
		view(_VIEW_COUNT_304)
		view(_VIEW_COUNT_SIZE)
		view(_VIEW_LIST_VIEWS)
	*/
}

func InsertAccess(access *apachelogs.AccessLog) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("Panic caught: %s", r))
		}
	}()

	_, err = db.Exec(_SQL_INSERT_ACCESS,
		access.IP,
		access.Method,
		access.ProcTime,
		access.Protocol,
		access.QueryString,
		access.Referrer,
		access.Size,
		access.Status.I,
		access.Status.Title(),
		access.Status.Description(),
		access.DateTime,
		access.URI,
		access.UserAgent,
		access.UserID,
		access.FileName,
	)
	return
}

/*
func Query(sql string) (rows *sql.Rows, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("Panic caught: %s", r))
		}
	}()

	rows, err = db.Query(sql)

	return
}*/

func Dump(filename string) {
	//go-sqlite3.
}

func Close() {
	db.Close()
}
