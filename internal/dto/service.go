package dto

// SERVICES

type RequestServiceWithID struct {
	ID int32 `uri:"id" binding:"required,gt=0"`
}

type CreateServiceRequest struct {
	RequiredBayTypeID  int32  `json:"required_bay_type_id" binding:"required,gt=0"`
	Name               string `json:"name" binding:"required,not_blank,min=1,max=255"`
	AnticipatedMinutes int32  `json:"anticipated_minutes" binding:"required,gte=30"`
}

type UpdateServiceRequest struct {
	RequiredBayTypeID  *int32  `json:"required_bay_type_id" binding:"omitempty,gt=0"`
	Name               *string `json:"name" binding:"omitempty,not_blank,min=1,max=255"`
	AnticipatedMinutes *int32  `json:"anticipated_minutes" binding:"omitempty,gt=0"`
}

type SearchServicesQuery struct {
	Name string `form:"name" binding:"required,not_blank,min=1,max=255"`
}

// SERVICE REQUIREMENTS

type AddServiceRequirementsRequest struct {
	SkillIds []int32 `json:"skill_ids" binding:"required,min=1,dive,gt=0"`
}

type RemoveServiceRequirementsRequest struct {
	SkillIds []int32 `json:"skill_ids" binding:"required,min=1,dive,gt=0"`
}

// RESPONSES

type ServiceResponseHTTP struct {
	ID                 int32   `json:"id"`
	RequiredBayTypeID  int32   `json:"required_bay_type_id"`
	TypeName           *string `json:"type_name"`
	ServiceName        string  `json:"service_name"`
	AnticipatedMinutes int32   `json:"anticipated_minutes"`
	CreatedAt          string  `json:"created_at"`
	UpdatedAt          string  `json:"updated_at"`
}

type ServiceDetailResponseHTTP struct {
	ID                  int32  `json:"id"`
	RequiredBayTypeID   int32  `json:"required_bay_type_id"`
	RequiredBayTypeName string `json:"required_bay_type_name"`
	Name                string `json:"name"`
	AnticipatedMinutes  int32  `json:"anticipated_minutes"`
	CreatedAt           string `json:"created_at"`
	UpdatedAt           string `json:"updated_at"`
	RequiredSkillNames  any    `json:"required_skill_names"`
}
