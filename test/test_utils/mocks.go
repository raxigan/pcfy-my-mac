package test_utils

import (
	"time"
)

type FakeTimeProvider struct {
}

func (tp FakeTimeProvider) Now() time.Time {
	parse, _ := time.Parse("2006-01-02 15:04:05", "2023-09-27 12:30:00")
	return parse
}
