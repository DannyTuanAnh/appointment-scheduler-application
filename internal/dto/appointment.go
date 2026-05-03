package dto

import "time"

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

// APPOINTMENTS MANAGEMENT

type RequestAppointmentWithID struct {
	ID int32 `uri:"id" binding:"required,gt=0"`
}

// Use for list endpoints that require a time window.
// Expected format: RFC3339 (ex: 2026-05-03T09:00:00+07:00)
// Note: binding tag uses 'datetime=2006-01-02T15:04:05Z07:00'.

type ListAppointmentsInTimeRangeQuery struct {
	FromTime string `form:"from_time" binding:"required,datetime=2006-01-02T15:04:05Z07:00"`
	ToTime   string `form:"to_time" binding:"required,datetime=2006-01-02T15:04:05Z07:00"`
}

type ListAppointmentsByDealershipQuery struct {
	DealershipID int32 `form:"dealership_id" binding:"required,gt=0"`
	ListAppointmentsInTimeRangeQuery
}

type ListAppointmentsByTechnicianQuery struct {
	TechnicianID int32 `form:"technician_id" binding:"required,gt=0"`
	ListAppointmentsInTimeRangeQuery
}

type ListAppointmentsByServiceBayQuery struct {
	BayID int32 `form:"bay_id" binding:"required,gt=0"`
	ListAppointmentsInTimeRangeQuery
}

type ListAppointmentsByServiceQuery struct {
	ServiceID int32 `form:"service_id" binding:"required,gt=0"`
	ListAppointmentsInTimeRangeQuery
}

type SearchAppointmentsQuery struct {
	DealershipID int32  `form:"dealership_id" binding:"required,gt=0"`
	CustomerName string `form:"customer_name" binding:"required,not_blank,min=1,max=255"`
}

type CreateAppointmentRequest struct {
	DealershipID int32  `json:"dealership_id" binding:"required,gt=0"`
	ServiceID    int32  `json:"service_id" binding:"required,gt=0"`
	BayTypeID    int32  `json:"bay_type_id" binding:"required,gt=0"`
	CustomerName string `json:"customer_name" binding:"required,not_blank,min=1,max=255"`
	StartTime    string `json:"start_time" binding:"required,datetime=2006-01-02T15:04:05Z07:00"`
	EndTime      string `json:"end_time" binding:"required,datetime=2006-01-02T15:04:05Z07:00"`
}

type PairKey struct {
	BayID        int32
	TechnicianID int32
}

type CreateAppointmentParams struct {
	DealershipID      int32
	ServiceID         int32
	PairKeyBayAndTech []PairKey
	CustomerName      string
	StartTime         time.Time
	EndTime           time.Time
}

// Status values are validated in service layer against known enum values.
// Using string in DTO for easier client usage.

type UpdateAppointmentStatusRequest struct {
	Status string `json:"status" binding:"required,not_blank,oneof=in_progress completed no_show"`
}

type CancelledAppointmentRequest struct {
	Status string `json:"status" binding:"required,not_blank,oneof=in_progress completed no_show"`
}

type MarkNoShowAppointmentsRequest struct {
	AppointmentIds []int32 `json:"appointment_ids" binding:"required,min=1,dive,gt=0"`
}

// Unified response object for list/search endpoints

type AppointmentResponseHTTP struct {
	ID                  int32  `json:"id"`
	DealershipID        int32  `json:"dealership_id"`
	DealershipName      string `json:"dealership_name"`
	DealershipOpenTime  string `json:"dealership_open_time"`
	DealershipCloseTime string `json:"dealership_close_time"`
	ServiceID           int32  `json:"service_id"`
	ServiceName         string `json:"service_name"`
	AnticipatedMinutes  int32  `json:"anticipated_minutes"`
	BayID               int32  `json:"bay_id"`
	BayName             string `json:"bay_name"`
	TechnicianID        int32  `json:"technician_id"`
	TechnicianName      string `json:"technician_name"`
	CustomerName        string `json:"customer_name"`
	Status              string `json:"status"`
	StartTime           string `json:"start_time"`
	EndTime             string `json:"end_time"`
	CreatedAt           string `json:"created_at"`
	UpdatedAt           string `json:"updated_at"`
}
