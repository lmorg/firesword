package sqlite

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
