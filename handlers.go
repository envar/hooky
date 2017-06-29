package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	mgo "gopkg.in/mgo.v2"
)

const (
	defaultBranch = "master"
)

type GithubWebhookHandler struct {
	Client *GithubClient
	DB     *mgo.Database
}

func (h *GithubWebhookHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// TODO: verify header X-Hub-Signature

	var payload GithubPushPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// TODO: check ref against repo branch
	// TODO: consider using bulk operations
	// TODO: implement add and delete operations

	for _, path := range payload.HeadCommit.Modified {
		name := GetRuleSetName(path)
		if name == "" {
			continue
		}

		fmt.Printf("updating: %s", name)

		contents, _, err := h.Client.GetFileContent(payload.HeadCommit.ID, path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// TODO: deploy updated document to db
		selector := map[string]interface{}{
			"name": name,
		}
		var update interface{}
		if err := json.Unmarshal([]byte(contents), &update); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := h.DB.C("rulesets").Update(selector, update); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}

type RepoHandler struct {
	Client *GithubClient
}

func (h *RepoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if path == "" {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	path = path[1:len(path)] // strip off leading /

	// TODO: allow ref to be specified in query parameters?
	ref := defaultBranch

	// TODO: make sure users can only push to authorized branches
	// TODO: make directory if not exists

	switch r.Method {
	case http.MethodGet:
		// retrieve file from github and return file contents
		contents, _, err := h.Client.GetFileContent(ref, path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write([]byte(contents))
		w.WriteHeader(http.StatusOK)
	case http.MethodPut:
		message := fmt.Sprintf("updated %s", path)

		_, sha, err := h.Client.GetFileContent(ref, path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		contents, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := h.Client.UpdateFileContent(path, message, sha, ref, contents); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	case http.MethodPost:
		message := fmt.Sprintf("added %s", path)

		contents, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := h.Client.CreateFileContent(path, message, contents, ref); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	case http.MethodDelete:
		message := fmt.Sprintf("deleted %s", path)

		_, sha, err := h.Client.GetFileContent(ref, path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := h.Client.DeleteFile(path, message, sha, ref); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}
