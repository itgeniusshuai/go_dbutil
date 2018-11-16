package main

import (
	"../datasource"
	"fmt"
)

func main(){
	mainFunc()
}

type User struct{
	Id int `json:"id"`
	Name string `json:"name"`
}
func mainFunc(){
	ds := datasource.GetDataSource("root","3108220w","localhost",3306,"test")
	sql := "select *from user where id = ?"
	//resMap,err := ds.QueryOneMap(sql,1)
	//if err != nil{
	//	fmt.Println(err)
	//}
	var user User
	//res,_ := ds.QueryOne(sql,&user,1)
	res1,_ := ds.QueryMany(sql,&user,1)
	//fmt.Println(res)
	fmt.Println(res1)
}
