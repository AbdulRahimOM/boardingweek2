package domain

// General error response...............................................
type ErrorRes struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Error   string `json:"error"`
}

// General success response...............................................
type SuccessRes struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
}

//Specifid Success Responses...............................................

type UserCreatedRes struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	UserId  int    `json:"user_id"`
}

type UserUpdatedRes struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	User    User   `json:"user"`
}

type GetUserRes struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	User    User   `json:"user"`
}

type GetUserNamesRes struct {
	Status bool     `json:"status"`
	Message string  `json:"message"`
	Names  []string `json:"names"`
}