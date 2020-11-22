package main

import (
	"github.com/rhinoman/couchdb-go"
	"time"
)

var conn *couchdb.Connection

func createDB() *couchdb.Database {
	var err error
	if conn == nil {
		var timeout = 10 * time.Second
		conn, err = couchdb.NewConnection("0.0.0.0", 5984, timeout)
	}
	if err != nil {
		panic(err)
	}
	auth := couchdb.BasicAuth{Username: "ginta", Password: "12345"}
	conn.CreateDB("mydb", &auth)
	return conn.SelectDB("mydb", &auth)
}
