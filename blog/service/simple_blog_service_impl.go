package service

import (
	"context"
	"github.com/rshlin/go-blog-api-assesment/blog/model"
	"github.com/rshlin/go-blog-api-assesment/blog/repository"
)

type SimpleBlogServiceImpl struct {
	repository repository.BlogRepository
}

func NewSimpleBlogService(repository repository.BlogRepository) *SimpleBlogServiceImpl {
	return &SimpleBlogServiceImpl{
		repository: repository,
	}
}

func (s *SimpleBlogServiceImpl) FindAll(ctx context.Context, page int, pageSize int) (*model.PaginatedPosts, error) {
	return s.repository.FindAll(ctx, page, pageSize)
}

func (s *SimpleBlogServiceImpl) FindById(ctx context.Context, id int) (*model.Post, error) {
	return s.repository.FindById(ctx, id)
}

func (s *SimpleBlogServiceImpl) Create(ctx context.Context, post model.Post) (*model.Post, error) {
	return s.repository.Create(ctx, post)
}

func (s *SimpleBlogServiceImpl) Update(ctx context.Context, author model.Author, post model.Post) (*model.Post, error) {
	return s.repository.Update(ctx, author, post)
}

func (s *SimpleBlogServiceImpl) Delete(ctx context.Context, author model.Author, id int) error {
	return s.repository.Delete(ctx, author, id)
}
