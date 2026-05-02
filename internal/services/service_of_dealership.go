package services

import (
	"context"
	"encoding/json"

	"github.com/DannyTuanAnh/appointment-scheduler-application/internal/db/sqlc"
	"github.com/DannyTuanAnh/appointment-scheduler-application/internal/dto"
	"github.com/DannyTuanAnh/appointment-scheduler-application/internal/repositories"
	"github.com/DannyTuanAnh/appointment-scheduler-application/internal/utils"
)

type serviceService struct {
	repo repositories.ServiceRepository
}

func NewServiceService(repo repositories.ServiceRepository) ServiceService {
	return &serviceService{repo: repo}
}

func (s *serviceService) CreateService(ctx context.Context, req dto.CreateServiceRequest) error {
	return s.repo.CreateService(ctx, sqlc.CreateServiceParams{
		RequiredBayTypeID:  req.RequiredBayTypeID,
		Name:               req.Name,
		AnticipatedMinutes: req.AnticipatedMinutes,
	})
}

func (s *serviceService) UpdateServiceByID(ctx context.Context, id int32, req dto.UpdateServiceRequest) error {
	if req.RequiredBayTypeID == nil && req.Name == nil && req.AnticipatedMinutes == nil {
		return nil
	}
	return s.repo.UpdateServiceByID(ctx, sqlc.UpdateServiceByIDParams{
		ID:                 id,
		RequiredBayTypeID:  req.RequiredBayTypeID,
		Name:               req.Name,
		AnticipatedMinutes: req.AnticipatedMinutes,
	})
}

func (s *serviceService) DeleteServiceByID(ctx context.Context, id int32) error {
	_, err := s.repo.DeleteServiceByID(ctx, id)
	return err
}

func (s *serviceService) ListServices(ctx context.Context) ([]dto.ServiceResponseHTTP, error) {
	rows, err := s.repo.ListServices(ctx)
	if err != nil {
		return nil, err
	}
	res := make([]dto.ServiceResponseHTTP, len(rows))
	for i, r := range rows {
		res[i] = dto.ServiceResponseHTTP{
			ID:                 r.ID,
			RequiredBayTypeID:  r.RequiredBayTypeID,
			TypeName:           r.TypeName,
			ServiceName:        r.ServiceName,
			AnticipatedMinutes: r.AnticipatedMinutes,
			CreatedAt:          r.CreatedAt.Format("15:04:05 02-01-2006"),
			UpdatedAt:          r.UpdatedAt.Format("15:04:05 02-01-2006"),
		}
	}
	return res, nil
}

func (s *serviceService) GetServiceDetailByID(ctx context.Context, id int32) (dto.ServiceDetailResponseHTTP, error) {
	row, err := s.repo.GetServiceDetailByID(ctx, id)
	if err != nil {
		return dto.ServiceDetailResponseHTTP{}, err
	}

	return dto.ServiceDetailResponseHTTP{
		ID:                  row.ID,
		RequiredBayTypeID:   row.RequiredBayTypeID,
		RequiredBayTypeName: row.RequiredBayTypeName,
		Name:                row.Name,
		AnticipatedMinutes:  row.AnticipatedMinutes,
		CreatedAt:           row.CreatedAt.Format("15:04:05 02-01-2006"),
		UpdatedAt:           row.UpdatedAt.Format("15:04:05 02-01-2006"),
		RequiredSkillNames:  row.RequiredSkillsName,
	}, nil
}

func (s *serviceService) SearchServicesByName(ctx context.Context, name string) ([]dto.ServiceResponseHTTP, error) {
	rows, err := s.repo.SearchServicesByName(ctx, name)
	if err != nil {
		return nil, err
	}
	res := make([]dto.ServiceResponseHTTP, len(rows))
	for i, r := range rows {
		res[i] = dto.ServiceResponseHTTP{
			ID:                 r.ID,
			RequiredBayTypeID:  r.RequiredBayTypeID,
			TypeName:           r.TypeName,
			ServiceName:        r.ServiceName,
			AnticipatedMinutes: r.AnticipatedMinutes,
			CreatedAt:          r.CreatedAt.Format("15:04:05 02-01-2006"),
			UpdatedAt:          r.UpdatedAt.Format("15:04:05 02-01-2006"),
		}
	}
	return res, nil
}

// SERVICE REQUIREMENTS

func (s *serviceService) AddSkillRequirementsToService(ctx context.Context, serviceID int32, skillIDs []int32) error {
	return s.repo.AddSkillRequirementsToService(ctx, serviceID, skillIDs)
}

func (s *serviceService) RemoveSkillRequirementsFromService(ctx context.Context, serviceID int32, skillIDs []int32) (int64, error) {
	return s.repo.RemoveSkillRequirementsFromService(ctx, serviceID, skillIDs)
}

func decodeInt32Array(raw interface{}) ([]int32, error) {
	if raw == nil {
		return []int32{}, nil
	}
	switch v := raw.(type) {
	case []byte:
		var out []int32
		if err := json.Unmarshal(v, &out); err != nil {
			return nil, utils.WrapError(err, "failed to parse required_skill_ids", utils.ErrCodeInternal)
		}
		return out, nil
	case string:
		var out []int32
		if err := json.Unmarshal([]byte(v), &out); err != nil {
			return nil, utils.WrapError(err, "failed to parse required_skill_ids", utils.ErrCodeInternal)
		}
		return out, nil
	default:
		b, err := json.Marshal(v)
		if err != nil {
			return nil, utils.WrapError(err, "failed to parse required_skill_ids", utils.ErrCodeInternal)
		}
		var out []int32
		if err := json.Unmarshal(b, &out); err != nil {
			return nil, utils.WrapError(err, "failed to parse required_skill_ids", utils.ErrCodeInternal)
		}
		return out, nil
	}
}
