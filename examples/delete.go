package examples

func deletes() {
	_ = client.Users.DeleteById(post.User.Id)
	_ = client.Posts.Delete(post)
}
