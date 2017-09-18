package db

import (
	"database/sql"
	"fmt"

	"github.com/SleepSpotify/SleepSpotify/config"
	_ "github.com/go-sql-driver/mysql" // the mysql driver
)

// DB the db driver to control the database
var DB *sql.DB

// InitDB function to init the db to be call in the main
func InitDB(config config.Config) error {
	var err error
	DB, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", config.DB.Username, config.DB.Password, config.DB.Host, config.DB.Port, config.DB.Name))
	return err
}
