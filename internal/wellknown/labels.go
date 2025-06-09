// SPDX-FileCopyrightText: Copyright (C) SchedMD LLC.
// SPDX-License-Identifier: Apache-2.0

package wellknown

const (
	// LabelSlurmNodeName indicates the Slurm NodeName which corresponds to the
	// labeled Kubernetes node.
	LabelSlurmNodeName = "slinky.slurm.net/slurm-nodename"

	// LabelPlaceholderJobId indicates the Slurm JobId which corresponds to the
	// the pod's placeholder job.
	LabelPlaceholderJobId = "scheduler.slinky.slurm.net/slurm-jobid"
)
