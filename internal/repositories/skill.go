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

type skillRepository struct {
	q sqlc.Querier
}

func NewSkillRepository(q sqlc.Querier) SkillRepository {
	return &skillRepository{q: q}
}

func (sr *skillRepository) CreateSkill(ctx context.Context, name string) (sqlc.Skill, error) {
	skill, err := sr.q.CreateSkill(ctx, name)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return sqlc.Skill{}, utils.NewError(fmt.Sprintf("skill with name '%s' already exists", name), utils.ErrCodeConflict)
		}
		return sqlc.Skill{}, utils.WrapError(err, "failed to create skill", utils.ErrCodeInternal)
	}
	return skill, nil
}

func (sr *skillRepository) ListSkills(ctx context.Context) ([]sqlc.Skill, error) {
	skills, err := sr.q.ListSkills(ctx)
	if err != nil {
		return nil, utils.WrapError(err, "failed to list skills", utils.ErrCodeInternal)
	}
	return skills, nil
}

func (sr *skillRepository) GetSkillByID(ctx context.Context, id int32) (sqlc.Skill, error) {
	skill, err := sr.q.GetSkillByID(ctx, id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "22P02" {
			return sqlc.Skill{}, utils.NewError(fmt.Sprintf("invalid skill id '%d'", id), utils.ErrCodeBadRequest)
		}
		if errors.Is(err, sql.ErrNoRows) {
			return sqlc.Skill{}, utils.NewError(fmt.Sprintf("skill with id '%d' is not found", id), utils.ErrCodeNotFound)
		}
		return sqlc.Skill{}, utils.WrapError(err, "failed to get skill by id", utils.ErrCodeInternal)
	}
	return skill, nil
}

func (sr *skillRepository) SearchSkillsByName(ctx context.Context, name string) ([]sqlc.Skill, error) {
	nameArg := name
	skills, err := sr.q.SearchSkillsByName(ctx, &nameArg)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []sqlc.Skill{}, nil
		}
		return nil, utils.WrapError(err, "failed to search skills by name", utils.ErrCodeInternal)
	}
	return skills, nil
}

func (sr *skillRepository) UpdateSkillNameByID(ctx context.Context, id int32, name string) (sqlc.Skill, error) {
	skill, err := sr.q.UpdateSkillNameByID(ctx, sqlc.UpdateSkillNameByIDParams{ID: id, Name: &name})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return sqlc.Skill{}, utils.NewError(fmt.Sprintf("skill with name '%s' already exists", name), utils.ErrCodeConflict)
		}
		if errors.Is(err, sql.ErrNoRows) {
			return sqlc.Skill{}, utils.NewError(fmt.Sprintf("skill with id '%d' is not found", id), utils.ErrCodeNotFound)
		}
		return sqlc.Skill{}, utils.WrapError(err, "failed to update skill", utils.ErrCodeInternal)
	}
	return skill, nil
}

func (sr *skillRepository) DeleteSkillByID(ctx context.Context, id int32) error {
	rows, err := sr.q.DeleteSkillByID(ctx, id)
	if rows == 0 {
		return utils.NewError(fmt.Sprintf("skill with id '%d' is not found", id), utils.ErrCodeNotFound)
	}
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "22P02" {
			return utils.NewError(fmt.Sprintf("invalid skill id '%d'", id), utils.ErrCodeBadRequest)
		}
		return utils.WrapError(err, "failed to delete skill", utils.ErrCodeInternal)
	}
	return nil
}
