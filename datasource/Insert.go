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


//func (this *dataSource) InsertOneObjSelective(sql string,obj interface{})(int64,error){
//	colNames,err := parse.ParseInsertSql(sql)
//	if err != nil{
//		return 0,err
//	}
//	var params = make([]interface{},0)
//	tagMap := common.GetStructTag(obj,"col")
//	if err != nil{
//		return 0,err
//	}
//	for _,colName := range colNames{
//		// 如果tag活着名称相同获取值否则赋值为空
//		common.IndexOfStrArr(colNames,colName)
//	}
//	return 0,nil
//}

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

