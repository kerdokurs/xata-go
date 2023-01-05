package examples

import (
	"fmt"
	xg "github.com/kerdokurs/xata-go"
)

func query() {
	// Note: Filter functions can be chained together
	query := client.Posts.
		Select("*", "user.*").
		Filter("likes", xg.Gt, 0)
	post, _ := query.GetFirst()
	fmt.Println(post.Likes) // Should print something greater than 0

	posts, _ := query.GetMany()
	fmt.Println(len(posts)) // Can contain many posts
}
