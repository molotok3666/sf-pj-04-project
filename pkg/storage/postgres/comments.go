package postgres

import (
	"APIGateway/pkg/storage"
	"context"
)

const CHECK_LIMIT = 10

func (s *DbStorage) AddComment(c storage.Comment) (int, error) {
	var id int
	err := s.db.QueryRow(context.Background(), `
			INSERT INTO comments (news_id, content, parent_id)
			VALUES ($1, $2, $3)
			RETURNING id;
		`,
		c.NewsId,
		c.Content,
		c.ParentId,
	).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *DbStorage) Comments(newsId uint64) ([]storage.Comment, error) {
	rows, err := s.db.Query(context.Background(), `
			SELECT *
			FROM comments
			WHERE news_id = $1
			ORDER BY id desc
			LIMIT $1
		`,
		newsId,
	)

	if err != nil {
		return nil, err
	}

	var comments []storage.Comment
	for rows.Next() {
		var c storage.Comment
		err = rows.Scan(
			&c.Id,
			&c.NewsId,
			&c.Content,
			&c.ParentId,
		)
		if err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}

	return comments, nil
}
