package domain

type NewUserReq struct {
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required,email"`
}

type UpdateUserReq struct {
	Id    int    `json:"id" validate:"required"`
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required,email"`
}

type DeleteUserReq struct {
	Id int `json:"id" validate:"required"`
}

type GetNamesReq struct {
	Method int32 `json:"method" validate:"required"`
	WaitTime int32 `json:"waitTime" validate:"required"`
}
