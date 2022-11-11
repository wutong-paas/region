package services

import (
	"testing"

	"github.com/wutong-paas/region/apis/core/v1alpha1"
	"sigs.k8s.io/yaml"
)

func Test(t *testing.T) {
	sysComponentConfigs := map[string]*SysComponentConfig{
		"wutong-observer": {
			Description: "梧桐容器平台可观测组件",
			InstallWay:  "helm",
			Namespace:   "wt-system",
			AvailableVersions: map[string]*v1alpha1.SysComponentVersionInfo{
				"v0.1.0": {
					HelmRepoName:  "wutong-observer",
					HelmChartName: "wutong-observer",
					HelmRepoUrl:   "https://wutong-paas.github.io/wutong-observer",
				},
				"v0.1.1": {
					HelmRepoName:  "wutong-observer",
					HelmChartName: "wutong-observer",
					HelmRepoUrl:   "https://wutong-paas.github.io/wutong-observer",
				},
			},
		},
		"wutong-gateway": {
			Description: "梧桐容器平台网关组件",
			InstallWay:  "apply",
			Namespace:   "wt-system",
			AvailableVersions: map[string]*v1alpha1.SysComponentVersionInfo{
				"v0.1.0": {
					ApplyFileUrl: "https://raw.githubusercontent.com/wutong-paas/wutong-gateway/v0.1.0/deploy/manifests.yaml",
				},
				"v1.0.0": {
					ApplyFileUrl: "https://raw.githubusercontent.com/wutong-paas/wutong-gateway/v1.0.0/deploy/manifests.yaml",
				},
			},
		},
		"test-redis": {
			Description: "测试系统组件 Redis",
			InstallWay:  "helm",
			Namespace:   "test",
			AvailableVersions: map[string]*v1alpha1.SysComponentVersionInfo{
				"17.3.8": {
					HelmRepoName:  "bitnami",
					HelmChartName: "redis",
					HelmRepoUrl:   "https://charts.bitnami.com/bitnami",
				},
			},
		},
	}

	out, err := yaml.Marshal(sysComponentConfigs)
	if err != nil {
		t.Log(err)
	}

	t.Log(string(out))
}
