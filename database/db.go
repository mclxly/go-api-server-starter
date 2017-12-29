package database

import (
	// "log"
	"database/sql"
	
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"  

	"github.com/spf13/viper"            // Application Config
)

func NewConnect() (*GlobalDB, error) {
	connStr := viper.GetString("db_connect_string")
	// log.Print(connStr)

	if dbConn, err := sqlx.Connect("postgres", connStr); err != nil {
		return nil, err
	} else {
		p := &GlobalDB{dbConn: dbConn}
		if err := p.dbConn.Ping(); err != nil {
			return nil, err
		}
		// if err := p.createTablesIfNotExist(); err != nil {
		// 	return nil, err
		// }
		// if err := p.prepareSqlStatements(); err != nil {
		// 	return nil, err
		// }

		return p, nil
	}
}

type GlobalDB struct {
	dbConn *sqlx.DB

	sqlSelectPeople *sqlx.Stmt
	sqlInsertPerson *sqlx.NamedStmt
	sqlSelectPerson *sql.Stmt
}
