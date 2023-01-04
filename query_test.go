package xatago

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestQuery_GetMany(t *testing.T) {
	assert := assert.New(t)
	posts, err := client.Posts.Select("*", "user.*").GetMany()
	assert.Nil(err, "Getting posts should not return an error")

	assert.Equal(3, len(posts), "Query should return 2 elements")

	hasUserIndex := -1
	for i, post := range posts {
		if post.User != nil {
			if hasUserIndex > -1 {
				assert.Fail("Only one post should have a user")
			}

			hasUserIndex = i
		}
	}

	assert.NotEqualf(-1, hasUserIndex, "One post should have a user")
}

func TestQuery_GetFirst(t *testing.T) {
	assert := assert.New(t)
	post, err := client.Posts.Select("*").GetFirst()
	assert.Nil(err, "Getting a single post should not return an error")
	assert.NotNil(post, "Post should not be nil")
}

func TestQuery_Filter(t *testing.T) {
	assert := assert.New(t)
	post, err := client.Posts.Select("*").Filter("title", Is, "Test").GetFirst()
	assert.Nil(err, "Filtering a single post should not return an error")
	assert.NotNil(post, "Filter query should return one element")
	assert.Equal("Test", post.Title, "Post titles should match")
}

func TestQuery_FilterMultiple(t *testing.T) {
	assert := assert.New(t)
	post, err := client.Posts.Select("*").
		Filter("title", Is, "Test").
		Filter("subtitle", Is, "Subtitle").
		GetFirst()
	assert.Nil(err, "Filtering a single post should not return an error")
	assert.NotNil(post, "Filter query should return one element")
	assert.Equal("Test", post.Title, "Post titles should match")
	assert.Equal("Subtitle", post.Subtitle, "Post subtitles should match")

	likes := post.Likes

	post, err = client.Posts.Select("*").
		Filter("likes", Gt, 0).
		GetFirst()
	assert.Nil(err, "Filtering a single post should not return an error")
	assert.NotNil(post, "Filter query should return one element")
	assert.Equal(likes, post.Likes, "Post like amounts should match")

	post, err = client.Posts.Select("*").
		Filter("likes", Lt, 0).
		GetFirst()
	assert.NotNil(err, "Filtering a non-existing item should return an error")
	assert.Nil(post, "Item returned should be nil")
}
