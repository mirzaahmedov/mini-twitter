package model

type Follow struct {
	FollowerID  string `json:"follower_id,omitempty"`
	FollowingID string `json:"following_id,omitempty"`
}
