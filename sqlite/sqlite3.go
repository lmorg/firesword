package sqlite

import (
	"github.com/lmorg/apachelogs"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

var (
	db *sql.DB
)

const (
	_SQL_CREATE_TABLE = `CREATE TABLE access (
							id 			integer PRIMARY KEY,
							ip          string,
							method      string,
							proc    	integer,
							proto		string,
							qs      	string,
							ref         string,
							size 		integer,
							status    	integer,
							stitle		string,
							sdesc		string,
							datetime    datetime,
							uri         string,
							ua    		string,
							uid   		string,
							file    	string
						);`

	_SQL_INSERT_ACCESS = `INSERT INTO access (
							ip,
							method,
							proc,
							proto,
							qs,
							ref,
							size,
							status,
							stitle,
							sdesc,
							datetime,
							uri,
							ua,
							uid,
							file
						) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`
)

func New( /*filename string*/ ) {
	// empty string == in memory
	//if filename == "" {
	filename := ":memory:"
	//filename := "logs.db"
	//}

	var err error
	db, err = sql.Open("sqlite3", fmt.Sprintf("file:%s?cache=shared&mode=rwc", filename))
	if err != nil {
		log.Fatalln("could not open database:", err)
	}

	_, err = db.Exec(_SQL_CREATE_TABLE)
	if err != nil {
		log.Fatalln("could not create table:", err)
	}

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

	//go insertEventHandler()
}

func InsertAccess(access apachelogs.AccessLog) (err error) {
	defer func() {
		if r := recover(); r != nil {
			//fmt.Println("Panic caught:", r)
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

func Query(sql string) (rows *sql.Rows, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("Panic caught: %s", r))
		}
	}()

	rows, err = db.Query(sql)

	return
}

func Close() {
	db.Close()
}
