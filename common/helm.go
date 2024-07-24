package common

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/golang-jwt/jwt/v5"
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

type userKubeAuthInfo struct {
	userName   string
	userType   string
	userAuthId string
	rolesInfo  map[string]clusterInfo
}

type clusterInfo struct {
	userType      string
	namespaceList []string
}

type KubeInfo struct {
	AimCluster   string
	AimNamespace string
	AimApiServer string
	AimToken     string
}

func InitKubeInfo(c *fiber.Ctx) (*KubeInfo, error) {
	namespace := c.Params("namespace")
	if strings.ToLower(namespace) == ALL_NAMESPACE {
		if c.Route().Name != LIST_RELEASES {
			// No other routes allow namespaces 'all' except list release
			return nil, fmt.Errorf(NAMESPACE_ALL_NOT_ALLOWED)
		}
		namespace = ""
	}

	kubeInfo := &KubeInfo{
		AimCluster:   c.Params("clusterId"),
		AimNamespace: namespace,
	}

	err := getKubeToken(c, kubeInfo)
	if err != nil {
		return nil, err
	}

	return kubeInfo, nil
}

func ActionConfigInit(c *fiber.Ctx) (*action.Configuration, error) {
	kubeInfo, err := InitKubeInfo(c)
	if err != nil {
		return nil, err
	}

	fmt.Println("kubeInfo:", kubeInfo)
	actionConfig := new(action.Configuration)

	settings.KubeAPIServer = kubeInfo.AimApiServer
	settings.KubeToken = kubeInfo.AimToken
	settings.KubeInsecureSkipTLSVerify = true
	settings.SetNamespace(kubeInfo.AimNamespace)

	log.Infof("SEND :: CLUSTER: %v, NAMESPACE: %v", kubeInfo.AimCluster, kubeInfo.AimNamespace)

	err = actionConfig.Init(settings.RESTClientGetter(), kubeInfo.AimNamespace, os.Getenv("HELM_DRIVER"), glog.Infof)
	if err != nil {
		glog.Errorf("%+v", err)
		return nil, err
	}

	return actionConfig, nil
}

func getKubeToken(c *fiber.Ctx, kubeInfo *KubeInfo) error {
	claims := c.Locals("user").(*jwt.Token).Claims.(jwt.MapClaims)
	userType := claims["userType"].(string)
	userAuthId := claims["userAuthId"].(string)

	//get cluster details
	if err := getClusterDetails(kubeInfo); err != nil {
		return err
	}

	switch userType {
	case AUTH_CLUSTER_ADMIN, AUTH_USER:
		userType = claims["rolesInfo"].(map[string]interface{})[kubeInfo.AimCluster].(map[string]interface{})["userType"].(string)
		if err := getUserToken(userType, userAuthId, kubeInfo); err != nil {
			return err
		}
	}

	return nil
}

func getClusterDetails(kubeInfo *KubeInfo) error {
	path := fmt.Sprintf("%v/%v", config.Env.VaultClusterPath, kubeInfo.AimCluster)
	data, err := read(path)
	if err != nil {
		return err
	}

	kubeInfo.AimApiServer = data["clusterApiUrl"].(string)
	kubeInfo.AimToken = data["clusterToken"].(string)

	return nil
}

func getUserToken(userType string, userAuthId string, kubeInfo *KubeInfo) error {
	path := fmt.Sprintf("%v/%v/%v", config.Env.VaultUserPath, userAuthId, kubeInfo.AimCluster)
	if userType == AUTH_USER {
		path = fmt.Sprintf("%v/%v", path, kubeInfo.AimNamespace)
	}
	fmt.Println("kubeInfo:", kubeInfo)
	fmt.Println("path:", path)
	data, err := read(path)
	if err != nil {
		return err
	}

	kubeInfo.AimToken = data["clusterToken"].(string)
	return nil
}
