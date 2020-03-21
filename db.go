package main

import (
	"strconv"
	"strings"

	_ "github.com/mattn/go-sqlite3"

	"database/sql"
)

type dbHelper struct {
	DBpath string
	db     *sql.DB
}

type ResultKeys struct {
	Keys    []ResultKey
	IsError bool
	Error   string
}

type ResultKey struct {
	KeyID     int
	Name      string
	Key       string
	Signature string
}

func (d *dbHelper) NewHelper(DBpath string) {
	d.DBpath = DBpath
	db, err := sql.Open("sqlite3", d.DBpath)
	if err != nil {
		panic(err)
	}

	d.db = db
}

func (d *dbHelper) CreateDatabase(DBpath string) error {
	db, _ := sql.Open("sqlite3", DBpath)

	stmt, err := db.Prepare("CREATE TABLE Keys(ID INTEGER PRIMARY KEY,Name TEXT,Key TEXT UNIQUE,Signature TEXT)")
	if err != nil {
		panic(err)
	}
	_, err = stmt.Exec()
	if err != nil {
		panic(err)
	}

	stmt, err = db.Prepare("CREATE TABLE Groups(ID INTEGER PRIMARY KEY,Name TEXT UNIQUE)")
	if err != nil {
		panic(err)
	}

	_, err = stmt.Exec()
	if err != nil {
		panic(err)
	}

	stmt, err = db.Prepare("CREATE TABLE GroupKeys(GroupName TEXT,KeyID INTERGER,FOREIGN KEY (GroupName) REFERENCES Groups (Name), FOREIGN KEY (KeyID)	REFERENCES Keys (ID) )")
	if err != nil {
		panic(err)
	}

	_, err = stmt.Exec()

	if err != nil {
		panic(err)
	}

	db.Close()
	return nil
}

func (d *dbHelper) DoesGroupExist(groupName string) (bool, error) {
	var result string

	err := d.db.QueryRow("SELECT Name FROM Groups WHERE Name=?", groupName).Scan(&result)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (d *dbHelper) GetGroupKeys(groupName string) ([]ResultKey, error) {
	var KeyIDs []int

	rows, err := d.db.Query("SELECT KeyID FROM GroupKeys WHERE GroupName=?", groupName)
	if err != nil {
		return []ResultKey{}, err
	}

	for rows.Next() {
		var id int
		rows.Scan(&id)

		KeyIDs = append(KeyIDs, id)
	}
	rows.Close()

	var stmtStr string
	stmtStr = "SELECT * FROM Keys WHERE ID IN ("
	for x := range KeyIDs {
		stmtStr += strconv.Itoa(KeyIDs[x]) + ","
	}

	stmtStr = strings.TrimSuffix(stmtStr, ",")
	stmtStr += ")"

	rows, err = d.db.Query(stmtStr)
	if err != nil {
		return []ResultKey{}, err
	}

	var r []ResultKey

	for rows.Next() {
		var x ResultKey
		rows.Scan(&x.KeyID, &x.Name, &x.Key, &x.Signature)

		r = append(r, x)
	}

	return r, nil
}
