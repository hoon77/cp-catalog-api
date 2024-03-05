package handler

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2/log"
	"go-api/common"
	"go-api/config"
	"helm.sh/helm/v3/pkg/cli"
	"io"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/apimachinery/pkg/apis/testapigroup/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
	sigyaml "sigs.k8s.io/yaml"
	"strings"
)

var (
	settings = cli.New()
)

func Settings() {
	settings.RepositoryConfig = config.Env.HelmRepoConfig
	settings.RepositoryCache = config.Env.HelmRepoCache
}

func GetResources(out string) []*v1.Carp {
	res, err := ParseManifests(out)
	if err != nil {
		res = append(res, &v1.Carp{
			TypeMeta: metav1.TypeMeta{Kind: "ManifestParseError"},
			ObjectMeta: metav1.ObjectMeta{
				Name: err.Error(),
			},
			Spec: v1.CarpSpec{},
			Status: v1.CarpStatus{
				Phase:   "BrokenManifest",
				Message: err.Error(),
			},
		})
		//_ = c.AbortWithError(http.StatusInternalServerError, err)
		//return
	}
	return res
}

func ParseManifests(out string) ([]*v1.Carp, error) {
	dec := yaml.NewYAMLOrJSONDecoder(strings.NewReader(out), 4096)
	res := make([]*v1.Carp, 0)
	var tmp interface{}
	for {
		err := dec.Decode(&tmp)
		if err == io.EOF {
			break
		}

		if err != nil {
			return res, err
		}

		jsoned, err := json.Marshal(tmp)
		if err != nil {
			return res, err
		}

		var doc v1.Carp
		err = json.Unmarshal(jsoned, &doc)
		if err != nil {
			return res, err
		}

		if doc.Kind == "" {
			log.Warnf("Manifest piece is not k8s resource: %s", jsoned)
			continue
		}

		res = append(res, &doc)
	}
	return res, nil
}

func ConvertYAML(results map[string]interface{}) string {
	obj, err := sigyaml.Marshal(results)
	if err != nil {
		return common.EMPTY_STR
	}
	return string(obj)
}
