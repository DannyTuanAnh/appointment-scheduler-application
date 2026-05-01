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

type dealershipRepository struct {
	dealership_repo sqlc.Querier
}

func NewDealershipRepository(dealership_repo sqlc.Querier) DealershipRepository {
	return &dealershipRepository{
		dealership_repo: dealership_repo,
	}
}

func (dr *dealershipRepository) GetAllDealerships(ctx context.Context) ([]sqlc.Dealership, error) {
	dealerships, err := dr.dealership_repo.ListDealerships(ctx)
	if err != nil {
		return []sqlc.Dealership{}, utils.WrapError(err, "failed to get all dealerships", utils.ErrCodeInternal)
	}

	return dealerships, nil
}

func (dr *dealershipRepository) GetDealershipByID(ctx context.Context, dealership_id int32) (sqlc.Dealership, error) {
	dealership, err := dr.dealership_repo.GetDealershipByID(ctx, dealership_id)
	if err != nil {

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "22P02" {
			return sqlc.Dealership{}, utils.NewError(fmt.Sprintf("invalid dealership id '%d'", dealership_id), utils.ErrCodeBadRequest)

		} else if errors.Is(err, sql.ErrNoRows) {
			return sqlc.Dealership{}, utils.NewError(fmt.Sprintf("dealership with id '%d' is not found", dealership_id), utils.ErrCodeNotFound)

		}

		return sqlc.Dealership{}, utils.WrapError(err, "failed to get dealership by id", utils.ErrCodeInternal)
	}

	return dealership, nil
}

func (dr *dealershipRepository) UpdateDealershipByID(ctx context.Context, id int32, name string, openTime, closeTime time.Time) error {
	_, err := dr.dealership_repo.UpdateDealershipByID(ctx, sqlc.UpdateDealershipByIDParams{
		ID:        id,
		Name:      name,
		OpenTime:  utils.ConvertTimeToPgTypeTime(openTime),
		CloseTime: utils.ConvertTimeToPgTypeTime(closeTime),
	})

	if err != nil {

		var pgErr *pgconn.PgError
		if errors.Is(err, sql.ErrNoRows) {
			return utils.NewError(fmt.Sprintf("dealership with id '%d' is not found", id), utils.ErrCodeNotFound)

		} else if errors.As(err, &pgErr) && pgErr.Code == "22P02" {
			return utils.NewError(fmt.Sprintf("Postgres Error: %s. Detail: %s", pgErr.Message, pgErr.Detail), utils.ErrCodeBadRequest)

		} else if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return utils.NewError(fmt.Sprintf("dealership with name '%s' already exists", name), utils.ErrCodeConflict)
		}

		return utils.WrapError(err, "failed to update dealership", utils.ErrCodeInternal)
	}

	return nil
}

func (dr *dealershipRepository) CreateDealership(ctx context.Context, name string, openTime, closeTime time.Time) error {
	_, err := dr.dealership_repo.CreateDealership(ctx, sqlc.CreateDealershipParams{
		Name:      name,
		OpenTime:  utils.ConvertTimeToPgTypeTime(openTime),
		CloseTime: utils.ConvertTimeToPgTypeTime(closeTime),
	})

	if err != nil {

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return utils.NewError(fmt.Sprintf("dealership with name '%s' already exists", name), utils.ErrCodeConflict)

		} else if errors.As(err, &pgErr) && pgErr.Code == "22P02" {
			return utils.NewError(fmt.Sprintf("Postgres Error: %s. Detail: %s", pgErr.Message, pgErr.Detail), utils.ErrCodeBadRequest)

		}

		return utils.WrapError(err, "failed to create dealership", utils.ErrCodeInternal)
	}

	return nil
}

func (dr *dealershipRepository) DeleteDealershipByID(ctx context.Context, id int32) error {
	rowAffected, err := dr.dealership_repo.DeleteDealershipByID(ctx, id)

	if rowAffected == 0 {
		return utils.NewError(fmt.Sprintf("dealership with id '%d' is not found", id), utils.ErrCodeNotFound)
	}

	if err != nil {

		var pgErr *pgconn.PgError
		if errors.Is(err, sql.ErrNoRows) {
			return utils.NewError(fmt.Sprintf("dealership with id '%d' is not found", id), utils.ErrCodeNotFound)

		} else if errors.As(err, &pgErr) && pgErr.Code == "22P02" {
			return utils.NewError(fmt.Sprintf("invalid dealership id '%d'", id), utils.ErrCodeBadRequest)
		}

		return utils.WrapError(err, "failed to delete dealership", utils.ErrCodeInternal)
	}

	return nil
}

func (dr *dealershipRepository) SearchDealershipsByName(ctx context.Context, name string) ([]sqlc.Dealership, error) {
	nameArg := name
	dealerships, err := dr.dealership_repo.SearchDealershipsByName(ctx, &nameArg)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []sqlc.Dealership{}, nil
		}

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "22P02" {
			return nil, utils.NewError(fmt.Sprintf("Postgres Error: %s. Detail: %s", pgErr.Message, pgErr.Detail), utils.ErrCodeBadRequest)
		}
		return nil, utils.WrapError(err, "failed to search dealerships by name", utils.ErrCodeInternal)
	}

	return dealerships, nil
}
