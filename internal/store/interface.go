package store

import (
	"twitter/internal/model"
)

type Store interface {
	User() UserStore
	Tweet() TweetStore
	Follow() FollowStore
	Open() error
	Close() error
}

type UserStore interface {
	Create(payload *model.UserCreatePayload) (*model.User, error)
	Search(term, userID string) ([]model.User, error)
	FindByID(id string) (*model.User, error)
	FindByUsername(username string) (*model.User, error)
}

type TweetStore interface {
	Create(payload *model.TweetCreatePayload) (*model.Tweet, error)
	Delete(id, userId string) error
	Update(id, userID string, payload *model.Tweet) (*model.Tweet, error)
	FindByID(id string) (*model.Tweet, error)
	GetTweetsByUser(userID string) ([]model.Tweet, error)
	GetTweetsFromFollowedUsers(userID string) ([]model.Tweet, error)
}

type FollowStore interface {
	Create(payload *model.Follow) (*model.Follow, error)
	Delete(followerID, followingID string) error
	FindAllWithFollowerID(id string) ([]model.User, error)
}
