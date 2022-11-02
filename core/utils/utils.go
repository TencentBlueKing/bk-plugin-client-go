package utils

import "strconv"

func CovertStrInt(strContent string) int {
	result, err := strconv.Atoi(strContent)
	if err == nil {
		return result
	}
	return 0
}
