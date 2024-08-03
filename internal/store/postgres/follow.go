package postgres

import (
	"database/sql"
	"errors"

	"twitter/internal/model"
	"twitter/internal/store"
)

type FollowStore struct {
	db *sql.DB
}

func (ps *PostgresStore) Follow() store.FollowStore {
	return &FollowStore{ps.db}
}

func (s *FollowStore) Create(payload *model.Follow) (*model.Follow, error) {
	follow := new(model.Follow)
	err := s.db.
		QueryRow(`
			INSERT INTO follows (follower_id, following_id) 
			VALUES ($1, $2) 
			RETURNING follower_id, following_id`, payload.FollowerID, payload.FollowingID,
			).Scan(&follow.FollowerID, &follow.FollowingID)
	if err != nil {
		return nil, err
	}
	return follow, nil
}

func (s *FollowStore) Delete(followerID, followingID string) error {
	result, err := s.db.Exec(`--sql 
		DELETE FROM follows 
		WHERE follower_id = $1 
		AND following_id = $2`, 
		followerID, followingID,
	)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("no rows effected")
	}
	return nil
}

func (s *FollowStore) FindAllWithFollowerID(id string) ([]model.User, error) {
	follows := make([]model.User, 0)

	rows, err := s.db.Query(`
		SELECT u.id, u.name, u.username, u.bio, u.profile_picture, u.created_at 
		FROM follows f
		LEFT JOIN users u
		ON u.id = f.following_id 
		WHERE f.follower_id = $1`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		user := new(model.User)
		err := rows.Scan(&user.ID, &user.Name, &user.Username, &user.Bio, &user.ProfilePicture, &user.CreatedAt)
		if err != nil {
			return nil, err
		}
		follows = append(follows, *user)
	}

	return follows, nil
}
