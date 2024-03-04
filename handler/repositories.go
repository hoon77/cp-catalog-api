package handler

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
	"go-api/common"
	"helm.sh/helm/v3/pkg/helmpath"
	"helm.sh/helm/v3/pkg/repo"
	"os"
	"path/filepath"
)

type repositoryElement struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// AddRepo
// @Summary Add Repository
// @Tags Repository
// @Accept json
// @Produce json
// @Router /api/repositories [Post]
func AddRepo(c *fiber.Ctx) error {

	return c.JSON(fiber.Map{
		"message": "success",
		"data":    "this is test",
	})
}

// ListRepos
// @Summary List Repository
// @Tags Repository
// @Accept json
// @Produce json
// @Router /api/repositories [Get]
func ListRepos(c *fiber.Ctx) error {
	repositories, err := repo.LoadFile(settings.RepositoryConfig)
	if err != nil {
		return common.RespErr(c, err)
	}

	repos := make([]repositoryElement, 0, len(repositories.Repositories))
	for _, re := range repositories.Repositories {
		repos = append(repos, repositoryElement{Name: re.Name, URL: re.URL})
	}
	return common.RespOK(c, repos)
}

// RemoveRepo
// @Summary Remove Repository
// @Tags Repository
// @Accept json
// @Produce json
// @Router /api/repositories/:repositories [Delete]
func RemoveRepo(c *fiber.Ctx) error {
	repoName := c.Params("repositories")
	repoFile, err := repo.LoadFile(settings.RepositoryConfig)
	if err != nil {
		return common.RespErr(c, err)
	}

	if !repoFile.Has(repoName) {
		return common.RespErr(c, fmt.Errorf(common.REPO_NO_NAMED_FOUND))
	}

	if !repoFile.Remove(repoName) {
		return common.RespErr(c, err)
	}

	if err := repoFile.WriteFile(settings.RepositoryConfig, 0600); err != nil {
		return common.RespErr(c, err)
	}

	if err := removeRepoCache(settings.RepositoryCache, repoName); err != nil {
		return common.RespErr(c, err)
	}
	return common.RespOK(c, nil)
}

func removeRepoCache(root, name string) error {
	idx := filepath.Join(root, helmpath.CacheChartsFile(name))
	if _, err := os.Stat(idx); err == nil {
		os.Remove(idx)
	}

	idx = filepath.Join(root, helmpath.CacheIndexFile(name))
	if _, err := os.Stat(idx); os.IsNotExist(err) {
		return nil
	} else if err != nil {
		return errors.Wrapf(err, "can't remove index file %s", idx)
	}
	return os.Remove(idx)
}
