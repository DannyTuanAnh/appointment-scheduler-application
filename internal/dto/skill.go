package dto

type RequestSkillWithID struct {
	ID int32 `uri:"id" binding:"required,gt=0"`
}

type CreateSkillRequest struct {
	Name string `json:"name" binding:"required,not_blank,min=1,max=100"`
}

type UpdateSkillRequest struct {
	Name *string `json:"name" binding:"omitempty,not_blank,min=1,max=100"`
}

type SearchSkillsQuery struct {
	Name string `form:"name" binding:"required,not_blank,min=1,max=100"`
}

type SkillResponseHTTP struct {
	ID        int32  `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
