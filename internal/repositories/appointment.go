package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/DannyTuanAnh/appointment-scheduler-application/internal/db"
	"github.com/DannyTuanAnh/appointment-scheduler-application/internal/db/sqlc"
	"github.com/DannyTuanAnh/appointment-scheduler-application/internal/dto"
	"github.com/DannyTuanAnh/appointment-scheduler-application/internal/utils"
	"github.com/jackc/pgx/v5"
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

func (ar *appointmentRepository) GetServiceByID(ctx context.Context, id int32) (sqlc.Service, error) {
	row, err := ar.appointment_repo.GetServiceByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return sqlc.Service{}, utils.NewError(fmt.Sprintf("service with id '%d' is not found", id), utils.ErrCodeNotFound)
		}
		return sqlc.Service{}, utils.WrapError(err, "failed to get service detail", utils.ErrCodeInternal)
	}
	return row, nil
}

func (ar *appointmentRepository) GetSkillRequirementIDs(ctx context.Context, serviceID int32) ([]int32, error) {
	skillIDs, err := ar.appointment_repo.GetSkillRequirementIDs(ctx, serviceID)
	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && pgErr.Code == "22P02" {
			return nil, utils.NewError(fmt.Sprintf("invalid service id '%d'", serviceID), utils.ErrCodeBadRequest)
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

func (ar *appointmentRepository) ListAppointments(ctx context.Context) ([]sqlc.ListAppointmentsRow, error) {
	rows, err := ar.appointment_repo.ListAppointments(ctx)
	if err != nil {
		return nil, utils.WrapError(err, "failed to list appointments", utils.ErrCodeInternal)
	}
	return rows, nil
}

func (ar *appointmentRepository) ListAppointmentsByDealershipInTimeRange(ctx context.Context, arg sqlc.ListAppointmentsByDealershipInTimeRangeParams) ([]sqlc.ListAppointmentsByDealershipInTimeRangeRow, error) {
	rows, err := ar.appointment_repo.ListAppointmentsByDealershipInTimeRange(ctx, arg)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "22P02" {
			return nil, utils.NewError("invalid input parameters", utils.ErrCodeBadRequest)
		}
		return nil, utils.WrapError(err, "failed to list appointments by dealership in time range", utils.ErrCodeInternal)
	}
	return rows, nil
}

func (ar *appointmentRepository) ListAppointmentsByTechnicianInTimeRange(ctx context.Context, arg sqlc.ListAppointmentsByTechnicianInTimeRangeParams) ([]sqlc.ListAppointmentsByTechnicianInTimeRangeRow, error) {
	rows, err := ar.appointment_repo.ListAppointmentsByTechnicianInTimeRange(ctx, arg)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "22P02" {
			return nil, utils.NewError("invalid input parameters", utils.ErrCodeBadRequest)
		}
		return nil, utils.WrapError(err, "failed to list appointments by technician in time range", utils.ErrCodeInternal)
	}
	return rows, nil
}

func (ar *appointmentRepository) ListAppointmentsByServiceBayInTimeRange(ctx context.Context, arg sqlc.ListAppointmentsByServiceBayInTimeRangeParams) ([]sqlc.ListAppointmentsByServiceBayInTimeRangeRow, error) {
	rows, err := ar.appointment_repo.ListAppointmentsByServiceBayInTimeRange(ctx, arg)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "22P02" {
			return nil, utils.NewError("invalid input parameters", utils.ErrCodeBadRequest)
		}
		return nil, utils.WrapError(err, "failed to list appointments by service bay in time range", utils.ErrCodeInternal)
	}
	return rows, nil
}

func (ar *appointmentRepository) ListAppointmentsByServiceInTimeRange(ctx context.Context, arg sqlc.ListAppointmentsByServiceInTimeRangeParams) ([]sqlc.ListAppointmentsByServiceInTimeRangeRow, error) {
	rows, err := ar.appointment_repo.ListAppointmentsByServiceInTimeRange(ctx, arg)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "22P02" {
			return nil, utils.NewError("invalid input parameters", utils.ErrCodeBadRequest)
		}
		return nil, utils.WrapError(err, "failed to list appointments by service in time range", utils.ErrCodeInternal)
	}
	return rows, nil
}

