package handler

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"go-api/common"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/downloader"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/release"
	"helm.sh/helm/v3/pkg/storage/driver"
	helmtime "helm.sh/helm/v3/pkg/time"
	"k8s.io/kubectl/pkg/cmd/get"
	"sigs.k8s.io/yaml"
	"strconv"
	"time"
)

var defaultTimeout = "5m0s"

type releaseElement struct {
	Name         string      `json:"name"`
	Namespace    string      `json:"namespace"`
	Repo         string      `json:"repo"`
	Revision     string      `json:"revision"`
	Updated      string      `json:"updated"`
	Status       string      `json:"status"`
	Chart        string      `json:"chart"`
	ChartVersion string      `json:"chart_version"`
	AppVersion   string      `json:"app_version"`
	Icon         string      `json:"icon"`
	Notes        string      `json:"notes,omitempty"`
	Values       string      `json:"values"`
	Resources    interface{} `json:"resources"`
	Mainifest    string      `json:"manifest"`
}

type releaseInfo struct {
	Revision    int           `json:"revision"`
	Updated     helmtime.Time `json:"updated"`
	Status      string        `json:"status"`
	Chart       string        `json:"chart"`
	AppVersion  string        `json:"app_version"`
	Description string        `json:"description"`
}
type releaseHistory []releaseInfo

// ListReleases
// @Summary List Releases
// @Tags Releases
// @Accept json
// @Produce json
// @Router /api/clusters/:clusterId/namespaces/:namespace/releases [Get]
func ListReleases(c *fiber.Ctx) error {
	lse, err := ListSearchCheck(c)
	if err != nil {
		return common.RespErr(c, err)
	}

	actionConfig, err := common.ActionConfigInit(c)
	if err != nil {
		return common.RespErr(c, err)
	}

	client := action.NewList(actionConfig)
	client.All = true
	client.SortReverse = true
	results, err := client.Run()
	if err != nil {
		return common.RespErr(c, err)
	}

	elements := make([]interface{}, 0, len(results))
	for _, r := range results {
		elements = append(elements, constructReleaseElement(r, false))
	}

	itemCount, resultData := ResourceListProcessing(elements, lse)
	return common.ListRespOK(c, itemCount, resultData)
}

// GetReleaseInfo
// @Summary Get Releases Info
// @Tags Releases
// @Accept json
// @Produce json
// @Router /api/clusters/:clusterId/namespaces/:namespace/releases/:release [Get]
func GetReleaseInfo(c *fiber.Ctx) error {
	name := c.Params("release")
	actionConfig, err := common.ActionConfigInit(c)
	if err != nil {
		return common.RespErr(c, err)
	}

	client := action.NewGet(actionConfig)
	results, err := client.Run(name)
	if err != nil {
		if errors.Is(err, driver.ErrReleaseNotFound) {
			return common.RespErr(c, fmt.Errorf(common.RELEASE_NOT_FOUND))
		}
		return common.RespErr(c, err)
	}

	releaseElement := constructReleaseInfoElement(results)

	return common.RespOK(c, releaseElement)
}

// InstallRelease
// @Summary Install Release
// @Tags Releases
// @Accept json
// @Produce json
// @Router /api/clusters/:clusterId/namespaces/:namespace/releases/:release [Post]
func InstallRelease(c *fiber.Ctx) error {
	newRelease := new(releaseElement)
	if err := c.BodyParser(newRelease); err != nil {
		return common.RespErr(c, err)
	}
	newRelease.Name = c.Params("release")
	newRelease.Namespace = c.Params("namespace")

	if newRelease.Chart == "" {
		return common.RespErr(c, fmt.Errorf(common.CHART_INFO_INVALID))
	}

	if err := runInstall(c, newRelease); err != nil {
		return common.RespErr(c, err)
	}
	return common.RespOK(c, nil)
}

