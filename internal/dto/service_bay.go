package dto

// SERVICE BAYS

type CreateServiceBayRequest struct {
	DealershipID int32  `json:"dealership_id" binding:"required,gt=0"`
	BayTypeID    int32  `json:"bay_type_id" binding:"required,gt=0"`
	Name         string `json:"name" binding:"required,not_blank,min=1,max=255"`
}

type UpdateServiceBayRequest struct {
	DealershipID *int32  `json:"dealership_id" binding:"omitempty,gt=0"`
	BayTypeID    *int32  `json:"bay_type_id" binding:"omitempty,gt=0"`
	Name         *string `json:"name" binding:"omitempty,not_blank,min=1,max=255"`
	IsActive     *bool   `json:"is_active" binding:"omitempty"`
}

type RequestServiceBayWithID struct {
	ID int32 `uri:"id" binding:"required,gt=0"`
}

type SearchServiceBaysQuery struct {
	Name         string `form:"name" binding:"required,not_blank,min=1,max=255"`
	DealershipID *int32 `form:"dealership_id" binding:"omitempty,gt=0"`
	BayTypeID    *int32 `form:"bay_type_id" binding:"omitempty,gt=0"`
}

type ListServiceBaysQuery struct {
	DealershipID *int32 `form:"dealership_id" binding:"omitempty,gt=0"`
	BayTypeID    *int32 `form:"bay_type_id" binding:"omitempty,gt=0"`
}

type ServiceBayResponseHTTP struct {
	ID             int32   `json:"id"`
	DealershipID   int32   `json:"dealership_id"`
	BayTypeID      int32   `json:"bay_type_id"`
	DealershipName string  `json:"dealership_name"`
	TypeName       *string `json:"type_name"`
	ServiceBayName string  `json:"service_bay_name"`
	IsActive       bool    `json:"is_active"`
	CreatedAt      string  `json:"created_at"`
	UpdatedAt      string  `json:"updated_at"`
}

// SERVICE BAY TYPES

type CreateServiceBayTypeRequest struct {
	Name string `json:"name" binding:"required,not_blank,min=1,max=255"`
}

type UpdateServiceBayTypeRequest struct {
	Name *string `json:"name" binding:"omitempty,not_blank,min=1,max=255"`
}

type RequestServiceBayTypeWithID struct {
	ID int32 `uri:"id" binding:"required,gt=0"`
}

type SearchServiceBayTypesQuery struct {
	Name string `form:"name" binding:"required,not_blank,min=1,max=255"`
}

type ServiceBayTypeResponseHTTP struct {
	ID        int32  `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
