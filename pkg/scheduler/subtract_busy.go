package scheduler

func SubtractBusy(work Interval, busy []Interval) []Interval {
	if len(busy) == 0 {
		return []Interval{work}
	}

	var free []Interval
	cursor := work.Start

	for _, b := range busy {
		if b.Start.After(cursor) {
			free = append(free, Interval{
				Start: cursor,
				End:   b.Start,
			})
		}
		if b.End.After(cursor) {
			cursor = b.End
		}
	}

	if cursor.Before(work.End) {
		free = append(free, Interval{
			Start: cursor,
			End:   work.End,
		})
	}

	return free
}
