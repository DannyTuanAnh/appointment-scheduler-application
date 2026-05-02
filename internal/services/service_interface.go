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

type SkillService interface {
	CreateSkill(ctx context.Context, name string) (dto.SkillResponseHTTP, error)
	ListSkills(ctx context.Context) ([]dto.SkillResponseHTTP, error)
	GetSkillByID(ctx context.Context, id int32) (dto.SkillResponseHTTP, error)
	SearchSkillsByName(ctx context.Context, name string) ([]dto.SkillResponseHTTP, error)
	UpdateSkillNameByID(ctx context.Context, id int32, name *string) (dto.SkillResponseHTTP, error)
	DeleteSkillByID(ctx context.Context, id int32) error
}

type ServiceBayService interface {
	// SERVICE BAYS
	ListServiceBays(ctx context.Context) ([]dto.ServiceBayResponseHTTP, error)
	ListServiceBaysByTypeID(ctx context.Context, bayTypeID int32) ([]dto.ServiceBayResponseHTTP, error)
	ListServiceBaysByDealershipID(ctx context.Context, dealershipID int32) ([]dto.ServiceBayResponseHTTP, error)
	ListServiceBaysByDealershipIDAndTypeID(ctx context.Context, dealershipID int32, bayTypeID int32) ([]dto.ServiceBayResponseHTTP, error)

	SearchServiceBaysByName(ctx context.Context, name string) ([]dto.ServiceBayResponseHTTP, error)
	SearchServiceBaysByNameAndDealershipID(ctx context.Context, name string, dealershipID int32) ([]dto.ServiceBayResponseHTTP, error)
	SearchServiceBaysByNameAndTypeID(ctx context.Context, name string, bayTypeID int32) ([]dto.ServiceBayResponseHTTP, error)
	SearchServiceBaysByNameDealershipIDAndTypeID(ctx context.Context, name string, dealershipID int32, bayTypeID int32) ([]dto.ServiceBayResponseHTTP, error)

	GetServiceBayByID(ctx context.Context, id int32) (dto.ServiceBayResponseHTTP, error)
	CreateServiceBay(ctx context.Context, req dto.CreateServiceBayRequest) error
	UpdateServiceBayByID(ctx context.Context, id int32, req dto.UpdateServiceBayRequest) error
	DeleteServiceBayByID(ctx context.Context, id int32) error

	// SERVICE BAY TYPES
	ListServiceBayTypes(ctx context.Context) ([]dto.ServiceBayTypeResponseHTTP, error)
	GetServiceBayTypeByID(ctx context.Context, id int32) (dto.ServiceBayTypeResponseHTTP, error)
	SearchServiceBayTypesByName(ctx context.Context, name string) ([]dto.ServiceBayTypeResponseHTTP, error)
	CreateServiceBayType(ctx context.Context, name string) (dto.ServiceBayTypeResponseHTTP, error)
	UpdateServiceBayTypeByID(ctx context.Context, id int32, name *string) (dto.ServiceBayTypeResponseHTTP, error)
	DeleteServiceBayTypeByID(ctx context.Context, id int32) error
}

type TechnicianService interface {
	CreateTechnician(ctx context.Context, req dto.CreateTechnicianRequest) error
	SetTechnicianOnLeave(ctx context.Context, id int32) (dto.TechnicianResponseHTTP, error)
	SetTechnicianBackToWork(ctx context.Context, id int32) (dto.TechnicianResponseHTTP, error)
	TransferTechnicianDealership(ctx context.Context, id int32, dealershipID int32) (dto.TechnicianResponseHTTP, error)
	UpdateTechnicianInfoByID(ctx context.Context, id int32, req dto.UpdateTechnicianRequest) error
	DeleteTechnicianByID(ctx context.Context, id int32) error
	DeleteTechnicianIfInactiveOverOneMonth(ctx context.Context, id int32) error

	AddSkillsToTechnician(ctx context.Context, technicianID int32, skillIDs []int32) error
	RemoveSkillsFromTechnician(ctx context.Context, technicianID int32, skillIDs []int32) (int64, error)

	ListTechniciansByDealershipID(ctx context.Context, dealershipID int32) ([]dto.TechnicianResponseHTTP, error)
	SearchTechniciansByName(ctx context.Context, name string) ([]dto.TechnicianResponseHTTP, error)
	SearchTechniciansByNameAndDealershipID(ctx context.Context, dealershipID int32, name string) ([]dto.TechnicianResponseHTTP, error)
	FindActiveTechniciansByDealershipWithRequiredSkills(ctx context.Context, dealershipID int32, skillIDs []int32) ([]int32, error)
	GetDetailTechnicianByID(ctx context.Context, id int32) (dto.TechnicianDetailResponseHTTP, error)
	GetTechnicianByID(ctx context.Context, id int32) (dto.TechnicianResponseHTTP, error)
}
