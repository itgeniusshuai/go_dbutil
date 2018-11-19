package datasource

import (
	"github.com/itgeniusshuai/go_common/common"
	"fmt"
	"github.com/kataras/iris/core/errors"
	"database/sql"
	"sync"
)

var txMap map[int]*sql.Tx
var txLock sync.RWMutex

func GetCurrentTX()*sql.Tx{
	txLock.Lock()
	defer txLock.Unlock()
	return txMap[common.GetGoroutineId()]
}

func PutTX(routineId int,tx *sql.Tx){
	txLock.Lock()
	defer txLock.Unlock()
	txMap[routineId] = tx
}

func DeleteTX(routineId int){
	txLock.Lock()
	defer txLock.Unlock()
	delete(txMap,routineId)
}

func (this *dataSource)TranFuncExec( s interface{},params ...interface{})(res []interface{},err error){
	tx,err := this.Begin()
	routineId := common.GetGoroutineId()
	PutTX(routineId,tx)
	defer DeleteTX(routineId)
	defer func(){
		if err1 := recover(); err1 != nil {
			tx.Rollback()
			err = errors.New(fmt.Sprintf("%v",err1))
		}
	}()
	// tx
	if err != nil{
		return nil,err
	}
	res,err = common.FuncReflectExec(s,params...)
	if err != nil{
		tx.Rollback()
		return nil,err
	}
	tx.Commit()
	return res,err
}



