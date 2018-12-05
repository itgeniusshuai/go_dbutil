package datasource

import "database/sql"

type DataSource struct{
	*sql.DB
	Name string
	User string
	Pwd string
	Url string
}

type ModelPtr interface{}

type PROPAGATION int
const (
	PROPAGATION_REQUIRED PROPAGATION = iota
	PROPAGATION_NEW
	PROPAGATION_NESTED
	)

const (
	DB_STRUCT_TAG string = "db"
)

