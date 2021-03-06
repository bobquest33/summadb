package database

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/fiatjaf/sublevel"
	"github.com/syndtr/goleveldb/leveldb/filter"
	"github.com/syndtr/goleveldb/leveldb/opt"

	settings "github.com/fiatjaf/summadb/settings"
)

const (
	DOC_STORE     = "d"
	REV_STORE     = "r"
	LOCAL_STORE   = "l"
	BY_SEQ        = "s"
	USER_STORE    = "u"
	PATH_METADATA = "m"

	UPDATE_SEQ_KEY = "_update_seq_key"
)

var db *sublevel.AbstractLevel

func Start() {
	dbfile := settings.DBFILE
	var err error

	// try to open an existing database
	db, err = sublevel.Open(dbfile, &opt.Options{
		Filter:         filter.NewBloomFilter(10),
		ErrorIfMissing: true,
	})
	if err != nil {
		// database is missing, create it and do initial setup
		db, err = sublevel.Open(dbfile, &opt.Options{
			Filter:       filter.NewBloomFilter(10),
			ErrorIfExist: true,
		})

		// admin party
		SetRulesAt("", map[string]interface{}{
			"_read":  "*",
			"_write": "*",
			"_admin": "*",
		})
	}

	if err != nil {
		log.WithFields(log.Fields{
			"error":  err,
			"DBFILE": settings.DBFILE,
		}).Fatal("couldn't open database file.")
	}
}

func End() {
	db.Close()
}

func Erase() error {
	dbfile := settings.DBFILE
	err := os.RemoveAll(dbfile)
	return err
}
