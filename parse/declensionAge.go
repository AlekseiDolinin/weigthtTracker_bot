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
	year := age % 10
	return y[year]
}
