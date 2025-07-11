package models

import (
	"time"
)

type AvgRecordsPeriod struct {
	weight float64
	t      time.Time
}

func (a AvgRecordsPeriod) GetWeight() float64 {
	return a.weight
}

func (a AvgRecordsPeriod) GetTime() time.Time {
	return a.t
}

// возвращает экземпляр записи
func NewAvgRecordsPeriod(weight float64, t time.Time) AvgRecordsPeriod {
	return AvgRecordsPeriod{weight: weight, t: t}
}
