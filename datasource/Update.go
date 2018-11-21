package datasource

import "database/sql"

func (this *dataSource)Update(sql string,params ...interface{})(int64,error) {
	var err error
	var res sql.Result
	tx := GetFistUnNilTX()
	if tx == nil{
		res,err = tx.Exec(sql,params...)
	}else {
		res, err = this.Exec(sql, params...)
	}
	if err !=  nil{
		return 0,err
	}
	return res.RowsAffected()
}