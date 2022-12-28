package xatago

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestQuery_GetMany(t *testing.T) {
	assert := assert.New(t)
	posts, err := client.Posts.Select("*", "user.*").GetMany()
	assert.Nil(err, "Getting posts should not return an error")

	assert.Equal(2, len(posts), "Query should return 2 elements")

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
