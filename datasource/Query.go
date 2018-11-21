package datasource

import (
 	_ "github.com/go-sql-driver/mysql"
	"database/sql"
	"reflect"
	"github.com/itgeniusshuai/go_common/common"
	"fmt"
)

func GetDataSource(user string,pwd string,ip string,port int,dbName string) *dataSource{
	ds := dataSource{}
	db,err := sql.Open("mysql",user+":"+pwd+"@tcp("+ip+":"+common.IntToStr(port)+")/"+dbName)
	if err != nil{
		panic(err.Error())
	}
	ds.DB = db
	return &ds
}

// 查询单个
func (this *dataSource)QueryOne(sql string,obj ModelPtr,params ...interface{})(interface{},error){
	rows,err := this.Query(sql,params...)
	if err != nil{
		return nil,err
	}

	if rows.Next(){
		resMap,err := fullMap(rows)
		if err != nil{
			return nil,err
		}
		return mapToObj(resMap,obj),nil
	}
	return nil,nil
}

// 查询单个
func (this *dataSource)QueryOneMap(sql string,params ...interface{})(interface{},error){

	rows,err := this.Query(sql,params...)
	if err != nil{
		return nil,err
	}
	var resMap map[string]interface{}
	if rows.Next(){
		resMap,err = fullMap(rows)
		if err != nil{
			return nil,err
		}
	}
	return resMap,nil
}


// 查询多个
func (this *dataSource)QueryMany(sql string,obj ModelPtr,params ...interface{})(interface{},error){
	rows,err := this.Query(sql,params...)
	if err != nil{
		return nil,err
	}
	var resList []ModelPtr
	for rows.Next(){
		m := fullObj(obj,rows)
		fmt.Println(m)
		resList = append(resList, m)
	}
	return resList,nil
}

// 查询多个
func (this *dataSource)QueryManyMap(sql string,params ...interface{})(interface{},error){
	rows,err := this.Query(sql,params)
	if err != nil{
		return nil,err
	}
	var resList []map[string]interface{}
	for rows.Next(){
		m,err := fullMap(rows)
		if err != nil{
			return nil,err
		}
		resList = append(resList, m)
	}
	return resList,nil
}

// 查询单个
func (this *dataSource)QueryOneWitchTX(tx sql.Tx,sql string,obj ModelPtr,params ...interface{})(interface{},error){
	rows,err := tx.Query(sql,params...)
	if err != nil{
		return nil,err
	}

	if rows.Next(){
		resMap,err := fullMap(rows)
		if err != nil{
			return nil,err
		}
		return mapToObj(resMap,obj),nil
	}
	return nil,nil
}

// 查询单个
func (this *dataSource)QueryOneMapWithTX(tx *sql.Tx,sql string,params ...interface{})(interface{},error){

	rows,err := tx.Query(sql,params...)
	if err != nil{
		return nil,err
	}
	var resMap map[string]interface{}
	if rows.Next(){
		resMap,err = fullMap(rows)
		if err != nil{
			return nil,err
		}
	}
	return resMap,nil
}


// 查询多个
func (this *dataSource)QueryManyWithTx(tx *sql.Tx,sql string,obj ModelPtr,params ...interface{})(interface{},error){
	rows,err := tx.Query(sql,params...)
	if err != nil{
		return nil,err
	}
	var resList []ModelPtr
	for rows.Next(){
		m := fullObj(obj,rows)
		fmt.Println(m)
		resList = append(resList, m)
	}
	return resList,nil
}

// 查询多个
func (this *dataSource)QueryManyMapWithTx(tx *sql.Tx,sql string,params ...interface{})(interface{},error){
	rows,err := tx.Query(sql,params)
	if err != nil{
		return nil,err
	}
	var resList []map[string]interface{}
	for rows.Next(){
		m,err := fullMap(rows)
		if err != nil{
			return nil,err
		}
		resList = append(resList, m)
	}
	return resList,nil
}


func fullObj(obj ModelPtr,rows *sql.Rows) ModelPtr{
	resMap,err := fullMap(rows)
	if err != nil{
		return nil
	}
	return mapToObj(resMap,obj)
}

func fullMap(rows *sql.Rows) (map[string]interface{},error){
	columns,err := rows.Columns()
	if err != nil{
		return nil,err
	}
	fmt.Println(columns)
	var columnNum = len(columns)
	var resList = make([]interface{},columnNum)
	values := make([]interface{}, len(columns))
	for i := range values {
		resList[i] = &values[i]
	}
	rows.Scan(resList...)
	var resMap = make(map[string]interface{})
	for i,e := range columns{
		resMap[e] = resList[i]
	}
	return resMap,nil
}

// 将map中与结构体中对应的field对应上
func mapToObj(m map[string]interface{},obj ModelPtr) ModelPtr{
	objType := reflect.TypeOf(obj)
	newObj := reflect.New(objType.Elem()).Interface()
	newObjValue := reflect.ValueOf(newObj).Elem()

	for k,v := range m{
		f := newObjValue.FieldByName(k)
		if f.IsValid(){
			f.Set(reflect.ValueOf(v))
		}
	}

	return newObj
}
