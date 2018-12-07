package datasource

import (
	"github.com/itgeniusshuai/go_common/common"
	"database/sql"
	"github.com/itgeniusshuai/go_dbutil/parse"
)


func (this *DataSource) InsertOne(sql string,params ...interface{})(int64,error){
	res,err := this.Exec(sql,params...)
	if err != nil{
		return -1,err
	}
	return res.LastInsertId()
}

func (this *DataSource) InsertOneWithTX(tx *sql.Tx,sql string,params ...interface{})(int64,error){
	res,err := tx.Exec(sql,params...)
	if err != nil{
		return -1,err
	}
	return res.LastInsertId()
}


func (this *DataSource) InsertOneObj(sql string,obj interface{})(int64,error){
	paramMap := common.ObjToMapByTagName(obj,DB_STRUCT_TAG)
	return this.InsertOneMap(sql,paramMap)
}

func (this *DataSource) InsertOneMap(sql string,paramMap map[string]interface{})(int64,error){
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

func (this *DataSource) InsertOneObjWithTX(tx *sql.Tx,sql string,obj interface{})(int64,error){
	paramMap := common.ObjToMapByTagName(obj,DB_STRUCT_TAG)
	return this.InsertOneMapWithTX(tx,sql,paramMap)
}

func (this *DataSource) InsertOneMapWithTX(tx *sql.Tx,sql string,paramMap map[string]interface{})(int64,error){
	colNames,err := parse.ParseInsertSql(sql)
	if err != nil{
		return -1,err
	}
	var params = make([]interface{},0)
	for _,colName := range colNames{
		params = append(params,paramMap[colName])
	}
	return this.InsertOneWithTX(tx,sql,params...)
}


