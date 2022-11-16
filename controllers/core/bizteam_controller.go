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
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	"k8s.io/klog/v2"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	corev1alpha1 "github.com/wutong-paas/region/apis/core/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// BizTeamReconciler reconciles a BizTeam object
type BizTeamReconciler struct {
	client.Client
	Scheme *runtime.Scheme

	// recorder is used to record events
	Recorder record.EventRecorder
}

//+kubebuilder:rbac:groups=core.wutong.io,resources=bizteams,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core.wutong.io,resources=bizteams/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=core.wutong.io,resources=bizteams/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the BizTeam object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.12.2/pkg/reconcile
func (r *BizTeamReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	klog.Info("=> Reconciling BizTeam")
	instance := &corev1alpha1.BizTeam{}
	err := r.Get(ctx, client.ObjectKey{Name: req.Name}, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request - return and don't requeue
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return ctrl.Result{}, err
	}

	if instance.Status.Phase == "" {
		instance.Status.Phase = corev1alpha1.BizTeamCreating
	}

	klog.Infof("BizTeam [%s] current status: %s", instance.Name, instance.Status.Phase)
	switch instance.Status.Phase {
	case corev1alpha1.BizTeamCreating:
		namespace := namespaceForTeam(instance)
		if err := controllerutil.SetControllerReference(instance, namespace, r.Scheme); err != nil {
			return ctrl.Result{}, err
		}

		found := &corev1.Namespace{}
		err = r.Get(ctx, client.ObjectKey{Name: namespace.Name}, found)
		if err != nil {
			if errors.IsNotFound(err) {
				r.Recorder.Event(instance, corev1.EventTypeNormal, "NamespaceCreating", "Creating namespace "+namespace.Name)
				err = r.Create(ctx, namespace)
				if err != nil {
					return ctrl.Result{}, err
				}
				r.Recorder.Event(instance, corev1.EventTypeNormal, "NamespaceCreated", "Created namespace "+namespace.Name)
			} else {
				return ctrl.Result{}, err
			}
		}
		if found.Status.Phase == corev1.NamespaceActive || found.Status.Phase == corev1.NamespaceTerminating {
			instance.Status.Phase = corev1alpha1.BizTeamPhase(found.Status.Phase)
		}
	case corev1alpha1.BizTeamActive:
	case corev1alpha1.BizTeamTerminating:
		klog.Info(fmt.Sprintf("BizTeam [%s] status.phase: %s", instance.Name, instance.Status.Phase))
		return ctrl.Result{}, nil
	}
	if err = r.Status().Update(ctx, instance); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func namespaceForTeam(team *corev1alpha1.BizTeam) *corev1.Namespace {
	namespace := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: team.Name,
			Labels: map[string]string{
				"creator":    "wutong",
				"controller": team.Name,
			},
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(team, corev1alpha1.SchemeGroupVersion.WithKind("BizTeam")),
			},
		},
	}
	return namespace
}

// SetupWithManager sets up the controller with the Manager.
func (r *BizTeamReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&corev1alpha1.BizTeam{}).
		Complete(r)
}
