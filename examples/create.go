package examples

import "fmt"

func create() {
	// Set the columns you want to set
	// Keep in mind that if the column has been set to:
	// 1. `not null`, errors may be returned when values are not set
	// 2. `unique`, errors may be returned when existing values are set
	post, _ = client.Posts.Create(&Post{
		Title: "Post title!",
	})
	fmt.Println(post.Title) // Post title!

	user, _ := client.Users.Create(&User{
		Name: "Xata",
	})
	fmt.Println(user.Name) // Xata

	postWithUser, _ := client.Posts.Create(&Post{
		Title: "Post with user!",
		User: &User{
			Id: user.Id, // Only Id should be set to connect the post with a user
		},
	})
	fmt.Println(postWithUser.User.Name) // Xata
}
