package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
	"go-api/common"
	"io"
	"net/http/httptest"
	"testing"
)

var (
	app = fiber.New()
	c   = app.AcquireCtx(&fasthttp.RequestCtx{})
)

func setTestConfig() {
	settings.RepositoryConfig = "../testdata/repositories.yaml"
	settings.RepositoryCache = "../testdata/cache"
}

func TestAddRepo(t *testing.T) {
	setTestConfig()
	var result common.ResultStatus
	mcPostBody := map[string]interface{}{
		"name": "bitnami",
		"url":  "https://charts.bitnami.com/bitnami",
	}
	postBody, _ := json.Marshal(mcPostBody)
	app.Post("/api/repositories", AddRepo)
	req := httptest.NewRequest("POST", "/api/repositories", bytes.NewReader(postBody))
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	res, _ := app.Test(req, -1)
	body, _ := io.ReadAll(res.Body)
	json.Unmarshal(body, &result)
	assert.Equal(t, 200, result.HttpStatusCode)
}

func TestListRepoCharts(t *testing.T) {
	var result common.ResultStatus
	setTestConfig()

	mcPostBody := map[string]interface{}{
		"name": "bitnami",
		"url":  "https://charts.bitnami.com/bitnami",
	}
	postBody, _ := json.Marshal(mcPostBody)
	app.Post("/api/repositories", AddRepo)
	req := httptest.NewRequest("POST", "/api/repositories", bytes.NewReader(postBody))

	app.Get("/api/repositories/:repositories/charts", ListRepoCharts)
	req = httptest.NewRequest("GET", "/api/repositories/bitnami/charts", nil)
	res, err := app.Test(req, -1)
	body, _ := io.ReadAll(res.Body)
	json.Unmarshal(body, &result)
	assert.Nil(t, err)
	assert.Equal(t, 200, result.HttpStatusCode)
}

func TestListRepos(t *testing.T) {
	var result common.ResultStatus
	setTestConfig()
	app.Get("/api/repositories", ListRepos)
	req := httptest.NewRequest("GET", "/api/repositories", nil)
	res, err := app.Test(req, -1)
	body, _ := io.ReadAll(res.Body)
	json.Unmarshal(body, &result)
	assert.Nil(t, err)
	assert.Equal(t, 200, result.HttpStatusCode)
}

func TestRemoveRepo(t *testing.T) {
	var result common.ResultStatus
	setTestConfig()
	mcPostBody := map[string]interface{}{
		"name": "bitnami",
		"url":  "https://charts.bitnami.com/bitnami",
	}
	postBody, _ := json.Marshal(mcPostBody)
	app.Post("/api/repositories", AddRepo)
	req := httptest.NewRequest("POST", "/api/repositories", bytes.NewReader(postBody))

	app.Delete("/api/repositories/:repositories", RemoveRepo)
	req = httptest.NewRequest("DELETE", "/api/repositories/bitnami", nil)
	res, err := app.Test(req, -1)
	body, _ := io.ReadAll(res.Body)
	json.Unmarshal(body, &result)
	assert.Nil(t, err)
	assert.Equal(t, 200, result.HttpStatusCode)
}

func TestUpdateRepo(t *testing.T) {
	var result common.ResultStatus
	setTestConfig()

	app.Put("/api/repositories/:repositories", UpdateRepo)
	req := httptest.NewRequest("PUT", "/api/repositories/grafana", nil)
	res, err := app.Test(req, -1)
	body, _ := io.ReadAll(res.Body)
	json.Unmarshal(body, &result)
	assert.Nil(t, err)
	assert.Equal(t, 200, result.HttpStatusCode)
}

func Test_syncRepoLock(t *testing.T) {
	type args struct {
		repoFile string
	}
	tests := []struct {
		name    string
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.wantErr(t, syncRepoLock(tt.args.repoFile), fmt.Sprintf("syncRepoLock(%v)", tt.args.repoFile))
		})
	}
}
