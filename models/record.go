package models

import (
	"time"
)

type Record struct {
	id      int
	weight  float64
	t       time.Time
	deleted int //0 - exist, 1 - deleted
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

func (r Record) GetStatus() int {
	return r.deleted
}

func (r *Record) SetStatus(delete int) {
	r.deleted = delete
}

// возвращает экземпляр записи
func NewRecord(id int, weight float64, t time.Time, deleted int) Record {
	return Record{id: id, weight: weight, t: t, deleted: deleted}
}
