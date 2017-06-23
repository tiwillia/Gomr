package gomr

import (
	"database/sql"
	"fmt"
	"github.com/go-gorp/gorp"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang/glog"
)

// Create and return a Database object
func InitDB(host, port, user, password, dbname string) (db *gorp.DbMap, err error) {
	connectionString := fmt.Sprintf("%s:%s@%s([%s]:%s)/%s",
		user, password, "tcp", host, port, dbname)

	// Create a connection with the database
	var dbCon *sql.DB
	dbCon, err = sql.Open("mysql", connectionString)
	if err != nil {
		return
	}

	// Verify that the database connection works
	err = dbCon.Ping()
	if err != nil {
		return
	}

	// Set up gorp mappings
	db = &gorp.DbMap{Db: dbCon, Dialect: gorp.MySQLDialect{"InnoDB", "UTF8"}}

	defineTables(db)
	if err := db.CreateTablesIfNotExists(); err != nil {
		glog.Fatalln("Unable to create tables:", err)
	}

	return db, err
}

func defineTables(Dbm *gorp.DbMap) {
	// Column sizes and options are defined on the database table structs,
	//   there is no reason to set it here.
	_ = Dbm.AddTableWithName(Karma{}, "karma").SetKeys(true, "Id")
	_ = Dbm.AddTableWithName(Factoid{}, "factoids").SetKeys(true, "Id")
}
