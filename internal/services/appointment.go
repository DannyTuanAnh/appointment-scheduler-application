package services

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/DannyTuanAnh/appointment-scheduler-application/internal/db/sqlc"
	"github.com/DannyTuanAnh/appointment-scheduler-application/internal/dto"
	"github.com/DannyTuanAnh/appointment-scheduler-application/internal/repositories"
	"github.com/DannyTuanAnh/appointment-scheduler-application/internal/utils"
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

	startTime, err := formatTimeToQuery(req.PreferenceTime, openTime)
	if err != nil {
		return dto.AvailabilityResponse{}, fmt.Errorf("failed to format open time (service layer): %w", err)
	}

	endTime, err := formatTimeToQuery(req.PreferenceTime, closeTime)
	if err != nil {
		return dto.AvailabilityResponse{}, fmt.Errorf("failed to format close time (service layer): %w", err)
	}

	reqGet := sqlc.GetAppointmentsOfBayOrTechnicianInTimeRangeParams{
		BayTypeID:    req.BayTypeID,
		DealershipID: req.DealershipID,
		SkillIds:     ids,
		FromTime:     startTime,
		ToTime:       endTime,
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
func formatTimeToQuery(dateStr string, timeObj time.Time) (time.Time, error) {
	loc, _ := time.LoadLocation("Asia/Ho_Chi_Minh")

	fullDateTimeStr := fmt.Sprintf("%s %02d:%02d:00", dateStr, timeObj.Hour(), timeObj.Minute())

	t, err := time.ParseInLocation("2006-01-02 15:04:05", fullDateTimeStr, loc)
	if err != nil {
		return time.Time{}, err
	}

	return t, nil
}

func normalizeDBTimeToHHmm(t time.Time) time.Time {
	tt := t.In(time.FixedZone("UTC+7", 7*60*60))
	return time.Date(1, time.January, 1, tt.Hour(), tt.Minute(), 0, 0, tt.Location())
}

// APPOINTMENT MANAGEMENT

func (as *appointmentService) ListAppointments(ctx context.Context) ([]dto.AppointmentResponseHTTP, error) {
	rows, err := as.repo.ListAppointments(ctx)
	if err != nil {
		return nil, err
	}
	return mapAppointmentsListRows(rows), nil
}

func (as *appointmentService) ListAppointmentsByDealershipInTimeRange(ctx context.Context, req dto.ListAppointmentsByDealershipQuery) ([]dto.AppointmentResponseHTTP, error) {
	fromTime, toTime, err := parseRFC3339Range(req.FromTime, req.ToTime)
	if err != nil {
		return nil, err
	}
	rows, err := as.repo.ListAppointmentsByDealershipInTimeRange(ctx, sqlc.ListAppointmentsByDealershipInTimeRangeParams{
		DealershipID: req.DealershipID,
		FromTime:     fromTime,
		ToTime:       toTime,
	})
	if err != nil {
		return nil, err
	}
	return mapAppointmentsByDealershipRows(rows), nil
}

func (as *appointmentService) ListAppointmentsByTechnicianInTimeRange(ctx context.Context, req dto.ListAppointmentsByTechnicianQuery) ([]dto.AppointmentResponseHTTP, error) {
	fromTime, toTime, err := parseRFC3339Range(req.FromTime, req.ToTime)
	if err != nil {
		return nil, err
	}
	rows, err := as.repo.ListAppointmentsByTechnicianInTimeRange(ctx, sqlc.ListAppointmentsByTechnicianInTimeRangeParams{
		TechnicianID: req.TechnicianID,
		FromTime:     fromTime,
		ToTime:       toTime,
	})
	if err != nil {
		return nil, err
	}
	return mapAppointmentsByTechnicianRows(rows), nil
}

func (as *appointmentService) ListAppointmentsByServiceBayInTimeRange(ctx context.Context, req dto.ListAppointmentsByServiceBayQuery) ([]dto.AppointmentResponseHTTP, error) {
	fromTime, toTime, err := parseRFC3339Range(req.FromTime, req.ToTime)
	if err != nil {
		return nil, err
	}
	rows, err := as.repo.ListAppointmentsByServiceBayInTimeRange(ctx, sqlc.ListAppointmentsByServiceBayInTimeRangeParams{
		BayID:    req.BayID,
		FromTime: fromTime,
		ToTime:   toTime,
	})
	if err != nil {
		return nil, err
	}
	return mapAppointmentsByBayRows(rows), nil
}

func (as *appointmentService) ListAppointmentsByServiceInTimeRange(ctx context.Context, req dto.ListAppointmentsByServiceQuery) ([]dto.AppointmentResponseHTTP, error) {
	fromTime, toTime, err := parseRFC3339Range(req.FromTime, req.ToTime)
	if err != nil {
		return nil, err
	}
	rows, err := as.repo.ListAppointmentsByServiceInTimeRange(ctx, sqlc.ListAppointmentsByServiceInTimeRangeParams{
		ServiceID: req.ServiceID,
		FromTime:  fromTime,
		ToTime:    toTime,
	})
	if err != nil {
		return nil, err
	}
	return mapAppointmentsByServiceRows(rows), nil
}

func (as *appointmentService) SearchAppointments(ctx context.Context, req dto.SearchAppointmentsQuery) ([]dto.AppointmentResponseHTTP, error) {
	name := strings.TrimSpace(req.CustomerName)
	if name == "" {
		return []dto.AppointmentResponseHTTP{}, nil
	}
	rows, err := as.repo.SearchAppointmentsByCustomerNameAndDealershipID(ctx, sqlc.SearchAppointmentsByCustomerNameAndDealershipIDParams{
		DealershipID: req.DealershipID,
		CustomerName: &name,
	})
	if err != nil {
		return nil, err
	}
	return mapAppointmentsSearchRows(rows), nil
}

func (as *appointmentService) CreateAppointment(ctx context.Context, req dto.CreateAppointmentRequest) error {
	start, err := time.Parse(time.RFC3339, req.StartTime)
	if err != nil {
		return utils.NewError("invalid start_time (expected RFC3339)", utils.ErrCodeBadRequest)
	}

	end, err := time.Parse(time.RFC3339, req.EndTime)
	if err != nil {
		return utils.NewError("invalid end_time (expected RFC3339)", utils.ErrCodeBadRequest)
	}

	if !end.After(start) {
		return utils.NewError("end_time must be after start_time", utils.ErrCodeBadRequest)
	}

	openTime, closeTime, err := as.repo.GetWorkHoursOfDealership(ctx, req.DealershipID)
	if err != nil {
		return err
	}

	if toMinutes(start) < toMinutes(openTime) || toMinutes(end) > toMinutes(closeTime) {
		return utils.NewError(fmt.Sprintf("appointment time must be within dealership work hours (%s - %s)", openTime.Format("15:04"), closeTime.Format("15:04")), utils.ErrCodeBadRequest)
	}

	service, err := as.repo.GetServiceByID(ctx, req.ServiceID) // check if service exists
	if err != nil {
		return err
	}

	if toMinutes(end)-toMinutes(start) < int(service.AnticipatedMinutes) {
		return utils.NewError(fmt.Sprintf("appointment duration must be at least %d minutes based on the service's anticipated duration", service.AnticipatedMinutes), utils.ErrCodeBadRequest)
	}

	ids, err := as.repo.GetSkillRequirementIDs(ctx, req.ServiceID)
	if err != nil {
		return err
	}

	bayIDs, err := as.repo.GetServiceBayIDsByDealershipIDAndTypeID(ctx, req.DealershipID, req.BayTypeID)
	if err != nil {
		return err
	}

	techIDs, err := as.repo.FindTechniciansByDealershipWithRequiredSkills(ctx, req.DealershipID, ids)
	if err != nil {
		return err
	}

	reqGet := sqlc.GetAppointmentsOfBayOrTechnicianInTimeRangeParams{
		BayTypeID:    req.BayTypeID,
		DealershipID: req.DealershipID,
		SkillIds:     ids,
		FromTime:     start,
		ToTime:       end,
	}

	var free []dto.PairKey
	busy := make(map[dto.PairKey]bool)

	resp, err := as.repo.GetAppointment(ctx, reqGet)
	if err != nil {
		return err
	}

	log.Println("Appointments found in requested slot:", len(resp))
	log.Println("Resp: ", resp)

	for _, r := range resp {
		key := dto.PairKey{BayID: r.BayID, TechnicianID: r.TechnicianID}
		busy[key] = true
	}

	log.Println("Busy bay-tech combinations for requested slot:", busy)

	for _, bayID := range bayIDs {
		for _, techID := range techIDs {
			key := dto.PairKey{BayID: bayID, TechnicianID: techID}
			if !busy[key] {
				free = append(free, key)
			}
		}
	}

	log.Println("Free bay-tech combinations for requested slot:", free)

	return as.repo.CreateAppointment(ctx, dto.CreateAppointmentParams{
		DealershipID:      req.DealershipID,
		ServiceID:         req.ServiceID,
		PairKeyBayAndTech: free,
		CustomerName:      req.CustomerName,
		StartTime:         start,
		EndTime:           end,
	})
}

func toMinutes(t time.Time) int {
	return t.Hour()*60 + t.Minute()
}

func (as *appointmentService) UpdateAppointmentStatusByID(ctx context.Context, appointmentID int32, status string) error {
	statusTrim := strings.ToLower(strings.TrimSpace(status))
	if statusTrim == "" {
		return utils.NewError("status is required", utils.ErrCodeBadRequest)
	}

	// Keep in sync with DB enum status_type
	switch status {
	case "booked", "in_progress", "completed", "cancelled", "no_show":
		// ok
	default:
		return utils.NewError("invalid status", utils.ErrCodeBadRequest)
	}

	return as.repo.UpdateAppointmentStatusByID(ctx, sqlc.UpdateAppointmentStatusByIDParams{
		Status:        statusTrim,
		AppointmentID: appointmentID,
	})
}

func (as *appointmentService) MarkNoShowAppointmentsForDealershipInTimeRange(ctx context.Context, appointmentIds []int32) error {
	return as.repo.MarkNoShowAppointmentsForDealershipInTimeRange(ctx, appointmentIds)
}

func parseRFC3339Range(fromStr, toStr string) (time.Time, time.Time, error) {
	fromTime, err := time.Parse(time.RFC3339, fromStr)
	if err != nil {
		return time.Time{}, time.Time{}, utils.NewError("invalid from_time (expected RFC3339)", utils.ErrCodeBadRequest)
	}
	toTime, err := time.Parse(time.RFC3339, toStr)
	if err != nil {
		return time.Time{}, time.Time{}, utils.NewError("invalid to_time (expected RFC3339)", utils.ErrCodeBadRequest)
	}
	if !toTime.After(fromTime) {
		return time.Time{}, time.Time{}, utils.NewError("to_time must be after from_time", utils.ErrCodeBadRequest)
	}
	return fromTime, toTime, nil
}

func mapAppointmentsListRows(rows []sqlc.ListAppointmentsRow) []dto.AppointmentResponseHTTP {
	out := make([]dto.AppointmentResponseHTTP, len(rows))
	for i, r := range rows {
		out[i] = mapAppointmentCommon(
			r.ID,
			r.DealershipID,
			r.DealershipName,
			utils.ConvertPgTypeTimeToTime(r.DealershipOpenTime).Format("15:04"),
			utils.ConvertPgTypeTimeToTime(r.DealershipCloseTime).Format("15:04"),
			r.ServiceID,
			r.ServiceName,
			r.AnticipatedMinutes,
			r.BayID,
			strOrEmpty(r.BayName),
			r.TechnicianID,
			strOrEmpty(r.TechnicianName),
			r.CustomerName,
			fmt.Sprint(r.Status),
			fmt.Sprint(r.StartTime),
			fmt.Sprint(r.EndTime),
			r.CreatedAt,
			r.UpdatedAt,
		)
	}
	return out
}

func mapAppointmentsByDealershipRows(rows []sqlc.ListAppointmentsByDealershipInTimeRangeRow) []dto.AppointmentResponseHTTP {
	out := make([]dto.AppointmentResponseHTTP, len(rows))
	for i, r := range rows {
		out[i] = mapAppointmentCommon(
			r.ID,
			r.DealershipID,
			r.DealershipName,
			utils.ConvertPgTypeTimeToTime(r.DealershipOpenTime).Format("15:04"),
			utils.ConvertPgTypeTimeToTime(r.DealershipCloseTime).Format("15:04"),
			r.ServiceID,
			r.ServiceName,
			r.AnticipatedMinutes,
			r.BayID,
			strOrEmpty(r.BayName),
			r.TechnicianID,
			strOrEmpty(r.TechnicianName),
			r.CustomerName,
			fmt.Sprint(r.Status),
			fmt.Sprint(r.StartTime),
			fmt.Sprint(r.EndTime),
			r.CreatedAt,
			r.UpdatedAt,
		)
	}
	return out
}

func mapAppointmentsByTechnicianRows(rows []sqlc.ListAppointmentsByTechnicianInTimeRangeRow) []dto.AppointmentResponseHTTP {
	out := make([]dto.AppointmentResponseHTTP, len(rows))
	for i, r := range rows {
		out[i] = mapAppointmentCommon(
			r.ID,
			r.DealershipID,
			r.DealershipName,
			utils.ConvertPgTypeTimeToTime(r.DealershipOpenTime).Format("15:04"),
			utils.ConvertPgTypeTimeToTime(r.DealershipCloseTime).Format("15:04"),
			r.ServiceID,
			r.ServiceName,
			r.AnticipatedMinutes,
			r.BayID,
			strOrEmpty(r.BayName),
			r.TechnicianID,
			strOrEmpty(r.TechnicianName),
			r.CustomerName,
			fmt.Sprint(r.Status),
			fmt.Sprint(r.StartTime),
			fmt.Sprint(r.EndTime),
			r.CreatedAt,
			r.UpdatedAt,
		)
	}
	return out
}

func mapAppointmentsByBayRows(rows []sqlc.ListAppointmentsByServiceBayInTimeRangeRow) []dto.AppointmentResponseHTTP {
	out := make([]dto.AppointmentResponseHTTP, len(rows))
	for i, r := range rows {
		out[i] = mapAppointmentCommon(
			r.ID,
			r.DealershipID,
			r.DealershipName,
			utils.ConvertPgTypeTimeToTime(r.DealershipOpenTime).Format("15:04"),
			utils.ConvertPgTypeTimeToTime(r.DealershipCloseTime).Format("15:04"),
			r.ServiceID,
			r.ServiceName,
			r.AnticipatedMinutes,
			r.BayID,
			strOrEmpty(r.BayName),
			r.TechnicianID,
			strOrEmpty(r.TechnicianName),
			r.CustomerName,
			fmt.Sprint(r.Status),
			fmt.Sprint(r.StartTime),
			fmt.Sprint(r.EndTime),
			r.CreatedAt,
			r.UpdatedAt,
		)
	}
	return out
}

func mapAppointmentsByServiceRows(rows []sqlc.ListAppointmentsByServiceInTimeRangeRow) []dto.AppointmentResponseHTTP {
	out := make([]dto.AppointmentResponseHTTP, len(rows))
	for i, r := range rows {
		out[i] = mapAppointmentCommon(
			r.ID,
			r.DealershipID,
			r.DealershipName,
			utils.ConvertPgTypeTimeToTime(r.DealershipOpenTime).Format("15:04"),
			utils.ConvertPgTypeTimeToTime(r.DealershipCloseTime).Format("15:04"),
			r.ServiceID,
			r.ServiceName,
			r.AnticipatedMinutes,
			r.BayID,
			strOrEmpty(r.BayName),
			r.TechnicianID,
			strOrEmpty(r.TechnicianName),
			r.CustomerName,
			fmt.Sprint(r.Status),
			fmt.Sprint(r.StartTime),
			fmt.Sprint(r.EndTime),
			r.CreatedAt,
			r.UpdatedAt,
		)
	}
	return out
}

func mapAppointmentsSearchRows(rows []sqlc.SearchAppointmentsByCustomerNameAndDealershipIDRow) []dto.AppointmentResponseHTTP {
	out := make([]dto.AppointmentResponseHTTP, len(rows))
	for i, r := range rows {
		out[i] = mapAppointmentCommon(
			r.ID,
			r.DealershipID,
			r.DealershipName,
			utils.ConvertPgTypeTimeToTime(r.DealershipOpenTime).Format("15:04"),
			utils.ConvertPgTypeTimeToTime(r.DealershipCloseTime).Format("15:04"),
			r.ServiceID,
			r.ServiceName,
			r.AnticipatedMinutes,
			r.BayID,
			strOrEmpty(r.BayName),
			r.TechnicianID,
			strOrEmpty(r.TechnicianName),
			r.CustomerName,
			fmt.Sprint(r.Status),
			fmt.Sprint(r.StartTime),
			fmt.Sprint(r.EndTime),
			r.CreatedAt,
			r.UpdatedAt,
		)
	}
	return out
}

func mapAppointmentCommon(
	id int32,
	dealershipID int32,
	dealershipName string,
	openTime string,
	closeTime string,
	serviceID int32,
	serviceName string,
	anticipatedMinutes int32,
	bayID int32,
	bayName string,
	technicianID int32,
	technicianName string,
	customerName string,
	status string,
	startTime string,
	endTime string,
	createdAt time.Time,
	updatedAt time.Time,
) dto.AppointmentResponseHTTP {
	return dto.AppointmentResponseHTTP{
		ID:                  id,
		DealershipID:        dealershipID,
		DealershipName:      dealershipName,
		DealershipOpenTime:  openTime,
		DealershipCloseTime: closeTime,
		ServiceID:           serviceID,
		ServiceName:         serviceName,
		AnticipatedMinutes:  anticipatedMinutes,
		BayID:               bayID,
		BayName:             bayName,
		TechnicianID:        technicianID,
		TechnicianName:      technicianName,
		CustomerName:        customerName,
		Status:              status,
		StartTime:           startTime,
		EndTime:             endTime,
		CreatedAt:           createdAt.Format("15:04:05 02-01-2006"),
		UpdatedAt:           updatedAt.Format("15:04:05 02-01-2006"),
	}
}

func strOrEmpty(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
