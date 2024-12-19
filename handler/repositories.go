package handler

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofrs/flock"
	"github.com/pkg/errors"
	"go-api/common"
	"go-api/config"
	"helm.sh/helm/v3/cmd/helm/search"
	"helm.sh/helm/v3/pkg/getter"
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

/*func addRepoVaildCheck(newRepo *addRepositoryElement) error {
	if newRepo.Name == "" || newRepo.URL == "" {
		return fmt.Errorf(common.REPO_NAME_URL_REQUIRED)
	}

	match, _ := regexp.MatchString(common.REPO_NAME_REGEXP_PATTERN, newRepo.Name)
	if !match {
		return fmt.Errorf(common.REPO_NAME_PATTERN_NOT_ALLOWED)
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
*/
// AddRepo
// @Summary Add Repository
// @Tags Repository
// @Accept json
// @Produce json
// @Router /api/repositories [Post]
func AddRepo(c *fiber.Ctx) error {
	repoFile := settings.RepositoryConfig
	newRepo := new(addRepositoryElement)
	if err := c.BodyParser(newRepo); err != nil {
		return common.RespErr(c, err)
	}
	/*	if err := addRepoVaildCheck(newRepo); err != nil {
			return common.RespErr(c, err)
		}
	*/
	log.Infof("Add repo :: name: %s, url: %s", newRepo.Name, newRepo.URL)
	/*	if err := getRepoConnectionStatus(newRepo.URL); err != nil {
			return common.RespErr(c, err)
		}
	*/
	// Ensure the file directory exists as it is required for file locking
	err := os.MkdirAll(filepath.Dir(repoFile), os.ModePerm)
	if err != nil && !os.IsExist(err) {
		return common.RespErr(c, err)
	}

	// Acquire a file lock for process synchronization
	if err := syncRepoLock(repoFile); err != nil {

		return common.RespErr(c, err)
	}

	f, err := repo.LoadFile(repoFile)
	if err != nil {
		log.Errorf("AddRepo:: faild load file :: %v", err)
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
	caFilePath := ""
	if len(newRepo.CaBase64) > 0 {
		caFile := fmt.Sprintf("%v_%v.crt", newRepo.Name, generatingId())
		caFilePath = filepath.Join(config.Env.HelmRepoCA, caFile)
		if err := os.MkdirAll(config.Env.HelmRepoCA, os.ModePerm); err != nil && !os.IsExist(err) {
			return common.RespErr(c, err)
		}
		/*	if err := saveRepoCaFile(caFilePath, newRepo.CaBase64); err != nil {
			return common.RespErr(c, err)
		}*/
		repoEntry.CAFile = caFilePath
	}

	r, err := repo.NewChartRepository(&repoEntry, getter.All(settings))
	if err != nil {
		log.Errorf("NewChartRepository ::  %s", err.Error())
		_ = RemoveFile(caFilePath)
		return common.RespErr(c, err)
	}

	// set cache path
	if settings.RepositoryCache != "" {
		r.CachePath = settings.RepositoryCache
	}

	if _, err := r.DownloadIndexFile(); err != nil {
		log.Errorf("DownloadIndexFile ::  %s", err.Error())
		_ = RemoveFile(caFilePath)
		return common.RespErr(c, err)
	}

	f.Update(&repoEntry)

	if err := f.WriteFile(repoFile, 0600); err != nil {
		log.Errorf("Write Repofile ::  %s", err.Error())
		_ = RemoveFile(caFilePath)
		return common.RespErr(c, err)
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
	lse, err := ListSearchCheck(c)
	if err != nil {
		return common.RespErr(c, err)
	}

	repositories, err := repo.LoadFile(settings.RepositoryConfig)
	if err != nil {
		log.Errorf("ListRepos:: faild load file :: %v", err)
		return common.RespErr(c, fmt.Errorf(common.REPO_FAILED_LOADING_FILE))
	}

	repos := make([]interface{}, 0, len(repositories.Repositories))
	for i := len(repositories.Repositories) - 1; i >= 0; i-- {
		re := repositories.Repositories[i]
		repos = append(repos, repositoryElement{Name: re.Name, URL: re.URL})
	}

	itemCount, resultData := ResourceListProcessing(repos, lse)
	return common.ListRespOK(c, itemCount, resultData)
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
		log.Errorf("RemoveRepo:: faild load file :: %v", err)
		return common.RespErr(c, fmt.Errorf(common.REPO_FAILED_LOADING_FILE))
	}

	if !repoFile.Has(repoName) {
		return common.RespErr(c, fmt.Errorf(common.REPO_NO_NAMED_FOUND))
	}
	removeRepo := repoFile.Get(repoName)

	if !repoFile.Remove(repoName) {
		return common.RespErr(c, err)
	}

	if err := repoFile.WriteFile(settings.RepositoryConfig, 0600); err != nil {
		return common.RespErr(c, err)
	}

	/*	if err := removeRepoCache(settings.RepositoryCache, repoName); err != nil {
			return common.RespErr(c, err)
		}
	*/
	// delete repo ca.crt
	_ = RemoveFile(removeRepo.CAFile)

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
	log.Infof("Update repo (name: %s, url: %s)", updateRepo.Name, updateRepo.URL)
	//err = updateChart(updateRepo)
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
	lse, err := ListSearchCheck(c)
	if err != nil {
		return common.RespErr(c, err)
	}

	repoName := c.Params("repositories")
	version := ">0.0.0"
	index, err := buildSearchIndex(repoName)
	if err != nil {
		return common.RespErr(c, err)
	}

	var res []*search.Result
	res = index.All()
	search.SortScore(res)
	data, err := applyConstraint(version, false, res)
	if err != nil {
		return common.RespErr(c, err)
	}

	chartList := make([]interface{}, 0, len(data))
	for _, v := range data {
		chartList = append(chartList, repoChartElement{
			Name:        strings.Replace(v.Chart.Name, repoName+"/", "", 1),
			Version:     v.Chart.Version,
			AppVersion:  v.Chart.AppVersion,
			Description: v.Chart.Description,
			Icon:        v.Chart.Icon,
			RepoName:    repoName,
			Deprecated:  v.Chart.Deprecated,
		})
	}

	itemCount, resultData := ResourceListProcessing(chartList, lse)
	return common.ListRespOK(c, itemCount, resultData)
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
