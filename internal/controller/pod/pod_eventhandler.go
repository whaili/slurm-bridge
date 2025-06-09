// SPDX-FileCopyrightText: Copyright (C) SchedMD LLC.
// SPDX-License-Identifier: Apache-2.0

package pod

import (
	"context"
	"fmt"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/util/workqueue"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type podEventHandler struct {
	client.Reader
	SchedulerName string
}

// Create implements handler.EventHandler.
func (h *podEventHandler) Create(ctx context.Context, evt event.CreateEvent, q workqueue.TypedRateLimitingInterface[reconcile.Request]) {
	pod, ok := evt.Object.(*corev1.Pod)
	if !ok {
		utilruntime.HandleError(fmt.Errorf("event object is not a pod %#v", evt.Object))
		return
	}
	if h.isManagedPod(pod) {
		enqueuePod(q, pod)
	}
}

// Delete implements handler.EventHandler.
func (h *podEventHandler) Delete(ctx context.Context, evt event.DeleteEvent, q workqueue.TypedRateLimitingInterface[reconcile.Request]) {
	pod, ok := evt.Object.(*corev1.Pod)
	if !ok {
		utilruntime.HandleError(fmt.Errorf("event object is not a pod %#v", evt.Object))
		return
	}
	if h.isManagedPod(pod) {
		enqueuePod(q, pod)
	}
}

// Generic implements handler.EventHandler.
func (h *podEventHandler) Generic(ctx context.Context, evt event.GenericEvent, q workqueue.TypedRateLimitingInterface[reconcile.Request]) {
	pod, ok := evt.Object.(*corev1.Pod)
	if !ok {
		utilruntime.HandleError(fmt.Errorf("event object is not a pod %#v", evt.Object))
		return
	}
	enqueuePod(q, pod)
}

// Update implements handler.EventHandler.
func (h *podEventHandler) Update(ctx context.Context, evt event.UpdateEvent, q workqueue.TypedRateLimitingInterface[reconcile.Request]) {
	podOld, ok := evt.ObjectOld.(*corev1.Pod)
	if !ok {
		utilruntime.HandleError(fmt.Errorf("event object (new) is not a pod %#v", evt.ObjectNew))
		return
	}
	podNew, ok := evt.ObjectNew.(*corev1.Pod)
	if !ok {
		utilruntime.HandleError(fmt.Errorf("event object (old) is not a pod %#v", evt.ObjectOld))
		return
	}
	if h.isManagedPod(podNew) {
		enqueuePod(q, podNew)
	} else if h.isManagedPod(podOld) {
		enqueuePod(q, podOld)
	}
}

func (h *podEventHandler) isManagedPod(pod *corev1.Pod) bool {
	if pod == nil {
		return false
	}
	return pod.Spec.SchedulerName == h.SchedulerName
}

var _ handler.EventHandler = &podEventHandler{}

func enqueuePod(q workqueue.TypedRateLimitingInterface[reconcile.Request], pod *corev1.Pod) {
	enqueuePodAfter(q, pod, 0)
}

func enqueuePodAfter(q workqueue.TypedRateLimitingInterface[reconcile.Request], pod *corev1.Pod, duration time.Duration) {
	if pod == nil {
		return
	}
	req := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Namespace: pod.GetNamespace(),
			Name:      pod.GetName(),
		},
	}
	q.AddAfter(req, duration)
}
