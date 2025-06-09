// SPDX-FileCopyrightText: Copyright (C) SchedMD LLC.
// SPDX-License-Identifier: Apache-2.0

package slurmjobir

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	jobset "sigs.k8s.io/jobset/api/jobset/v1alpha2"
)

var (
	// Ref: https://jobset.sigs.k8s.io/docs/
	jobSet_v1alpha2 = metav1.TypeMeta{APIVersion: "jobset.x-k8s.io/v1alpha2", Kind: "JobSet"}
)

// fromJobSet will translate a pod from a JobSet into a SlurmJobIR.
func (t *translator) fromJobSet(pod *corev1.Pod, rootPOM *metav1.PartialObjectMetadata) (*SlurmJobIR, error) {
	jobSet := &jobset.JobSet{}
	key := client.ObjectKey{Namespace: rootPOM.GetNamespace(), Name: rootPOM.GetName()}
	if err := t.Get(t.ctx, key, jobSet); err != nil {
		return nil, err
	}

	// Construct the rootPOM representing the job for this
	// pod and use fromJob to populate slurmJobIR.
	jobRootPOM := &metav1.PartialObjectMetadata{
		TypeMeta: job_v1,
		ObjectMeta: metav1.ObjectMeta{
			Name:      pod.Labels["job-name"],
			Namespace: rootPOM.GetNamespace(),
		},
	}
	slurmJobIR, err := t.fromJob(pod, jobRootPOM)
	if err != nil {
		return nil, err
	}

	return slurmJobIR, nil
}
