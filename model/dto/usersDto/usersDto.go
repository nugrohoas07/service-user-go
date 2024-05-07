package usersDto

type (
	LoginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	CreateUserRequest struct {
		FullName string `json:"fullname" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}
)
