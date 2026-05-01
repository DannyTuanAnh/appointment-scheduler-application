package services

import (
	"context"
	"fmt"
	"time"

	"github.com/DannyTuanAnh/appointment-scheduler-application/internal/dto"
	"github.com/DannyTuanAnh/appointment-scheduler-application/internal/repositories"
	"github.com/DannyTuanAnh/appointment-scheduler-application/internal/utils"
)

type dealershipService struct {
	repo repositories.DealershipRepository
}

func NewDealershipService(repo repositories.DealershipRepository) DealershipService {
	return &dealershipService{
		repo: repo,
	}
}

func (ds *dealershipService) GetAllDealerships(ctx context.Context) ([]dto.DealershipResponseHTTP, error) {
	dealerships, err := ds.repo.GetAllDealerships(ctx)
	if err != nil {
		return nil, err
	}

	response := make([]dto.DealershipResponseHTTP, len(dealerships))
	for i, dealership := range dealerships {
		response[i] = dto.DealershipResponseHTTP{
			ID:        dealership.ID,
			Name:      dealership.Name,
			OpenTime:  utils.ConvertPgTypeTimeToTime(dealership.OpenTime).Format("15:04"),
			CloseTime: utils.ConvertPgTypeTimeToTime(dealership.CloseTime).Format("15:04"),
			CreatedAt: dealership.CreatedAt.Format("15:04:05 02-01-2006"),
			UpdatedAt: dealership.UpdatedAt.Format("15:04:05 02-01-2006"),
		}
	}

	return response, nil
}

func (ds *dealershipService) GetDealershipByID(ctx context.Context, id int32) (dto.DealershipResponseHTTP, error) {
	dealership, err := ds.repo.GetDealershipByID(ctx, id)
	if err != nil {
		return dto.DealershipResponseHTTP{}, err
	}

	response := dto.DealershipResponseHTTP{
		ID:        dealership.ID,
		Name:      dealership.Name,
		OpenTime:  utils.ConvertPgTypeTimeToTime(dealership.OpenTime).Format("15:04"),
		CloseTime: utils.ConvertPgTypeTimeToTime(dealership.CloseTime).Format("15:04"),
		CreatedAt: dealership.CreatedAt.Format("15:04:05 02-01-2006"),
		UpdatedAt: dealership.UpdatedAt.Format("15:04:05 02-01-2006"),
	}

	return response, nil
}

func (ds *dealershipService) CreateDealership(ctx context.Context, name string, openTime, closeTime time.Time) error {
	if closeTime.Equal(openTime) {
		return utils.NewError("close_time must be different from open_time", utils.ErrCodeBadRequest)
	}

	if closeTime.Before(openTime) {
		return utils.NewError("close_time must be after open_time", utils.ErrCodeBadRequest)
	}

	duration := closeTime.Sub(openTime)
	if duration < 8*time.Hour {
		return utils.NewError(fmt.Sprintf("business hours must be at least 8 hours (current: %s)", formatDurationToHM(duration)), utils.ErrCodeBadRequest)
	}

	err := ds.repo.CreateDealership(ctx, name, openTime, closeTime)
	if err != nil {
		return err
	}

	return nil
}

func formatDurationToHM(duration time.Duration) string {
	hours := int(duration.Hours())
	minutes := int(duration.Minutes()) % 60
	return fmt.Sprintf("%dh %dm", hours, minutes)
}

func (ds *dealershipService) UpdateDealership(ctx context.Context, id int32, name *string, openTime, closeTime *time.Time) error {
	dealership, err := ds.repo.GetDealershipByID(ctx, id)
	if err != nil {
		return err
	}

	var finalName string
	if name != nil {
		finalName = *name
	} else {
		finalName = dealership.Name
	}

	var finalOpen, finalClose time.Time

	if openTime != nil {
		finalOpen = utils.NormalizeToSameDate(*openTime)
	} else {
		finalOpen = utils.NormalizeToSameDate(utils.ConvertPgTypeTimeToTime(dealership.OpenTime))
	}

	if closeTime != nil {
		finalClose = utils.NormalizeToSameDate(*closeTime)
	} else {
		finalClose = utils.NormalizeToSameDate(utils.ConvertPgTypeTimeToTime(dealership.CloseTime))
	}

	if finalClose.Equal(finalOpen) {
		return utils.NewError("close_time must be different from open_time", utils.ErrCodeBadRequest)
	}

	if finalClose.Before(finalOpen) {
		return utils.NewError("close_time must be after open_time", utils.ErrCodeBadRequest)
	}

	duration := finalClose.Sub(finalOpen)
	if duration < 8*time.Hour {
		return utils.NewError(fmt.Sprintf("business hours must be at least 8 hours (current: %s)", formatDurationToHM(duration)), utils.ErrCodeBadRequest)
	}

	err = ds.repo.UpdateDealershipByID(ctx, id, finalName, finalOpen, finalClose)
	if err != nil {
		return err
	}

	return nil
}

func (ds *dealershipService) DeleteDealershipByID(ctx context.Context, id int32) error {
	err := ds.repo.DeleteDealershipByID(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (ds *dealershipService) SearchDealershipsByName(ctx context.Context, name string) ([]dto.DealershipResponseHTTP, error) {
	dealerships, err := ds.repo.SearchDealershipsByName(ctx, name)
	if err != nil {
		return nil, err
	}

	response := make([]dto.DealershipResponseHTTP, len(dealerships))
	for i, dealership := range dealerships {
		response[i] = dto.DealershipResponseHTTP{
			ID:        dealership.ID,
			Name:      dealership.Name,
			OpenTime:  utils.ConvertPgTypeTimeToTime(dealership.OpenTime).Format("15:04"),
			CloseTime: utils.ConvertPgTypeTimeToTime(dealership.CloseTime).Format("15:04"),
			CreatedAt: dealership.CreatedAt.Format("15:04:05 02-01-2006"),
			UpdatedAt: dealership.UpdatedAt.Format("15:04:05 02-01-2006"),
		}
	}

	return response, nil
}
