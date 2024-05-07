package usersEntity

type (
	UserData struct {
		ID       string `json:"id"`
		FullName string `json:"fullname"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
)