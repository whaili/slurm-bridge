// SPDX-FileCopyrightText: Copyright (C) SchedMD LLC.
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"reflect"
	"testing"

	corev1 "k8s.io/api/core/v1"
)

func TestMergeTolerations(t *testing.T) {
	toleration := NewTolerationNodeBridged("foo")
	type args struct {
		tolerations []corev1.Toleration
		toleration  corev1.Toleration
	}
	tests := []struct {
		name string
		args args
		want []corev1.Toleration
	}{
		{
			name: "Toleration added",
			args: args{
				tolerations: []corev1.Toleration{},
				toleration:  *toleration,
			},
			want: []corev1.Toleration{*toleration},
		},
		{
			name: "Toleration already exists",
			args: args{
				tolerations: []corev1.Toleration{*toleration},
				toleration:  *toleration,
			},
			want: []corev1.Toleration{*toleration},
		},
		{
			name: "Toleration added and other toleration remains",
			args: args{
				tolerations: []corev1.Toleration{{Key: "foo"}, *toleration},
				toleration:  *toleration,
			},
			want: []corev1.Toleration{{Key: "foo"}, *toleration},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MergeTolerations(tt.args.tolerations, tt.args.toleration); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MergeTolerations() = %v, want %v", got, tt.want)
			}
		})
	}
}
