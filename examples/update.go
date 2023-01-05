package examples

import (
	"fmt"
	xg "github.com/kerdokurs/xata-go"
)

func update() {
	updatedPost, _ := client.Posts.Update(&Post{
		Id:    post.Id,       // Id has to be set
		Likes: xg.SetPtr(10), // Pointer value cannot be set directly, use this utility instead
	})
	fmt.Println(updatedPost.Likes) // 10
}
