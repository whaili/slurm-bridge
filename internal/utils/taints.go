// SPDX-FileCopyrightText: Copyright (C) SchedMD LLC.
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	corev1 "k8s.io/api/core/v1"
)

const (
	// TaintKeyBridgedNode will be added when nodes detected as bridged,
	// containing co-located kubelet and slurmd services, and removed
	// when detected as no longer bridged.
	TaintKeyBridgedNode = "slinky.slurm.net/managed-node"
)

var (
	// TaintNodeBridged will mark a node such that:
	// - Any already-running pods that do not tolerate the taint will be evicted.
	//   Currently enforced by the Kubernetes NodeController.
	// - Indicates the node is being managed by the named slurm-bridge scheduler.
	TaintNodeBridged = corev1.Taint{
		Key:    TaintKeyBridgedNode,
		Effect: corev1.TaintEffectNoExecute,
	}
)

func NewTaintNodeBridged(schedulerName string) *corev1.Taint {
	taint := TaintNodeBridged
	taint.Value = schedulerName
	return &taint
}

var (
	// TolerationNodeBridged is used to mark pods such that they can run on the
	// slurm-bridge marked nodes.
	TolerationNodeBridged = corev1.Toleration{
		Key:      TaintNodeBridged.Key,
		Operator: corev1.TolerationOpEqual,
		Effect:   TaintNodeBridged.Effect,
	}
)

func NewTolerationNodeBridged(schedulerName string) *corev1.Toleration {
	toleration := TolerationNodeBridged
	toleration.Value = schedulerName
	return &toleration
}
