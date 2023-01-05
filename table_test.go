package xatago

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTableImpl_CreateAndDeleteSimple(t *testing.T) {
	assert := assert.New(t)

	item, err := client.Posts.Create(&post{
		Title:    "New Post!",
		Subtitle: "Post subtitle",
		Likes:    10,
	})
	assert.Nil(err, "Creating item should not return an error")
	assert.NotNil(item, "Creating item should return the created object")
	assert.NotEmpty(item.ID(), "Creating item should return its id")

	post, err := client.Posts.
		Filter("id", Is, item.ID()).
		GetFirst()
	assert.Nil(err, "Filtering the created item should not return an error")
	assert.NotNil(post, "Filtering the created item should return it")
	assert.Equal(item.ID(), post.ID(), "IDs of the items should match")

	err = client.Posts.Delete(post)
	assert.Nil(err, "Deleting item should not return error")
}

func TestTableImpl_CreateAndDeleteComplex(t *testing.T) {
	assert := assert.New(t)

	author, err := client.Users.Create(&user{
		Name: "Author",
	})
	assert.Nil(err, "Creating should work")
	assert.NotNil(author, "Object should be returned")
	assert.NotEmpty(author.ID(), "ID should be set")

	item, err := client.Posts.Create(&post{
		Title:    "Post with user",
		Subtitle: "Nothing",
		Likes:    0,
		User: &user{
			Id: author.ID(),
		},
	})
	assert.Nil(err, "Creating a linked post should work")
	assert.NotNil(item, "Object should be returned")
	assert.NotEmpty(item.ID(), "ID to the post should be returned")

	createdItem, err := client.Posts.Select("*", "user.*").
		Filter("subtitle", Is, "Nothing").
		GetFirst()
	assert.Nil(err, "Created post should exist")
	assert.NotNil(createdItem, "Created post should exist")
	assert.NotNil(createdItem.User, "Created post should have user")
	assert.Equal(author.Name, createdItem.User.Name, "Created post user names should match")

	err = client.Posts.Delete(createdItem)
	assert.Nil(err, "Deleting the post should not return an error")
	err = client.Users.DeleteById(createdItem.User.ID())
	assert.Nil(err, "Deleting the user should not return an error")
}

func TestTableImpl_Update(t *testing.T) {
	assert := assert.New(t)

	item, err := client.Posts.Create(&post{
		Title:    "Updatable post",
		Subtitle: "Nothing here",
	})
	assert.Nil(err, "Creating should not return an error")
	assert.NotNil(item, "Creating should return ID")

	updatedItem, err := client.Posts.Update(&post{
		Id:       item.ID(),
		Subtitle: "New subtitle!",
	})
	assert.Nil(err, "Updating item should not return an error")
	assert.NotNil(updatedItem, "Updated item should be returned")
	assert.Equal(item.ID(), updatedItem.ID(), "Updated item should have the same ID")
	assert.Equal(item.Title, updatedItem.Title, "Unchanged fields should remain the same")
	assert.Equal("New subtitle!", updatedItem.Subtitle, "Updated fields should be reflected in the result")

	err = client.Posts.Delete(updatedItem)
	assert.Nil(err, "Deleting the updated item should not return an error")
}
