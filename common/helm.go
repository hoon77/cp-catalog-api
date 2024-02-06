package common

import (
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

func InitKubeInformation(c *fiber.Ctx) *KubeInformation {
	namespace := c.Params("namespace")
	if strings.ToLower(namespace) == ALL_NAMESPACE {
		namespace = ""
	}

	return &KubeInformation{
		AimCluster:   c.Params("clusterId"),
		AimNamespace: namespace,
		AimApiServer: config.Env.K8sApiServer,
		AimToken:     config.Env.K8sToken,
	}
}

func ActionConfigInit(c *fiber.Ctx) (*action.Configuration, error) {
	kubeInfo := InitKubeInformation(c)
	actionConfig := new(action.Configuration)

	settings.KubeAPIServer = kubeInfo.AimApiServer
	settings.KubeToken = kubeInfo.AimToken
	settings.KubeInsecureSkipTLSVerify = true

	err := actionConfig.Init(settings.RESTClientGetter(), kubeInfo.AimNamespace, os.Getenv("HELM_DRIVER"), glog.Infof)
	if err != nil {
		glog.Errorf("%+v", err)
		return nil, err
	}

	return actionConfig, nil
}
