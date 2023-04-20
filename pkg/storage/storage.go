package storage

const NEWS_PAGE_LIMIT = 15

// News - новость.
type News struct {
	Id      uint64
	GUID    string
	Title   string
	Content string
	PubTime uint64
	Link    string
}

// Comment - комментарий.
type Comment struct {
	Id       uint64
	NewsId   uint64
	Content  string
	ParentId uint64
}

// RequestLog - лог запроса.
type RequestLog struct {
	Id         uint64
	Ip         string
	Timestamp  uint64
	StatusCode string
	RequestId  string
}

type NewsInterface interface {
	NewsDetail(uint64) (News, error)
	News(int, string) ([]News, int, error)
}

type CommentsInterface interface {
	AddComment(Comment) (int, error)
	Comments(uint64) ([]Comment, error)
}

type RequestLogsInterface interface {
	AddLog(RequestLog)
}

// Конструктор новости
func NewNews(guid string, title string, content string, pubTime uint64, link string) News {
	return News{
		0,
		guid,
		title,
		content,
		pubTime,
		link,
	}
}

// Конструктор комментария
func NewComments(newsId uint64, content string, parentId uint64) Comment {
	return Comment{
		0,
		newsId,
		content,
		parentId,
	}
}

func NewRequestLog(ip string, timestamp uint64, statusCode string, requestId string) RequestLog {
	return RequestLog{
		0,
		ip,
		timestamp,
		statusCode,
		requestId,
	}
}
