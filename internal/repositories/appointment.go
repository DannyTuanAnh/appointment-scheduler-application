package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/DannyTuanAnh/appointment-scheduler-application/internal/db/sqlc"
	"github.com/DannyTuanAnh/appointment-scheduler-application/internal/utils"
	"github.com/jackc/pgx/v5/pgconn"
)

type appointmentRepository struct {
	appointment_repo sqlc.Querier
}

func NewAppointmentRepository(appointment_repo sqlc.Querier) AppointmentRepository {
	return &appointmentRepository{
		appointment_repo: appointment_repo,
	}
}

func (ar *appointmentRepository) GetAppointment(ctx context.Context, params sqlc.GetAppointmentsOfBayOrTechnicianInTimeRangeParams) ([]sqlc.GetAppointmentsOfBayOrTechnicianInTimeRangeRow, error) {
	resp, err := ar.appointment_repo.GetAppointmentsOfBayOrTechnicianInTimeRange(ctx, params)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "22P02" {
			return nil, utils.NewError("invalid input parameters", utils.ErrCodeBadRequest)
		}

		return nil, utils.WrapError(err, "failed to get appointments of bay or technician in time range", utils.ErrCodeInternal)
	}

	return resp, nil
}

func (ar *appointmentRepository) GetWorkHoursOfDealership(ctx context.Context, dealershipID int32) (time.Time, time.Time, error) {
	dealership, err := ar.appointment_repo.GetDealershipByID(ctx, dealershipID)
	if err != nil {

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "22P02" {
			return time.Time{}, time.Time{}, utils.NewError(fmt.Sprintf("invalid dealership id '%d'", dealershipID), utils.ErrCodeBadRequest)

		} else if errors.Is(err, sql.ErrNoRows) {
			return time.Time{}, time.Time{}, utils.NewError(fmt.Sprintf("dealership with id '%d' is not found", dealershipID), utils.ErrCodeNotFound)

		}

		return time.Time{}, time.Time{}, utils.WrapError(err, "failed to get work hours of dealership", utils.ErrCodeInternal)
	}

	return utils.ConvertPgTypeTimeToTime(dealership.OpenTime), utils.ConvertPgTypeTimeToTime(dealership.CloseTime), nil

}

func (ar *appointmentRepository) GetSkillRequirementIDs(ctx context.Context, erviceID int32) ([]int32, error) {
	skillIDs, err := ar.appointment_repo.GetSkillRequirementIDs(ctx, erviceID)
	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && pgErr.Code == "22P02" {
			return nil, utils.NewError(fmt.Sprintf("invalid service id '%d'", erviceID), utils.ErrCodeBadRequest)
		}

		return nil, utils.WrapError(err, "failed to get skill requirement ids by service id", utils.ErrCodeInternal)
	}

	return skillIDs, nil
}

func (ar *appointmentRepository) GetServiceBayIDsByDealershipIDAndTypeID(ctx context.Context, dealershipID, bayTypeID int32) ([]int32, error) {
	bays, err := ar.appointment_repo.ListServiceBaysByDealershipIDAndTypeID(ctx, sqlc.ListServiceBaysByDealershipIDAndTypeIDParams{
		DealershipID: dealershipID,
		BayTypeID:    bayTypeID,
	})

	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && pgErr.Code == "22P02" {
			return nil, utils.NewError(fmt.Sprintf("invalid input parameters: dealership_id '%d' or bay_type_id '%d'", dealershipID, bayTypeID), utils.ErrCodeBadRequest)
		}

		return nil, utils.WrapError(err, "failed to get service bays by dealership id and type id", utils.ErrCodeInternal)
	}

	var bay_ids []int32
	for _, bay := range bays {
		bay_ids = append(bay_ids, bay.ID)
	}

	return bay_ids, nil
}

func (ar *appointmentRepository) FindTechniciansByDealershipWithRequiredSkills(ctx context.Context, dealershipID int32, skillIDs []int32) ([]int32, error) {
	techs, err := ar.appointment_repo.FindActiveTechniciansByDealershipWithRequiredSkills(ctx, sqlc.FindActiveTechniciansByDealershipWithRequiredSkillsParams{
		DealershipID: dealershipID,
		SkillIds:     skillIDs,
	})

	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && pgErr.Code == "22P02" {
			return nil, utils.NewError(fmt.Sprintf("invalid input parameters: dealership_id '%d' or skill_ids '%v'", dealershipID, skillIDs), utils.ErrCodeBadRequest)
		}

		return nil, utils.WrapError(err, "failed to find technicians by dealership with required skills", utils.ErrCodeInternal)
	}

	return techs, nil
}
