package handler

import (
	"fmt"
	"github.com/Masterminds/semver/v3"
	"github.com/gofiber/fiber/v2"
	"go-api/common"
	"helm.sh/helm/v3/cmd/helm/search"
	"helm.sh/helm/v3/pkg/repo"
)

// searchMaxScore suggests that any score higher than this is not considered a match.
const searchMaxScore = 25

type repoChartElement struct {
	Name       string `json:"name"`
	Version    string `json:"version"`
	AppVersion string `json:"app_version"`
}

type repoChartList []repoChartElement

// GetChartVersions
// @Summary Get Chart sVersions
// @Tags Charts
// @Accept json
// @Produce json
// @Router /api/repositories/:repositories/charts/:charts/versions [Get]
func GetChartVersions(c *fiber.Ctx) error {
	repoName := c.Params("repositories")
	charts := c.Params("charts") // search keyword
	version := c.Query("version")
	// default stable
	if version == "" {
		version = ">0.0.0"
	}

	index, err := buildSearchIndex(repoName)
	if err != nil {
		return common.RespErr(c, err)
	}

	var res []*search.Result
	res, err = index.Search(fmt.Sprintf("%s/%s[^-]", repoName, charts), searchMaxScore, true)
	if err != nil {
		return common.RespErr(c, err)
	}

	search.SortScore(res)
	data, err := applyConstraint(version, true, res)
	if err != nil {
		return common.RespErr(c, err)
	}

	chartList := make(repoChartList, 0, len(data))
	for _, v := range data {
		chartList = append(chartList, repoChartElement{
			Name:       v.Name,
			Version:    v.Chart.Version,
			AppVersion: v.Chart.AppVersion,
		})
	}

	return common.RespOK(c, chartList)
}

func buildSearchIndex(repoName string) (*search.Index, error) {
	index := search.NewIndex()

	path := fmt.Sprintf("%s/%s-index.yaml", settings.RepositoryCache, repoName)
	fmt.Println("path:", path)
	indexFile, err := repo.LoadIndexFile(path)
	if err != nil {
		return nil, fmt.Errorf("레파지토리가 존재하지 않습니다.", repoName)
	}

	index.AddRepo(repoName, indexFile, true)
	return index, nil
}

func applyConstraint(version string, versions bool, res []*search.Result) ([]*search.Result, error) {
	if len(version) == 0 {
		return res, nil
	}

	constraint, err := semver.NewConstraint(version)
	if err != nil {
		return res, fmt.Errorf("an invalid version/constraint format")
	}

	data := res[:0]
	foundNames := map[string]bool{}
	for _, r := range res {
		// if not returning all versions and already have found a result,
		// you're done!
		if !versions && foundNames[r.Name] {
			continue
		}
		v, err := semver.NewVersion(r.Chart.Version)
		if err != nil {
			continue
		}
		if constraint.Check(v) {
			data = append(data, r)
			foundNames[r.Name] = true
		}
	}

	return data, nil
}
