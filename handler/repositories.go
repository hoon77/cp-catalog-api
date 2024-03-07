package handler

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofrs/flock"
	"github.com/pkg/errors"
	"go-api/common"
	"go-api/config"
	"helm.sh/helm/v3/cmd/helm/search"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/helmpath"
	"helm.sh/helm/v3/pkg/repo"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Repositories that have been permanently deleted and no longer work
var deprecatedRepos = map[string]string{
	"//kubernetes-charts.storage.googleapis.com":           "https://charts.helm.sh/stable",
	"//kubernetes-charts-incubator.storage.googleapis.com": "https://charts.helm.sh/incubator",
}

type repositoryElement struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type addRepositoryElement struct {
	Name     string `json:"name"`
	URL      string `json:"url"`
	Username string `json:"username"`
	Password string `json:"password"`
	CaBase64 string `json:"ca_base64"`
}

func addRepoVaildCheck(newRepo *addRepositoryElement) error {
	if newRepo.Name == "" || newRepo.URL == "" {
		return fmt.Errorf(common.REPO_NAME_URL_REQUIRED)
	}
	if strings.Contains(newRepo.Name, "/") {
		return fmt.Errorf(common.REPO_NAME_CONTAINS_SC)
	}
	// Block deprecated repos
	for oldURL, newURL := range deprecatedRepos {
		if strings.Contains(newRepo.URL, oldURL) {
			return fmt.Errorf("repo %q is no longer available; try %q instead", newRepo.URL, newURL)
		}
	}

	if (newRepo.Username != "" && newRepo.Password == "") || (newRepo.Username == "" && newRepo.Password != "") {
		return errors.New(common.REPO_USERNAME_PASSWD_REQUIRED)
	}
	return nil
}

