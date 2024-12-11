package handler

import (
	"github.com/gofiber/fiber/v2"
	"helm.sh/helm/v3/pkg/repo"
	"reflect"
	"testing"
)

func TestAddRepo(t *testing.T) {
	type args struct {
		c *fiber.Ctx
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := AddRepo(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("AddRepo() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClearRepoCache(t *testing.T) {
	type args struct {
		c *fiber.Ctx
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ClearRepoCache(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("ClearRepoCache() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestListRepoCharts(t *testing.T) {
	type args struct {
		c *fiber.Ctx
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ListRepoCharts(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("ListRepoCharts() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestListRepos(t *testing.T) {
	type args struct {
		c *fiber.Ctx
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ListRepos(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("ListRepos() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRemoveRepo(t *testing.T) {
	type args struct {
		c *fiber.Ctx
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := RemoveRepo(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("RemoveRepo() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUpdateRepo(t *testing.T) {
	type args struct {
		c *fiber.Ctx
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := UpdateRepo(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("UpdateRepo() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUpdateRepoAll(t *testing.T) {
	type args struct {
		repoFile *repo.File
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := UpdateRepoAll(tt.args.repoFile); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UpdateRepoAll() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_addRepoVaildCheck(t *testing.T) {
	type args struct {
		newRepo *addRepositoryElement
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := addRepoVaildCheck(tt.args.newRepo); (err != nil) != tt.wantErr {
				t.Errorf("addRepoVaildCheck() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_getRepoConnectionStatus(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := getRepoConnectionStatus(tt.args.url); (err != nil) != tt.wantErr {
				t.Errorf("getRepoConnectionStatus() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_isNotExist(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isNotExist(tt.args.err); got != tt.want {
				t.Errorf("isNotExist() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_removeRepoCache(t *testing.T) {
	type args struct {
		root string
		name string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := removeRepoCache(tt.args.root, tt.args.name); (err != nil) != tt.wantErr {
				t.Errorf("removeRepoCache() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_saveRepoCaFile(t *testing.T) {
	type args struct {
		caFilePath string
		base64CA   string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := saveRepoCaFile(tt.args.caFilePath, tt.args.base64CA); (err != nil) != tt.wantErr {
				t.Errorf("saveRepoCaFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_syncRepoLock(t *testing.T) {
	type args struct {
		repoFile string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := syncRepoLock(tt.args.repoFile); (err != nil) != tt.wantErr {
				t.Errorf("syncRepoLock() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_updateChart(t *testing.T) {
	type args struct {
		repoEntry *repo.Entry
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := updateChart(tt.args.repoEntry); (err != nil) != tt.wantErr {
				t.Errorf("updateChart() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
