package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/DannyTuanAnh/appointment-scheduler-application/internal/db/sqlc"
	"github.com/DannyTuanAnh/appointment-scheduler-application/internal/dto"
	"github.com/DannyTuanAnh/appointment-scheduler-application/internal/repositories"
	"github.com/DannyTuanAnh/appointment-scheduler-application/internal/utils"
	"github.com/jackc/pgx/v5/pgtype"
)

type technicianService struct {
	repo repositories.TechnicianRepository
}

func NewTechnicianService(repo repositories.TechnicianRepository) TechnicianService {
	return &technicianService{repo: repo}
}

func (s *technicianService) CreateTechnician(ctx context.Context, req dto.CreateTechnicianRequest) error {
	lvl, err := parseTechnicianLevel(req.Level)
	if err != nil {
		return err
	}

	err = s.repo.CreateTechnician(ctx, sqlc.CreateTechnicianParams{
		DealershipID: req.DealershipID,
		Name:         req.Name,
		Level:        lvl,
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *technicianService) SetTechnicianOnLeave(ctx context.Context, id int32) (dto.TechnicianResponseHTTP, error) {
	tech, err := s.repo.SetTechnicianOnLeave(ctx, id)
	if err != nil {
		return dto.TechnicianResponseHTTP{}, err
	}
	return dto.TechnicianResponseHTTP{
		ID:             tech.ID,
		DealershipID:   tech.DealershipID,
		TechnicianName: tech.Name,
		Level:          string(tech.Level),
		IsActive:       tech.IsActive,
		InactiveSince:  formatPgTimestamptzPtr(tech.InactiveSince),
		CreatedAt:      tech.CreatedAt.Format("15:04:05 02-01-2006"),
		UpdatedAt:      tech.UpdatedAt.Format("15:04:05 02-01-2006"),
	}, nil
}

func (s *technicianService) SetTechnicianBackToWork(ctx context.Context, id int32) (dto.TechnicianResponseHTTP, error) {
	tech, err := s.repo.SetTechnicianBackToWork(ctx, id)
	if err != nil {
		return dto.TechnicianResponseHTTP{}, err
	}
	return dto.TechnicianResponseHTTP{
		ID:             tech.ID,
		DealershipID:   tech.DealershipID,
		TechnicianName: tech.Name,
		Level:          string(tech.Level),
		IsActive:       tech.IsActive,
		InactiveSince:  formatPgTimestamptzPtr(tech.InactiveSince),
		CreatedAt:      tech.CreatedAt.Format("15:04:05 02-01-2006"),
		UpdatedAt:      tech.UpdatedAt.Format("15:04:05 02-01-2006"),
	}, nil
}

func (s *technicianService) TransferTechnicianDealership(ctx context.Context, id int32, dealershipID int32) (dto.TechnicianResponseHTTP, error) {
	tech, err := s.repo.TransferTechnicianDealership(ctx, id, dealershipID)
	if err != nil {
		return dto.TechnicianResponseHTTP{}, err
	}
	return dto.TechnicianResponseHTTP{
		ID:             tech.ID,
		DealershipID:   tech.DealershipID,
		TechnicianName: tech.Name,
		Level:          string(tech.Level),
		IsActive:       tech.IsActive,
		InactiveSince:  formatPgTimestamptzPtr(tech.InactiveSince),
		CreatedAt:      tech.CreatedAt.Format("15:04:05 02-01-2006"),
		UpdatedAt:      tech.UpdatedAt.Format("15:04:05 02-01-2006"),
	}, nil
}

func (s *technicianService) UpdateTechnicianInfoByID(ctx context.Context, id int32, req dto.UpdateTechnicianRequest) error {
	if req.Name == nil && req.Level == nil {
		return utils.NewError("nothing to update", utils.ErrCodeBadRequest)
	}

	var lvl *sqlc.TechnicianLevel
	if req.Level != nil {
		parsed, err := parseTechnicianLevel(*req.Level)
		if err != nil {
			return err
		}
		lvl = &parsed
	}

	err := s.repo.UpdateTechnicianInfoByID(ctx, id, req.Name, lvl)
	if err != nil {
		return err
	}

	return nil
}

func (s *technicianService) DeleteTechnicianIfInactiveOverOneMonth(ctx context.Context, id int32) error {
	return s.repo.DeleteTechnicianIfInactiveOverOneMonth(ctx, id)
}

func (s *technicianService) DeleteTechnicianByID(ctx context.Context, id int32) error {
	return s.repo.DeleteTechnicianByID(ctx, id)
}

func (s *technicianService) AddSkillsToTechnician(ctx context.Context, technicianID int32, skillIDs []int32) error {
	return s.repo.AddSkillsToTechnician(ctx, technicianID, skillIDs)
}

func (s *technicianService) RemoveSkillsFromTechnician(ctx context.Context, technicianID int32, skillIDs []int32) (int64, error) {
	return s.repo.RemoveSkillsFromTechnician(ctx, technicianID, skillIDs)
}

func (s *technicianService) ListTechniciansByDealershipID(ctx context.Context, dealershipID int32) ([]dto.TechnicianResponseHTTP, error) {
	rows, err := s.repo.ListTechniciansByDealershipID(ctx, dealershipID)
	if err != nil {
		return nil, err
	}

	res := make([]dto.TechnicianResponseHTTP, len(rows))
	for i, r := range rows {
		res[i] = mapTechnicianListRow(r)
	}
	return res, nil
}

func (s *technicianService) SearchTechniciansByName(ctx context.Context, name string) ([]dto.TechnicianResponseHTTP, error) {
	rows, err := s.repo.SearchTechniciansByName(ctx, name)
	if err != nil {
		return nil, err
	}

	res := make([]dto.TechnicianResponseHTTP, len(rows))
	for i, r := range rows {
		res[i] = mapTechnicianSearchRow(r)
	}
	return res, nil
}

func (s *technicianService) SearchTechniciansByNameAndDealershipID(ctx context.Context, dealershipID int32, name string) ([]dto.TechnicianResponseHTTP, error) {
	rows, err := s.repo.SearchTechniciansByNameAndDealershipID(ctx, dealershipID, name)
	if err != nil {
		return nil, err
	}

	res := make([]dto.TechnicianResponseHTTP, len(rows))
	for i, r := range rows {
		res[i] = mapTechnicianSearchByDealershipRow(r)
	}
	return res, nil
}

func (s *technicianService) FindActiveTechniciansByDealershipWithRequiredSkills(ctx context.Context, dealershipID int32, skillIDs []int32) ([]int32, error) {
	return s.repo.FindActiveTechniciansByDealershipWithRequiredSkills(ctx, dealershipID, skillIDs)
}

func (s *technicianService) GetDetailTechnicianByID(ctx context.Context, id int32) (dto.TechnicianDetailResponseHTTP, error) {
	row, err := s.repo.GetDetailTechnicianByID(ctx, id)
	if err != nil {
		return dto.TechnicianDetailResponseHTTP{}, err
	}

	skills, err := decodeSkills(row.Skills)
	if err != nil {
		return dto.TechnicianDetailResponseHTTP{}, err
	}

	return dto.TechnicianDetailResponseHTTP{
		TechnicianResponseHTTP: dto.TechnicianResponseHTTP{
			ID:             row.ID,
			DealershipID:   row.DealershipID,
			DealershipName: row.DealershipName,
			TechnicianName: row.TechnicianName,
			Level:          string(row.Level),
			IsActive:       row.IsActive,
			InactiveSince:  formatPgTimestamptzPtr(row.InactiveSince),
			CreatedAt:      row.CreatedAt.Format("15:04:05 02-01-2006"),
			UpdatedAt:      row.UpdatedAt.Format("15:04:05 02-01-2006"),
		},
		Skills: skills,
	}, nil
}

func (s *technicianService) GetTechnicianByID(ctx context.Context, id int32) (dto.TechnicianResponseHTTP, error) {
	row, err := s.repo.GetTechnicianByID(ctx, id)
	if err != nil {
		return dto.TechnicianResponseHTTP{}, err
	}

	return dto.TechnicianResponseHTTP{
		ID:             row.ID,
		DealershipID:   row.DealershipID,
		DealershipName: row.DealershipName,
		TechnicianName: row.TechnicianName,
		Level:          string(row.Level),
		IsActive:       row.IsActive,
		InactiveSince:  formatPgTimestamptzPtr(row.InactiveSince),
		CreatedAt:      row.CreatedAt.Format("15:04:05 02-01-2006"),
		UpdatedAt:      row.UpdatedAt.Format("15:04:05 02-01-2006"),
	}, nil
}

func parseTechnicianLevel(level string) (sqlc.TechnicianLevel, error) {
	switch level {
	case string(sqlc.TechnicianLevelFresher), string(sqlc.TechnicianLevelJunior), string(sqlc.TechnicianLevelSenior):
		return sqlc.TechnicianLevel(level), nil
	default:
		return "", utils.NewError(fmt.Sprintf("invalid technician level '%s'. Must be one of: '%s', '%s', '%s'", level, sqlc.TechnicianLevelFresher, sqlc.TechnicianLevelJunior, sqlc.TechnicianLevelSenior), utils.ErrCodeBadRequest)
	}
}

func formatPgTimestamptzPtr(ts pgtype.Timestamptz) *string {
	if !ts.Valid {
		return nil
	}
	str := ts.Time.Format(time.RFC3339)
	return &str
}

func decodeSkills(raw interface{}) ([]string, error) {
	if raw == nil {
		return []string{}, nil
	}

	// most drivers return []byte for jsonb
	switch v := raw.(type) {
	case []byte:
		var out []string
		if err := json.Unmarshal(v, &out); err != nil {
			return nil, utils.WrapError(err, "failed to parse skills", utils.ErrCodeInternal)
		}
		return out, nil
	case string:
		var out []string
		if err := json.Unmarshal([]byte(v), &out); err != nil {
			return nil, utils.WrapError(err, "failed to parse skills", utils.ErrCodeInternal)
		}
		return out, nil
	default:
		b, err := json.Marshal(v)
		if err != nil {
			return nil, utils.WrapError(err, "failed to parse skills", utils.ErrCodeInternal)
		}
		var out []string
		if err := json.Unmarshal(b, &out); err != nil {
			return nil, utils.WrapError(err, "failed to parse skills", utils.ErrCodeInternal)
		}
		return out, nil
	}
}

func mapTechnicianListRow(r sqlc.ListTechniciansByDealershipIDRow) dto.TechnicianResponseHTTP {
	return dto.TechnicianResponseHTTP{
		ID:             r.ID,
		DealershipID:   r.DealershipID,
		DealershipName: r.DealershipName,
		TechnicianName: r.TechnicianName,
		Level:          string(r.Level),
		IsActive:       r.IsActive,
		InactiveSince:  formatPgTimestamptzPtr(r.InactiveSince),
		CreatedAt:      r.CreatedAt.Format("15:04:05 02-01-2006"),
		UpdatedAt:      r.UpdatedAt.Format("15:04:05 02-01-2006"),
	}
}

func mapTechnicianSearchRow(r sqlc.SearchTechniciansByNameRow) dto.TechnicianResponseHTTP {
	return dto.TechnicianResponseHTTP{
		ID:             r.ID,
		DealershipID:   r.DealershipID,
		DealershipName: r.DealershipName,
		TechnicianName: r.TechnicianName,
		Level:          string(r.Level),
		IsActive:       r.IsActive,
		InactiveSince:  formatPgTimestamptzPtr(r.InactiveSince),
		CreatedAt:      r.CreatedAt.Format("15:04:05 02-01-2006"),
		UpdatedAt:      r.UpdatedAt.Format("15:04:05 02-01-2006"),
	}
}

func mapTechnicianSearchByDealershipRow(r sqlc.SearchTechniciansByNameAndDealershipIDRow) dto.TechnicianResponseHTTP {
	return dto.TechnicianResponseHTTP{
		ID:             r.ID,
		DealershipID:   r.DealershipID,
		DealershipName: r.DealershipName,
		TechnicianName: r.TechnicianName,
		Level:          string(r.Level),
		IsActive:       r.IsActive,
		InactiveSince:  formatPgTimestamptzPtr(r.InactiveSince),
		CreatedAt:      r.CreatedAt.Format("15:04:05 02-01-2006"),
		UpdatedAt:      r.UpdatedAt.Format("15:04:05 02-01-2006"),
	}
}
