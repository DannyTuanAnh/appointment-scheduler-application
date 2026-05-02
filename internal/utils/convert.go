package utils

import (
	"regexp"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

var (
	matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
	matchAllCap   = regexp.MustCompile("([a-z0-9])([A-Z])")
)

func CamelToSnake(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

func NormalizeToSameDate(t time.Time) time.Time {
	return time.Date(0, 1, 1, t.Hour(), t.Minute(), 0, 0, time.UTC)
}

func ConvertTimeToPgTypeTime(input time.Time) pgtype.Time {
	if input.IsZero() {
		return pgtype.Time{
			Valid: false,
		}
	}

	nano := time.Duration(input.Hour())*time.Hour +
		time.Duration(input.Minute())*time.Minute +
		time.Duration(input.Second())*time.Second +
		time.Duration(input.Nanosecond())

	return pgtype.Time{
		Microseconds: nano.Microseconds(),
		Valid:        true,
	}
}

func ConvertPgTypeTimeToTime(input pgtype.Time) time.Time {
	if !input.Valid {
		return time.Time{}
	}

	d := time.Duration(input.Microseconds) * time.Microsecond

	zeroTime := time.Date(1, 1, 1, 0, 0, 0, 0, time.FixedZone("UTC+7", 7*60*60))
	return zeroTime.Add(d)
}

func ConvertPgTypeTimeToString(input pgtype.Time) string {
	if !input.Valid {
		return ""
	}

	d := time.Duration(input.Microseconds) * time.Microsecond

	zeroTime := time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC)
	t := zeroTime.Add(d)

	return t.Format("15:04")
}
