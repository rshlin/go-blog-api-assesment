package repository

import (
	"context"
	"github.com/rshlin/go-blog-api-assesment/blog/model"
	"testing"
)

func TestInMemoryBlogRepository(t *testing.T) {
	repo := NewInMemoryBlogRepository()

	// Test Create
	post := model.Post{Title: "Test Title", Content: "Test Content", Author: "Test Author"}
	createdPost, err := repo.Create(context.Background(), post)
	if err != nil {
		t.Fatalf("failed to Create post: %v", err)
	}
	if createdPost.Id == 0 {
		t.Fatalf("expected post ID to be set, got %v", createdPost.Id)
	}

	// Test find by ID
	foundPost, err := repo.FindById(context.Background(), createdPost.Id)
	if err != nil {
		t.Fatalf("failed to find post: %v", err)
	}
	if foundPost.Title != "Test Title" {
		t.Fatalf("expected post title to be 'Test Title', got %v", foundPost.Title)
	}

	// Test Update by another author
	updatedPost := *foundPost
	updatedPost.Title = "Updated Title"
	_, err = repo.Update(context.Background(), "Another Author", updatedPost)
	if err == nil {
		t.Fatalf("expected not authorized error, got nil")
	}

	// Test Update by the same author
	_, err = repo.Update(context.Background(), "Test Author", updatedPost)
	if err != nil {
		t.Fatalf("failed to Update post: %v", err)
	}
	foundPost, err = repo.FindById(context.Background(), createdPost.Id)
	if err != nil {
		t.Fatalf("failed to find post: %v", err)
	}
	if foundPost.Title != "Updated Title" {
		t.Fatalf("expected post title to be 'Updated Title', got %v", foundPost.Title)
	}

	// Test Delete by another author
	err = repo.Delete(context.Background(), "Another Author", createdPost.Id)
	if err == nil {
		t.Fatalf("expected not authorized error, got nil")
	}

	// Test Delete by the same author
	err = repo.Delete(context.Background(), "Test Author", createdPost.Id)
	if err != nil {
		t.Fatalf("failed to Delete post: %v", err)
	}

	// Test find by ID after deletion
	_, err = repo.FindById(context.Background(), createdPost.Id)
	if err == nil {
		t.Fatalf("expected post not found error, got nil")
	}

	// Test FindAll with invalid page and pageSize
	_, err = repo.FindAll(context.Background(), 0, -1)
	if err == nil {
		t.Fatalf("expected invalid page or pageSize error, got nil")
	}
}
