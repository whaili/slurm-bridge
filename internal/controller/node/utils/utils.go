// SPDX-FileCopyrightText: Copyright (C) SchedMD LLC.
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"context"

	corev1 "k8s.io/api/core/v1"

	"github.com/SlinkyProject/slurm-bridge/internal/wellknown"
)

func MakeNodeNameMap(ctx context.Context, nodeList *corev1.NodeList) map[string]string {
	nodeNameMap := make(map[string]string, len(nodeList.Items))
	for _, node := range nodeList.Items {
		slurmNodeName := GetSlurmNodeName(&node)
		nodeNameMap[slurmNodeName] = node.GetName()
	}
	return nodeNameMap
}

func GetSlurmNodeName(node *corev1.Node) string {
	slurmNodeName, ok := node.Labels[wellknown.LabelSlurmNodeName]
	if ok {
		return slurmNodeName
	}
	return node.GetName()
}
