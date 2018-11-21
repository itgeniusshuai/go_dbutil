package datasource

import (
	"../parse"
	"github.com/itgeniusshuai/go_common/common"
	"database/sql"
)


func (this *dataSource) InsertOne(sql string,params ...interface{})(int64,error){
	tx := GetFistUnNilTX()
	if tx != nil{
		return this.InsertOneWithTX(tx,sql,params...)
	}
	res,err := this.Exec(sql,params...)
	if err != nil{
		return -1,err
	}
	return res.LastInsertId()
}

func (this *dataSource) InsertOneWithTX(tx *sql.Tx,sql string,params ...interface{})(int64,error){
	res,err := tx.Exec(sql,params...)
	if err != nil{
		return -1,err
	}
	return res.LastInsertId()
}


func (this *dataSource) InsertOneObj(sql string,obj interface{})(int64,error){
	sqlColNames,err := parse.ParseInsertSql(sql)
	if err != nil{
		return -1,err
	}
	var params = make([]interface{},0)
	if err != nil{
		return -1,err
	}
	// 参数对象值
	for _,sqlColName := range sqlColNames{
		// 如果tag活着名称相同获取值否则赋值为空
		fv := common.GetValueByField(obj,sqlColName)
		if fv == nil{
			fv = common.GetFieldValueByFieldTag(obj,"col",sqlColName)
		}
		params = append(params, fv)
	}
	return this.InsertOne(sql,params...)
}

func (this *dataSource) InsertOneMap(sql string,paramMap map[string]interface{})(int64,error){
	colNames,err := parse.ParseInsertSql(sql)
	if err != nil{
		return -1,err
	}
	var params = make([]interface{},0)
	for _,colName := range colNames{
		params = append(params,paramMap[colName])
	}
	return this.InsertOne(sql,params...)
}

