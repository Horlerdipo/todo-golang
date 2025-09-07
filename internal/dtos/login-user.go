package dtos

type LoginUserDTO struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginUserResponseDto struct {
	Email     string       `json:"email"`
	FirstName string       `json:"first_name"`
	LastName  string       `json:"last_name"`
	Token     TokenDetails `json:"token"`
}

type TokenDetails struct {
	Token string `json:"token"`
	Exp   string `json:"exp"`
}
