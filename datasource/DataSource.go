package datasource

import "database/sql"

type dataSource struct{
	Name string
	User string
	Pwd string
	Url string
	db *sql.DB
}

