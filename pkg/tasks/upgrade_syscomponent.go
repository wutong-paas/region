/*
Copyright 2022 The Wutong Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package tasks

import (
	"fmt"
	"os/exec"

	"github.com/wutong-paas/region/apis/core/v1alpha1"
)

type upgradeSysComponentTask struct {
	Instance *v1alpha1.SysComponent
}

func NewUpgradeSysComponentTask(instance *v1alpha1.SysComponent) *upgradeSysComponentTask {
	return &upgradeSysComponentTask{Instance: instance}
}

func (t *upgradeSysComponentTask) Run() error {
	ver := t.Instance.Spec.AvailableVersions[t.Instance.Spec.CurrentVersion]

	// Step 1: helm repo add
	err := exec.Command("helm", "repo", "add", ver.HelmRepoName, ver.HelmRepoName).Run()
	if err != nil {
		t.Instance.Status.Phase = v1alpha1.SysComponentAbnormal
		t.Instance.Status.Message = "helm repo add failed"
		return err
	}

	// Step 2: helm repo update
	err = exec.Command("helm", "repo", "update").Run()
	if err != nil {
		t.Instance.Status.Phase = v1alpha1.SysComponentAbnormal
		t.Instance.Status.Message = "helm repo update failed"
		return err
	}

	// Step 3: helm upgrade
	err = exec.Command("helm", "upgrade", t.Instance.Name, fmt.Sprintf("%s/%s", ver.HelmRepoName, ver.HelmChartName), "--namespace", t.Instance.Namespace, "--version", t.Instance.Spec.CurrentVersion).Run()
	if err != nil {
		t.Instance.Status.Phase = v1alpha1.SysComponentAbnormal
		t.Instance.Status.Message = "helm install failed"
		return err
	}

	t.Instance.Status.Phase = v1alpha1.SysComponentInstalled
	t.Instance.Status.Message = "install completed"

	return nil
}