// UpgradeRelease
// @Summary Upgrade Release
// @Tags Releases
// @Accept json
// @Produce json
// @Router /api/clusters/:clusterId/namespaces/:namespace/releases/:release [Put]
func UpgradeRelease(c *fiber.Ctx) error {
	upgradeRelease := new(releaseElement)
	if err := c.BodyParser(upgradeRelease); err != nil {
		return common.RespErr(c, err)
	}
	upgradeRelease.Name = c.Params("release")
	upgradeRelease.Namespace = c.Params("namespace")
	if upgradeRelease.Chart == "" {
		return common.RespErr(c, fmt.Errorf(common.CHART_INFO_INVALID))
	}

	vals, err := mergeValues(upgradeRelease.Values)
	if err != nil {
		return common.RespErr(c, err)
	}

	actionConfig, err := common.ActionConfigInit(c)
	if err != nil {
		return common.RespErr(c, err)
	}

	client := action.NewUpgrade(actionConfig)
	client.Namespace = upgradeRelease.Namespace
	client.Version = upgradeRelease.ChartVersion

	aimChart := fmt.Sprintf("%s/%s", upgradeRelease.Repo, upgradeRelease.Chart)
	cp, err := client.ChartPathOptions.LocateChart(aimChart, settings)
	if err != nil {
		return common.RespErr(c, err)
	}

	chartRequested, err := loader.Load(cp)
	if err != nil {
		return common.RespErr(c, err)
	}

	if req := chartRequested.Metadata.Dependencies; req != nil {
		if err := action.CheckDependencies(chartRequested, req); err != nil {
			return common.RespErr(c, err)
		}
	}

	_, err = client.Run(upgradeRelease.Name, chartRequested, vals)
	if err != nil {
		return common.RespErr(c, err)
	}

	return common.RespOK(c, nil)
}

// RollbackRelease
// @Summary Rollback Release
// @Tags Releases
// @Accept json
// @Produce json
// @Router /api/clusters/:clusterId/namespaces/:namespace/releases/:release/versions/:reversion [Put]
func RollbackRelease(c *fiber.Ctx) error {
	name := c.Params("release")
	reversionStr := c.Params("reversion")

	reversion, err := strconv.Atoi(reversionStr)
	if err != nil {
		return common.RespErr(c, fmt.Errorf(common.REVERSION_NUMBER_INVALID))
	}

	actionConfig, err := common.ActionConfigInit(c)
	if err != nil {
		return common.RespErr(c, err)
	}

	client := action.NewRollback(actionConfig)
	client.Version = reversion

	err = client.Run(name)
	if err != nil {
		return common.RespErr(c, err)
	}

	return common.RespOK(c, nil)
}

// UninstallRelease
// @Summary Uninstall Release
// @Tags Releases
// @Accept json
// @Produce json
// @Router /api/clusters/:clusterId/namespaces/:namespace/releases/:release [Delete]
func UninstallRelease(c *fiber.Ctx) error {
	name := c.Params("release")
	actionConfig, err := common.ActionConfigInit(c)
	if err != nil {
		return common.RespErr(c, err)
	}

	client := action.NewUninstall(actionConfig)
	_, err = client.Run(name)
	if err != nil {
		return common.RespErr(c, err)
	}
	return common.RespOK(c, nil)
}

// GetReleaseHistories
// @Summary Get Release Histories
// @Tags Releases
// @Accept json
// @Produce json
// @Router /api/clusters/:clusterId/namespaces/:namespace/releases/:release/histories [Get]
func GetReleaseHistories(c *fiber.Ctx) error {
	name := c.Params("release")
	actionConfig, err := common.ActionConfigInit(c)
	if err != nil {
		return common.RespErr(c, err)
	}

	client := action.NewHistory(actionConfig)
	results, err := client.Run(name)
	if err != nil {
		return common.RespErr(c, err)
	}

	if len(results) == 0 {
		return common.RespOK(c, &releaseHistory{})
	}

	return common.RespOK(c, getReleaseHistory(results))
}

