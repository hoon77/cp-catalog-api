package handler

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"go-api/common"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/release"
	"strconv"
)

type releaseElement struct {
	Name         string `json:"name"`
	Namespace    string `json:"namespace"`
	Revision     string `json:"revision"`
	Updated      string `json:"updated"`
	Status       string `json:"status"`
	Chart        string `json:"chart"`
	ChartVersion string `json:"chart_version"`
	AppVersion   string `json:"app_version"`
	Icon         string `json:"icon"`
	Notes        string `json:"notes,omitempty"`
}

// ListReleases godoc
// @Summary List Releases
// @Accept json
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

// GetRelease godoc
// @Summary Get Release
// @Accept json
// @Produce json
// @Router /api/clusters/:clusterId/namespaces/:namespace/releases/:release [Get]
func GetRelease(c *fiber.Ctx) error {
	name := c.Param("release")
	namespace := c.Param("namespace")
	info := c.Query("info")
	kubeConfig := c.Query("kube_config")
	if info == "" {
		info = "values"
	}
	kubeContext := c.Query("kube_context")
	infos := []string{"hooks", "manifest", "notes", "values"}
	infoMap := map[string]bool{}
	for _, i := range infos {
		infoMap[i] = true
	}
	if _, ok := infoMap[info]; !ok {
		respErr(c, fmt.Errorf("bad info %s, release info only support hooks/manifest/notes/values", info))
		return
	}
	actionConfig, err := actionConfigInit(InitKubeInformation(namespace, kubeContext, kubeConfig))
	if err != nil {
		respErr(c, err)
		return
	}

	if info == "values" {
		output := c.Query("output")
		// get values output format
		if output == "" {
			output = "json"
		}
		if output != "json" && output != "yaml" {
			respErr(c, fmt.Errorf("invalid format type %s, output only support json/yaml", output))
			return
		}

		client := action.NewGetValues(actionConfig)
		results, err := client.Run(name)
		if err != nil {
			respErr(c, err)
			return
		}
		if output == "yaml" {
			obj, err := yaml.Marshal(results)
			if err != nil {
				respErr(c, err)
				return
			}
			respOK(c, string(obj))
			return
		}
		respOK(c, results)
		return
	}

	client := action.NewGet(actionConfig)
	results, err := client.Run(name)
	if err != nil {
		respErr(c, err)
		return
	}
	// TODO: support all
	if info == "hooks" {
		if len(results.Hooks) < 1 {
			respOK(c, []*release.Hook{})
			return
		}
		respOK(c, results.Hooks)
		return
	} else if info == "manifest" {
		respOK(c, results.Manifest)
		return
	} else if info == "notes" {
		respOK(c, results.Info.Notes)
		return
	}
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
