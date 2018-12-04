package datasource

import (
 	_ "github.com/go-sql-driver/mysql"
	"database/sql"
	"reflect"
	"fmt"
	"strings"
	"github.com/itgeniusshuai/go_common/common"
	"time"
)

func GetDataSource(user string,pwd string,ip string,port int,dbName string) *DataSource{
	ds := DataSource{}
	db,err := sql.Open("mysql",user+":"+pwd+"@tcp("+ip+":"+common.IntToStr(port)+")/"+dbName)
	if err != nil{
		panic(err.Error())
	}
	ds.DB = db
	return &ds
}

// 查询单个
func (this *DataSource)QueryOne(sql string,obj ModelPtr,params ...interface{})(interface{},error){
	tx := GetFistUnNilTX()
	if tx != nil{
		return this.QueryOneWitchTX(tx,sql,obj,params...)
	}
	rows,err := this.Query(sql,params...)
	if err != nil{
		return nil,err
	}
	types,_ := rows.ColumnTypes()
	if rows.Next(){
		resMap,err := fullMap(rows,types)
		if err != nil{
			return nil,err
		}
		return mapToObj(resMap,obj),nil
	}
	return nil,nil
}

// 查询单个
func (this *DataSource)QueryOneMap(sql string,params ...interface{})(interface{},error){
	tx := GetFistUnNilTX()
	if tx != nil{
		return this.QueryOneMapWithTX(tx,sql,params...)
	}
	rows,err := this.Query(sql,params...)
	if err != nil{
		return nil,err
	}
	types,_ := rows.ColumnTypes()
	var resMap map[string]interface{}
	if rows.Next(){
		resMap,err = fullMap(rows,types)
		if err != nil{
			return nil,err
		}
	}
	return resMap,nil
}


// 查询多个
func (this *DataSource)QueryMany(sql string,obj ModelPtr,params ...interface{})(interface{},error){
	tx := GetFistUnNilTX()
	if tx != nil{
		return this.QueryManyWithTx(tx,sql,obj,params...)
	}
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
func (this *DataSource)QueryManyMap(sql string,params ...interface{})([]map[string]interface{},error){
	tx := GetFistUnNilTX()
	if tx != nil{
		return this.QueryManyMapWithTx(tx,sql,params...)
	}
	rows,err := this.Query(sql,params...)
	if err != nil{
		return nil,err
	}
	types,_ := rows.ColumnTypes()
	var resList []map[string]interface{}
	for rows.Next(){
		m,err := fullMap(rows,types)
		if err != nil{
			return nil,err
		}
		resList = append(resList, m)
	}
	return resList,nil
}

// 查询单个
func (this *DataSource)QueryOneWitchTX(tx *sql.Tx,sql string,obj ModelPtr,params ...interface{})(interface{},error){
	rows,err := tx.Query(sql,params...)
	if err != nil{
		return nil,err
	}
	types,_ := rows.ColumnTypes()
	if rows.Next(){
		resMap,err := fullMap(rows,types)
		if err != nil{
			return nil,err
		}
		return mapToObj(resMap,obj),nil
	}
	return nil,nil
}

// 查询单个
func (this *DataSource)QueryOneMapWithTX(tx *sql.Tx,sql string,params ...interface{})(map[string]interface{},error){

	rows,err := tx.Query(sql,params...)
	if err != nil{
		return nil,err
	}
	var resMap map[string]interface{}
	types,_ := rows.ColumnTypes()
	if rows.Next(){
		resMap,err = fullMap(rows,types)
		if err != nil{
			return nil,err
		}
	}
	return resMap,nil
}


// 查询多个
func (this *DataSource)QueryManyWithTx(tx *sql.Tx,sql string,obj ModelPtr,params ...interface{})(interface{},error){
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
func (this *DataSource)QueryManyMapWithTx(tx *sql.Tx,sql string,params ...interface{})([]map[string]interface{},error){
	rows,err := tx.Query(sql,params)
	if err != nil{
		return nil,err
	}
	var resList []map[string]interface{}
	types,_ := rows.ColumnTypes()
	for rows.Next(){
		m,err := fullMap(rows,types)
		if err != nil{
			return nil,err
		}
		resList = append(resList, m)
	}
	return resList,nil
}

func (this *DataSource)QueryNum(sqlStr string,params ...interface{})(int64,error){
	tx := GetFistUnNilTX()
	var err error
	var num int64
	var rows *sql.Rows
	if tx != nil{
		rows,err = this.Query(sqlStr,params...)
	}else{
		rows,err = this.Query(sqlStr,params...)
	}
	if rows.Next(){
		rows.Scan(&num)
	}
	return num,err
}


func fullObj(obj ModelPtr,rows *sql.Rows) ModelPtr{
	types,_ := rows.ColumnTypes()
	resMap,err := fullMap(rows,types)
	if err != nil{
		return nil
	}
	return mapToObj(resMap,obj)
}

func fullMap(rows *sql.Rows,types []*sql.ColumnType) (map[string]interface{},error){
	columns,err := rows.Columns()
	if err != nil{
		return nil,err
	}
	var columnNum = len(columns)
	var resList = make([]interface{},columnNum)
	values := make([]interface{}, len(columns))

	for i,_ := range values{
		typeName := types[i].ScanType().String()
		if typeName == "sql.RawBytes"{
			var str = ""
			values[i] = &str
		}else if typeName=="mysql.NullTime"{
			var t time.Time
			values[i] = &t
		}else if strings.HasPrefix(typeName,"sql.NullInt"){
			var t  interface{}
			values[i] = &t
		}else if typeName == "int8"{
			var t int8
			values[i] = &t
		}else if typeName == "int16"{
			var t int16
			values[i] = &t
		}else if typeName == "int32"{
			var t int32
			values[i] = &t
		}else if typeName == "int64"{
			var t int64
			values[i] = &t
		}else if typeName == "int"{
			var t int8
			values[i] = &t
		}else if typeName == "uint"{
			var t uint
			values[i] = &t
		}else if typeName == "uint8"{
			var t uint8
			values[i] = &t
		}else if typeName == "uint16"{
			var t uint16
			values[i] = &t
		}else if typeName == "uint32"{
			var t int32
			values[i] = &t
		}else if typeName == "uint64"{
			var t int64
			values[i] = &t
		}else if typeName == "float32"{
			var t float32
			values[i] = &t
		}else if typeName == "float64"{
			var t float64
			values[i] = &t
		}else{
			values[i] = reflect.New(types[i].ScanType()).Interface()
		}
	}
	for i := range values {
		resList[i] = values[i]
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

func InterfaceToValue(v interface{},t *sql.ColumnType) interface{}{
	var result = v
	fmt.Println(t.DatabaseTypeName())
	dbType := strings.ToLower(t.DatabaseTypeName())
	switch dbType {
	case "int":result = v.(int);break;
	case "int8":result = v.(int8);break;
	case "int16":result = v.(int16);break;
	case "int32":result = v.(int32);break;
	case "int64":result = v.(int64);break;
	case "uint":result = v.(uint);break;
	case "uint8":result = v.(uint8);break;
	case "uint16":result = v.(uint16);break;
	case "uint32":result = v.(uint32);break;
	case "uint64":result = v.(uint64);break;
	case "string":result = v.(string);break;
	case "float32":result = v.(float32);break;
	case "float64":result = v.(float64);break;
	}
	return result
}

