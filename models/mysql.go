package models

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type MysqlConfig struct {
	Type string

	User        string
	Password    string
	Host        string
	Port        string
	Database    string
	TablePrefix string
}

func (cfg MysqlConfig) String() string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local", cfg.User, cfg.Password, cfg.Host, cfg.Database)
}

func DefaultMysqlConfig() MysqlConfig {
	return MysqlConfig{
		Type:     "mysql",
		User:     "baloo",
		Password: "junglebook",
		Host:     "localhost",
		Port:     "3306",
		Database: "lenslocked",
	}
}

// Open will open a SQL connnection with the provided
// Mysql database. cakkers of Open need to ensure
// db.Close() method
func Open(config MysqlConfig) (*sql.DB, error) {
	db, err := sql.Open("mysql", config.String())

	if err != nil {
		return nil, fmt.Errorf("Open: %w", err)
	}
	return db, nil
}
