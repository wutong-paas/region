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
	"os/exec"

	"github.com/wutong-paas/region/apis/core/v1alpha1"
)

type uninstallSysComponentTask struct {
	Instance *v1alpha1.SysComponent
}

func NewUninstallSysComponentTask(instance *v1alpha1.SysComponent) *uninstallSysComponentTask {
	return &uninstallSysComponentTask{Instance: instance}
}

func (t *uninstallSysComponentTask) Run() error {
	err := exec.Command("helm", "uninstall", t.Instance.Name, "--namespace", t.Instance.Namespace).Run()
	if err != nil {
		t.Instance.Status.Phase = v1alpha1.SysComponentAbnormal
		t.Instance.Status.Message = "helm uninstall failed"
		return err
	}

	t.Instance.Status.Phase = v1alpha1.SysComponentUnInstalled
	t.Instance.Status.Message = "uninstall completed"

	return nil
}
