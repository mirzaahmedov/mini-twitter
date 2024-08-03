package postgres

import (
	"database/sql"

	"twitter/internal/model"
	"twitter/internal/store"
)

type UserStore struct {
	db *sql.DB
}

func (ps *PostgresStore) User() store.UserStore {
	return &UserStore{ps.db}
}

func (s *UserStore) Create(payload *model.UserCreatePayload) (*model.User, error) {
	user := new(model.User)
	err := s.db.
		QueryRow(`
			INSERT INTO users (name, username, bio, password, profile_picture) 
			VALUES ($1, $2, $3, $4, $5) 
			RETURNING id, name, username, bio, password, profile_picture, created_at`, 
			payload.Name, payload.Username, payload.Bio, payload.Password, payload.ProfilePicture,
		).Scan(&user.ID, &user.Name, &user.Username, &user.Bio, &user.Password, &user.ProfilePicture, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserStore) FindByID(id string) (*model.User, error) {
	user := new(model.User)
	err := s.db.
		QueryRow(`
			SELECT id, name, username, bio, password, profile_picture, created_at 
			FROM users 
			WHERE id = $1`, 
			id).Scan(&user.ID, &user.Name, &user.Username, &user.Bio, &user.Password, &user.ProfilePicture, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserStore) FindByUsername(username string) (*model.User, error) {
	user := new(model.User)
	err := s.db.
		QueryRow(`
			SELECT id, name, username, bio, password, profile_picture, created_at 
			FROM users 
			WHERE username = $1`, 
			username).Scan(&user.ID, &user.Name, &user.Username, &user.Bio, &user.Password, &user.ProfilePicture, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserStore) Search(term, userID string) ([]model.User, error) {
	users := make([]model.User, 0)
	rows, err := s.db.
		Query(`
			SELECT id, name, username, bio, profile_picture, created_at 
			FROM users 
			WHERE id != $1 
			AND id 
			NOT IN (
				SELECT following_id 
				FROM follows f 
				WHERE f.follower_id = $1) 
				AND (username ILIKE '%' || $2 || '%' OR name ILIKE '%' || $2 || '%') 
				LIMIT 20`, 
				userID, 
				term)
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
		users = append(users, *user)
	}

	return users, nil
}
