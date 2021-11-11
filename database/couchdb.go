package database

import (
	_ "github.com/go-kivik/couchdb/v4"
	kivik "github.com/go-kivik/kivik/v4"
)

var client *kivik.Client

func init() {
	var err error
	client, err = kivik.New("couch", "http://localhost:5984/")
	if err != nil {
		panic(err)
	}
}

func Db(name string) *kivik.DB {
	return client.DB(name)
}
