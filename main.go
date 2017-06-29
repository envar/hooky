package main

import (
	"net/http"
	"os"

	"log"
)

func main() {
	// fetch configuration from environment
	httpAddr := os.Getenv("HTTP_ADDR")
	if httpAddr == "" {
		httpAddr = ":8080"
	}
	secret := os.Getenv("GITHUB_SECRET")
	if secret == "" {
		log.Fatal("GITHUB_SECRET must be set and non-empty")
	}
	owner := os.Getenv("GITHUB_OWNER")
	if owner == "" {
		log.Fatal("GITHUB_OWNER must be set and non-empty")
	}
	repo := os.Getenv("GITHUB_REPO")
	if owner == "" {
		log.Fatal("GITHUB_REPO must be set and non-empty")
	}
	token := os.Getenv("GITHUB_API_TOKEN")
	if token == "" {
		log.Fatal("GITHUB_API_TOKEN must be set and non-empty")
	}

	// create a github client
	client := NewGithubClient(token, secret, owner, repo)

	// setup routes
	http.Handle("/gh", &GithubWebhookHandler{client, nil})
	http.Handle("/", &RepoHandler{client})

	// start the server
	log.Printf("listening on %s...\n", httpAddr)
	log.Fatal(http.ListenAndServe(httpAddr, nil))
}
