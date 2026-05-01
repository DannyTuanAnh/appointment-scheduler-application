package services

import (
	"context"

	"github.com/DannyTuanAnh/appointment-scheduler-application/internal/dto"
	"github.com/DannyTuanAnh/appointment-scheduler-application/internal/repositories"
	"github.com/DannyTuanAnh/appointment-scheduler-application/internal/utils"
)

type skillService struct {
	repo repositories.SkillRepository
}

func NewSkillService(repo repositories.SkillRepository) SkillService {
	return &skillService{repo: repo}
}

func (ss *skillService) CreateSkill(ctx context.Context, name string) (dto.SkillResponseHTTP, error) {
	skill, err := ss.repo.CreateSkill(ctx, name)
	if err != nil {
		return dto.SkillResponseHTTP{}, err
	}

	return dto.SkillResponseHTTP{
		ID:        skill.ID,
		Name:      skill.Name,
		CreatedAt: skill.CreatedAt.Format("15:04:05 02-01-2006"),
		UpdatedAt: skill.UpdatedAt.Format("15:04:05 02-01-2006"),
	}, nil
}

func (ss *skillService) ListSkills(ctx context.Context) ([]dto.SkillResponseHTTP, error) {
	skills, err := ss.repo.ListSkills(ctx)
	if err != nil {
		return nil, err
	}

	resp := make([]dto.SkillResponseHTTP, len(skills))
	for i, s := range skills {
		resp[i] = dto.SkillResponseHTTP{
			ID:        s.ID,
			Name:      s.Name,
			CreatedAt: s.CreatedAt.Format("15:04:05 02-01-2006"),
			UpdatedAt: s.UpdatedAt.Format("15:04:05 02-01-2006"),
		}
	}
	return resp, nil
}

func (ss *skillService) GetSkillByID(ctx context.Context, id int32) (dto.SkillResponseHTTP, error) {
	s, err := ss.repo.GetSkillByID(ctx, id)
	if err != nil {
		return dto.SkillResponseHTTP{}, err
	}
	return dto.SkillResponseHTTP{
		ID:        s.ID,
		Name:      s.Name,
		CreatedAt: s.CreatedAt.Format("15:04:05 02-01-2006"),
		UpdatedAt: s.UpdatedAt.Format("15:04:05 02-01-2006"),
	}, nil
}

func (ss *skillService) SearchSkillsByName(ctx context.Context, name string) ([]dto.SkillResponseHTTP, error) {
	skills, err := ss.repo.SearchSkillsByName(ctx, name)
	if err != nil {
		return nil, err
	}
	resp := make([]dto.SkillResponseHTTP, len(skills))
	for i, s := range skills {
		resp[i] = dto.SkillResponseHTTP{
			ID:        s.ID,
			Name:      s.Name,
			CreatedAt: s.CreatedAt.Format("15:04:05 02-01-2006"),
			UpdatedAt: s.UpdatedAt.Format("15:04:05 02-01-2006"),
		}
	}
	return resp, nil
}

func (ss *skillService) UpdateSkillNameByID(ctx context.Context, id int32, name *string) (dto.SkillResponseHTTP, error) {
	if name == nil {
		return dto.SkillResponseHTTP{}, utils.NewError("no update fields", utils.ErrCodeBadRequest)
	}

	updated, err := ss.repo.UpdateSkillNameByID(ctx, id, *name)
	if err != nil {
		return dto.SkillResponseHTTP{}, err
	}

	return dto.SkillResponseHTTP{
		ID:        updated.ID,
		Name:      updated.Name,
		CreatedAt: updated.CreatedAt.Format("15:04:05 02-01-2006"),
		UpdatedAt: updated.UpdatedAt.Format("15:04:05 02-01-2006"),
	}, nil
}

func (ss *skillService) DeleteSkillByID(ctx context.Context, id int32) error {
	return ss.repo.DeleteSkillByID(ctx, id)
}