// GetReleaseResources
// @Summary Get Release Resources
// @Tags Releases
// @Accept json
// @Produce json
// @Router /api/clusters/:clusterId/namespaces/:namespace/releases/:release/resources [Get]
func GetReleaseResources(c *fiber.Ctx) error {
	name := c.Params("release")
	actionConfig, err := common.ActionConfigInit(c)
	if err != nil {
		return common.RespErr(c, err)
	}

	client := action.NewStatus(actionConfig)
	client.ShowResources = true
	client.ShowResourcesTable = true
	status, err := client.Run(name)
	if err != nil {
		return common.RespErr(c, err)
	}

	buf := new(bytes.Buffer)
	if status.Info.Resources != nil && len(status.Info.Resources) > 0 {
		printFlags := get.NewHumanPrintFlags()
		typePrinter, _ := printFlags.ToPrinter("")
		printer := &get.TablePrinter{Delegate: typePrinter}

		var keys []string
		for key := range status.Info.Resources {
			keys = append(keys, key)
		}

		for _, t := range keys {
			_, _ = fmt.Fprintf(buf, "==> %s\n", t)

			vk := status.Info.Resources[t]
			for _, resource := range vk {
				if err := printer.PrintObj(resource, buf); err != nil {
					_, _ = fmt.Fprintf(buf, "failed to print object type %s: %v\n", t, err)
				}
			}

			buf.WriteString("\n")
		}

	}

	return common.RespOK(c, buf.String())
}

func runInstall(c *fiber.Ctx, release *releaseElement) (err error) {
	vals, err := mergeValues(release.Values)
	if err != nil {
		return
	}

	actionConfig, err := common.ActionConfigInit(c)
	if err != nil {
		return
	}

	client := action.NewInstall(actionConfig)
	client.ReleaseName = release.Name
	client.Namespace = release.Namespace
	client.Version = release.ChartVersion

	aimChart := fmt.Sprintf("%s/%s", release.Repo, release.Chart)

	cp, err := client.ChartPathOptions.LocateChart(aimChart, settings)
	if err != nil {
		return
	}

	chartRequested, err := loader.Load(cp)
	if err != nil {
		return
	}

	validInstallableChart, err := isChartInstallable(chartRequested)
	if !validInstallableChart {
		return
	}

	if req := chartRequested.Metadata.Dependencies; req != nil {
		// If CheckDependencies returns an error, we have unfulfilled dependencies.
		// As of Helm 2.4.0, this is treated as a stopping condition:
		// https://github.com/helm/helm/issues/2209
		if err = action.CheckDependencies(chartRequested, req); err != nil {
			if client.DependencyUpdate {
				man := &downloader.Manager{
					ChartPath:        cp,
					Keyring:          client.ChartPathOptions.Keyring,
					SkipUpdate:       false,
					Getters:          getter.All(settings),
					RepositoryConfig: settings.RepositoryConfig,
					RepositoryCache:  settings.RepositoryCache,
				}
				if err = man.Update(); err != nil {
					return
				}
			} else {
				return
			}
		}
	}

	_, err = client.Run(chartRequested, vals)
	if err != nil {
		return
	}

	return nil
}

func isChartInstallable(ch *chart.Chart) (bool, error) {
	switch ch.Metadata.Type {
	case "", "application":
		return true, nil
	}

	return false, fmt.Errorf("charts are not installable")
}

func constructReleaseElement(r *release.Release, showStatus bool) releaseElement {
	element := releaseElement{
		Name:         r.Name,
		Namespace:    r.Namespace,
		Revision:     strconv.Itoa(r.Version),
		Status:       r.Info.Status.String(),
		Chart:        r.Chart.Metadata.Name,
		ChartVersion: r.Chart.Metadata.Version,
		AppVersion:   r.Chart.Metadata.AppVersion,
		Icon:         r.Chart.Metadata.Icon,
		Resources:    make([]string, 0),
	}
	if showStatus {
		element.Notes = r.Info.Notes
	}
	t := "-"
	if tspb := r.Info.LastDeployed; !tspb.IsZero() {
		t = tspb.Format(time.DateTime)
	}
	element.Updated = t

	return element
}

