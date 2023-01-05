package examples

import (
	xg "github.com/kerdokurs/xata-go"
)

// Keep in mind that fields should have `omitempty` JSON tag to make single field updates, creations and queries work.
// Note: If using any numeric field, a pointer should be used. Otherwise, unset values will default to 0 and partial update
// when not setting the field will set the field in the database to 0.

type Post struct {
	Id    string `json:"id,omitempty"`
	Title string `json:"title,omitempty"`
	Likes *int   `json:"likes,omitempty"`
	User  *User  `json:"user,omitempty"`
}

func (p Post) ID() string {
	return p.Id
}

type User struct {
	Id   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

func (u User) ID() string {
	return u.Id
}

type MyAPIClient struct {
	*xg.Client
	Posts xg.Table[Post]
	Users xg.Table[User]
}

func buildXataClient(accessToken, databaseURL string) *MyAPIClient {
	xataClient := xg.NewClient(accessToken, databaseURL)
	myClient := &MyAPIClient{
		Client: xataClient,
		Posts:  xg.NewTableImpl[Post](xataClient, "Posts"),
		Users:  xg.NewTableImpl[User](xataClient, "Users"),
	}
	return myClient
}

var (
	client = buildXataClient("API_KEY", "DATABASE_URL")
	post   *Post
)

func Examples() {
	create()
	update()
	deletes()
	query()
}
