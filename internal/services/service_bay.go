package services

import (
	"context"

	"github.com/DannyTuanAnh/appointment-scheduler-application/internal/db/sqlc"
	"github.com/DannyTuanAnh/appointment-scheduler-application/internal/dto"
	"github.com/DannyTuanAnh/appointment-scheduler-application/internal/repositories"
	"github.com/DannyTuanAnh/appointment-scheduler-application/internal/utils"
)

type serviceBayService struct {
	repo repositories.ServiceBayRepository
}

func NewServiceBayService(repo repositories.ServiceBayRepository) ServiceBayService {
	return &serviceBayService{repo: repo}
}

// SERVICE BAYS

func (s *serviceBayService) ListServiceBays(ctx context.Context) ([]dto.ServiceBayResponseHTTP, error) {
	bays, err := s.repo.ListServiceBays(ctx)
	if err != nil {
		return []dto.ServiceBayResponseHTTP{}, err
	}

	resp := make([]dto.ServiceBayResponseHTTP, len(bays))
	for i, b := range bays {
		resp[i] = dto.ServiceBayResponseHTTP{
			ID:             b.ID,
			DealershipID:   b.DealershipID,
			BayTypeID:      b.BayTypeID,
			DealershipName: b.DealershipName,
			TypeName:       b.TypeName,
			ServiceBayName: b.ServiceBayName,
			IsActive:       b.IsActive,
			CreatedAt:      b.CreatedAt.Format("15:04:05 02-01-2006"),
			UpdatedAt:      b.UpdatedAt.Format("15:04:05 02-01-2006"),
		}
	}

	return resp, nil
}

func (s *serviceBayService) ListServiceBaysByTypeID(ctx context.Context, bayTypeID int32) ([]dto.ServiceBayResponseHTTP, error) {
	bays, err := s.repo.ListServiceBaysByTypeID(ctx, bayTypeID)
	if err != nil {
		return []dto.ServiceBayResponseHTTP{}, err
	}

	resp := make([]dto.ServiceBayResponseHTTP, len(bays))
	for i, b := range bays {
		resp[i] = dto.ServiceBayResponseHTTP{
			ID:             b.ID,
			DealershipID:   b.DealershipID,
			BayTypeID:      b.BayTypeID,
			DealershipName: b.DealershipName,
			TypeName:       b.TypeName,
			ServiceBayName: b.ServiceBayName,
			IsActive:       b.IsActive,
			CreatedAt:      b.CreatedAt.Format("15:04:05 02-01-2006"),
			UpdatedAt:      b.UpdatedAt.Format("15:04:05 02-01-2006"),
		}
	}

	return resp, nil

}

func (s *serviceBayService) ListServiceBaysByDealershipID(ctx context.Context, dealershipID int32) ([]dto.ServiceBayResponseHTTP, error) {
	bays, err := s.repo.ListServiceBaysByDealershipID(ctx, dealershipID)
	if err != nil {
		return []dto.ServiceBayResponseHTTP{}, err
	}

	resp := make([]dto.ServiceBayResponseHTTP, len(bays))
	for i, b := range bays {
		resp[i] = dto.ServiceBayResponseHTTP{
			ID:             b.ID,
			DealershipID:   b.DealershipID,
			BayTypeID:      b.BayTypeID,
			DealershipName: b.DealershipName,
			TypeName:       b.TypeName,
			ServiceBayName: b.ServiceBayName,
			IsActive:       b.IsActive,
			CreatedAt:      b.CreatedAt.Format("15:04:05 02-01-2006"),
			UpdatedAt:      b.UpdatedAt.Format("15:04:05 02-01-2006"),
		}
	}

	return resp, nil
}

func (s *serviceBayService) ListServiceBaysByDealershipIDAndTypeID(ctx context.Context, dealershipID int32, bayTypeID int32) ([]dto.ServiceBayResponseHTTP, error) {
	bays, err := s.repo.ListServiceBaysByDealershipIDAndTypeID(ctx, dealershipID, bayTypeID)

	if err != nil {
		return nil, err
	}

	resp := make([]dto.ServiceBayResponseHTTP, len(bays))
	for i, b := range bays {
		resp[i] = dto.ServiceBayResponseHTTP{
			ID:             b.ID,
			DealershipID:   b.DealershipID,
			BayTypeID:      b.BayTypeID,
			DealershipName: b.DealershipName,
			TypeName:       b.TypeName,
			ServiceBayName: b.ServiceBayName,
			IsActive:       b.IsActive,
			CreatedAt:      b.CreatedAt.Format("15:04:05 02-01-2006"),
			UpdatedAt:      b.UpdatedAt.Format("15:04:05 02-01-2006"),
		}
	}

	return resp, nil
}

