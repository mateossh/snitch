package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gopkg.in/go-ini/ini.v1"
	"net/http"
	"os"
	"os/user"
	"path"
	"strconv"
)

// GithubCredentials contains PersonalToken for GitHub API authorization
type GithubCredentials struct {
	PersonalToken string
}

func (creds GithubCredentials) query(method, url string, jsonBody map[string]interface{}) (map[string]interface{}, error) {
	bodyBuffer := new(bytes.Buffer)
	err := json.NewEncoder(bodyBuffer).Encode(jsonBody)

	req, err := http.NewRequest(method, url, bodyBuffer)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "token "+creds.PersonalToken)
	req.Header.Add("Content-Type", "application/json")

	return QueryHTTP(req)
}

func (creds GithubCredentials) getIssue(repo string, todo Todo) (map[string]interface{}, error) {
	json, err := creds.query(
		"GET",
		// FIXME(#59): possible GitHub API injection attack
		"https://api.github.com/repos/"+repo+"/issues/"+(*todo.ID)[1:],
		nil)

	if err != nil {
		return nil, err
	}

	return json, nil
}

func (creds GithubCredentials) postIssue(repo string, todo Todo, body string) (Todo, error) {
	json, err := creds.query(
		"POST",
		"https://api.github.com/repos/"+repo+"/issues",
		map[string]interface{}{
			"title": todo.Title,
			"body":  body,
		})
	if err != nil {
		return todo, err
	}

	id := "#" + strconv.Itoa(int(json["number"].(float64)))
	todo.ID = &id

	return todo, err
}

// GithubCredentialsFromFile gets GithubCredentials from a filepath
func GithubCredentialsFromFile(filepath string) (GithubCredentials, error) {
	cfg, err := ini.Load(filepath)
	if err != nil {
		return GithubCredentials{}, err
	}

	return GithubCredentials{
		PersonalToken: cfg.Section("github").Key("personal_token").String(),
	}, nil
}

// GithubCredentialsFromToken returns a GithubCredentials from a string token
func GithubCredentialsFromToken(token string) GithubCredentials {
	return GithubCredentials{
		PersonalToken: token,
	}
}
