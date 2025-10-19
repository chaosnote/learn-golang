package domain

import "time"

// 只是一個範例 domain model + repository interface
type Article struct {
	ID      int
	Title   string
	Content string
	NodeID  string
	Updated time.Time
}

type ArticleRepository interface {
	Create(a *Article) (int, error)
	GetByID(id int) (*Article, error)
	ListRecent(limit int) ([]*Article, error)
}