func constructReleaseInfoElement(r *release.Release) releaseElement {
	element := releaseElement{
		Name:         r.Name,
		Namespace:    r.Namespace,
		Revision:     strconv.Itoa(r.Version),
		Status:       r.Info.Status.String(),
		Chart:        r.Chart.Metadata.Name,
		ChartVersion: r.Chart.Metadata.Version,
		AppVersion:   r.Chart.Metadata.AppVersion,
		Icon:         r.Chart.Metadata.Icon,
		Notes:        r.Info.Notes,
		Values:       ConvertYAML(r.Chart.Values),
		Resources:    GetResources(r.Manifest),
		Mainifest:    r.Manifest,
	}

	t := "-"
	if tspb := r.Info.LastDeployed; !tspb.IsZero() {
		t = tspb.Format(time.DateTime)
	}
	element.Updated = t

	return element
}

// MergeValues merges values from files specified via -f/--values and directly
// via --set-json, --set, --set-string, or --set-file, marshaling them to YAML
func mergeValues(values string) (map[string]interface{}, error) {
	byts := []byte(values)
	vals := map[string]interface{}{}

	if err := yaml.Unmarshal(byts, &vals); err != nil {
		return nil, fmt.Errorf(common.FAILED_TO_PARSE_VALUES)
	}
	return vals, nil
}

func getReleaseHistory(rls []*release.Release) (history releaseHistory) {
	for i := len(rls) - 1; i >= 0; i-- {
		r := rls[i]
		c := formatChartname(r.Chart)
		s := r.Info.Status.String()
		v := r.Version
		d := r.Info.Description
		a := formatAppVersion(r.Chart)

		rInfo := releaseInfo{
			Revision:    v,
			Status:      s,
			Chart:       c,
			AppVersion:  a,
			Description: d,
		}
		if !r.Info.LastDeployed.IsZero() {
			rInfo.Updated = r.Info.LastDeployed

		}
		history = append(history, rInfo)
	}

	return history
}

func formatChartname(c *chart.Chart) string {
	if c == nil || c.Metadata == nil {
		// This is an edge case that has happened in prod, though we don't
		// know how: https://github.com/helm/helm/issues/1347
		return "MISSING"
	}
	return fmt.Sprintf("%s-%s", c.Name(), c.Metadata.Version)
}

func formatAppVersion(c *chart.Chart) string {
	if c == nil || c.Metadata == nil {
		// This is an edge case that has happened in prod, though we don't
		// know how: https://github.com/helm/helm/issues/1347
		return "MISSING"
	}
	return c.AppVersion()
}
func GetReleaseOld(c *fiber.Ctx) error {
	infos := []string{"hooks", "manifest", "notes", "values"}

	name := c.Params("release")
	info := c.Query("info")

	if info == "" {
		info = "values"
	}

	infoMap := map[string]bool{}
	for _, i := range infos {
		infoMap[i] = true
	}
	if _, ok := infoMap[info]; !ok {
		return common.RespErr(c, fmt.Errorf("bad info %s, release info only support hooks/manifest/notes/values", info))
	}

	actionConfig, err := common.ActionConfigInit(c)
	if err != nil {
		return common.RespErr(c, err)
	}

	// values
	if info == "values" {
		output := c.Query("output")
		// get values output format
		if output == "" {
			output = "json"
		}
		if output != "json" && output != "yaml" {
			return common.RespErr(c, fmt.Errorf("invalid format type %s, output only support json/yaml", output))
		}

		client := action.NewGetValues(actionConfig)
		results, err := client.Run(name)
		if err != nil {
			return common.RespErr(c, err)
		}

		if output == "yaml" {
			obj, err := yaml.Marshal(results)
			if err != nil {
				return common.RespErr(c, err)
			}
			return common.RespOK(c, string(obj))
		}
		return common.RespOK(c, results)
	}

	client := action.NewGet(actionConfig)
	results, err := client.Run(name)
	if err != nil {
		return common.RespErr(c, err)
	}

	// TODO: support all
	if info == "hooks" {
		if len(results.Hooks) < 1 {
			return common.RespOK(c, []*release.Hook{})
		}
		return common.RespOK(c, results.Hooks)

	} else if info == "manifest" {
		return common.RespOK(c, results.Manifest)
	} else if info == "notes" {
		return common.RespOK(c, results.Info.Notes)

	}

	return common.RespOK(c, nil)
}
