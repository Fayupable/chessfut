package chesscom

import (
	"regexp"
	"strconv"
	"time"
)

var archiveDateRe = regexp.MustCompile(`/(\d{4})/(\d{2})$`)

func parseArchiveDate(archiveURL string) (time.Time, bool) {
	match := archiveDateRe.FindStringSubmatch(archiveURL)
	if len(match) < 3 {
		return time.Time{}, false
	}
	year, err := strconv.Atoi(match[1])
	if err != nil {
		return time.Time{}, false
	}
	month, err := strconv.Atoi(match[2])
	if err != nil {
		return time.Time{}, false
	}
	return time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC), true
}

func monthFloor(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.UTC)
}
