package models

import (
	"time"
)

type Record struct {
	id     int
	weight float64
	t      time.Time
}

func (r Record) GetId() int {
	return r.id
}

func (r Record) GetWeight() float64 {
	return r.weight
}

func (r Record) GetTime() time.Time {
	return r.t
}

// возвращает экземпляр записи
func NewRecord(id int, weight float64, t time.Time) Record {
	return Record{id: id, weight: weight, t: t}
}
