package datasource

import "database/sql"

type dataSource struct{
	*sql.DB
	Name string
	User string
	Pwd string
	Url string
}

type ModelPtr interface{}

