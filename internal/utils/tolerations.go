// SPDX-FileCopyrightText: Copyright (C) SchedMD LLC.
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	corev1 "k8s.io/api/core/v1"
)

func MergeTolerations(tolerations []corev1.Toleration, toleration corev1.Toleration) []corev1.Toleration {
	newtolerations := tolerations

	tolerationFound := false
	// Using range to iterate over the slice
	for _, curtoleration := range tolerations {
		if curtoleration.MatchToleration(&toleration) {
			tolerationFound = true
			break
		}
	}

	if !tolerationFound {
		newtolerations = append(newtolerations, toleration)
	}

	return newtolerations
}
