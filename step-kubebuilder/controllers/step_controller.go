/*


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

package controllers

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	v1 "github.com/programming-kubernetes/cnat/step-kubebuilder/api/v1"
	webappv1 "github.com/programming-kubernetes/cnat/step-kubebuilder/api/v1"
)

// StepReconciler reconciles a Step object
type StepReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=webapp.example.com,resources=steps,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=webapp.example.com,resources=steps/status,verbs=get;update;patch

func (r *StepReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	now := metav1.Now().UTC()
	reqLogger := r.Log.WithValues("step", req.NamespacedName)
	// your logic here
	reqLogger.Info("============================ RECONCILER RUNNING==========================")

	instance := &v1.Step{}
	err := r.Get(context.TODO(), req.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request - return and don't requeue:
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request:
		return ctrl.Result{}, err
	}
	reqLogger.Info("=======INSTANCE OF STEP ====", "lastupdate", instance.Status.LastUpdate.UTC(), "now time ", now.String())
	retryAfter := instance.Status.LastUpdate
	if retryAfter.IsZero() || now.Sub(retryAfter.Time.UTC()) >= time.Duration(time.Minute.Nanoseconds()*2) {
		reqLogger.Info("=======TIME TO RUN THE STEP CR =================")
		// As long as we haven't executed the command yet, we need to check if it's time already to act:
		//reqLogger.Info("Checking schedule")
		podList := &corev1.PodList{}
		err = r.List(context.TODO(), podList, client.InNamespace(instance.Spec.PodNamespace))
		if err != nil {
			if errors.IsNotFound(err) {
				// Request object not found, could have been deleted after reconcile request - return and don't requeue:
				return ctrl.Result{}, nil
			}
			// Error reading the object - requeue the request:
			return ctrl.Result{}, err
		}
		pod := newPodForCR(instance, podList)
		// Set At instance as the owner and controller
		reqLogger.Info("Pod Created", "name", pod.Name)
		if err := controllerutil.SetControllerReference(instance, pod, r.Scheme); err != nil {
			// requeue with error
			return ctrl.Result{}, err
		}
		reqLogger.Info("Pod Assigned to controller", "name", pod.Name)
		Dupfound := &corev1.Pod{}
		err = r.Get(context.TODO(), types.NamespacedName{Name: pod.Name, Namespace: pod.Namespace}, Dupfound)
		// Try to see if the pod already exists and if not
		// (which we expect) then create a one-shot pod as per spec:
		if err != nil && errors.IsNotFound(err) {
			reqLogger.Info("Going to create a new pod since not found", "name", pod.Name)
			err = r.Create(context.TODO(), pod)
			if err != nil {
				// requeue with error
				return ctrl.Result{}, err
			}
			reqLogger.Info("Pod launched", "name", pod.Name)
		} else if err != nil {
			// requeue with error
			reqLogger.Info("Requeue with error")
			return ctrl.Result{}, err
		} else if Dupfound.Status.Phase == corev1.PodFailed || Dupfound.Status.Phase == corev1.PodSucceeded {
			reqLogger.Info("Container terminated", "reason", Dupfound.Status.Reason, "message", Dupfound.Status.Message)

		} else {
			// don't requeue because it will happen automatically when the pod status changes
			return ctrl.Result{}, nil
		}
		// case v1.PhaseDone:
		// 	reqLogger.Info("Phase: DONE")
		// 	return ctrl.Result{}, nil
		// default:
		// 	reqLogger.Info("NOP")
		// 	return ctrl.Result{}, nil
		// }

		instance.Status.Phase = v1.PhaseDone
		instance.Status.LastUpdate = metav1.Now()
		err = r.Status().Update(ctx, instance)
		if err != nil {
			// Error reading the object - requeue the request:
			return ctrl.Result{}, err
		}
		reqLogger.Info("UPDATED THE POD STATUS", "LAST UPDATE", instance.Status.LastUpdate.UTC().String())
		return ctrl.Result{RequeueAfter: time.Duration(time.Minute.Nanoseconds() * 2)}, nil
	}
	abc := fmt.Sprintf("%f", now.Sub(retryAfter.Time.UTC()).Minutes())
	reqLogger.Info("Still TIme to perform the operation", "duration", abc)
	return ctrl.Result{RequeueAfter: now.Sub(retryAfter.Time.UTC())}, nil
}

func (r *StepReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&webappv1.Step{}).
		Complete(r)
}

func newPodForCR(cr *v1.Step, pl *corev1.PodList) *corev1.Pod {
	labels := map[string]string{
		"app": cr.Name,
	}
	rand.Seed(time.Now().UnixNano())
	chars := []rune("abcdefghijklmnopqrstuvwxyz" +
		"0123456789")
	length := 8
	var b strings.Builder
	for i := 0; i < length; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}
	str := b.String()
	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name + "-pod" + str,
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:    "busybox",
					Image:   "busybox",
					Command: strings.Split(getPodList(pl), " "),
				},
			},
			RestartPolicy: corev1.RestartPolicyOnFailure,
		},
	}
}

func getPodList(pl *corev1.PodList) string {
	s := "echo "
	for _, p := range pl.Items {
		s = s + p.Name + "\n"
	}
	return s
}
