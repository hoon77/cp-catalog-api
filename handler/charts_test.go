package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"go-api/common"
	"helm.sh/helm/v3/cmd/helm/search"
	"io"
	"net/http/httptest"
	"testing"
)

func TestGetChartInfo(t *testing.T) {
	setTestConfig()
	var result common.ResultStatus
	app.Get("/api/repositories/:repositories/charts/:charts/info", GetChartInfo)
	req := httptest.NewRequest("GET", "/api/repositories/grafana/charts/grafana/info", nil)

	res, _ := app.Test(req, -1)
	body, _ := io.ReadAll(res.Body)
	json.Unmarshal(body, &result)
	assert.Equal(t, 200, result.HttpStatusCode)
}
func TestGetChartInfoValues(t *testing.T) {
	setTestConfig()
	var result common.ResultStatus

	app.Get("/api/repositories/:repositories/charts/:charts/info", GetChartInfo)
	req := httptest.NewRequest("GET", "/api/repositories/grafana/charts/grafana/info?info=values", nil)

	res, _ := app.Test(req, -1)
	body, _ := io.ReadAll(res.Body)
	json.Unmarshal(body, &result)
	assert.Equal(t, 200, result.HttpStatusCode)
}

func TestGetChartInfoInfoReadme(t *testing.T) {
	setTestConfig()
	var result common.ResultStatus
	app.Get("/api/repositories/:repositories/charts/:charts/info", GetChartInfo)
	req := httptest.NewRequest("GET", "/api/repositories/grafana/charts/grafana/info?info=readme", nil)

	res, _ := app.Test(req, -1)
	body, _ := io.ReadAll(res.Body)
	json.Unmarshal(body, &result)
	assert.Equal(t, 200, result.HttpStatusCode)
}

func TestGetChartInfoInfoChart(t *testing.T) {
	setTestConfig()
	var result common.ResultStatus

	app.Get("/api/repositories/:repositories/charts/:charts/info", GetChartInfo)
	req := httptest.NewRequest("GET", "/api/repositories/grafana/charts/grafana/info?info=chart", nil)

	res, _ := app.Test(req, -1)
	body, _ := io.ReadAll(res.Body)
	json.Unmarshal(body, &result)
	assert.Equal(t, 200, result.HttpStatusCode)
}
func TestGetChartVersions(t *testing.T) {
	setTestConfig()
	var result common.ResultStatus
	mcPostBody := map[string]interface{}{
		"name": "repo",
		"url":  "https://charts.bitnami.com/bitnami",
	}
	postBody, _ := json.Marshal(mcPostBody)
	app.Post("/api/repositories", AddRepo)
	req := httptest.NewRequest("POST", "/api/repositories", bytes.NewReader(postBody))
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	app.Get("/api/charts/:charts/versions", GetChartVersions)
	req = httptest.NewRequest("GET", "/api/charts/nginx/versions?repo=bitnami&version=", nil)

	res, _ := app.Test(req, -1)
	body, _ := io.ReadAll(res.Body)
	json.Unmarshal(body, &result)
	assert.Equal(t, 200, result.HttpStatusCode)
}

func Test_applyConstraint(t *testing.T) {
	type args struct {
		version  string
		versions bool
		res      []*search.Result
	}
	tests := []struct {
		name    string
		args    args
		want    []*search.Result
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := applyConstraint(tt.args.version, tt.args.versions, tt.args.res)
			if !tt.wantErr(t, err, fmt.Sprintf("applyConstraint(%v, %v, %v)", tt.args.version, tt.args.versions, tt.args.res)) {
				return
			}
			assert.Equalf(t, tt.want, got, "applyConstraint(%v, %v, %v)", tt.args.version, tt.args.versions, tt.args.res)
		})
	}
}

func Test_buildSearchIndex(t *testing.T) {
	type args struct {
		repoName string
	}
	tests := []struct {
		name    string
		args    args
		want    *search.Index
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := buildSearchIndex(tt.args.repoName)
			if !tt.wantErr(t, err, fmt.Sprintf("buildSearchIndex(%v)", tt.args.repoName)) {
				return
			}
			assert.Equalf(t, tt.want, got, "buildSearchIndex(%v)", tt.args.repoName)
		})
	}
}
