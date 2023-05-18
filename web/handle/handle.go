package handle

import (
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"log"
	"nobozo/etc"
	"os"
)

type handle_t struct {
	dbh *sql.DB
	dsn string
}

var Handle handle_t

func Connect() {
	var err error
	var cfg mysql.Config
	var dsn string

	fmt.Println("in Connect")
	cfg = mysql.Config{
		User:   etc.DatabaseUser,
		Passwd: etc.DatabasePassword,
		Net:    "tcp",
		Addr:   etc.DatabaseHost + ":" + etc.DatabasePort,
		DBName: etc.DatabaseName,
	}

	if etc.DatabaseType == "Oracle" {
		os.Setenv("NLS_LANG", "AMERICAN_AMERICA.AL32UTF8")
		os.Setenv("NLS_NCHAR", "AL32UTF8")
	}

	dsn = cfg.FormatDSN()
	Handle.dbh, err = sql.Open(etc.DatabaseType, dsn)
	Handle.dsn = dsn

	if err != nil {
		log.Printf("Db err = %v\n", err)
		fmt.Fprintf(os.Stderr, "Db err = %v\n", err)
	}
}
