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
