package parse

import (
	"testing"
	"time"
)

func TestParseMonth(t *testing.T) {
	cases := []struct {
		// имя теста
		name string
		// значения на вход тестируемой функции
		values []string
		// желаемый результат
		want []string
	}{
		// тестовые данные № 1
		{
			name: "correct values",
			values: []string{"2025-01-09T21:56:18+04:00",
				"2025-02-09T21:56:18+04:00",
				"2025-03-09T21:56:18+04:00",
				"2025-04-09T21:56:18+04:00",
				"2025-05-09T21:56:18+04:00",
				"2025-06-09T21:56:18+04:00",
				"2025-07-09T21:56:18+04:00",
				"2025-08-09T21:56:18+04:00",
				"2025-09-09T21:56:18+04:00",
				"2025-10-09T21:56:18+04:00",
				"2025-11-09T21:56:18+04:00",
				"2025-12-09T21:56:18+04:00",
			},
			want: []string{
				"января",
				"февраля",
				"марта",
				"апреля",
				"мая",
				"июня",
				"июля",
				"августа",
				"сентября",
				"октября",
				"ноября",
				"декабря",
			},
		},
	}

	for _, newCase := range cases {
		newCase := newCase

		t.Run(newCase.name, func(t *testing.T) {
			for i := 0; i < len(newCase.values); i++ {
				value, err := time.Parse(time.RFC3339, newCase.values[i])
				if err != nil {
					t.Fatalf("Failed to parse expected time: %v", err)
				}

				got := ParseMonth(value.Month())

				if got != newCase.want[i] {
					t.Errorf("ParseMonth(value.Month()) mismatch: got %v, want %v", got, newCase.want[i])
				}
			}
		})
	}

}
