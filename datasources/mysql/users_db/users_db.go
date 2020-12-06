package users_db

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
)

const (
	mysqlUsername = "mysql_username"
	mysqlPassword = "mysql_password"
	mysqlHost     = "mysql_host"
	mysqlSchema   = "mysql_schema"
)
var (
	/* Capitalized variable names are exported for access in other packages */
	Db *sql.DB // database connection pool.

	username = os.Getenv(mysqlUsername)
	password = os.Getenv(mysqlPassword)
	host = os.Getenv(mysqlHost)
	schema = os.Getenv(mysqlSchema)
)

// By importing users_id package anywhere, init function executes
func init()  {
	datasourceName := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8",
		username,
		password,
		host,
		schema,
	)
	var err error
	Db, err = sql.Open("mysql", datasourceName)
	if err != nil {
		panic(err) // not going to start the application
	}
	if err := Db.Ping(); err != nil {
		panic(err)
	}
	log.Println("database successfully configured")
}
