package common

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang/glog"
	"go-api/config"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/cli"
	"os"
	"strings"
)

var (
	settings = cli.New()
)

type KubeInformation struct {
	AimCluster   string
	AimNamespace string
	AimApiServer string
	AimToken     string
}

func InitKubeInformation(c *fiber.Ctx) (*KubeInformation, error) {
	namespace := c.Params("namespace")
	if strings.ToLower(namespace) == ALL_NAMESPACE {
		if c.Route().Name != LIST_RELEASES {
			// No other routes allow namespaces 'all' except list release
			return nil, fmt.Errorf(NAMESPACE_ALL_NOT_ALLOWED)
		}
		namespace = ""
	}

	return &KubeInformation{
		AimCluster:   c.Params("clusterId"),
		AimNamespace: namespace,
		AimApiServer: config.Env.K8sApiServer,
		AimToken:     config.Env.K8sToken,
	}, nil
}

func ActionConfigInit(c *fiber.Ctx) (*action.Configuration, error) {
	kubeInfo, err := InitKubeInformation(c)
	if err != nil {
		return nil, err
	}
	actionConfig := new(action.Configuration)

	settings.KubeAPIServer = kubeInfo.AimApiServer
	settings.KubeToken = kubeInfo.AimToken
	settings.KubeInsecureSkipTLSVerify = true

	err = actionConfig.Init(settings.RESTClientGetter(), kubeInfo.AimNamespace, os.Getenv("HELM_DRIVER"), glog.Infof)
	if err != nil {
		glog.Errorf("%+v", err)
		return nil, err
	}

	return actionConfig, nil
}