func (s *serviceBayService) SearchServiceBaysByName(ctx context.Context, name string) ([]dto.ServiceBayResponseHTTP, error) {
	bays, err := s.repo.SearchServiceBaysByName(ctx, name)
	if err != nil {
		return nil, err
	}

	resp := make([]dto.ServiceBayResponseHTTP, len(bays))
	for i, b := range bays {
		resp[i] = dto.ServiceBayResponseHTTP{
			ID:             b.ID,
			DealershipID:   b.DealershipID,
			BayTypeID:      b.BayTypeID,
			DealershipName: b.DealershipName,
			TypeName:       b.TypeName,
			ServiceBayName: b.ServiceBayName,
			IsActive:       b.IsActive,
			CreatedAt:      b.CreatedAt.Format("15:04:05 02-01-2006"),
			UpdatedAt:      b.UpdatedAt.Format("15:04:05 02-01-2006"),
		}
	}

	return resp, nil
}

func (s *serviceBayService) SearchServiceBaysByNameAndDealershipID(ctx context.Context, name string, dealershipID int32) ([]dto.ServiceBayResponseHTTP, error) {
	bays, err := s.repo.SearchServiceBaysByNameAndDealershipID(ctx, name, dealershipID)
	if err != nil {
		return nil, err
	}

	resp := make([]dto.ServiceBayResponseHTTP, len(bays))
	for i, b := range bays {
		resp[i] = dto.ServiceBayResponseHTTP{
			ID:             b.ID,
			DealershipID:   b.DealershipID,
			BayTypeID:      b.BayTypeID,
			DealershipName: b.DealershipName,
			TypeName:       b.TypeName,
			ServiceBayName: b.ServiceBayName,
			IsActive:       b.IsActive,
			CreatedAt:      b.CreatedAt.Format("15:04:05 02-01-2006"),
			UpdatedAt:      b.UpdatedAt.Format("15:04:05 02-01-2006"),
		}
	}

	return resp, nil
}

func (s *serviceBayService) SearchServiceBaysByNameAndTypeID(ctx context.Context, name string, bayTypeID int32) ([]dto.ServiceBayResponseHTTP, error) {
	bays, err := s.repo.SearchServiceBaysByNameAndTypeID(ctx, name, bayTypeID)
	if err != nil {
		return nil, err
	}

	resp := make([]dto.ServiceBayResponseHTTP, len(bays))
	for i, b := range bays {
		resp[i] = dto.ServiceBayResponseHTTP{
			ID:             b.ID,
			DealershipID:   b.DealershipID,
			BayTypeID:      b.BayTypeID,
			DealershipName: b.DealershipName,
			TypeName:       b.TypeName,
			ServiceBayName: b.ServiceBayName,
			IsActive:       b.IsActive,
			CreatedAt:      b.CreatedAt.Format("15:04:05 02-01-2006"),
			UpdatedAt:      b.UpdatedAt.Format("15:04:05 02-01-2006"),
		}
	}

	return resp, nil
}

func (s *serviceBayService) SearchServiceBaysByNameDealershipIDAndTypeID(ctx context.Context, name string, dealershipID int32, bayTypeID int32) ([]dto.ServiceBayResponseHTTP, error) {
	bays, err := s.repo.SearchServiceBaysByNameDealershipIDAndTypeID(ctx, name, dealershipID, bayTypeID)
	if err != nil {
		return nil, err
	}

	resp := make([]dto.ServiceBayResponseHTTP, len(bays))
	for i, b := range bays {
		resp[i] = dto.ServiceBayResponseHTTP{
			ID:             b.ID,
			DealershipID:   b.DealershipID,
			BayTypeID:      b.BayTypeID,
			DealershipName: b.DealershipName,
			TypeName:       b.TypeName,
			ServiceBayName: b.ServiceBayName,
			IsActive:       b.IsActive,
			CreatedAt:      b.CreatedAt.Format("15:04:05 02-01-2006"),
			UpdatedAt:      b.UpdatedAt.Format("15:04:05 02-01-2006"),
		}
	}

	return resp, nil
}

func (s *serviceBayService) GetServiceBayByID(ctx context.Context, id int32) (dto.ServiceBayResponseHTTP, error) {
	b, err := s.repo.GetServiceBayByID(ctx, id)
	if err != nil {
		return dto.ServiceBayResponseHTTP{}, err
	}

	return dto.ServiceBayResponseHTTP{
		ID:             b.ID,
		DealershipID:   b.DealershipID,
		BayTypeID:      b.BayTypeID,
		DealershipName: b.DealershipName,
		TypeName:       b.TypeName,
		ServiceBayName: b.ServiceBayName,
		IsActive:       b.IsActive,
		CreatedAt:      b.CreatedAt.Format("15:04:05 02-01-2006"),
		UpdatedAt:      b.UpdatedAt.Format("15:04:05 02-01-2006"),
	}, nil
}

