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

type installSysComponentTask struct {
	sysComponentTask
}

func NewInstallSysComponentTask(instance *v1alpha1.SysComponent) *installSysComponentTask {
	return &installSysComponentTask{sysComponentTask: sysComponentTask{Instance: instance}}
}

func (t *installSysComponentTask) Run() error {
	switch t.Instance.Spec.InstallWay {
	case v1alpha1.InstallHelm:
		return t.byHelm()
	case v1alpha1.InstallApply:
		return t.byApply()
	default:
		return errors.New("install way is not supported")
	}
}

func (t *installSysComponentTask) byHelm() error {
	ver := t.Instance.Spec.AvailableVersions[t.Instance.Spec.CurrentVersion]

	// Step 1: helm add repo
	err := helm.RepoAdd(ver.HelmRepoName, ver.HelmRepoUrl)
	if err != nil {
		t.setErrorStatus(err)
		return err
	}

	// Step 2: helm install
	r, err := helm.Install(t.Instance.Name, fmt.Sprintf("%s/%s", ver.HelmRepoName, ver.HelmChartName), t.Instance.Namespace, t.Instance.Spec.CurrentVersion, map[string]string{})
	if err != nil {
		t.setErrorStatus(err)
		return err
	}

	t.Instance.Status.Phase = v1alpha1.SysComponentInstalling

	switch r.Info.Status {
	case release.StatusDeployed:
		t.Instance.Status.Phase = v1alpha1.SysComponentInstalled
		// case release.StatusUnknown:
		// 	t.Instance.Status.Phase = v1alpha1.SysComponentUnknown
		// case release.StatusFailed:
		// 	t.Instance.Status.Phase = v1alpha1.SysComponentAbnormal
	}

	t.Instance.Status.Message = r.Info.Description
	return nil
}

func (t *installSysComponentTask) byApply() error {
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

	err := exector.Apply()

	if err != nil {
		t.setErrorStatus(err)
		return err
	}
	t.Instance.Status.Phase = v1alpha1.SysComponentInstalled
	t.Instance.Status.Message = "resouces applied."
	return nil
}
