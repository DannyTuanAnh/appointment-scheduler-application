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

type serviceBayRepository struct {
	q sqlc.Querier
}

func NewServiceBayRepository(q sqlc.Querier) ServiceBayRepository {
	return &serviceBayRepository{q: q}
}

// SERVICE BAYS

func (r *serviceBayRepository) ListServiceBays(ctx context.Context) ([]sqlc.ListServiceBaysRow, error) {
	bays, err := r.q.ListServiceBays(ctx)

	if err != nil {
		return nil, utils.WrapError(err, "failed to list service bays", utils.ErrCodeInternal)
	}

	return bays, nil
}

func (r *serviceBayRepository) ListServiceBaysByTypeID(ctx context.Context, bayTypeID int32) ([]sqlc.ListServiceBaysByTypeIDRow, error) {
	bays, err := r.q.ListServiceBaysByTypeID(ctx, bayTypeID)

	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && pgErr.Code == "22P02" {
			return []sqlc.ListServiceBaysByTypeIDRow{}, utils.NewError(fmt.Sprintf("invalid bay type id '%d'", bayTypeID), utils.ErrCodeBadRequest)
		}

		return nil, utils.WrapError(err, "failed to list service bays by type id", utils.ErrCodeInternal)
	}

	return bays, nil
}

func (r *serviceBayRepository) ListServiceBaysByDealershipID(ctx context.Context, dealershipID int32) ([]sqlc.ListServiceBaysByDealershipIDRow, error) {
	bays, err := r.q.ListServiceBaysByDealershipID(ctx, dealershipID)

	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && pgErr.Code == "22P02" {
			return []sqlc.ListServiceBaysByDealershipIDRow{}, utils.NewError(fmt.Sprintf("invalid dealership id '%d'", dealershipID), utils.ErrCodeBadRequest)
		}

		return nil, utils.WrapError(err, "failed to list service bays by dealership id", utils.ErrCodeInternal)
	}

	return bays, nil
}

func (r *serviceBayRepository) ListServiceBaysByDealershipIDAndTypeID(ctx context.Context, dealershipID int32, bayTypeID int32) ([]sqlc.ListServiceBaysByDealershipIDAndTypeIDRow, error) {
	bays, err := r.q.ListServiceBaysByDealershipIDAndTypeID(ctx, sqlc.ListServiceBaysByDealershipIDAndTypeIDParams{
		DealershipID: dealershipID,
		BayTypeID:    bayTypeID,
	})

	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && pgErr.Code == "22P02" {
			return []sqlc.ListServiceBaysByDealershipIDAndTypeIDRow{}, utils.NewError(fmt.Sprintf("invalid dealership id '%d' or bay type id '%d'", dealershipID, bayTypeID), utils.ErrCodeBadRequest)
		}

		return nil, utils.WrapError(err, "failed to list service bays by dealership id and type id", utils.ErrCodeInternal)
	}

	return bays, nil
}

func (r *serviceBayRepository) SearchServiceBaysByName(ctx context.Context, name string) ([]sqlc.SearchServiceBaysByNameRow, error) {
	bays, err := r.q.SearchServiceBaysByName(ctx, name)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "22P02" {
			return []sqlc.SearchServiceBaysByNameRow{}, utils.NewError(fmt.Sprintf("invalid name format '%s'", name), utils.ErrCodeBadRequest)
		}

		return nil, utils.WrapError(err, "failed to search service bays by name", utils.ErrCodeInternal)
	}

	return bays, nil
}

func (r *serviceBayRepository) SearchServiceBaysByNameAndDealershipID(ctx context.Context, name string, dealershipID int32) ([]sqlc.SearchServiceBaysByNameAndDealershipIDRow, error) {
	bays, err := r.q.SearchServiceBaysByNameAndDealershipID(ctx, sqlc.SearchServiceBaysByNameAndDealershipIDParams{
		DealershipID:   dealershipID,
		ServiceBayName: name,
	})

	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && pgErr.Code == "22P02" {
			return []sqlc.SearchServiceBaysByNameAndDealershipIDRow{}, utils.NewError(fmt.Sprintf("invalid dealership id '%d' or name format '%s'", dealershipID, name), utils.ErrCodeBadRequest)
		}

		return nil, utils.WrapError(err, "failed to search service bays by name and dealership id", utils.ErrCodeInternal)
	}

	return bays, nil
}