func (s *serviceBayService) CreateServiceBay(ctx context.Context, req dto.CreateServiceBayRequest) error {
	err := s.repo.CreateServiceBay(ctx, sqlc.CreateServiceBayParams{
		DealershipID: req.DealershipID,
		BayTypeID:    req.BayTypeID,
		Name:         req.Name,
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *serviceBayService) UpdateServiceBayByID(ctx context.Context, id int32, req dto.UpdateServiceBayRequest) error {
	if req.DealershipID == nil && req.BayTypeID == nil && req.Name == nil && req.IsActive == nil {
		return utils.NewError("no updated fields", utils.ErrCodeBadRequest)
	}

	arg := sqlc.UpdateServiceBayByIDParams{
		ID:           id,
		DealershipID: req.DealershipID,
		BayTypeID:    req.BayTypeID,
		Name:         req.Name,
		IsActive:     req.IsActive,
	}

	err := s.repo.UpdateServiceBayByID(ctx, arg)
	if err != nil {
		return err
	}

	return nil
}

func (s *serviceBayService) DeleteServiceBayByID(ctx context.Context, id int32) error {
	return s.repo.DeleteServiceBayByID(ctx, id)
}

// SERVICE BAY TYPES
func (s *serviceBayService) ListServiceBayTypes(ctx context.Context) ([]dto.ServiceBayTypeResponseHTTP, error) {
	types, err := s.repo.ListServiceBayTypes(ctx)
	if err != nil {
		return nil, err
	}
	resp := make([]dto.ServiceBayTypeResponseHTTP, len(types))
	for i, t := range types {
		resp[i] = dto.ServiceBayTypeResponseHTTP{
			ID:        t.ID,
			Name:      t.Name,
			CreatedAt: t.CreatedAt.Format("15:04:05 02-01-2006"),
			UpdatedAt: t.UpdatedAt.Format("15:04:05 02-01-2006"),
		}
	}
	return resp, nil
}

func (s *serviceBayService) CreateServiceBayType(ctx context.Context, name string) (dto.ServiceBayTypeResponseHTTP, error) {
	t, err := s.repo.CreateServiceBayType(ctx, name)
	if err != nil {
		return dto.ServiceBayTypeResponseHTTP{}, err
	}

	return dto.ServiceBayTypeResponseHTTP{
		ID:        t.ID,
		Name:      t.Name,
		CreatedAt: t.CreatedAt.Format("15:04:05 02-01-2006"),
		UpdatedAt: t.UpdatedAt.Format("15:04:05 02-01-2006"),
	}, nil
}

func (s *serviceBayService) GetServiceBayTypeByID(ctx context.Context, id int32) (dto.ServiceBayTypeResponseHTTP, error) {
	t, err := s.repo.GetServiceBayTypeByID(ctx, id)
	if err != nil {
		return dto.ServiceBayTypeResponseHTTP{}, err
	}

	return dto.ServiceBayTypeResponseHTTP{
		ID:        t.ID,
		Name:      t.Name,
		CreatedAt: t.CreatedAt.Format("15:04:05 02-01-2006"),
		UpdatedAt: t.UpdatedAt.Format("15:04:05 02-01-2006"),
	}, nil
}

func (s *serviceBayService) SearchServiceBayTypesByName(ctx context.Context, name string) ([]dto.ServiceBayTypeResponseHTTP, error) {
	types, err := s.repo.SearchServiceBayTypesByName(ctx, name)
	if err != nil {
		return nil, err
	}

	resp := make([]dto.ServiceBayTypeResponseHTTP, len(types))
	for i, t := range types {
		resp[i] = dto.ServiceBayTypeResponseHTTP{
			ID:        t.ID,
			Name:      t.Name,
			CreatedAt: t.CreatedAt.Format("15:04:05 02-01-2006"),
			UpdatedAt: t.UpdatedAt.Format("15:04:05 02-01-2006"),
		}
	}
	return resp, nil
}

func (s *serviceBayService) UpdateServiceBayTypeByID(ctx context.Context, id int32, name *string) (dto.ServiceBayTypeResponseHTTP, error) {
	if name == nil {
		return dto.ServiceBayTypeResponseHTTP{}, utils.NewError("no update fields", utils.ErrCodeBadRequest)
	}

	updated, err := s.repo.UpdateServiceBayTypeByID(ctx, id, *name)
	if err != nil {
		return dto.ServiceBayTypeResponseHTTP{}, err
	}

	return dto.ServiceBayTypeResponseHTTP{
		ID:        updated.ID,
		Name:      updated.Name,
		CreatedAt: updated.CreatedAt.Format("15:04:05 02-01-2006"),
		UpdatedAt: updated.UpdatedAt.Format("15:04:05 02-01-2006"),
	}, nil
}

func (s *serviceBayService) DeleteServiceBayTypeByID(ctx context.Context, id int32) error {
	return s.repo.DeleteServiceBayTypeByID(ctx, id)
}
