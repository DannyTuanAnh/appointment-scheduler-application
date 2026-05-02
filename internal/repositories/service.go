package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/DannyTuanAnh/appointment-scheduler-application/internal/db/sqlc"
	"github.com/DannyTuanAnh/appointment-scheduler-application/internal/utils"
	"github.com/jackc/pgx/v5/pgconn"
)

type serviceRepository struct {
	q sqlc.Querier
}

func NewServiceRepository(q sqlc.Querier) ServiceRepository {
	return &serviceRepository{q: q}
}

// SERVICES

func (r *serviceRepository) CreateService(ctx context.Context, arg sqlc.CreateServiceParams) error {
	_, err := r.q.CreateService(ctx, arg)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "22P02" {
				return utils.NewError("invalid service payload", utils.ErrCodeBadRequest)
			}
			if pgErr.Code == "23503" {
				return utils.NewError("required_bay_type not found", utils.ErrCodeBadRequest)
			}
			if pgErr.Code == "23505" {
				return utils.NewError("service already exists", utils.ErrCodeConflict)
			}
		}
		return utils.WrapError(err, "failed to create service", utils.ErrCodeInternal)
	}
	return nil
}

func (r *serviceRepository) UpdateServiceByID(ctx context.Context, arg sqlc.UpdateServiceByIDParams) error {
	_, err := r.q.UpdateServiceByID(ctx, arg)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return utils.NewError(fmt.Sprintf("service with id '%d' is not found", arg.ID), utils.ErrCodeNotFound)
		}
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23503" {
				return utils.NewError("required_bay_type not found", utils.ErrCodeBadRequest)
			}
			if pgErr.Code == "22P02" {
				return utils.NewError("invalid service payload", utils.ErrCodeBadRequest)
			}
			if pgErr.Code == "23505" {
				return utils.NewError("service already exists", utils.ErrCodeConflict)
			}
		}
		return utils.WrapError(err, "failed to update service", utils.ErrCodeInternal)
	}
	return nil
}

func (r *serviceRepository) DeleteServiceByID(ctx context.Context, id int32) (int64, error) {
	rows, err := r.q.DeleteServiceByID(ctx, id)
	if err != nil {
		return 0, utils.WrapError(err, "failed to delete service", utils.ErrCodeInternal)
	}
	if rows == 0 {
		return 0, utils.NewError(fmt.Sprintf("service with id '%d' is not found", id), utils.ErrCodeNotFound)
	}
	return rows, nil
}

func (r *serviceRepository) ListServices(ctx context.Context) ([]sqlc.ListServicesRow, error) {
	rows, err := r.q.ListServices(ctx)
	if err != nil {
		return nil, utils.WrapError(err, "failed to list services", utils.ErrCodeInternal)
	}
	return rows, nil
}

func (r *serviceRepository) GetServiceDetailByID(ctx context.Context, id int32) (sqlc.GetServiceDetailByIDRow, error) {
	row, err := r.q.GetServiceDetailByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return sqlc.GetServiceDetailByIDRow{}, utils.NewError(fmt.Sprintf("service with id '%d' is not found", id), utils.ErrCodeNotFound)
		}
		return sqlc.GetServiceDetailByIDRow{}, utils.WrapError(err, "failed to get service detail", utils.ErrCodeInternal)
	}
	return row, nil
}

func (r *serviceRepository) SearchServicesByName(ctx context.Context, name string) ([]sqlc.SearchServicesByNameRow, error) {
	rows, err := r.q.SearchServicesByName(ctx, name)
	if err != nil {
		return nil, utils.WrapError(err, "failed to search services", utils.ErrCodeInternal)
	}
	return rows, nil
}

// SERVICE REQUIREMENTS

func (r *serviceRepository) AddSkillRequirementsToService(ctx context.Context, serviceID int32, skillIDs []int32) error {
	err := r.q.AddSkillRequirementsToService(ctx, sqlc.AddSkillRequirementsToServiceParams{ServiceID: serviceID, SkillIds: skillIDs})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23503" {
			return utils.NewError("service or skill not found", utils.ErrCodeBadRequest)
		}
		return utils.WrapError(err, "failed to add skill requirements", utils.ErrCodeInternal)
	}
	return nil
}

func (r *serviceRepository) RemoveSkillRequirementsFromService(ctx context.Context, serviceID int32, skillIDs []int32) (int64, error) {
	rows, err := r.q.RemoveSkillRequirementsFromService(ctx, sqlc.RemoveSkillRequirementsFromServiceParams{ServiceID: serviceID, SkillIds: skillIDs})
	if err != nil {
		return 0, utils.WrapError(err, "failed to remove skill requirements", utils.ErrCodeInternal)
	}
	return rows, nil
}
