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

	UpdateUserRequest struct {
		ID       string `json:"id" binding:"required,uuid4"`
		FullName string `json:"fullname" binding:"omitempty"`
		Password string `json:"password" binding:"omitempty"`
	}

	Param struct {
		ID string `uri:"id" binding:"required,uuid4"`
	}

	Query struct {
		Email    string `form:"email" binding:"omitempty,email"`
		Fullname string `form:"fullname" binding:"omitempty"`
		Page     int    `form:"page" binding:"omitempty,number"`
		Size     int    `form:"size" binding:"omitempty,number"`
	}
)
