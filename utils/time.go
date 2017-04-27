package utils

import (
	"time"
)

var Weekdays = [...]string{
	"日",
	"一",
	"二",
	"三",
	"四",
	"五",
	"六",
}

func ConcatTime(date time.Time, clock time.Time) time.Time {
	return time.Date(date.Year(), date.Month(), date.Day(), clock.Hour(), clock.Minute(),
		clock.Second(), clock.Nanosecond(), time.Local)
}

func BeginOfDay(tm time.Time) time.Time {
	return time.Date(tm.Year(), tm.Month(), tm.Day(), 0, 0, 0, 0, tm.Location())
}

func BeginOfYesterday(tm time.Time) time.Time {
	return BeginOfDay(tm.Add(-24 * time.Hour))
}

func BeginOfTomorrow(tm time.Time) time.Time {
	return BeginOfDay(tm.Add(24 * time.Hour))
}
