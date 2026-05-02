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

type technicianRepository struct {
	q sqlc.Querier
}

func NewTechnicianRepository(q sqlc.Querier) TechnicianRepository {
	return &technicianRepository{q: q}
}

func (r *technicianRepository) CreateTechnician(ctx context.Context, arg sqlc.CreateTechnicianParams) error {
	_, err := r.q.CreateTechnician(ctx, arg)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "22P02" {
				return utils.NewError("invalid technician payload", utils.ErrCodeBadRequest)
			}
			if pgErr.Code == "23503" {
				return utils.NewError("dealership not found", utils.ErrCodeBadRequest)
			}
		}
		return utils.WrapError(err, "failed to create technician", utils.ErrCodeInternal)
	}
	return nil
}

func (r *technicianRepository) SetTechnicianOnLeave(ctx context.Context, id int32) (sqlc.Technician, error) {
	tech, err := r.q.SetTechnicianOnLeave(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return sqlc.Technician{}, utils.NewError(fmt.Sprintf("technician with id '%d' is not found", id), utils.ErrCodeNotFound)
		}
		return sqlc.Technician{}, utils.WrapError(err, "failed to set technician on leave", utils.ErrCodeInternal)
	}
	return tech, nil
}

func (r *technicianRepository) SetTechnicianBackToWork(ctx context.Context, id int32) (sqlc.Technician, error) {
	tech, err := r.q.SetTechnicianBackToWork(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return sqlc.Technician{}, utils.NewError(fmt.Sprintf("technician with id '%d' is not found", id), utils.ErrCodeNotFound)
		}
		return sqlc.Technician{}, utils.WrapError(err, "failed to set technician back to work", utils.ErrCodeInternal)
	}
	return tech, nil
}

func (r *technicianRepository) TransferTechnicianDealership(ctx context.Context, id int32, dealershipID int32) (sqlc.Technician, error) {
	tech, err := r.q.TransferTechnicianDealership(ctx, sqlc.TransferTechnicianDealershipParams{ID: id, DealershipID: dealershipID})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return sqlc.Technician{}, utils.NewError(fmt.Sprintf("technician with id '%d' is not found", id), utils.ErrCodeNotFound)
		}
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23503" {
			return sqlc.Technician{}, utils.NewError(fmt.Sprintf("dealership with id '%d' is not found", dealershipID), utils.ErrCodeBadRequest)
		}
		return sqlc.Technician{}, utils.WrapError(err, "failed to transfer technician dealership", utils.ErrCodeInternal)
	}
	return tech, nil
}

func (r *technicianRepository) UpdateTechnicianInfoByID(ctx context.Context, id int32, name *string, level *sqlc.TechnicianLevel) error {
	var nl sqlc.NullTechnicianLevel
	if level != nil {
		nl = sqlc.NullTechnicianLevel{TechnicianLevel: *level, Valid: true}
	}

	_, err := r.q.UpdateTechnicianInfoByID(ctx, sqlc.UpdateTechnicianInfoByIDParams{ID: id, Name: name, Level: nl})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return utils.NewError(fmt.Sprintf("technician with id '%d' is not found", id), utils.ErrCodeNotFound)
		}
		return utils.WrapError(err, "failed to update technician", utils.ErrCodeInternal)
	}
	return nil
}

func (r *technicianRepository) DeleteTechnicianIfInactiveOverOneMonth(ctx context.Context, id int32) error {
	if err := r.q.DeleteTechnicianIfInactiveOverOneMonth(ctx, id); err != nil {
		return utils.WrapError(err, "failed to delete technician", utils.ErrCodeInternal)
	}
	return nil
}

func (r *technicianRepository) DeleteTechnicianByID(ctx context.Context, id int32) error {
	rows, err := r.q.DeleteTechnicianByID(ctx, id)
	if rows == 0 {
		return utils.NewError(fmt.Sprintf("technician with id '%d' is not found", id), utils.ErrCodeNotFound)
	}

	if err != nil {
		return utils.WrapError(err, "failed to delete technician", utils.ErrCodeInternal)
	}

	return nil
}

