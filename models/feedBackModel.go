package models

import "time"

type FeedBack struct {
	date    time.Time
	userID  int64
	message string
}

func (f FeedBack) GetUseID() int64 {
	return f.userID
}

func (f FeedBack) GetMsg() string {
	return f.message
}

func (f FeedBack) GetTime() time.Time {
	return f.date
}

func NewFeedBack(d time.Time, id int64, msg string) FeedBack {
	return FeedBack{date: d, userID: id, message: msg}
}
