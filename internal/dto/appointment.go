package dto

type GetAppointmentRequest struct {
	ServiceID           int32  `form:"service_id" binding:"required,gt=0"`
	BayTypeID           int32  `form:"bay_type_id" binding:"required,gt=0"`
	DealershipID        int32  `form:"dealership_id" binding:"required,gt=0"`
	PreferenceTime      string `form:"preference_time" binding:"required,datetime=2006-01-02"`
	AnticipatedDuration int32  `form:"anticipated_duration" binding:"required,gt=0"`
}

type BusyIntervalResponse struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

type AvailabilityResponse struct {
	Date            string                 `json:"date"`
	WorkStart       string                 `json:"work_start"`
	WorkEnd         string                 `json:"work_end"`
	DurationMinutes int                    `json:"duration_minutes"`
	Busy            []BusyIntervalResponse `json:"busy"`
}
