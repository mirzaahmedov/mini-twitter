package model

type User struct {
	ID             string `json:"id,omitempty"`
	Name           string `json:"name,omitempty"`
	Username       string `json:"username,omitempty"`
	Bio            string `json:"bio"`
	Password       string `json:"-"`
	ProfilePicture string `json:"profile_picture,omitempty"`
	CreatedAt      string `json:"created_at,omitempty"`
}

type UserCreatePayload struct {
	Name           string `json:"name,omitempty" validate:"required"`
	Username       string `json:"username,omitempty" validate:"required"`
	Bio            string `json:"bio"`
	Password       string `json:"password" validate:"required"`
	ProfilePicture string `json:"profile_picture,omitempty"`
}

type UserLoginPayload struct {
	Username string `json:"username,omitempty" validate:"required"`
	Password string `json:"password" validate:"required"`
}
