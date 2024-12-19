package handler

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	"testing"
)

func TestGetReleaseHistories(t *testing.T) {
	type args struct {
		c *fiber.Ctx
	}
	tests := []struct {
		name    string
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.wantErr(t, GetReleaseHistories(tt.args.c), fmt.Sprintf("GetReleaseHistories(%v)", tt.args.c))
		})
	}
}

func TestGetReleaseInfo(t *testing.T) {
	type args struct {
		c *fiber.Ctx
	}
	tests := []struct {
		name    string
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.wantErr(t, GetReleaseInfo(tt.args.c), fmt.Sprintf("GetReleaseInfo(%v)", tt.args.c))
		})
	}
}

func TestGetReleaseResources(t *testing.T) {
	type args struct {
		c *fiber.Ctx
	}
	tests := []struct {
		name    string
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.wantErr(t, GetReleaseResources(tt.args.c), fmt.Sprintf("GetReleaseResources(%v)", tt.args.c))
		})
	}
}

func TestInstallRelease(t *testing.T) {
	type args struct {
		c *fiber.Ctx
	}
	tests := []struct {
		name    string
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.wantErr(t, InstallRelease(tt.args.c), fmt.Sprintf("InstallRelease(%v)", tt.args.c))
		})
	}
}

func TestListReleases(t *testing.T) {
	type args struct {
		c *fiber.Ctx
	}
	tests := []struct {
		name    string
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.wantErr(t, ListReleases(tt.args.c), fmt.Sprintf("ListReleases(%v)", tt.args.c))
		})
	}
}

func TestRollbackRelease(t *testing.T) {
	type args struct {
		c *fiber.Ctx
	}
	tests := []struct {
		name    string
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.wantErr(t, RollbackRelease(tt.args.c), fmt.Sprintf("RollbackRelease(%v)", tt.args.c))
		})
	}
}

func TestUninstallRelease(t *testing.T) {
	type args struct {
		c *fiber.Ctx
	}
	tests := []struct {
		name    string
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.wantErr(t, UninstallRelease(tt.args.c), fmt.Sprintf("UninstallRelease(%v)", tt.args.c))
		})
	}
}

func TestUpgradeRelease(t *testing.T) {
	type args struct {
		c *fiber.Ctx
	}
	tests := []struct {
		name    string
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.wantErr(t, UpgradeRelease(tt.args.c), fmt.Sprintf("UpgradeRelease(%v)", tt.args.c))
		})
	}
}

func Test_getHistory(t *testing.T) {
	type args struct {
		client *action.History
		name   string
	}
	tests := []struct {
		name    string
		args    args
		want    releaseHistory
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getHistory(tt.args.client, tt.args.name)
			if !tt.wantErr(t, err, fmt.Sprintf("getHistory(%v, %v)", tt.args.client, tt.args.name)) {
				return
			}
			assert.Equalf(t, tt.want, got, "getHistory(%v, %v)", tt.args.client, tt.args.name)
		})
	}
}

func Test_isChartInstallable(t *testing.T) {
	type args struct {
		ch *chart.Chart
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := isChartInstallable(tt.args.ch)
			if !tt.wantErr(t, err, fmt.Sprintf("isChartInstallable(%v)", tt.args.ch)) {
				return
			}
			assert.Equalf(t, tt.want, got, "isChartInstallable(%v)", tt.args.ch)
		})
	}
}
