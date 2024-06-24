package model

type PageMetadata struct {
	Page  int
	Size  int
	Total int
}

type Author = string

type Post struct {
	Author  Author
	Content string
	Id      int
	Title   string
}

type PaginatedPosts struct {
	PageMetadata PageMetadata
	Data         []Post
}
