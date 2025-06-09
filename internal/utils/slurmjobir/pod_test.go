// SPDX-FileCopyrightText: Copyright (C) SchedMD LLC.
// SPDX-License-Identifier: Apache-2.0

package slurmjobir

import (
	"context"
	"testing"

	corev1 "k8s.io/api/core/v1"
	apiequality "k8s.io/apimachinery/pkg/api/equality"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func Test_translator_fromPod(t *testing.T) {
	type fields struct {
		Reader client.Reader
		ctx    context.Context
	}
	type args struct {
		pod *corev1.Pod
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *SlurmJobIR
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &translator{
				Reader: tt.fields.Reader,
				ctx:    tt.fields.ctx,
			}
			got, err := tr.fromPod(tt.args.pod)
			if (err != nil) != tt.wantErr {
				t.Errorf("translator.fromPod() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !apiequality.Semantic.DeepEqual(got, tt.want) {
				t.Errorf("translator.fromPod() = %v, want %v", got, tt.want)
			}
		})
	}
}
