package repository

import (
	"context"
	"github.com/rshlin/go-blog-api-assesment/blog/model"
)

type BlogRepository interface {
	FindAll(ctx context.Context, page int, pageSize int) (*model.PaginatedPosts, error)
	FindById(ctx context.Context, id int) (*model.Post, error)
	Create(ctx context.Context, post model.Post) (*model.Post, error)
	Update(ctx context.Context, author model.Author, post model.Post) (*model.Post, error)
	Delete(ctx context.Context, author model.Author, id int) error
}
