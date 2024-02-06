package handler

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"go-api/common"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/release"
	"sigs.k8s.io/yaml"
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
