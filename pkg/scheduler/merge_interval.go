package scheduler

import "sort"

func MergeIntervals(intervals []Interval) []Interval {
	if len(intervals) == 0 {
		return nil
	}

	sort.Slice(intervals, func(i, j int) bool {
		return intervals[i].Start.Before(intervals[j].Start)
	})

	merged := []Interval{intervals[0]}

	for _, curr := range intervals[1:] {
		last := &merged[len(merged)-1]

		if !curr.Start.After(last.End) { // overlap or touching
			if curr.End.After(last.End) {
				last.End = curr.End
			}
		} else {
			merged = append(merged, curr)
		}
	}

	return merged
}
