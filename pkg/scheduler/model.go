package scheduler

import (
	"time"
)

type Interval struct {
	Start time.Time
	End   time.Time
}

type BusyRecord struct {
	BayID        int32
	TechnicianID int32
	Start        time.Time
	End          time.Time
}

type Slot struct {
	Start        time.Time
	End          time.Time
	BayID        int32
	TechnicianID int32
}
