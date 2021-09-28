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