// AddRepo
// @Summary Add Repository
// @Tags Repository
// @Accept json
// @Produce json
// @Router /api/repositories/:repositories [Post]
func AddRepo(c *fiber.Ctx) error {
	repoFile := settings.RepositoryConfig
	newRepo := new(addRepositoryElement)
	if err := c.BodyParser(newRepo); err != nil {
		return common.RespErr(c, err)
	}
	newRepo.Name = c.Params("repositories")

	if err := addRepoVaildCheck(newRepo); err != nil {
		return err
	}

	// Ensure the file directory exists as it is required for file locking
	err := os.MkdirAll(filepath.Dir(repoFile), os.ModePerm)
	if err != nil && !os.IsExist(err) {
		return err
	}

	// Acquire a file lock for process synchronization
	if err := syncRepoLock(repoFile); err != nil {
		return err
	}

	f, err := repo.LoadFile(repoFile)
	if err != nil {
		return common.RespErr(c, fmt.Errorf(common.REPO_FAILED_LOADING_FILE))
	}

	repoEntry := repo.Entry{
		Name:     newRepo.Name,
		URL:      newRepo.URL,
		Username: newRepo.Username,
		Password: newRepo.Password,
	}

	if f.Has(newRepo.Name) {
		existing := f.Get(newRepo.Name)
		if repoEntry != *existing {
			return common.RespErr(c, errors.Errorf(common.REPO_NAME_ALREADY_EXISTS))
		}
		// The add is idempotent so do nothing
		return common.RespErr(c, errors.Errorf(common.REPO_SAME_CONF_ALREADY_EXISTS))
	}

	// save ca.crt
	caFilePath := fmt.Sprintf("%s%s.crt", config.Env.RepoCertPath, newRepo.Name)
	if len(newRepo.CaBase64) > 0 {
		if err := os.MkdirAll(config.Env.RepoCertPath, os.ModePerm); err != nil && !os.IsExist(err) {
			return common.RespErr(c, err)
		}
		if err := saveRepoCaFile(caFilePath, newRepo.CaBase64); err != nil {
			return common.RespErr(c, err)
		}
		repoEntry.CAFile = caFilePath
	}

	r, err := repo.NewChartRepository(&repoEntry, getter.All(settings))
	if err != nil {
		return common.RespErr(c, err)
	}

	if _, err := r.DownloadIndexFile(); err != nil {
		return common.RespErr(c, errors.Errorf(common.REPO_CANNOT_BE_REACHED))
	}

	f.Update(&repoEntry)

	if err := f.WriteFile(repoFile, 0600); err != nil {
		return err
	}

	return common.RespOK(c, nil)
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
		return common.RespErr(c, fmt.Errorf(common.REPO_FAILED_LOADING_FILE))
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

// UpdateRepo
// @Summary Update Repository
// @Tags Repository
// @Accept json
// @Produce json
// @Router /api/repositories/:repositories [Put]
func UpdateRepo(c *fiber.Ctx) error {
	repoName := c.Params("repositories")
	repoFile, err := repo.LoadFile(settings.RepositoryConfig)
	if err != nil {
		return common.RespErr(c, fmt.Errorf(common.REPO_FAILED_LOADING_FILE))
	}
	if !repoFile.Has(repoName) {
		return common.RespErr(c, fmt.Errorf(common.REPO_NO_NAMED_FOUND))
	}

	updateRepo := repoFile.Get(repoName)
	err = updateChart(updateRepo)
	if err != nil {
		log.Errorf("Failed to update repo.. %s", err.Error())
		return common.RespErr(c, fmt.Errorf(common.REPO_UNABLE_UPDATE))
	}

	return common.RespOK(c, nil)
}

// ListRepoCharts
// @Summary List Repository Charts
// @Tags Repository
// @Accept json
// @Produce json
// @Router /api/repositories/:repositories/charts [Get]
func ListRepoCharts(c *fiber.Ctx) error {
	repoName := c.Params("repositories")
	version := ">0.0.0"
	index, err := buildSearchIndex(repoName)
	if err != nil {
		return common.RespErr(c, err)
	}

	var res []*search.Result
	res, err = index.Search(fmt.Sprintf("%s/", repoName), searchMaxScore, false)
	if err != nil {
		return common.RespErr(c, err)
	}

	search.SortScore(res)
	data, err := applyConstraint(version, false, res)
	if err != nil {
		return common.RespErr(c, err)
	}

	chartList := make(repoChartList, 0, len(data))
	for _, v := range data {
		chartList = append(chartList, repoChartElement{
			Name:        v.Name,
			Version:     v.Chart.Version,
			AppVersion:  v.Chart.AppVersion,
			Description: v.Chart.Description,
			Icon:        v.Chart.Icon,
		})
	}

	return common.RespOK(c, chartList)
}

func syncRepoLock(repoFile string) error {
	repoFileExt := filepath.Ext(repoFile)
	var lockPath string
	if len(repoFileExt) > 0 && len(repoFileExt) < len(repoFile) {
		lockPath = strings.TrimSuffix(repoFile, repoFileExt) + ".lock"
	} else {
		lockPath = repoFile + ".lock"
	}
	fileLock := flock.New(lockPath)
	lockCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	locked, err := fileLock.TryLockContext(lockCtx, time.Second)
	if err == nil && locked {
		defer fileLock.Unlock()
	}
	if err != nil {
		return err
	}

	return nil
}

func updateChart(repoEntry *repo.Entry) error {
	chartRepository, err := repo.NewChartRepository(repoEntry, getter.All(settings))
	if err != nil {
		return err
	}
	if _, err := chartRepository.DownloadIndexFile(); err != nil {
		return err
	}

	return nil
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

func saveRepoCaFile(caFilePath string, base64CA string) error {
	// decode CA
	origCA, err := base64.StdEncoding.DecodeString(base64CA)
	if err != nil {
		return fmt.Errorf(common.REPO_CA_INVALID)
	}
	err = os.WriteFile(caFilePath, origCA, 0644)
	if err != nil {
		return fmt.Errorf(common.REPO_CA_FAILED_SAVE)
	}

	return nil
}