func (r *technicianRepository) AddSkillsToTechnician(ctx context.Context, technicianID int32, skillIDs []int32) error {
	err := r.q.AddSkillsToTechnician(ctx, sqlc.AddSkillsToTechnicianParams{TechnicianID: technicianID, SkillIds: skillIDs})
	if err != nil {
		return utils.WrapError(err, "failed to add skills to technician", utils.ErrCodeInternal)
	}
	return nil
}

func (r *technicianRepository) RemoveSkillsFromTechnician(ctx context.Context, technicianID int32, skillIDs []int32) (int64, error) {
	rows, err := r.q.RemoveSkillsFromTechnician(ctx, sqlc.RemoveSkillsFromTechnicianParams{TechnicianID: technicianID, SkillIds: skillIDs})
	if err != nil {
		return 0, utils.WrapError(err, "failed to remove skills from technician", utils.ErrCodeInternal)
	}
	return rows, nil
}

func (r *technicianRepository) ListTechniciansByDealershipID(ctx context.Context, dealershipID int32) ([]sqlc.ListTechniciansByDealershipIDRow, error) {
	rows, err := r.q.ListTechniciansByDealershipID(ctx, dealershipID)
	if err != nil {
		return nil, utils.WrapError(err, "failed to list technicians", utils.ErrCodeInternal)
	}
	return rows, nil
}

func (r *technicianRepository) SearchTechniciansByName(ctx context.Context, name string) ([]sqlc.SearchTechniciansByNameRow, error) {
	rows, err := r.q.SearchTechniciansByName(ctx, name)
	if err != nil {
		return nil, utils.WrapError(err, "failed to search technicians", utils.ErrCodeInternal)
	}
	return rows, nil
}

func (r *technicianRepository) SearchTechniciansByNameAndDealershipID(ctx context.Context, dealershipID int32, name string) ([]sqlc.SearchTechniciansByNameAndDealershipIDRow, error) {
	rows, err := r.q.SearchTechniciansByNameAndDealershipID(ctx, sqlc.SearchTechniciansByNameAndDealershipIDParams{DealershipID: dealershipID, TechnicianName: name})
	if err != nil {
		return nil, utils.WrapError(err, "failed to search technicians by dealership", utils.ErrCodeInternal)
	}
	return rows, nil
}

func (r *technicianRepository) FindActiveTechniciansByDealershipWithRequiredSkills(ctx context.Context, dealershipID int32, skillIDs []int32) ([]int32, error) {
	ids, err := r.q.FindActiveTechniciansByDealershipWithRequiredSkills(ctx, sqlc.FindActiveTechniciansByDealershipWithRequiredSkillsParams{DealershipID: dealershipID, SkillIds: skillIDs})
	if err != nil {
		return nil, utils.WrapError(err, "failed to find active technicians", utils.ErrCodeInternal)
	}
	return ids, nil
}

func (r *technicianRepository) GetDetailTechnicianByID(ctx context.Context, id int32) (sqlc.GetDetailTechnicianByIDRow, error) {
	row, err := r.q.GetDetailTechnicianByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return sqlc.GetDetailTechnicianByIDRow{}, utils.NewError(fmt.Sprintf("technician with id '%d' is not found", id), utils.ErrCodeNotFound)
		}
		return sqlc.GetDetailTechnicianByIDRow{}, utils.WrapError(err, "failed to get technician detail", utils.ErrCodeInternal)
	}
	return row, nil
}

func (r *technicianRepository) GetTechnicianByID(ctx context.Context, id int32) (sqlc.GetTechnicianByIDRow, error) {
	row, err := r.q.GetTechnicianByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return sqlc.GetTechnicianByIDRow{}, utils.NewError(fmt.Sprintf("technician with id '%d' is not found", id), utils.ErrCodeNotFound)
		}
		return sqlc.GetTechnicianByIDRow{}, utils.WrapError(err, "failed to get technician", utils.ErrCodeInternal)
	}
	return row, nil
}
