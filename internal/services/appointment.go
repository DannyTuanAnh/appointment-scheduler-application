package services

import (
	"context"
	"fmt"
	"time"

	"github.com/DannyTuanAnh/appointment-scheduler-application/internal/db/sqlc"
	"github.com/DannyTuanAnh/appointment-scheduler-application/internal/dto"
	"github.com/DannyTuanAnh/appointment-scheduler-application/internal/repositories"
	"github.com/DannyTuanAnh/appointment-scheduler-application/pkg/scheduler"
)

type appointmentService struct {
	repo repositories.AppointmentRepository
}

func NewAppointmentService(repo repositories.AppointmentRepository) AppointmentService {
	return &appointmentService{
		repo: repo,
	}
}

func (as *appointmentService) GetAppointment(ctx context.Context, req dto.GetAppointmentRequest) (dto.AvailabilityResponse, error) {
	requestedDate, err := time.Parse("2006-01-02", req.PreferenceTime)
	if err != nil {
		return dto.AvailabilityResponse{}, fmt.Errorf("invalid date format for preference_time: %w", err)
	}

	openTime, closeTime, err := as.repo.GetWorkHoursOfDealership(ctx, req.DealershipID)
	if err != nil {
		return dto.AvailabilityResponse{}, err
	}

	ids, err := as.repo.GetSkillRequirementIDs(ctx, req.ServiceID)
	if err != nil {
		return dto.AvailabilityResponse{}, err
	}

	bayIDs, err := as.repo.GetServiceBayIDsByDealershipIDAndTypeID(ctx, req.DealershipID, req.BayTypeID)
	if err != nil {
		return dto.AvailabilityResponse{}, err
	}

	techIDs, err := as.repo.FindTechniciansByDealershipWithRequiredSkills(ctx, req.DealershipID, ids)
	if err != nil {
		return dto.AvailabilityResponse{}, err
	}

	reqGet := sqlc.GetAppointmentsOfBayOrTechnicianInTimeRangeParams{
		BayTypeID:    req.BayTypeID,
		DealershipID: req.DealershipID,
		SkillIds:     ids,
		FromTime:     formatTimeToQuery(req.PreferenceTime, openTime),
		ToTime:       formatTimeToQuery(req.PreferenceTime, closeTime),
	}

	resp, err := as.repo.GetAppointment(ctx, reqGet)
	if err != nil {
		return dto.AvailabilityResponse{}, err
	}

	serviceDuration := time.Duration(req.AnticipatedDuration) * time.Minute

	var busy []scheduler.BusyRecord
	for _, r := range resp {
		busy = append(busy, scheduler.BusyRecord{
			BayID:        r.BayID,
			TechnicianID: r.TechnicianID,
			Start:        normalizeDBTimeToHHmm(r.Duration.Lower.Time),
			End:          normalizeDBTimeToHHmm(r.Duration.Upper.Time),
		})
	}

	globalAvailable := BuildGlobalAvailability(bayIDs, techIDs, busy, openTime, closeTime, serviceDuration)

	result := BuildBusyIntervalsResponse(requestedDate, globalAvailable, openTime, closeTime, serviceDuration)

	return result, nil
}

func BuildGlobalAvailability(bayIDs []int32, techIDs []int32, busy []scheduler.BusyRecord, workStart, workEnd time.Time, serviceDuration time.Duration) []scheduler.Interval {
	work := scheduler.Interval{Start: workStart, End: workEnd}

	bayBusy := map[int32][]scheduler.Interval{}
	techBusy := map[int32][]scheduler.Interval{}

	for _, b := range busy {
		iv := scheduler.Interval{Start: b.Start, End: b.End}
		bayBusy[b.BayID] = append(bayBusy[b.BayID], iv)
		techBusy[b.TechnicianID] = append(techBusy[b.TechnicianID], iv)
	}

	// free per bay
	var allBayFree []scheduler.Interval
	for _, bayID := range bayIDs {
		free := scheduler.SubtractBusy(work, bayBusy[bayID])
		allBayFree = append(allBayFree, free...)
	}

	// free per tech
	var allTechFree []scheduler.Interval
	for _, techID := range techIDs {
		free := scheduler.SubtractBusy(work, techBusy[techID])
		allTechFree = append(allTechFree, free...)
	}

	// merge each side
	globalBayFree := scheduler.MergeIntervals(allBayFree)
	globalTechFree := scheduler.MergeIntervals(allTechFree)

	// intersect bay + tech
	globalAvailable := scheduler.Intersect(globalBayFree, globalTechFree)

	// filter too-short intervals
	var result []scheduler.Interval
	for _, iv := range globalAvailable {
		if iv.End.Sub(iv.Start) >= serviceDuration {
			result = append(result, iv)
		}
	}

	return result
}

func BuildBusyIntervalsResponse(date time.Time, globalAvailable []scheduler.Interval, workStart time.Time, workEnd time.Time, serviceDuration time.Duration) dto.AvailabilityResponse {
	work := scheduler.Interval{Start: workStart, End: workEnd}
	busy := scheduler.SubtractBusy(work, globalAvailable)

	resp := dto.AvailabilityResponse{
		Date:            date.Format("02-01-2006"),
		WorkStart:       workStart.Format("15:04"),
		WorkEnd:         workEnd.Format("15:04"),
		DurationMinutes: int(serviceDuration.Minutes()),
		Busy:            make([]dto.BusyIntervalResponse, 0, len(busy)),
	}

	for _, b := range busy {
		resp.Busy = append(resp.Busy, dto.BusyIntervalResponse{
			Start: b.Start.Format("15:04"),
			End:   b.End.Format("15:04"),
		})
	}

	return resp
}

// true format is "YYYY-MM-DD HH:mm:ss+07"
func formatTimeToQuery(dateStr string, timeStr time.Time) string {
	return fmt.Sprintf("%s %02d:%02d:00+07", dateStr, timeStr.Hour(), timeStr.Minute())
}

func normalizeDBTimeToHHmm(t time.Time) time.Time {
	tt := t.In(time.FixedZone("UTC+7", 7*60*60))
	return time.Date(1, time.January, 1, tt.Hour(), tt.Minute(), 0, 0, tt.Location())
}
