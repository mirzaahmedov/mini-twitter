package postgres

import (
	"database/sql"
	"errors"

	"twitter/internal/model"
	"twitter/internal/store"
)

type TweetStore struct {
	db *sql.DB
}

func (ps *PostgresStore) Tweet() store.TweetStore {
	return &TweetStore{ps.db}
}

func (s *TweetStore) Create(payload *model.TweetCreatePayload) (*model.Tweet, error) {
	tweet := new(model.Tweet)
	err := s.db.
		QueryRow(`--sql
			INSERT INTO tweets (content, author_id) 
			VALUES ($1, $2) 
			RETURNING id, content, author_id, attachment`,
			payload.Content, payload.AuthorID,
		).Scan(&tweet.ID, &tweet.Content, &tweet.AuthorID, &tweet.Attachment)
	if err != nil {
		return nil, err
	}
	return tweet, nil
}

func (s *TweetStore) Delete(id, userID string) error {
	result, err := s.db.Exec(`--sql
		DELETE FROM tweets 
		WHERE id = $1 
		AND author_id = $2`,
		id, userID)
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

func (s *TweetStore) Update(id, userID string, payload *model.Tweet) (*model.Tweet, error) {
	tweet := new(model.Tweet)

	err := s.db.QueryRow(`--sql
		UPDATE tweets 
		SET content = $1 
		WHERE id=$2 
		AND author_id=$3 
		RETURNING id, content, author_id, attachment`,
		payload.Content, id, userID,
	).Scan(&tweet.ID, &tweet.Content, &tweet.AuthorID, &tweet.Attachment)
	if err != nil {
		return nil, err
	}

	return tweet, nil
}

func (s *TweetStore) FindByID(id string) (*model.Tweet, error) {
	tweet := new(model.Tweet)

	err := s.db.QueryRow(`--sql
		SELECT id, content, author_id, attachment 
		FROM tweets 
		WHERE id = $1`,
		id).Scan(&tweet.ID, &tweet.Content, &tweet.AuthorID, &tweet.Attachment)
	if err != nil {
		return nil, err
	}

	return tweet, nil
}

func (s *TweetStore) GetTweetsByUser(userID string) ([]model.Tweet, error) {
	tweets := make([]model.Tweet, 0)

	rows, err := s.db.Query(`--sql
		SELECT t.id, t.content, t.author_id, t.attachment, t.created_at, t.updated_at, u.name, u.username, u.profile_picture  
		FROM tweets t
		INNER JOIN users u
		ON u.id = t.author_id 
		WHERE t.author_id = $1`,
		userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		tweet := new(model.Tweet)
		author := new(model.User)
		err = rows.Scan(&tweet.ID, &tweet.Content, &tweet.AuthorID, &tweet.Attachment, &tweet.CreatedAt, &tweet.UpdatedAt, &author.Name, &author.Username, &author.ProfilePicture)
		if err != nil {
			return nil, err
		}
		tweet.Author = *author
		tweets = append(tweets, *tweet)
	}

	return tweets, nil
}

func (s *TweetStore) GetTweetsFromFollowedUsers(userID string) ([]model.Tweet, error) {
	tweets := make([]model.Tweet, 0)

	rows, err := s.db.Query(`
		SELECT t.id, t.content, t.author_id, t.attachment, t.created_at, t.updated_at, u.name, u.username, u.profile_picture 
		FROM tweets t
		INNER JOIN users u
		ON u.id = t.author_id 
		WHERE t.author_id != $1 
		AND t.author_id 
		IN (
			SELECT following_id 
			FROM follows 
			WHERE follower_id = $1
		) LIMIT 30`,
		userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		tweet := new(model.Tweet)
		author := new(model.User)
		err := rows.Scan(&tweet.ID, &tweet.Content, &tweet.AuthorID, &tweet.Attachment, &tweet.CreatedAt, &tweet.UpdatedAt, &author.Name, &author.Username, &author.ProfilePicture)
		if err != nil {
			return nil, err
		}
		tweet.Author = *author
		tweets = append(tweets, *tweet)
	}

	return tweets, nil
}
