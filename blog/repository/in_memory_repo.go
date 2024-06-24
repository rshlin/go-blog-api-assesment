package repository

import (
	"context"
	"errors"
	error2 "github.com/rshlin/go-blog-api-assesment/blog/error"
	"github.com/rshlin/go-blog-api-assesment/blog/model"
	"sync"
	"sync/atomic"
)

type Option func(repository *InMemoryBlogRepository)

type InMemoryBlogRepository struct {
	store   *sync.Map
	counter int32
}

func NewInMemoryBlogRepository(opts ...Option) *InMemoryBlogRepository {
	r := &InMemoryBlogRepository{
		store: new(sync.Map),
	}

	for _, opt := range opts {
		opt(r)
	}

	return r
}

func (r *InMemoryBlogRepository) FindAll(ctx context.Context, page int, pageSize int) (*model.PaginatedPosts, error) {
	if page < 1 || pageSize < 1 {
		return nil, errors.New("invalid page or pageSize")
	}

	// Create a slice to store all posts
	allPosts := make([]model.Post, 0)

	// Iterate over the sync.Map and add each post to the slice
	r.store.Range(func(key interface{}, value interface{}) bool {
		// Check if the context is done. If so, stop iteration.
		select {
		case <-ctx.Done():
			return false
		default:
		}

		post, ok := value.(model.Post)
		if !ok {
			return false // continue iteration in case of error
		}
		allPosts = append(allPosts, post)
		return true
	})

	// Calculate total count
	totalCount := len(allPosts)

	// Calculate start and end indices for the slice
	startIndex := (page - 1) * pageSize
	if startIndex >= totalCount {
		return &model.PaginatedPosts{
			PageMetadata: model.PageMetadata{
				Page:  page,
				Size:  pageSize,
				Total: totalCount,
			},
			Data: make([]model.Post, 0),
		}, nil
	}
	endIndex := startIndex + pageSize
	if endIndex > totalCount {
		endIndex = totalCount
	}

	// Slice the allPosts slice to get the posts for the requested page
	pagePosts := allPosts[startIndex:endIndex]

	return &model.PaginatedPosts{
		PageMetadata: model.PageMetadata{
			Page:  page,
			Size:  pageSize,
			Total: totalCount,
		},
		Data: pagePosts,
	}, nil
}

func (r *InMemoryBlogRepository) FindById(ctx context.Context, id int) (*model.Post, error) {
	value, ok := r.store.Load(id)
	if !ok {
		return nil, error2.NewError(error2.NotFound, "post not found", nil)
	}
	post, ok := value.(model.Post)
	if !ok {
		return nil, errors.New("failed to type assert post")
	}
	return &post, nil
}

func (r *InMemoryBlogRepository) Create(ctx context.Context, post model.Post) (*model.Post, error) {
	post.Id = int(atomic.AddInt32(&r.counter, 1)) // Atomically increment the counter and assign it as the post's ID.
	r.store.Store(post.Id, post)
	return &post, nil
}

func (r *InMemoryBlogRepository) Update(ctx context.Context, author model.Author, post model.Post) (*model.Post, error) {
	value, ok := r.store.Load(post.Id)
	if !ok {
		return nil, error2.NewError(error2.NotFound, "post not found", nil)
	}
	existingPost, ok := value.(model.Post)
	if !ok {
		return nil, errors.New("failed to type assert post")
	}
	if existingPost.Author != author {
		return nil, error2.NewError(error2.NotAuthorized, "not authorized", nil)
	}
	r.store.Store(post.Id, post)
	return &post, nil
}

func (r *InMemoryBlogRepository) Delete(ctx context.Context, author model.Author, id int) error {
	value, ok := r.store.Load(id)
	if !ok {
		return error2.NewError(error2.NotFound, "post not found", nil)
	}
	existingPost, ok := value.(model.Post)
	if !ok {
		return errors.New("failed to type assert post")
	}
	if existingPost.Author != author {
		return error2.NewError(error2.NotAuthorized, "not authorized", nil)
	}
	r.store.Delete(id)
	return nil
}