func (r *serviceBayRepository) SearchServiceBaysByNameAndTypeID(ctx context.Context, name string, bayTypeID int32) ([]sqlc.SearchServiceBaysByNameAndTypeIDRow, error) {
	bays, err := r.q.SearchServiceBaysByNameAndTypeID(ctx, sqlc.SearchServiceBaysByNameAndTypeIDParams{
		BayTypeID:      bayTypeID,
		ServiceBayName: name,
	})

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "22P02" {
			return []sqlc.SearchServiceBaysByNameAndTypeIDRow{}, utils.NewError(fmt.Sprintf("invalid bay type id '%d' or name format '%s'", bayTypeID, name), utils.ErrCodeBadRequest)
		}

		return nil, utils.WrapError(err, "failed to search service bays by name and type id", utils.ErrCodeInternal)
	}

	return bays, nil
}

func (r *serviceBayRepository) SearchServiceBaysByNameDealershipIDAndTypeID(ctx context.Context, name string, dealershipID int32, bayTypeID int32) ([]sqlc.SearchServiceBaysByNameDealershipIDAndTypeIDRow, error) {
	bays, err := r.q.SearchServiceBaysByNameDealershipIDAndTypeID(ctx, sqlc.SearchServiceBaysByNameDealershipIDAndTypeIDParams{
		DealershipID:   dealershipID,
		BayTypeID:      bayTypeID,
		ServiceBayName: name,
	})

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "22P02" {
			return []sqlc.SearchServiceBaysByNameDealershipIDAndTypeIDRow{}, utils.NewError(fmt.Sprintf("invalid dealership id '%d' or bay type id '%d' or name format '%s'", dealershipID, bayTypeID, name), utils.ErrCodeBadRequest)
		}

		return nil, utils.WrapError(err, "failed to search service bays by name, dealership id and type id", utils.ErrCodeInternal)
	}

	return bays, nil
}

func (r *serviceBayRepository) GetServiceBayByID(ctx context.Context, id int32) (sqlc.GetServiceBayByIDRow, error) {
	bay, err := r.q.GetServiceBayByID(ctx, id)
	if err != nil {
		var pgErr *pgconn.PgError

		if errors.Is(err, sql.ErrNoRows) {
			return sqlc.GetServiceBayByIDRow{}, utils.NewError(fmt.Sprintf("service bay with id '%d' is not found", id), utils.ErrCodeNotFound)

		} else if errors.As(err, &pgErr) && pgErr.Code == "22P02" {
			return sqlc.GetServiceBayByIDRow{}, utils.NewError(fmt.Sprintf("invalid dealership id '%d'", id), utils.ErrCodeBadRequest)

		}

		return sqlc.GetServiceBayByIDRow{}, utils.WrapError(err, "failed to get service bay by id", utils.ErrCodeInternal)
	}
	return bay, nil
}

func (r *serviceBayRepository) CreateServiceBay(ctx context.Context, arg sqlc.CreateServiceBayParams) error {
	_, err := r.q.CreateServiceBay(ctx, arg)
	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return utils.NewError(fmt.Sprintf("service bay with name '%s' already exists in that dealership", arg.Name), utils.ErrCodeConflict)

		} else if errors.As(err, &pgErr) && pgErr.Code == "22P02" {
			return utils.NewError(fmt.Sprintf("invalid dealership id '%d' or bay type id '%d' or name format '%s'", arg.DealershipID, arg.BayTypeID, arg.Name), utils.ErrCodeBadRequest)
		}

		return utils.WrapError(err, "failed to create service bay", utils.ErrCodeInternal)
	}

	return nil
}

