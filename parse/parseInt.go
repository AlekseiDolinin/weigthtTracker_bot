package parse

import (
	"strconv"
	"strings"
)

// ищет в строке первое вхождение числа с плавающей точкой с разделителем: запятой или точкой
func ParseInt(imput string) (int64, error) {
	imputSplit := strings.Split(imput, "")
	var newStr []string
	for _, char := range imputSplit {
		num := string(char)
		if num == "0" || num == "1" || num == "2" || num == "3" || num == "4" || num == "5" || num == "6" || num == "7" || num == "8" || num == "9" {
			newStr = append(newStr, char)
		}
	}
	return strconv.ParseInt(strings.Join(newStr, ""), 10, 64)
}
