// SPDX-FileCopyrightText: Copyright (C) SchedMD LLC.
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"reflect"
	"testing"

	corev1 "k8s.io/api/core/v1"
)

func TestNewTaintNodeBridged(t *testing.T) {
	type args struct {
		schedulerName string
	}
	tests := []struct {
		name string
		args args
		want *corev1.Taint
	}{
		{
			name: "No scheduler name",
			args: args{
				schedulerName: "",
			},
			want: func() *corev1.Taint {
				t := TaintNodeBridged
				t.Value = ""
				return &t
			}(),
		},
		{
			name: "Taint with schedulerName foo",
			args: args{
				schedulerName: "foo",
			},
			want: func() *corev1.Taint {
				t := TaintNodeBridged
				t.Value = "foo"
				return &t
			}(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewTaintNodeBridged(tt.args.schedulerName)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewTaintNodeBridged() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewTolerationNodeBridged(t *testing.T) {
	type args struct {
		schedulerName string
	}
	tests := []struct {
		name string
		args args
		want *corev1.Toleration
	}{
		{
			name: "No scheduler name",
			args: args{
				schedulerName: "",
			},
			want: func() *corev1.Toleration {
				t := TolerationNodeBridged
				t.Value = ""
				return &t
			}(),
		},
		{
			name: "Taint with schedulerName foo",
			args: args{
				schedulerName: "foo",
			},
			want: func() *corev1.Toleration {
				t := TolerationNodeBridged
				t.Value = "foo"
				return &t
			}(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewTolerationNodeBridged(tt.args.schedulerName)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewTolerationNodeBridged() = %v, want %v", got, tt.want)
			}
		})
	}
}
