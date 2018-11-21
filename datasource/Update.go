package datasource

import (
	"database/sql"
	"bytes"
	"reflect"
	"fmt"
	"strings"
)

func (this *dataSource)Update(sqlStr string,params ...interface{})(int64,error) {
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

func (this *dataSource)UpdateSelectiveBySelect(tableName string,objMap map[string]interface{},whereMap map[string]interface{},params ...interface{})(int64,error) {
	var sqlBuf bytes.Buffer
	sqlBuf.WriteString(" udpate ")
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