func (ar *appointmentRepository) SearchAppointmentsByCustomerNameAndDealershipID(ctx context.Context, arg sqlc.SearchAppointmentsByCustomerNameAndDealershipIDParams) ([]sqlc.SearchAppointmentsByCustomerNameAndDealershipIDRow, error) {
	rows, err := ar.appointment_repo.SearchAppointmentsByCustomerNameAndDealershipID(ctx, arg)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "22P02" {
			return nil, utils.NewError("invalid input parameters", utils.ErrCodeBadRequest)
		}
		return nil, utils.WrapError(err, "failed to search appointments by customer name and dealership id", utils.ErrCodeInternal)
	}
	return rows, nil
}

func (ar *appointmentRepository) CreateAppointment(ctx context.Context, arg dto.CreateAppointmentParams) error {

	success := false

	for _, pair := range arg.PairKeyBayAndTech {

		tx, err := db.DBPool.BeginTx(ctx, pgx.TxOptions{
			IsoLevel:   pgx.ReadCommitted,
			AccessMode: pgx.ReadWrite,
		})
		if err != nil {
			return utils.WrapError(err, "failed to begin transaction for creating appointment", utils.ErrCodeInternal)
		}

		locked, err := tryLockPair(ctx, tx, pair.BayID, pair.TechnicianID)
		if err != nil {
			tx.Rollback(ctx)
			return utils.WrapError(err, "failed to acquire advisory lock for bay and technician pair", utils.ErrCodeInternal)
		}

		if !locked {
			tx.Rollback(ctx)
			continue
		}

		qtx := sqlc.New(tx)
		_, err = qtx.CreateAppointment(ctx, sqlc.CreateAppointmentParams{
			DealershipID: arg.DealershipID,
			CustomerName: arg.CustomerName,
			ServiceID:    arg.ServiceID,
			BayID:        pair.BayID,
			TechnicianID: pair.TechnicianID,
			StartTime:    arg.StartTime,
			EndTime:      arg.EndTime,
		})

		if err == nil {
			err = tx.Commit(ctx)
			if err != nil {
				return utils.WrapError(err, "failed to commit transaction for creating appointment", utils.ErrCodeInternal)
			}
			success = true
			break
		}

		tx.Rollback(ctx)
	}

	if !success {
		return utils.NewError("appointment time conflicts with an existing appointment for all provided bay and technician combinations", utils.ErrCodeConflict)
	}

	return nil
}

func (ar *appointmentRepository) InsertAppointment(ctx context.Context, arg sqlc.CreateAppointmentParams) error {
	_, err := ar.appointment_repo.CreateAppointment(ctx, arg)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "22P02":
				return utils.NewError("invalid input parameters", utils.ErrCodeBadRequest)
			case "23503":
				return utils.NewError("resource not found", utils.ErrCodeNotFound)
			case "23P01":
				// exclusion constraint conflict (overlapping tstzrange)
				return utils.NewError("appointment time conflicts with an existing appointment", utils.ErrCodeConflict)
			case "23505":
				return utils.NewError("conflict", utils.ErrCodeConflict)
			}
		}
		return utils.WrapError(err, "failed to create appointment", utils.ErrCodeInternal)
	}

	return nil
}

func tryLockPair(ctx context.Context, tx pgx.Tx, bayID, techID int32) (bool, error) {
	var locked bool
	err := tx.QueryRow(ctx,
		`SELECT pg_try_advisory_xact_lock($1, $2)`,
		bayID,
		techID,
	).Scan(&locked)

	return locked, err
}

func (ar *appointmentRepository) UpdateAppointmentStatusByID(ctx context.Context, arg sqlc.UpdateAppointmentStatusByIDParams) error {
	_, err := ar.appointment_repo.UpdateAppointmentStatusByID(ctx, arg)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "22P02":
				return utils.NewError("invalid input parameters", utils.ErrCodeBadRequest)
			case "23503":
				return utils.NewError("resource not found", utils.ErrCodeNotFound)
			}
		}
		if errors.Is(err, sql.ErrNoRows) {
			return utils.NewError("appointment not found", utils.ErrCodeNotFound)
		}
		return utils.WrapError(err, "failed to update appointment status", utils.ErrCodeInternal)
	}
	return nil
}

func (ar *appointmentRepository) MarkNoShowAppointmentsForDealershipInTimeRange(ctx context.Context, appointmentIds []int32) error {
	err := ar.appointment_repo.MarkNoShowAppointmentsForDealershipInTimeRange(ctx, appointmentIds)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "22P02" {
			return utils.NewError("invalid input parameters", utils.ErrCodeBadRequest)
		}
		return utils.WrapError(err, "failed to mark no-show appointments for dealership in time range", utils.ErrCodeInternal)
	}
	return nil
}
