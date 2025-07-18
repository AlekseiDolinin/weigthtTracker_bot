package parse

import (
	"testing"
)

func TestDeclensionAge(t *testing.T) {
	cases := []struct {
		name   string
		values []int
		want   []string
	}{
		// тестовые данные № 1
		{
			name:   "second decade",
			values: []int{20, 21, 22, 23, 24, 25, 26, 27, 28, 29},
			want: []string{
				"лет",
				"год",
				"года",
				"года",
				"года",
				"лет",
				"лет",
				"лет",
				"лет",
				"лет",
			},
		},
		// тестовые данные № 2
		{
			name:   "third decade",
			values: []int{30, 31, 32, 33, 34, 35, 36, 37, 38, 39},
			want: []string{
				"лет",
				"год",
				"года",
				"года",
				"года",
				"лет",
				"лет",
				"лет",
				"лет",
				"лет",
			},
		},
	}

	for _, newCase := range cases {
		newCase := newCase

		t.Run(newCase.name, func(t *testing.T) {
			for i := 0; i < len(newCase.values); i++ {
				value := newCase.values[i]
				got := DeclensionAge(value)
				if got != newCase.want[i] {
					t.Errorf("ParseAge(value) mismatch: got %v, want %v", got, newCase.want[i])
				}
			}
		})
	}
}
