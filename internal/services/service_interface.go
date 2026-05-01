package services

import (
	"context"
	"time"

	"github.com/DannyTuanAnh/appointment-scheduler-application/internal/dto"
)

type DealershipService interface {
	GetAllDealerships(ctx context.Context) ([]dto.DealershipResponseHTTP, error)
	GetDealershipByID(ctx context.Context, id int32) (dto.DealershipResponseHTTP, error)
	SearchDealershipsByName(ctx context.Context, name string) ([]dto.DealershipResponseHTTP, error)
	CreateDealership(ctx context.Context, name string, openTime, closeTime time.Time) error
	UpdateDealership(ctx context.Context, id int32, name *string, openTime, closeTime *time.Time) error
	DeleteDealershipByID(ctx context.Context, id int32) error
}
