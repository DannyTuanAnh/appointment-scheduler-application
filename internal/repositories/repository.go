package repositories

import (
	"context"
	"time"

	"github.com/DannyTuanAnh/appointment-scheduler-application/internal/db/sqlc"
)

type DealershipRepository interface {
	GetAllDealerships(ctx context.Context) ([]sqlc.Dealership, error)
	GetDealershipByID(ctx context.Context, dealership_id int32) (sqlc.Dealership, error)
	SearchDealershipsByName(ctx context.Context, name string) ([]sqlc.Dealership, error)
	UpdateDealershipByID(ctx context.Context, id int32, name string, openTime, closeTime time.Time) error
	CreateDealership(ctx context.Context, name string, openTime, closeTime time.Time) error
	DeleteDealershipByID(ctx context.Context, id int32) error
}

type SkillRepository interface {
	CreateSkill(ctx context.Context, name string) (sqlc.Skill, error)
	ListSkills(ctx context.Context) ([]sqlc.Skill, error)
	GetSkillByID(ctx context.Context, id int32) (sqlc.Skill, error)
	SearchSkillsByName(ctx context.Context, name string) ([]sqlc.Skill, error)
	UpdateSkillNameByID(ctx context.Context, id int32, name string) (sqlc.Skill, error)
	DeleteSkillByID(ctx context.Context, id int32) error
}

type ServiceBayRepository interface {
	// SERVICE BAYS
	ListServiceBays(ctx context.Context) ([]sqlc.ListServiceBaysRow, error)
	ListServiceBaysByDealershipID(ctx context.Context, dealershipID int32) ([]sqlc.ListServiceBaysByDealershipIDRow, error)
	ListServiceBaysByTypeID(ctx context.Context, bayTypeID int32) ([]sqlc.ListServiceBaysByTypeIDRow, error)
	ListServiceBaysByDealershipIDAndTypeID(ctx context.Context, dealershipID int32, bayTypeID int32) ([]sqlc.ListServiceBaysByDealershipIDAndTypeIDRow, error)

	SearchServiceBaysByName(ctx context.Context, name string) ([]sqlc.SearchServiceBaysByNameRow, error)
	SearchServiceBaysByNameAndDealershipID(ctx context.Context, name string, dealershipID int32) ([]sqlc.SearchServiceBaysByNameAndDealershipIDRow, error)
	SearchServiceBaysByNameAndTypeID(ctx context.Context, name string, bayTypeID int32) ([]sqlc.SearchServiceBaysByNameAndTypeIDRow, error)
	SearchServiceBaysByNameDealershipIDAndTypeID(ctx context.Context, name string, dealershipID int32, bayTypeID int32) ([]sqlc.SearchServiceBaysByNameDealershipIDAndTypeIDRow, error)

	GetServiceBayByID(ctx context.Context, id int32) (sqlc.GetServiceBayByIDRow, error)
	CreateServiceBay(ctx context.Context, arg sqlc.CreateServiceBayParams) error
	UpdateServiceBayByID(ctx context.Context, arg sqlc.UpdateServiceBayByIDParams) error
	DeleteServiceBayByID(ctx context.Context, id int32) error

	// SERVICE BAY TYPES
	CreateServiceBayType(ctx context.Context, name string) (sqlc.ServiceBayType, error)
	ListServiceBayTypes(ctx context.Context) ([]sqlc.ServiceBayType, error)
	GetServiceBayTypeByID(ctx context.Context, id int32) (sqlc.ServiceBayType, error)
	SearchServiceBayTypesByName(ctx context.Context, name string) ([]sqlc.ServiceBayType, error)
	UpdateServiceBayTypeByID(ctx context.Context, id int32, name string) (sqlc.ServiceBayType, error)
	DeleteServiceBayTypeByID(ctx context.Context, id int32) error
}
