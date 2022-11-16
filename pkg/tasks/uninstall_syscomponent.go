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
	"errors"
	"fmt"

	"github.com/wutong-paas/region/apis/core/v1alpha1"
	"github.com/wutong-paas/region/pkg/helm"
	"github.com/wutong-paas/region/pkg/kube"
	"helm.sh/helm/v3/pkg/release"
	controllerruntime "sigs.k8s.io/controller-runtime"
)

type uninstallSysComponentTask struct {
	sysComponentTask
}

func NewUninstallSysComponentTask(instance *v1alpha1.SysComponent) *uninstallSysComponentTask {
	return &uninstallSysComponentTask{sysComponentTask: sysComponentTask{Instance: instance}}
}

func (t *uninstallSysComponentTask) Run() error {
	switch t.Instance.Spec.InstallWay {
	case v1alpha1.InstallHelm:
		return t.byHelm()
	case v1alpha1.InstallApply:
		return t.byApply()
	default:
		return errors.New("install way is not supported")
	}
}

func (t *uninstallSysComponentTask) byHelm() error {
	r, err := helm.Uninstall(t.Instance.Name, t.Instance.Namespace)
	if err != nil {
		t.setErrorStatus(err)
		return err
	}

	if r == nil {
		t.Instance.Status.Phase = v1alpha1.SysComponentUnInstalled
		t.Instance.Status.Message = "sys component uninstalled"
		return nil
	}

	switch r.Info.Status {
	case release.StatusUninstalled:
		t.Instance.Status.Phase = v1alpha1.SysComponentUnInstalled
		// case release.StatusUnknown:
		// 	t.Instance.Status.Phase = v1alpha1.SysComponentUnknown
		// case release.StatusFailed:
		// 	t.Instance.Status.Phase = v1alpha1.SysComponentAbnormal
	}

	t.Instance.Status.Message = r.Info.Description
	return nil
}

func (t *uninstallSysComponentTask) byApply() error {
	ver := t.Instance.Spec.AvailableVersions[t.Instance.Spec.CurrentVersion]

	var content string
	if ver.ApplyFileContent != "" {
		content = ver.ApplyFileContent
	} else {
		if ver.ApplyFileUrl != "" {
			contentBytes, err := kube.Download(ver.ApplyFileUrl)
			if err != nil {
				t.setErrorStatus(err)
				return fmt.Errorf("download file from [%s] failed: %s", ver.ApplyFileUrl, err.Error())
			}
			content = string(contentBytes)
		} else {
			err := errors.New("sys component configurations is not valid")
			t.setErrorStatus(err)
			return err
		}
	}
	exector := kube.NewApplyExector(content, t.Instance.Namespace, controllerruntime.GetConfigOrDie())

	err := exector.Delete()

	if err != nil {
		t.setErrorStatus(err)
		return err
	}
	t.Instance.Status.Phase = v1alpha1.SysComponentInstalled
	t.Instance.Status.Message = "resouces applied."
	return nil
}
