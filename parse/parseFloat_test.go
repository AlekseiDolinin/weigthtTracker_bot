package parse

import "testing"

func TestParseFloat(t *testing.T) {
	cases := []struct {
		// имя теста
		name string
		// значения на вход тестируемой функции
		values []string
		// желаемый результат
		want float64
	}{
		// тестовые данные № 1
		{
			name:   "correct values",
			values: []string{"10.10", "10,10", "10,1", "10.1"},
			want:   10.10,
		},
		// тестовые данные № 2
		{
			name:   "mixed values",
			values: []string{"a10.10", "10,10v", "10,s1", "10f.1"},
			want:   10.10,
		},
		// тестовые данные № 3
		{
			name:   "negative values",
			values: []string{"-10.10", "-10,10", "-10,s1", "-10f.1"},
			want:   10.10,
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			for _, value := range tc.values {

				got, _ := ParseFloat(value)

				if got != tc.want {
					t.Errorf("ParseFloat(%v) = %v; want %v", value, got, tc.want)
				}
			}
		})
	}
}
