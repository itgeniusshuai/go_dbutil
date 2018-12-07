package datasource

import (
	"errors"
)

func ReplaceSqlPlaceHolders(sql string, paramMap map[string]interface{})(string,[]interface{},error){
	sql =  sql + "$"
	var stack []rune
	var paramList []interface{}
	var runes = []rune(sql)
	var i = 0
	var ch = runes[i]
	for ch != '$'{
		switch ch {
		case '#':
			if runes[i+1] != '{'{
				return "",nil,errors.New("#后面必须跟{")
			}
		case '{':
			var key []rune
			for ch != '$'{
				i++
				ch = runes[i]
				key = append(key, ch)
				if runes[i+1] == '}'{
					break
				}
			}
			v := paramMap[string(key)]
			if v == nil{
				return "",nil,errors.New("参数"+string(key)+"不能为空")
			}
			paramList = append(paramList, v)
			stack = append(stack, '?')
		case '}':
		default:
			stack = append(stack, ch)
		}
		i++
		ch = runes[i]
	}

	return string(stack),paramList,nil
}

//func InterfaceToRunes(v interface{}) ([]rune,error){
//	var runes []rune
//	switch v1 := v.(type) {
//	case string:
//		runes = append(runes, '"')
//		runes = append(runes, []rune(v1)...)
//		runes = append(runes, '"')
//	case int:runes = append(runes, []rune(common.IntToStr(v1))...)
//	case int8:runes = append(runes, []rune(common.Int8ToStr(v1))...)
//	case int16:runes = append(runes, []rune(common.Int16ToStr(v1))...)
//	case int32:runes = append(runes, []rune(common.Int32ToStr(v1))...)
//	case int64:runes = append(runes, []rune(common.Int64ToStr(v1))...)
//	case uint:runes = append(runes, []rune(common.UintToStr(v1))...)
//	case uint8:runes = append(runes, []rune(common.Uint8ToStr(v1))...)
//	case uint16:runes = append(runes, []rune(common.Uint16ToStr(v1))...)
//	case uint32:runes = append(runes, []rune(common.Uint32ToStr(v1))...)
//	case uint64:runes = append(runes, []rune(common.Uint64ToStr(v1))...)
//	case float32:runes = append(runes, []rune(common.Float32ToStr(v1))...)
//	case float64:runes = append(runes, []rune(common.Float64ToStr(v1))...)
//	default:
//		return nil,errors.New("不支持的类型"+reflect.TypeOf(v).Kind().String())
//
//	}
//	return runes,nil
//}