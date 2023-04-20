package postgres

import (
	"APIGateway/pkg/storage"
	"context"
)

func (s *DbStorage) AddLog(r storage.RequestLog) {
	s.db.Exec(context.Background(), `
			INSERT INTO request_logs (ip, timestamp, status_code, request_id)
			VALUES ($1, $2, $3, $4);
		`,
		r.Ip,
		r.Timestamp,
		r.StatusCode,
		r.RequestId,
	)
}
