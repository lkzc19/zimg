package upload

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"zimg/config"
	"zimg/utils"
)

func Gitee(filePath string) string {
	owner, ok := config.Get(config.GiteeOwner)
	if !ok {
		utils.Boom(errors.New(fmt.Sprintf("[%s] required", config.GiteeOwner)))
	}
	repo, ok := config.Get(config.GiteeRepo)
	if !ok {
		utils.Boom(errors.New(fmt.Sprintf("[%s] required", config.GiteeRepo)))
	}
	bucket, ok := config.Get(config.GiteeBucket)
	if !ok {
		bucket = "default"
	}
	token, ok := config.Get(config.GiteeToken)
	if !ok {
		utils.Boom(errors.New(fmt.Sprintf("[%s] required", config.GiteeToken)))
	}

	content := utils.GetBytes(filePath)
	path := fmt.Sprintf("%s/%s%s", bucket, utils.ToMd5(content), filepath.Ext(filePath))

	formData := &bytes.Buffer{}
	writer := multipart.NewWriter(formData)
	_ = writer.WriteField("message", "zimg")
	_ = writer.WriteField("content", utils.ToBase64(content))
	_ = writer.WriteField("access_token", token)
	err := writer.Close()
	utils.Boom(err)

	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf("https://gitee.com/api/v5/repos/%s/%s/contents/%s", owner, repo, path),
		formData,
	)
	utils.Boom(err)
	req.Header.Add("Content-Type", writer.FormDataContentType())

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

	if res.Status == "400 Bad Request" {
		return fmt.Sprintf("https://gitee.com/%s/%s/raw//master/%s", owner, repo, path)
	} else if res.Status == "401 Unauthorized" {
		utils.Boom(errors.New(m["message"].(string)))
	}

	return m["content"].(map[string]any)["download_url"].(string)
}
