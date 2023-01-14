package xatago

import (
	"os"
)

type post struct {
	Id       string `json:"id,omitempty"`
	Title    string `json:"title,omitempty"`
	Subtitle string `json:"subtitle,omitempty"`
	Likes    int    `json:"likes,omitempty"`

	User *user `json:"user,omitempty"`
}

func (p post) ID() string {
	return p.Id
}

type user struct {
	Id   string `json:"id,omitempty"`
	Name string `json:"name"`
}

func (u user) ID() string {
	return u.Id
}

type missing struct {
	Id string `json:"id,omitempty"`
}

func (m missing) ID() string {
	return m.Id
}

type apiClient struct {
	*Client
	Posts   Table[post]
	Users   Table[user]
	Missing Table[missing]
}

var (
	databaseURL string
	apiKey      string

	unauthenticatedClient *apiClient
	client                *apiClient
)

func init() {
	databaseURL = os.Getenv("DB_URL")
	apiKey = os.Getenv("API_KEY")

	unauthenticatedClient = buildXataClient("", databaseURL)
	client = buildXataClient(apiKey, databaseURL)
}

func buildXataClient(accessToken string, databaseURL string) *apiClient {
	xataClient := NewClient(accessToken, databaseURL)
	client := &apiClient{
		Client:  xataClient,
		Posts:   NewTableImpl[post](xataClient, "Post"),
		Users:   NewTableImpl[user](xataClient, "User"),
		Missing: NewTableImpl[missing](xataClient, "Missing"),
	}

	return client
}
