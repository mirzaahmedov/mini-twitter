package model

type Tweet struct {
	ID         string  `json:"id,omitempty"`
	Content    string  `json:"content,omitempty"`
	AuthorID   string  `json:"author_id,omitempty"`
	Author     User    `json:"author"`
	Attachment *string `json:"attachment,omitempty"`
	UpdatedAt  *string `json:"updated_at,omitempty"`
	CreatedAt  string  `json:"created_at,omitempty"`
}

type TweetCreatePayload struct {
	Content  string `json:"content,omitempty" validate:"required"`
	AuthorID string `json:"author_id,omitempty"`
}