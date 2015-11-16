package utils

import "time"

var Weekdays = [...]string{
	"日",
	"一",
	"二",
	"三",
	"四",
	"五",
	"六",
}

const (
	TIME_PATTERN  = "2006-01-02 15:04"
	DATE_PATTERN  = "2006-01-02"
	CLOCK_PATTERN = "15:04"
)

var (
	Location *time.Location
)
