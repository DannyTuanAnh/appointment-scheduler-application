package dto

// TECHNICIANS

type RequestTechnicianWithID struct {
	ID int32 `uri:"id" binding:"required,gt=0"`
}

type CreateTechnicianRequest struct {
	DealershipID int32  `json:"dealership_id" binding:"required,gt=0"`
	Name         string `json:"name" binding:"required,not_blank,min=1,max=255"`
	Level        string `json:"level" binding:"required,not_blank"`
}

type UpdateTechnicianRequest struct {
	Name  *string `json:"name" binding:"omitempty,not_blank,min=1,max=255"`
	Level *string `json:"level" binding:"omitempty,not_blank"`
}

type TransferTechnicianDealershipRequest struct {
	DealershipID int32 `json:"dealership_id" binding:"required,gt=0"`
}

type SearchTechniciansQuery struct {
	Name         string `form:"name" binding:"required,not_blank,min=1,max=255"`
	DealershipID *int32 `form:"dealership_id" binding:"omitempty,gt=0"`
}

type ListTechniciansQuery struct {
	DealershipID int32 `form:"dealership_id" binding:"required,gt=0"`
}

// Skills assignment

type AddSkillsToTechnicianRequest struct {
	SkillIds []int32 `json:"skill_ids" binding:"required,min=1,dive,gt=0"`
}

type RemoveSkillsFromTechnicianRequest struct {
	SkillIds []int32 `json:"skill_ids" binding:"required,min=1,dive,gt=0"`
}

// Used for FindActiveTechniciansByDealershipWithRequiredSkills

type FindActiveTechniciansQuery struct {
	DealershipID int32   `form:"dealership_id" binding:"required,gt=0"`
	SkillIDs     []int32 `form:"skill_ids[]" binding:"omitempty,dive,gt=0"`
}

// Responses

type TechnicianResponseHTTP struct {
	ID             int32   `json:"id"`
	DealershipID   int32   `json:"dealership_id"`
	DealershipName string  `json:"dealership_name"`
	TechnicianName string  `json:"technician_name"`
	Level          string  `json:"level"`
	IsActive       bool    `json:"is_active"`
	InactiveSince  *string `json:"inactive_since"`
	CreatedAt      string  `json:"created_at"`
	UpdatedAt      string  `json:"updated_at"`
}

type TechnicianDetailResponseHTTP struct {
	TechnicianResponseHTTP
	Skills []string `json:"skills"`
}
