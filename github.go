package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"
	"zimg/config"
	"zimg/utils"
)

func uploadGithub(filePath string) string {
	owner, ok := config.Get(config.GithubOwner)
	if !ok {
		utils.Boom(errors.New(fmt.Sprintf("[%s] required", config.GithubOwner)))
	}
	repo, ok := config.Get(config.GithubRepo)
	if !ok {
		utils.Boom(errors.New(fmt.Sprintf("[%s] required", config.GithubRepo)))
	}
	bucket, ok := config.Get(config.GithubBucket)
	if !ok {
		bucket = "default"
	}
	token, ok := config.Get(config.GithubToken)
	if !ok {
		utils.Boom(errors.New(fmt.Sprintf("[%s] required", config.GithubToken)))
	}

	content := utils.GetBytes(filePath)
	path := fmt.Sprintf("%s/%s%s", bucket, utils.ToMd5(content), filepath.Ext(filePath))

	payload, err := json.Marshal(
		map[string]string{
			"message": "upload",
			"content": utils.ToBase64(content),
		},
	)
	utils.Boom(err)
	req, err := http.NewRequest(
		"PUT",
		fmt.Sprintf("https://api.github.com/repos/%s/%s/contents/%s", owner, repo, path),
		strings.NewReader(string(payload)),
	)
	utils.Boom(err)
	req.Header.Add("Accept", "application/vnd.github+json")
	req.Header.Add("X-GitHub-Api-Version", "2022-11-28")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+token)

	client := &http.Client{}
	res, err := client.Do(req)
	utils.Boom(err)
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	utils.Boom(err)

	var m map[string]any
	err = json.Unmarshal(body, &m)
	if err != nil {
		utils.Boom(err)
	}

	if m["status"] == "422" {
		return fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/main/%s", owner, repo, path)
	} else if m["status"] == "401" {
		utils.Boom(errors.New(m["message"].(string)))
	}

	return m["content"].(map[string]any)["download_url"].(string)
}
