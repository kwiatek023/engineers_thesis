package utils

import "strconv"

func ParseStrToPositiveInt(str string) int {
	number, err := strconv.Atoi(str)
	if err != nil || number <= 0 {
		// TODO handle error
		panic(err)
	}

	return number
}
