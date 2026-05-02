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

type TechnicianRepository interface {
	CreateTechnician(ctx context.Context, arg sqlc.CreateTechnicianParams) error
	SetTechnicianOnLeave(ctx context.Context, id int32) (sqlc.Technician, error)
	SetTechnicianBackToWork(ctx context.Context, id int32) (sqlc.Technician, error)
	TransferTechnicianDealership(ctx context.Context, id int32, dealershipID int32) (sqlc.Technician, error)
	UpdateTechnicianInfoByID(ctx context.Context, id int32, name *string, level *sqlc.TechnicianLevel) error
	DeleteTechnicianByID(ctx context.Context, id int32) error
	DeleteTechnicianIfInactiveOverOneMonth(ctx context.Context, id int32) error

	AddSkillsToTechnician(ctx context.Context, technicianID int32, skillIDs []int32) error
	RemoveSkillsFromTechnician(ctx context.Context, technicianID int32, skillIDs []int32) (int64, error)

	ListTechniciansByDealershipID(ctx context.Context, dealershipID int32) ([]sqlc.ListTechniciansByDealershipIDRow, error)
	SearchTechniciansByName(ctx context.Context, name string) ([]sqlc.SearchTechniciansByNameRow, error)
	SearchTechniciansByNameAndDealershipID(ctx context.Context, dealershipID int32, name string) ([]sqlc.SearchTechniciansByNameAndDealershipIDRow, error)
	FindActiveTechniciansByDealershipWithRequiredSkills(ctx context.Context, dealershipID int32, skillIDs []int32) ([]int32, error)
	GetDetailTechnicianByID(ctx context.Context, id int32) (sqlc.GetDetailTechnicianByIDRow, error)
	GetTechnicianByID(ctx context.Context, id int32) (sqlc.GetTechnicianByIDRow, error)
}

type ServiceRepository interface {
	// SERVICES
	CreateService(ctx context.Context, arg sqlc.CreateServiceParams) error
	UpdateServiceByID(ctx context.Context, arg sqlc.UpdateServiceByIDParams) error
	DeleteServiceByID(ctx context.Context, id int32) (int64, error)
	ListServices(ctx context.Context) ([]sqlc.ListServicesRow, error)
	GetServiceDetailByID(ctx context.Context, id int32) (sqlc.GetServiceDetailByIDRow, error)
	SearchServicesByName(ctx context.Context, name string) ([]sqlc.SearchServicesByNameRow, error)

	// SERVICE REQUIREMENTS
	AddSkillRequirementsToService(ctx context.Context, serviceID int32, skillIDs []int32) error
	RemoveSkillRequirementsFromService(ctx context.Context, serviceID int32, skillIDs []int32) (int64, error)
}

type AppointmentRepository interface {
	GetAppointment(ctx context.Context, params sqlc.GetAppointmentsOfBayOrTechnicianInTimeRangeParams) ([]sqlc.GetAppointmentsOfBayOrTechnicianInTimeRangeRow, error)
	GetWorkHoursOfDealership(ctx context.Context, dealershipID int32) (time.Time, time.Time, error)
	GetSkillRequirementIDs(ctx context.Context, erviceID int32) ([]int32, error)
	GetServiceBayIDsByDealershipIDAndTypeID(ctx context.Context, dealershipID, bayTypeID int32) ([]int32, error)
	FindTechniciansByDealershipWithRequiredSkills(ctx context.Context, dealershipID int32, skillIDs []int32) ([]int32, error)
}
