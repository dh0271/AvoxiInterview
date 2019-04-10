package persistence

import (
	"log"

	"github.com/oschwald/maxminddb-golang"
)

var db *maxminddb.Reader

func Init(path string) {
	local, err := maxminddb.Open(path)

	if err != nil {
		log.Panicln("Failed to load database")
		panic(err)
	}

	db = local
}

func Close() {
	if db != nil {
		db.Close()
	}
}

func GetDB() *maxminddb.Reader {
	// This shouldn't happen but just in case
	if db == nil {
		log.Panic("Database is not initialized")
	}

	return db
}
