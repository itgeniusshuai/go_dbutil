package datasource

import (
	"../parse"
	"github.com/itgeniusshuai/go_common/common"
)


func (this *dataSource) InsertOne(sql string,params ...interface{})(int64,error){
	res,err := this.Exec(sql,params)
	if err != nil{
		return 0,err
	}
	return res.RowsAffected()
}


func (this *dataSource) InsertOneObj(sql string,obj interface{})(int64,error){
	sqlColNames,err := parse.ParseInsertSql(sql)
	if err != nil{
		return 0,err
	}
	var params = make([]interface{},0)
	if err != nil{
		return 0,err
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
	return 0,nil
}

func (this *dataSource) InsertOneMap(sql string,paramMap map[string]interface{})(int64,error){
	colNames,err := parse.ParseInsertSql(sql)
	if err != nil{
		return 0,err
	}
	var params = make([]interface{},0)
	for _,colName := range colNames{
		params = append(params,paramMap[colName])
	}
	return this.InsertOne(sql,params...)
}

