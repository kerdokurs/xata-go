package xatago

import (
	"os"
)

type post struct {
	Record
	Title    string `json:"title"`
	Subtitle string `json:"subtitle"`
	Likes    int    `json:"likes"`

	User *user `json:"user,omitempty"`
}

type user struct {
	Record
	Name string `json:"name"`
}

type missing struct {
	Record
}

type apiClient struct {
	*Client
	Posts   Table[post]
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
		Posts:   NewTableImpl[post](xataClient, "Posts"),
		Missing: NewTableImpl[missing](xataClient, "Missing"),
	}

	return client
}
