package persistence

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"idv/chris/config"
	"idv/chris/domain"
)

// A simple sqlx / sql wrapper return *sql.DB
func NewMariaDB(cfg *config.APPConfig, logger *zap.Logger) (*sql.DB, error) {
	dsn := cfg.Mariadb.DSN
	if dsn == "" {
		return nil, fmt.Errorf("mariadb.dsn is empty")
	}
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	// simple ping
	if err := db.Ping(); err != nil {
		return nil, err
	}
	logger.Info("mariadb connected")
	return db, nil
}

// Concrete implementation of domain.ArticleRepository using *sql.DB
type ArticleRepoMySQL struct {
	DB *sql.DB
}

func NewArticleRepoMySQL(db *sql.DB) domain.ArticleRepository {
	return &ArticleRepoMySQL{DB: db}
}

// Implement methods (minimal)
func (r *ArticleRepoMySQL) Create(a *domain.Article) (int, error) {
	res, err := r.DB.Exec("INSERT INTO articles (Title, Content, NodeID, UpdateDt) VALUES (?, ?, ?, NOW())", a.Title, a.Content, a.NodeID)
	if err != nil {
		return 0, err
	}
	id64, _ := res.LastInsertId()
	return int(id64), nil
}

func (r *ArticleRepoMySQL) GetByID(id int) (*domain.Article, error) {
	row := r.DB.QueryRow("SELECT RowID, Title, Content, NodeID, UpdateDt FROM articles WHERE RowID = ?", id)
	var a domain.Article
	if err := row.Scan(&a.ID, &a.Title, &a.Content, &a.NodeID, &a.Updated); err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *ArticleRepoMySQL) ListRecent(limit int) ([]*domain.Article, error) {
	rows, err := r.DB.Query("SELECT RowID, Title, Content, NodeID, UpdateDt FROM articles ORDER BY UpdateDt DESC LIMIT ?", limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []*domain.Article
	for rows.Next() {
		var a domain.Article
		if err := rows.Scan(&a.ID, &a.Title, &a.Content, &a.NodeID, &a.Updated); err != nil {
			return nil, err
		}
		out = append(out, &a)
	}
	return out, nil
}

func Module() fx.Option {
	return fx.Options(
		fx.Provide(NewMariaDB),
		// bind repo impl to domain.ArticleRepository
		fx.Provide(fx.Annotate(
			NewArticleRepoMySQL,
			// no result tags; provides domain.ArticleRepository
		)),
	)
}
