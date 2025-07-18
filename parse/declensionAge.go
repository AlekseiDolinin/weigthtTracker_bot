package parse

var y = map[int]string{
	1: "год",
	2: "года",
	3: "года",
	4: "года",
	5: "лет",
	6: "лет",
	7: "лет",
	8: "лет",
	9: "лет",
	0: "лет",
}

func DeclensionAge(age int) string {
	switch age {
	case 11, 12, 13, 14:
		return y[0]
	default:
		year := age % 10
		return y[year]
	}
}
