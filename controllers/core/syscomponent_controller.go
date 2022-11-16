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

package core

import (
	"context"

	"helm.sh/helm/v3/pkg/release"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	"k8s.io/klog/v2"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	corev1alpha1 "github.com/wutong-paas/region/apis/core/v1alpha1"
	"github.com/wutong-paas/region/pkg/helm"
	"github.com/wutong-paas/region/pkg/tasks"
	"k8s.io/apimachinery/pkg/api/errors"
)

// SysComponentReconciler reconciles a SysComponent object
type SysComponentReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	// Recorder is used to record events
	Recorder record.EventRecorder
}

//+kubebuilder:rbac:groups=core.wutong.io,resources=syscomponents,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core.wutong.io,resources=syscomponents/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=core.wutong.io,resources=syscomponents/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the SysComponent object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.12.2/pkg/reconcile
func (r *SysComponentReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	klog.Info("=> Reconciling SysComponent")
	instance := &corev1alpha1.SysComponent{}
	err := r.Get(ctx, req.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request - return and don't requeue
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return ctrl.Result{}, err
	}

	klog.Infof("SysComponent [%s/%s] current status: %s", instance.Name, instance.Namespace, instance.Status.Phase)
	requeue := false
	switch instance.Status.Phase {
	case "":
		instance.Status.Phase = corev1alpha1.SysComponentPendingInstall
		instance.Status.Message = "pending to install"
		requeue = true
	case corev1alpha1.SysComponentInstalled:
		if instance.Spec.InstallWay == "helm" {
			rel, ok := helm.Status(instance.Name, instance.Namespace)
			if !ok {
				instance.Status.Phase = corev1alpha1.SysComponentInstalling
				instance.Status.Message = "helm release not found, need to install again"
			} else {
				if rel != nil {
					if instance.Spec.CurrentVersion != rel.Chart.Metadata.Version {
						instance.Status.Phase = corev1alpha1.SysComponentPendingUpgrade
						instance.Status.Message = "helm release version not match, need to upgrade"
						requeue = true
						break
					}

					switch rel.Info.Status {
					case release.StatusDeployed:
						instance.Status.Phase = corev1alpha1.SysComponentInstalled
						// case release.StatusUnknown:
						// 	instance.Status.Phase = corev1alpha1.SysComponentUnknown
						// case release.StatusFailed:
						// 	instance.Status.Phase = corev1alpha1.SysComponentAbnormal
					}
					instance.Status.Message = rel.Info.Description
				} else {
					// instance.Status.Phase = corev1alpha1.SysComponentUnknown
					instance.Status.Message = "sys component installed, but status unknown"
				}
			}
		} else {
			// TODD: apply
		}
	case corev1alpha1.SysComponentUnInstalled:
		// delete the sys component
		r.Delete(ctx, instance)
		return ctrl.Result{}, nil
	case corev1alpha1.SysComponentPendingInstall:
		// helm install
		err = tasks.NewInstallSysComponentTask(instance).Run()
		if err != nil {
			klog.Errorf("install failed: %s", err)
		}
	case corev1alpha1.SysComponentUnInstalling:
		err = tasks.NewUninstallSysComponentTask(instance).Run()
		if err != nil {
			klog.Errorf("uninstall failed: %s", err)
		}
	case corev1alpha1.SysComponentPendingUpgrade:
		err = tasks.NewUpgradeSysComponentTask(instance).Run()
		if err != nil {
			klog.Errorf("upgrade failed: %s")
		}
	}

	if err = r.Status().Update(ctx, instance); err != nil {
		klog.Errorf("update status failed: %s", err)
		return ctrl.Result{Requeue: requeue}, err
	}

	return ctrl.Result{Requeue: requeue}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *SysComponentReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&corev1alpha1.SysComponent{}).
		Complete(r)
}
