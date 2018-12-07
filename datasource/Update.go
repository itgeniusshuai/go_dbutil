package datasource

import (
	"database/sql"
	"bytes"
	"reflect"
	"fmt"
	"strings"
	"github.com/itgeniusshuai/go_common/common"
)

func (this *DataSource)Update(sqlStr string,params ...interface{})(int64,error) {
	var err error
	var res sql.Result
	res, err = this.Exec(sqlStr, params...)
	if err !=  nil{
		return 0,err
	}
	return res.RowsAffected()
}

func (this *DataSource)UpdateBySelect(tableName string,objMap map[string]interface{},whereMap map[string]interface{},params ...interface{})(int64,error) {
	var sqlBuf bytes.Buffer
	sqlBuf.WriteString(" update ")
	sqlBuf.WriteString(tableName)
	sqlBuf.WriteString(" set ")
	var updateStrBuf []string
	var whereStrBuf []string
	for k,v := range objMap{
		if reflect.TypeOf(v).Kind().String() == "string"{
			updateStrBuf = append(updateStrBuf, fmt.Sprintf("%s=\"$v\"",k,v))
		}else{
			updateStrBuf = append(updateStrBuf, fmt.Sprintf("%s=%v",k,v))
		}
	}
	for k,v := range whereStrBuf{
		if reflect.TypeOf(v).Kind().String() == "string"{
			whereStrBuf = append(whereStrBuf, fmt.Sprintf("%s=\"$v\"",k,v))
		}else{
			whereStrBuf = append(whereStrBuf, fmt.Sprintf("%s=%v",k,v))
		}
	}
	sqlBuf.WriteString(strings.Join(updateStrBuf,","))
	sqlBuf.WriteString(strings.Join(whereStrBuf," and "))
	sqlStr := sqlBuf.String()
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

func (this *DataSource)UpdateByMap(sql string,params map[string]interface{}) (int64,error){
	// 替换#{}里面的内容
	reSql,paramList,err := ReplaceSqlPlaceHolders(sql,params)
	if err != nil{
		return 0,err
	}
	res,err := this.Exec(reSql,paramList...)
	if err != nil{
		return 0,err
	}
	return res.RowsAffected();
}

func (this *DataSource)UpdateByObj(sql string,obj interface{}) (int64,error){
	// 替换#{}里面的内容
	paramMap := common.ObjToMapByTagName(obj,DB_STRUCT_TAG)
	return this.UpdateByMap(sql,paramMap)
}

func (this *DataSource)UpdateWithTX(tx *sql.Tx,sqlStr string,params ...interface{})(int64,error) {
	var err error
	var res sql.Result
	res, err = tx.Exec(sqlStr, params...)
	if err !=  nil{
		return 0,err
	}
	return res.RowsAffected()
}

func (this *DataSource)UpdateBySelectWithTX(tx *sql.Tx,tableName string,objMap map[string]interface{},whereMap map[string]interface{},params ...interface{})(int64,error) {
	var sqlBuf bytes.Buffer
	sqlBuf.WriteString(" update ")
	sqlBuf.WriteString(tableName)
	sqlBuf.WriteString(" set ")
	var updateStrBuf []string
	var whereStrBuf []string
	for k,v := range objMap{
		if reflect.TypeOf(v).Kind().String() == "string"{
			updateStrBuf = append(updateStrBuf, fmt.Sprintf("%s=\"%v\"",k,v))
		}else{
			updateStrBuf = append(updateStrBuf, fmt.Sprintf("%s=%v",k,v))
		}
	}
	for k,v := range whereStrBuf{
		if reflect.TypeOf(v).Kind().String() == "string"{
			whereStrBuf = append(whereStrBuf, fmt.Sprintf("%s=\"$v\"",k,v))
		}else{
			whereStrBuf = append(whereStrBuf, fmt.Sprintf("%s=%v",k,v))
		}
	}
	sqlBuf.WriteString(strings.Join(updateStrBuf,","))
	sqlBuf.WriteString(strings.Join(whereStrBuf," and "))
	sqlStr := sqlBuf.String()
	var err error
	var res sql.Result
	res,err = tx.Exec(sqlStr,params...)
	if err !=  nil{
		return 0,err
	}
	return res.RowsAffected()
}

func (this *DataSource)UpdateByMapWithTX(tx *sql.Tx, sql string, params map[string]interface{}) (int64,error){
	// 替换#{}里面的内容
	reSql,paramList,err := ReplaceSqlPlaceHolders(sql,params)
	if err != nil{
		return 0,err
	}
	res,err := tx.Exec(reSql,paramList...)
	if err != nil{
		return 0,err
	}
	return res.RowsAffected();
}

func (this *DataSource)UpdateByObjWithTX(tx *sql.Tx, sql string,obj interface{}) (int64,error){
	// 替换#{}里面的内容
	paramMap := common.ObjToMapByTagName(obj,DB_STRUCT_TAG)
	return this.UpdateByMapWithTX(tx,sql,paramMap)
}







