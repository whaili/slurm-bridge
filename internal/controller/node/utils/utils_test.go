// SPDX-FileCopyrightText: Copyright (C) SchedMD LLC.
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"context"
	"testing"

	corev1 "k8s.io/api/core/v1"
	apiequality "k8s.io/apimachinery/pkg/api/equality"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/SlinkyProject/slurm-bridge/internal/wellknown"
)

func TestMakeNodeNameMap(t *testing.T) {
	ctx := context.Background()
	type args struct {
		ctx          context.Context
		kubeNodeList *corev1.NodeList
	}
	tests := []struct {
		name string
		args args
		want map[string]string
	}{
		{
			name: "No Labels",
			args: args{
				ctx: ctx,
				kubeNodeList: &corev1.NodeList{
					Items: []corev1.Node{
						{
							ObjectMeta: metav1.ObjectMeta{
								Name: "node-0",
							},
						},
						{
							ObjectMeta: metav1.ObjectMeta{
								Name: "node-1",
							},
						},
					},
				},
			},
			want: map[string]string{
				"node-0": "node-0",
				"node-1": "node-1",
			},
		},
		{
			name: "With Labels",
			args: args{
				ctx: ctx,
				kubeNodeList: &corev1.NodeList{
					Items: []corev1.Node{
						{
							ObjectMeta: metav1.ObjectMeta{
								Name: "node-0",
								Labels: map[string]string{
									wellknown.LabelSlurmNodeName: "slurm-0",
								},
							},
						},
						{
							ObjectMeta: metav1.ObjectMeta{
								Name: "node-1",
								Labels: map[string]string{
									wellknown.LabelSlurmNodeName: "slurm-1",
								},
							},
						},
					},
				},
			},
			want: map[string]string{
				"slurm-0": "node-0",
				"slurm-1": "node-1",
			},
		},
		{
			name: "Mixed",
			args: args{
				ctx: ctx,
				kubeNodeList: &corev1.NodeList{
					Items: []corev1.Node{
						{
							ObjectMeta: metav1.ObjectMeta{
								Name: "node-0",
								Labels: map[string]string{
									wellknown.LabelSlurmNodeName: "slurm-0",
								},
							},
						},
						{
							ObjectMeta: metav1.ObjectMeta{
								Name: "node-1",
							},
						},
					},
				},
			},
			want: map[string]string{
				"slurm-0": "node-0",
				"node-1":  "node-1",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MakeNodeNameMap(tt.args.ctx, tt.args.kubeNodeList); !apiequality.Semantic.DeepEqual(got, tt.want) {
				t.Errorf("MakeNodeNameMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetSlurmNodeName(t *testing.T) {
	type args struct {
		node *corev1.Node
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "No Labels",
			args: args{
				node: &corev1.Node{
					ObjectMeta: metav1.ObjectMeta{
						Name: "node-0",
					},
				},
			},
			want: "node-0",
		},
		{
			name: "With Label",
			args: args{
				node: &corev1.Node{
					ObjectMeta: metav1.ObjectMeta{
						Name: "node-0",
						Labels: map[string]string{
							wellknown.LabelSlurmNodeName: "slurm-0",
						},
					},
				},
			},
			want: "slurm-0",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetSlurmNodeName(tt.args.node); got != tt.want {
				t.Errorf("GetSlurmNodeName() = %v, want %v", got, tt.want)
			}
		})
	}
}
