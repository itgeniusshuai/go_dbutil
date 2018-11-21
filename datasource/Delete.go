package datasource

import "database/sql"

func (this *dataSource)Delete(sqlStr string,params ...interface{})(int64,error) {
	var err error
	var res sql.Result
	tx := GetFistUnNilTX()
	if tx == nil{
		res,err = tx.Exec(sqlStr,params...)
	}else {
		res, err = this.Exec(sqlStr, params...)
	}
	if err !=  nil{
		return 0,err
	}
	return res.RowsAffected()
}


