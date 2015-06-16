package sqlite

import (
	"database/sql"
	"fmt"
	"github.com/lmorg/apachelogs"
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
	_, err = db.Exec(_VIEW_ALL)
	if err != nil {
		log.Fatalln("could not create view:", err)
	}

	_, err = db.Exec(_VIEW_LATEST_NON_200)
	if err != nil {
		log.Fatalln("could not create view:", err)
	}

	_, err = db.Exec(_VIEW_LATEST_PROC)
	if err != nil {
		log.Fatalln("could not create view:", err)
	}

	_, err = db.Exec(_VIEW_LATEST_304)
	if err != nil {
		log.Fatalln("could not create view:", err)
	}

	_, err = db.Exec(_VIEW_COUNT_STATUS)
	if err != nil {
		log.Fatalln("could not create view:", err)
	}

	_, err = db.Exec(_VIEW_COUNT_304)
	if err != nil {
		log.Fatalln("could not create view:", err)
	}

	_, err = db.Exec(_VIEW_COUNT_SIZE)
	if err != nil {
		log.Fatalln("could not create view:", err)
	}

	_, err = db.Exec(_VIEW_LIST_VIEWS)
	if err != nil {
		log.Fatalln("could not create view:", err)
	}

}

func InsertAccess(access apachelogs.AccessLog, filename string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Pacnic caught:", r)
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
		filename,
	)
	return
}

func Query(sql string) (rows *sql.Rows, err error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Pacnic caught:", r)
		}
	}()

	rows, err = db.Query(sql)
	/*if err != nil {
		log.Fatalln("could not select data:", err)
	}*/
	return
}

func Close() {
	db.Close()
}
