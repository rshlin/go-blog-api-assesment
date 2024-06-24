package utils

import (
	"reflect"
	"testing"

	"github.com/rshlin/go-blog-api-assesment/api"
	blogModel "github.com/rshlin/go-blog-api-assesment/blog/model"
)

// TestMapToStruct tests the MapToStruct function
func TestMapToStruct(t *testing.T) {
	// setup domain model data
	domPosts := blogModel.PaginatedPosts{
		PageMetadata: blogModel.PageMetadata{
			Page:  1,
			Size:  10,
			Total: 100,
		},
		Data: []blogModel.Post{
			{
				Author:  "John Doe",
				Content: "This is a test post",
				Id:      1,
				Title:   "Test Post",
			},
			{
				Author:  "Jane Doe",
				Content: "This is another test post",
				Id:      2,
				Title:   "Test Post 2",
			},
		},
	}

	// setup expected api model data
	expectedApiPosts := api.PaginatedPosts{
		PageMetadata: &struct {
			Page  int `json:"page"`
			Size  int `json:"size"`
			Total int `json:"total"`
		}{
			Page:  1,
			Size:  10,
			Total: 100,
		},
		Data: []api.Post{
			{
				Author:  "John Doe",
				Content: "This is a test post",
				Id:      1,
				Title:   "Test Post",
			},
			{
				Author:  "Jane Doe",
				Content: "This is another test post",
				Id:      2,
				Title:   "Test Post 2",
			},
		},
	}

	// map domain model to api model
	var apiPosts api.PaginatedPosts
	err := MapToStruct(domPosts, &apiPosts)
	if err != nil {
		t.Errorf("Failed to map domain model to api model: %v", err)
	}

	// compare the mapped api model to the expected api model
	if !reflect.DeepEqual(apiPosts, expectedApiPosts) {
		t.Errorf("Mapped api model does not match expected api model")
	}
}
