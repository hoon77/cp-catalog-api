package config

import (
	"github.com/gofiber/fiber/v2/log"
	"github.com/spf13/viper"
)

var Env *envConfigs

func InitEnvConfigs() {
	Env = loadEnvVariables()
}

type envConfigs struct {
	ServerPort                string `mapstructure:"SERVER_PORT"`
	JwtSecret                 string `mapstructure:"JWT_SECRET"`
	K8sApiServer              string `mapstructure:"K8S_API_SERVER"`
	K8sToken                  string `mapstructure:"K8S_TOKEN"`
	HelmRepoConfig            string `mapstructure:"HELM_REPO_CONFIG"`
	HelmRepoCache             string `mapstructure:"HELM_REPO_CACHE"`
	HelmRepoCA                string `mapstructure:"HELM_REPO_CA"`
	ArtifactHubUrl            string `mapstructure:"ARTIFACT_HUB_API_URL"`
	ArtifactHubRepoSearch     string `mapstructure:"ARTIFACT_HUB_REPO_SEARCH"`
	ArtifactHubPackageSearch  string `mapstructure:"ARTIFACT_HUB_PACKAGE_SEARCH"`
	ArtifactHubPackageDetail  string `mapstructure:"ARTIFACT_HUB_PACKAGE_DETAIL"`
	ArtifactHubPackageValues  string `mapstructure:"ARTIFACT_HUB_PACKAGE_VALUES"`
	ArtifactHubPackageLogoUrl string `mapstructure:"ARTIFACT_HUB_PACKAGE_LOGO_URL"`
	VaultRoleName             string `mapstructure:"VAULT_ROLE_NAME"`
	VaultRoleId               string `mapstructure:"VAULT_ROLE_ID"`
	VaultSecretId             string `mapstructure:"VAULT_SECRET_ID"`
}

func loadEnvVariables() (config *envConfigs) {
	viper.AddConfigPath(".")
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("Error reading env file", err)
	}
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatal(err)
	}
	return
}
