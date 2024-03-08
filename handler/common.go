package handler

import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"go-api/common"
	"go-api/config"
	"helm.sh/helm/v3/pkg/cli"
	"io"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/apimachinery/pkg/apis/testapigroup/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
	"os"
	sigyaml "sigs.k8s.io/yaml"
	"strconv"
	"strings"
)

var (
	settings = cli.New()
)

type ListSearchElement struct {
	Offset int
	Limit  int
}

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

func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func RemoveFile(filename string) error {
	err := os.Remove(filename)
	if err != nil {
		return err
	}

	return nil
}

func ListSearchCheck(c *fiber.Ctx) (*ListSearchElement, error) {
	offset, err := strconv.Atoi(c.Query("offset"))
	if err != nil {
		offset = common.DEFAULT_OFFSET
	}
	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		limit = common.DEFAULT_LIMIT
	}

	if limit < 0 {
		return nil, fmt.Errorf(common.LIMIT_ILLEGAL_ARGUMENT)
	}
	if offset < 0 {
		return nil, fmt.Errorf(common.OFFSET_ILLEGAL_ARGUMENT)
	}

	if offset > 0 && limit == 0 {
		return nil, fmt.Errorf(common.OFFSET_REQUIRES_LIMIT_ILLEGAL_ARGUMENT)
	}

	lse := ListSearchElement{
		Offset: offset,
		Limit:  limit,
	}

	return &lse, nil
}

func ResourceListProcessing(list []interface{}, lse *ListSearchElement) []interface{} {
	allCounts := len(list)
	fmt.Println("lse:", lse)

	start := lse.Offset * lse.Limit

	fmt.Println("allCounts:", allCounts)
	fmt.Println("start:", start)

	if start > allCounts {
		return make([]interface{}, 0)
	}
	if (start + lse.Limit) > allCounts {
		return list[start:]
	}

	return list[start : start+lse.Limit]
}
