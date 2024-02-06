package handler

import (
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
