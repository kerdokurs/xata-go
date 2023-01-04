package xatago

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTableImpl_CreateAndDeleteSimple(t *testing.T) {
	assert := assert.New(t)
	item := post{
		Title:    "New Post!",
		Subtitle: "Post subtitle",
		Likes:    10,
	}
	id, err := client.Posts.Create(&item)
	assert.Nil(err, "Creating item should not return an error")
	assert.NotEmpty(id, "Creating item should return its id")

	post, err := client.Posts.
		Filter("id", Is, id).
		GetFirst()
	assert.Nil(err, "Filtering the created item should not return an error")
	assert.NotNil(post, "Filtering the created item should return it")
	assert.Equal(id, post.ID, "IDs of the items should match")

	err = client.Posts.Delete(id)
	assert.Nil(err, "Deleting item should not return error")
}

func TestTableImpl_CreateAndDeleteComplex(t *testing.T) {
	assert := assert.New(t)
	author := user{
		Name: "Author",
	}
	id, err := client.Users.Create(&author)
	assert.Nil(err, "Creating should work")
	assert.NotEmpty(id, "ID should be returned")

	item := post{
		Title:    "Post with user",
		Subtitle: "Nothing",
		Likes:    0,
		User: &user{
			ID: id,
		},
	}

	id, err = client.Posts.Create(&item)
	assert.Nil(err, "Creating a linked post should work")
	assert.NotEmpty(id, "ID to the post should be returned")

	createdItem, err := client.Posts.Select("*", "user.*").
		Filter("subtitle", Is, "Nothing").
		GetFirst()
	assert.Nil(err, "Created post should exist")
	assert.NotNil(createdItem, "Created post should exist")
	assert.NotNil(createdItem.User, "Created post should have user")
	assert.Equal(author.Name, createdItem.User.Name, "Created post user names should match")

	err = client.Posts.Delete(createdItem.ID)
	assert.Nil(err, "Deleting the post should not return an error")
	err = client.Users.Delete(createdItem.User.ID)
	assert.Nil(err, "Deleting the user should not return an error")
}
