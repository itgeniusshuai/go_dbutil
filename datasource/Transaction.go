package datasource

import (
	"github.com/itgeniusshuai/go_common/common"
	"fmt"
	"github.com/kataras/iris/core/errors"
	"database/sql"
	"sync"
	"time"
	"reflect"
)
// 暂不支持事务嵌套

func init(){
	// 防止泄露检查
	txMap = make(map[int][]*sql.Tx,0)
	go func(){
		for t := range time.Tick(time.Second * 2) {
			fmt.Println(t)
			deleteKeys := make([]int,0)
			for k,_ := range txMap{
				startTime := txTimeMap[k]
				currTime := time.Now().Second()
				if startTime == 0{
					deleteKeys = append(deleteKeys, k)
				}else if currTime - 60*60 > startTime{
					deleteKeys = append(deleteKeys, k)
					delete(txTimeMap,k)
				}
			}
			if len(deleteKeys) > 0 {
				for _,k := range deleteKeys{
					txs := txMap[k]
					if txs != nil{
						txs = nil
					}
					delete(txMap,k)
				}
			}
		}
	}()
}

var txMap map[int][]*sql.Tx
var txTimeMap map[int]int
var txLock sync.RWMutex

func GetCurrentTX()[]*sql.Tx{
	txLock.Lock()
	defer txLock.Unlock()
	return txMap[common.GetGoroutineId()]
}

func PutTX(routineId int,tx *sql.Tx) []*sql.Tx{
	txLock.Lock()
	defer txLock.Unlock()
	if txMap[routineId] == nil{
		txs := make([]*sql.Tx,0)
		txs = append(txs, tx)
		return txMap[routineId]
	}
	txTimeMap[routineId] = time.Now().Second()
	txMap[routineId] = append(txMap[routineId], tx)
	return txMap[routineId]
}

func DeleteTX(routineId int){
	txLock.Lock()
	defer txLock.Unlock()
	// 删除数组
	txs := txMap[routineId]
	if txs != nil && len(txs) > 1 {
		txs = txs[:len(txs)-1]
	}else{
		delete(txMap,routineId)
	}
}

func (this *dataSource)TranFuncExec( s interface{},propagation PROPAGATION,params ...interface{})(res []interface{},err error){
	// 如果是
	routineId := common.GetGoroutineId()
	tx,err  := this.getTxByPropagation(propagation)
	if err != nil{
		return nil,err
	}
	// 如果有事务，则共用事务
	defer DeleteTX(routineId)
	defer func(){
		if err1 := recover(); err1 != nil {
			if tx != nil{
				tx.Rollback()
			}
			err = errors.New(fmt.Sprintf("%v",err1))
		}
	}()
	res,err = common.FuncReflectExec(s,params...)
	if err != nil{
		if tx != nil{
			tx.Rollback()
		}
		return nil,err
	}
	if tx != nil{
		tx.Commit()
	}
	return res,err
}

func (this *dataSource)getTxByPropagation(propagation PROPAGATION) (*sql.Tx,error){
	routineId := common.GetGoroutineId()
	var tx *sql.Tx
	var err error
	switch propagation {
	case PROPAGATION_NEW:
		tx,err = this.Begin()
		PutTX(routineId,tx)
		break
	case PROPAGATION_REQUIRED:
		txs := GetCurrentTX()
		if txs == nil || len(txs) == 0{
			tx,err = this.Begin()
			PutTX(routineId,tx)
		}else{
			PutTX(routineId,nil)
		}
		break
	case PROPAGATION_NESTED:
		txs := GetCurrentTX()
		if txs == nil || len(txs) == 0{
			tx,err = this.Begin()
			PutTX(routineId,tx)
		}else{
			PutTX(routineId,nil)
		}
		break
	}
	return tx,err
}

func GetFistUnNilTX() *sql.Tx{
	txs := GetCurrentTX()
	// 如果不为空找到第一个不为空的事务执行
	if txs != nil && len(txs)>1{
		for i := len(txs)-1; i >0; i--{
			if txs[i] != nil{
				return txs[i]
			}
		}
	}
	return nil
}

func FullTranResult(res []interface{},fRes ...interface{})  {
	for i,e := range res{
		if e != nil{
			var fs = fRes[i]
			reflect.ValueOf(fs).Elem().Set(reflect.ValueOf(e))
		}
	}
}



