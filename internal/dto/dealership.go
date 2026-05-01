package dto

type CreateDealershipRequest struct {
	Name      string `json:"name" binding:"required,not_blank,min=2,max=255"`
	OpenTime  string `json:"open_time" binding:"required,datetime=15:04"`
	CloseTime string `json:"close_time" binding:"required,datetime=15:04"`
}

type RequestDealershipWithID struct {
	ID int32 `uri:"id" binding:"required,gt=0"`
}

type UpdateDealershipRequest struct {
	Name      *string `json:"name" binding:"omitempty,not_blank,min=2,max=255"`
	OpenTime  *string `json:"open_time" binding:"omitempty,datetime=15:04"`
	CloseTime *string `json:"close_time" binding:"omitempty,datetime=15:04"`
}

type DealershipResponseHTTP struct {
	ID        int32  `json:"id"`
	Name      string `json:"name"`
	OpenTime  string `json:"open_time"`
	CloseTime string `json:"close_time"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type SearchDealershipsQuery struct {
	Name string `form:"name" binding:"required,not_blank,min=2,max=255"`
}
