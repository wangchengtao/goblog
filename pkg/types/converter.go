package types

import (
	"goblog/pkg/logger"
	"strconv"
)

// Int64ToString 将 int64 转换为 string
func Int64ToString(num int64) string {
	return strconv.FormatInt(num, 10)
}

func Uint64ToString(num uint64) string {
	return strconv.FormatUint(num, 10)
}

func StringToInt(num string) int {
	i, err := strconv.Atoi(num)

	if err != nil {
		logger.LogError(err)
	}

	return i
}
