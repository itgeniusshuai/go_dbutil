package datasource

import "database/sql"

func (this *DataSource)Delete(sqlStr string,params ...interface{})(int64,error) {
	res, err := this.Exec(sqlStr, params...)
	if err !=  nil{
		return 0,err
	}
	return res.RowsAffected()
}

func (this *DataSource)DeleteWithTX(tx *sql.Tx,sqlStr string,params ...interface{})(int64,error) {
	res,err := tx.Exec(sqlStr,params...)
	if err !=  nil{
		return 0,err
	}
	return res.RowsAffected()
}


