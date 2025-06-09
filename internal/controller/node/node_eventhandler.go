// SPDX-FileCopyrightText: Copyright (C) SchedMD LLC.
// SPDX-License-Identifier: Apache-2.0

package node

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
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	nodeutils "github.com/SlinkyProject/slurm-bridge/internal/controller/node/utils"
)

type nodeEventHandler struct {
	client.Reader
}

// Create implements handler.EventHandler.
func (h *nodeEventHandler) Create(ctx context.Context, evt event.CreateEvent, q workqueue.TypedRateLimitingInterface[reconcile.Request]) {
	// Intentionally blank
}

// Delete implements handler.EventHandler.
func (h *nodeEventHandler) Delete(ctx context.Context, evt event.DeleteEvent, q workqueue.TypedRateLimitingInterface[reconcile.Request]) {
	// Intentionally blank
}

// Generic implements handler.EventHandler.
func (h *nodeEventHandler) Generic(ctx context.Context, evt event.GenericEvent, q workqueue.TypedRateLimitingInterface[reconcile.Request]) {
	logger := log.FromContext(ctx)

	node, ok := evt.Object.(*corev1.Node)
	if !ok {
		utilruntime.HandleError(fmt.Errorf("event object is not a node %#v", evt.Object))
		return
	}

	nodeList := &corev1.NodeList{}
	if err := h.List(ctx, nodeList); err != nil {
		logger.Error(err, "failed to list nodes")
		return
	}
	nodeNameMap := nodeutils.MakeNodeNameMap(ctx, nodeList)

	name, ok := nodeNameMap[node.GetName()]
	if !ok {
		name = node.GetName()
	}
	namespacedName := types.NamespacedName{
		Name: name,
	}
	if err := h.Get(ctx, namespacedName, node); err != nil {
		logger.Error(err, "failed to get node")
		return
	}
	enqueueNode(q, node)
}

// Update implements handler.EventHandler.
func (h *nodeEventHandler) Update(ctx context.Context, evt event.UpdateEvent, q workqueue.TypedRateLimitingInterface[reconcile.Request]) {
	// Intentionally blank
}

var _ handler.EventHandler = &nodeEventHandler{}

func enqueueNode(q workqueue.TypedRateLimitingInterface[reconcile.Request], node *corev1.Node) {
	enqueueNodeAfter(q, node, 0)
}

func enqueueNodeAfter(q workqueue.TypedRateLimitingInterface[reconcile.Request], node *corev1.Node, duration time.Duration) {
	if node == nil {
		return
	}
	req := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name: node.GetName(),
		},
	}
	q.AddAfter(req, duration)
}
