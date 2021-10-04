package controllers

type updateprofile struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

type updatePassword struct {
	Password        string `json:"password"`
	PasswordConfirm string `json:"password_confirm"`
}

type loginForm struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type registerForm struct {
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	PasswordConfirm string `json:"password_confirm"`
}

type userResponse struct {
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Email        string `json:"email"`
	IsAmbassador bool   `json:"is_ambassador"`
}