func (r *serviceBayRepository) UpdateServiceBayByID(ctx context.Context, arg sqlc.UpdateServiceBayByIDParams) error {
	_, err := r.q.UpdateServiceBayByID(ctx, arg)
	if err != nil {
		var pgErr *pgconn.PgError

		if errors.Is(err, sql.ErrNoRows) {
			return utils.NewError(fmt.Sprintf("service bay with id '%d' is not found", arg.ID), utils.ErrCodeNotFound)

		} else if errors.As(err, &pgErr) && pgErr.Code == "22P02" {
			return utils.NewError(fmt.Sprintf("invalid service bay id '%d' or dealership id '%d' or bay type id '%d' or name format '%s'", arg.ID, arg.DealershipID, arg.BayTypeID, *arg.Name), utils.ErrCodeBadRequest)

		} else if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return utils.NewError(fmt.Sprintf("service bay with name '%s' already exists in that dealership", *arg.Name), utils.ErrCodeConflict)
		}

		return utils.WrapError(err, "failed to update service bay", utils.ErrCodeInternal)
	}

	return nil
}

func (r *serviceBayRepository) DeleteServiceBayByID(ctx context.Context, id int32) error {
	rows, err := r.q.DeleteServiceBayByID(ctx, id)
	if rows == 0 {
		return utils.NewError(fmt.Sprintf("service bay with id '%d' is not found", id), utils.ErrCodeNotFound)
	}
	if err != nil {
		return utils.WrapError(err, "failed to delete service bay", utils.ErrCodeInternal)
	}
	return nil
}

// SERVICE BAY TYPES

func (r *serviceBayRepository) CreateServiceBayType(ctx context.Context, name string) (sqlc.ServiceBayType, error) {
	t, err := r.q.CreateServiceBayType(ctx, name)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return sqlc.ServiceBayType{}, utils.NewError(fmt.Sprintf("service bay type with name '%s' already exists", name), utils.ErrCodeConflict)
		}
		return sqlc.ServiceBayType{}, utils.WrapError(err, "failed to create service bay type", utils.ErrCodeInternal)
	}
	return t, nil
}

func (r *serviceBayRepository) ListServiceBayTypes(ctx context.Context) ([]sqlc.ServiceBayType, error) {
	types, err := r.q.ListServiceBayTypes(ctx)
	if err != nil {
		return nil, utils.WrapError(err, "failed to list service bay types", utils.ErrCodeInternal)
	}
	return types, nil
}

func (r *serviceBayRepository) GetServiceBayTypeByID(ctx context.Context, id int32) (sqlc.ServiceBayType, error) {
	t, err := r.q.GetServiceBayTypeByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return sqlc.ServiceBayType{}, utils.NewError(fmt.Sprintf("service bay type with id '%d' is not found", id), utils.ErrCodeNotFound)
		}
		return sqlc.ServiceBayType{}, utils.WrapError(err, "failed to get service bay type by id", utils.ErrCodeInternal)
	}
	return t, nil
}

func (r *serviceBayRepository) SearchServiceBayTypesByName(ctx context.Context, name string) ([]sqlc.ServiceBayType, error) {
	nameArg := name
	types, err := r.q.SearchServiceBayTypesByName(ctx, &nameArg)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []sqlc.ServiceBayType{}, nil
		}
		return nil, utils.WrapError(err, "failed to search service bay types by name", utils.ErrCodeInternal)
	}
	return types, nil
}

func (r *serviceBayRepository) UpdateServiceBayTypeByID(ctx context.Context, id int32, name string) (sqlc.ServiceBayType, error) {
	updated, err := r.q.UpdateServiceBayTypeByID(ctx, sqlc.UpdateServiceBayTypeByIDParams{ID: id, Name: &name})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return sqlc.ServiceBayType{}, utils.NewError(fmt.Sprintf("service bay type with id '%d' is not found", id), utils.ErrCodeNotFound)
		}
		return sqlc.ServiceBayType{}, utils.WrapError(err, "failed to update service bay type", utils.ErrCodeInternal)
	}
	return updated, nil
}

func (r *serviceBayRepository) DeleteServiceBayTypeByID(ctx context.Context, id int32) error {
	rows, err := r.q.DeleteServiceBayTypeByID(ctx, id)
	if rows == 0 {
		return utils.NewError(fmt.Sprintf("service bay type with id '%d' is not found", id), utils.ErrCodeNotFound)
	}
	if err != nil {
		return utils.WrapError(err, "failed to delete service bay type", utils.ErrCodeInternal)
	}
	return nil
}
