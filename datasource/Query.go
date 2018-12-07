package datasource

import (
 	_ "github.com/go-sql-driver/mysql"
	"database/sql"
	"reflect"
	"fmt"
	"strings"
	"github.com/itgeniusshuai/go_common/common"
	"time"
	"github.com/kataras/iris/core/errors"
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
func (this *DataSource)QueryOne(sql string,obj ModelPtr,params ...interface{})(ModelPtr,error){
	rows,err := this.Query(sql,params...)
	if err != nil{
		return nil,err
	}
	res,err := fullObjList(rows,obj)
	if err != nil{
		return nil,err
	}
	if len(res) > 1{
		return nil,errors.New("expect one but find many")
	}
	return res[0],nil
}

// 查询单个
func (this *DataSource)QueryOneMap(sql string,params ...interface{})(map[string]interface{},error){
	rows,err := this.Query(sql,params...)
	if err != nil{
		return nil,err
	}
	res,err := fullMapList(rows)
	if err != nil{
		return nil,err
	}
	if len(res) > 0 {
		return nil,errors.New("expect one but find many")
	}
	return res[0],nil
}


// 查询多个
func (this *DataSource)QueryMany(sql string,obj ModelPtr,params ...interface{})([]ModelPtr,error){

	rows,err := this.Query(sql,params...)
	if err != nil{
		return nil,err
	}
	return fullObjList(rows,obj)
}

// 查询多个
func (this *DataSource)QueryManyMap(sql string,params ...interface{})([]map[string]interface{},error){

	rows,err := this.Query(sql,params...)
	if err != nil{
		return nil,err
	}
	return fullMapList(rows)
}

// 查询单个
func (this *DataSource)QueryOneWitchTX(tx *sql.Tx,sql string,obj ModelPtr,params ...interface{})(interface{},error){
	rows,err := tx.Query(sql,params...)
	if err != nil{
		return nil,err
	}
	res,err := fullObjList(rows,obj)
	if err != nil{
		return nil,err
	}
	if len(res) > 1{
		return nil,errors.New("expect one but find many")
	}
	return res[0],nil
}

// 查询单个
func (this *DataSource)QueryOneMapWithTX(tx *sql.Tx,sql string,params ...interface{})(map[string]interface{},error){

	rows,err := tx.Query(sql,params...)
	if err != nil{
		return nil,err
	}
	res,err := fullMapList(rows)
	if err != nil{
		return nil,err
	}
	if len(res) > 1{
		return nil,errors.New("expect one but find many")
	}
	return res[0],err
}


// 查询多个
func (this *DataSource)QueryManyWithTx(tx *sql.Tx,sql string,obj ModelPtr,params ...interface{})([]ModelPtr,error){
	rows,err := tx.Query(sql,params...)
	if err != nil{
		return nil,err
	}
	return fullObjList(rows,obj)
}

// 查询多个
func (this *DataSource)QueryManyMapWithTx(tx *sql.Tx,sql string,params ...interface{})([]map[string]interface{},error){
	rows,err := tx.Query(sql,params)
	if err != nil{
		return nil,err
	}
	return fullMapList(rows)
}

func (this *DataSource)QueryNum(sqlStr string,params ...interface{})(int64,error){
	tx := GetFistUnNilTX()
	var err error
	var num int64
	var rows *sql.Rows
	if tx != nil{
		rows,err = tx.Query(sqlStr,params...)
	}else{
		rows,err = this.Query(sqlStr,params...)
	}
	if rows.Next(){
		rows.Scan(&num)
	}
	return num,err
}

func fullMapList(rows *sql.Rows)([]map[string]interface{},error){
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

func fullObjList(rows *sql.Rows,obj ModelPtr)([]ModelPtr,error){
	var resList []ModelPtr
	for rows.Next(){
		m := fullObj(obj,rows)
		fmt.Println(m)
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
	types,err := rows.ColumnTypes()
	columns,err := rows.Columns()
	if err != nil{
		return nil,err
	}
	var columnNum = len(columns)
	var resList = make([]interface{},columnNum)
	values := getValuesByScanType(types)
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
		}else{
			// 根据标签赋值
			fullByTagValue(objType,newObj,k,v)
		}
	}
	return newObj
}

func fullByTagValue(vt reflect.Type,vv interface{},tagValue string, v interface{}){
	vt1 := vt.Elem()
	fn := vt1.NumField()

	for i := 0; i < fn; i++{
		tag := vt1.Field(i).Tag
		curTagValue := tag.Get(DB_STRUCT_TAG)
		if curTagValue == tagValue{
			// 类型不同不尝试强转
			fv := common.InterfacePtrToInterface(v)
			reflect.ValueOf(vv).Elem().Field(i).Set(reflect.ValueOf(fv))
		}
	}
}

func getValuesByScanType(types []*sql.ColumnType)[]interface{}{
	values := make([]interface{}, len(types))

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
	return values
}

