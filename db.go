package main

import (
	"time"

	"github.com/rhinoman/couchdb-go"
)

var conn *couchdb.Connection

func createDB() *couchdb.Database {
	var err error
	if conn == nil {
		var timeout = time.Duration(500 * time.Millisecond)
		conn, err = couchdb.NewConnection("127.0.0.1", 5984, timeout)
	}
	if err != nil {
		panic(err)
	}
	auth := couchdb.BasicAuth{Username: "pinta", Password: "developer1"}
	conn.CreateDB("mydb", &auth)
	return conn.SelectDB("mydb", &auth)
}
