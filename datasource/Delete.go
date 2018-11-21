package datasource

func Delete(sql string,params ...interface{})error {
	tx := GetFistUnNilTX()
	if tx == nil{

	}
}
