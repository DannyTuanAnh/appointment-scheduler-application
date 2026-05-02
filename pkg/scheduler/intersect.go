package scheduler

import "time"

func Intersect(a, b []Interval) []Interval {
	var result []Interval
	i, j := 0, 0

	for i < len(a) && j < len(b) {
		start := maxTime(a[i].Start, b[j].Start)
		end := minTime(a[i].End, b[j].End)

		if start.Before(end) {
			result = append(result, Interval{Start: start, End: end})
		}

		if a[i].End.Before(b[j].End) {
			i++
		} else {
			j++
		}
	}

	return result
}

func minTime(a, b time.Time) time.Time {
	if a.Before(b) {
		return a
	}
	return b
}

func maxTime(a, b time.Time) time.Time {
	if a.After(b) {
		return a
	}
	return b
}
