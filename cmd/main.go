package main

import (
	"../datasource"
	"fmt"
	"reflect"
)

func main(){
	mainFunc()
}

type User struct{
	Id int `json:"id"`
	Name string `json:"name"`
}
func mainFunc(){
	var user User
	t := reflect.ValueOf(user)
	f := t.FieldByName("fsdf")
	fmt.Println(f.IsValid())


	ds := datasource.GetDataSource("root","123456","localhost",3306,"test")
	sql := "select *from user where id = ?"
	//resMap,err := ds.QueryOneMap(sql,1)
	//if err != nil{
	//	fmt.Println(err)
	//}
	//res,_ := ds.QueryOne(sql,&user,1)
	res1,_ := ds.QueryMany(sql,&user,1)
	//fmt.Println(res)
	fmt.Println(res1)
}
