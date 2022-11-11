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
	logger "sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/juju/errors"
	corev1alpha1 "github.com/wutong-paas/region/apis/core/v1alpha1"
	"github.com/wutong-paas/region/pkg/helm"
	"github.com/wutong-paas/region/pkg/tasks"
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
	log := logger.FromContext(ctx)
	log.Info("=> Reconciling SysComponent")
	instance := &corev1alpha1.SysComponent{}
	err := r.Get(ctx, req.NamespacedName, instance)
	if err != nil {
		if errors.Is(err, errors.NotFound) {
			// Request object not found, could have been deleted after reconcile request - return and don't requeue
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return ctrl.Result{}, err
	}

	if instance.Status.Phase == "" {
		instance.Status.Phase = corev1alpha1.SysComponentInstalling
	}

	klog.Infof("SysComponent [%s/%s] current status: %s", instance.Name, instance.Namespace, instance.Status.Phase)
	// log.Info("Phase: " + instance.Status.Phase)
	switch instance.Status.Phase {
	case "":
		instance.Status.Phase = corev1alpha1.SysComponentInstalling
		instance.Status.Message = "pending to install"
		err = r.Status().Update(ctx, instance)
		if err != nil {
			return ctrl.Result{}, err
		}
		return ctrl.Result{Requeue: true}, nil
	case corev1alpha1.SysComponentInstalled:
		if instance.Spec.InstallWay == "helm" {
			rel, ok := helm.Status(instance.Name, instance.Namespace)
			if !ok {
				instance.Status.Phase = corev1alpha1.SysComponentInstalling
				instance.Status.Message = "helm release not found, need to install again"
			} else {
				if rel != nil {
					switch rel.Info.Status {
					case release.StatusDeployed:
						instance.Status.Phase = corev1alpha1.SysComponentInstalled
					case release.StatusUnknown:
						instance.Status.Phase = corev1alpha1.SysComponentUnknown
					case release.StatusFailed:
						instance.Status.Phase = corev1alpha1.SysComponentAbnormal
					}
					instance.Status.Message = rel.Info.Description
				} else {
					instance.Status.Phase = corev1alpha1.SysComponentUnknown
					instance.Status.Message = "sys component installed, but status unknown"
				}
			}
		} else {
			// TODD: apply
		}
	case corev1alpha1.SysComponentUnInstalled:
		break
	case corev1alpha1.SysComponentInstalling:
		// helm install
		err = tasks.NewInstallSysComponentTask(instance).Run()
		if err != nil {
			log.Error(err, "install failed")
		}
	case corev1alpha1.SysComponentUnInstalling:
		err = tasks.NewUninstallSysComponentTask(instance).Run()
		if err != nil {
			log.Error(err, "uninstall failed")
		}
	case corev1alpha1.SysComponentUpgrading:
		err = tasks.NewUpgradeSysComponentTask(instance).Run()
		if err != nil {
			log.Error(err, "upgrade failed")
		}
	}

	if err = r.Status().Update(ctx, instance); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *SysComponentReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&corev1alpha1.SysComponent{}).
		Complete(r)
}
