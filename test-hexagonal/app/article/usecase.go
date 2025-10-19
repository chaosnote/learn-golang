package article

import (
	"context"

	"go.uber.org/fx"

	"idv/chris/domain"
)

type ArticleUseCase struct {
	repo domain.ArticleRepository
}

func NewArticleUseCase(repo domain.ArticleRepository) *ArticleUseCase {
	return &ArticleUseCase{repo: repo}
}

func (uc *ArticleUseCase) Create(ctx context.Context, a *domain.Article) (int, error) {
	return uc.repo.Create(a)
}

func (uc *ArticleUseCase) Get(ctx context.Context, id int) (*domain.Article, error) {
	return uc.repo.GetByID(id)
}

func Module() fx.Option {
	return fx.Options(
		fx.Provide(NewArticleUseCase),
	)
}
