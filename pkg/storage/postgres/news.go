package postgres

import (
	"APIGateway/pkg/storage"
	"context"
)

// Получет новости из БД
func (s *DbStorage) News(page int, search string) ([]storage.News, int, error) {
	var news []storage.News
	var pagesTotal int
	var err error
	if search == "" {
		news, pagesTotal, err = s.news(page)
	} else {
		news, pagesTotal, err = s.newsWithFilter(page, search)
	}

	if err != nil {
		return nil, 0, err
	}

	return news, pagesTotal, nil
}

func (s *DbStorage) news(page int) ([]storage.News, int, error) {
	var news []storage.News
	offset := (page - 1) * storage.NEWS_PAGE_LIMIT
	rows, err := s.db.Query(context.Background(), `
			SELECT *
			FROM news
			ORDER BY pub_time
			OFFSET $1
			LIMIT $2
		`,
		offset,
		storage.NEWS_PAGE_LIMIT,
	)

	if err != nil {
		return news, 0, nil
	}

	for rows.Next() {
		var n storage.News
		err = rows.Scan(
			&n.Id,
			&n.GUID,
			&n.Title,
			&n.Content,
			&n.PubTime,
			&n.Link,
		)

		if err != nil {
			return nil, 0, err
		}
		news = append(news, n)
	}

	row := s.db.QueryRow(context.Background(), `
			SELECT COUNT(id)
			FROM news
		`,
	)

	var pagesTotal int
	err = row.Scan(&pagesTotal)
	if err != nil {
		return news, 0, err
	}

	return news, pagesTotal, nil
}

func (s *DbStorage) newsWithFilter(page int, search string) ([]storage.News, int, error) {
	offset := (page - 1) * storage.NEWS_PAGE_LIMIT
	rows, err := s.db.Query(context.Background(), `
			SELECT *
			FROM news
			WHERE title ILIKE $1
			ORDER BY pub_time
			OFFSET $2
			LIMIT $3
		`,
		search,
		offset,
		storage.NEWS_PAGE_LIMIT,
	)

	var news []storage.News

	for rows.Next() {
		var n storage.News
		err = rows.Scan(
			&n.GUID,
			&n.Title,
			&n.Content,
			&n.PubTime,
			&n.Link,
		)

		if err != nil {
			return news, 0, err
		}
		news = append(news, n)
	}

	row := s.db.QueryRow(context.Background(), `
			SELECT COUNT(id)
			FROM news
			WHERE title ILIKE $1
		`,
	)

	var pagesTotal int
	err = row.Scan(&pagesTotal)
	if err != nil {
		return news, 0, err
	}

	return news, pagesTotal, nil
}

// Получет конкретную новость из БД
func (s *DbStorage) NewsDetail(id uint64) (storage.News, error) {
	var n storage.News
	row := s.db.QueryRow(context.Background(), `
			SELECT *
			FROM news
			WHERE id = $1
		`,
		id,
	)

	err := row.Scan(
		&n.Id,
		&n.GUID,
		&n.Title,
		&n.Content,
		&n.PubTime,
		&n.Link,
	)

	if err != nil {
		return n, err
	}

	return n, nil
}

// Добавляет новость в БД
func (s *DbStorage) AddNews(n storage.News) error {
	_, err := s.db.Exec(context.Background(), `
			INSERT INTO news (guid, title, content, pub_time, link)
			VALUES ($1, $2, $3, $4, $5)
			ON CONFLICT
			DO NOTHING
			RETURNING guid;
		`,
		n.GUID,
		n.Title,
		n.Content,
		n.PubTime,
		n.Link,
	)
	if err != nil {
		return err
	}

	return nil
}
