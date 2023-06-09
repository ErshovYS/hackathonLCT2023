package string_generator

import (
	"math"
	"strconv"
	"unicode/utf8"
)

func getString(val int) string {
	return strconv.FormatInt(int64(val), 10)
}

func floatToString(val float64) string {
	newVal := 0.0
	st := 0.0
	n := 0
	for i := 5; i >= 0; i-- {
		st = math.Pow(1000, float64(i))
		if val >= st {
			newVal = val / st
			n = i
			break
		}
	}
	intVal := int64(newVal)
	ost := val - float64(intVal)*st
	ost = ost / math.Pow(10, float64(n)*3-2)
	resVal := strconv.FormatInt(int64(newVal), 10)
	resOst := strconv.FormatInt(int64(ost), 10)
	return resVal + "." + resOst
}

func floatToStringFrom(val float64, to float64) string {
	newVal := 0.0
	st := 0.0
	n := 0
	for i := 5; i >= 0; i-- {
		st = math.Pow(1000, float64(i))
		if to >= st {
			newVal = val / st
			n = i
			break
		}
	}
	intVal := int64(newVal)
	ost := val - float64(intVal)*st
	ost = ost / math.Pow(10, float64(n)*3-2)
	resVal := strconv.FormatInt(int64(newVal), 10)
	resOst := strconv.FormatInt(int64(ost), 10)
	return resVal + "." + resOst
}

func getUnits(val float64) string {
	Units := []string{"", "тыс.", "млн.", "млрд.", "трлн.", "трлрд."}
	for i := 5; i >= 0; i-- {
		if val >= math.Pow(1000, float64(i)) {
			return " " + Units[i] + "руб"
		}
	}
	return "руб"
}

func parseString(str string) [4]string {
	res := [4]string{}

	strLen := utf8.RuneCountInString(str)

	localLen := 37

	if len(str) < localLen {
		res[0] = str
		return res
	}

	pos := 0
	oldPos := 0
	num := 0
	oldNum := 0
	newNum := 0

	for i, c := range str {
		if num > localLen {
			break
		}
		if c == ' ' {
			newNum = num
			pos = i
			oldPos = i
		}
		num++
	}
	res[0] = str[0:pos]

	if strLen-oldNum < localLen {
		res[1] = str[pos:]
		return res
	} else if pos < len(str) {
		num = 0
		oldNum = newNum
		for i, c := range str {
			if num > oldNum+localLen {
				break
			}
			if c == ' ' {
				newNum = num
				pos = i
			}
			num++
		}
		res[1] = str[oldPos:pos]
		oldPos = pos
	}

	if strLen-oldNum < localLen {
		res[2] = str[pos:]
		return res
	}
	if pos < len(str) {
		num = 0
		oldNum = newNum
		for i, c := range str {
			if num > oldNum+localLen {
				break
			}
			if c == ' ' {
				newNum = num
				pos = i
			}
			num++
		}
		res[2] = str[oldPos:pos]
		oldPos = pos
	}

	res[3] = str[pos:]

	return res
}
