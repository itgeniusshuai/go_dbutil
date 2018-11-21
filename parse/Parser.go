package parse

import (
	"regexp"
	"strings"
	"github.com/kataras/iris/core/errors"
)

func ParseInsertSql(sql string) ([]string,error){
	sql = strings.TrimSpace(sql)
	reg,_ := regexp.Compile(`(insert|INSERT)\s+(into|INTO)\s+([^\s]+)\(([^s]+)\)values\([^\s]+\)`)
	match := reg.FindStringSubmatch(sql)
	if len(match) <5{
		return nil,errors.New("非法sql")
	}
	insertColumns := match[4]
	reSep,_ := regexp.Compile("\\s*,\\s*")
	return reSep.Split(insertColumns,-1),nil
}




