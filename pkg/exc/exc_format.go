package exc

import (
	"drive/pkg/logger"
	"strconv"
)

// string to uint
func StrToUint(str string) (uint, error) {
	idUint, err := strconv.ParseUint(str, 10, 32)
	if err != nil {
		logger.Error("StrToUint类型转换失败")
		return 0, err
	}
	return uint(idUint), nil
}
