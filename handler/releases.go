package handler

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"go-api/common"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/downloader"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/release"
	"sigs.k8s.io/yaml"
	"strconv"
)

var (
	settings = cli.New()
)

type releaseElement struct {
	Name         string      `json:"name"`
	Namespace    string      `json:"namespace"`
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

// ListReleases godoc
// @Summary List Releases
// @Accept json-+
// @Produce json
// @Router /api/clusters/:clusterId/namespaces/:namespace/releases [Get]
func ListReleases(c *fiber.Ctx) error {
	actionConfig, err := common.ActionConfigInit(c)
	if err != nil {
		return common.RespErr(c, err)
	}

	client := action.NewList(actionConfig)
	client.Deployed = true
	results, err := client.Run()
	if err != nil {
		return common.RespErr(c, err)
	}

	elements := make([]releaseElement, 0, len(results))
	for _, r := range results {
		elements = append(elements, constructReleaseElement(r, false))
	}
	return common.RespOK(c, elements)
}

// GetReleaseInfo godoc
// @Summary Get Release Info
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
		return common.RespErr(c, err)
	}

	releaseElement := constructReleaseInfoElement(results)

	return common.RespOK(c, releaseElement)
}

// InstallRelease godoc
// @Summary Install Release
// @Accept json
// @Produce json
// @Router /api/clusters/:clusterId/namespaces/:namespace/releases/:release [Post]
func InstallRelease(c *fiber.Ctx) error {
	//name := c.Params("release")
	release := new(releaseElement)
	if err := c.BodyParser(release); err != nil {
		return common.RespErr(c, err)
	}

	if release.Chart == "" {
		return common.RespErr(c, fmt.Errorf("CHART_INFO_INVALID"))
	}

	if err = runInstall(name, namespace, kubeContext, aimChart, kubeConfig, options); err != nil {
		respErr(c, err)
		return
	}

	respOK(c, err)
	return
}

func runInstall(c *fiber.Ctx, name, namespace, kubeContext, aimChart, kubeConfig string, release releaseElement) (err error) {
	actionConfig, err := common.ActionConfigInit(c)
	if err != nil {
		return common.RespErr(c, err)
	}

	client := action.NewInstall(actionConfig)
	client.ReleaseName = name
	client.Namespace = namespace
	client.ChartPathOptions.Version = release.ChartVersion

	cp, err := client.ChartPathOptions.LocateChart(aimChart, settings)
	if err != nil {
		return common.RespErr(c, err)
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

	_, err = client.Run(chartRequested, nil)
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
		t = tspb.String()
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
		t = tspb.String()
	}
	element.Updated = t

	return element
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
