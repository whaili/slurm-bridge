// SPDX-FileCopyrightText: Copyright (C) SchedMD LLC.
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"slices"
	"testing"

	apiequality "k8s.io/apimachinery/pkg/api/equality"
)

func TestKeys(t *testing.T) {
	type args struct {
		items map[string]int32
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "Empty Map",
			args: args{
				items: map[string]int32{},
			},
			want: []string{},
		},
		{
			name: "One Item Map",
			args: args{
				items: map[string]int32{
					"foo": 0,
				},
			},
			want: []string{"foo"},
		},
		{
			name: "Two Item Map",
			args: args{
				items: map[string]int32{
					"foo": 0,
					"bar": 1,
				},
			},
			want: []string{"foo", "bar"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Keys(tt.args.items)
			slices.Sort(got)
			slices.Sort(tt.want)
			if !apiequality.Semantic.DeepEqual(got, tt.want) {
				t.Errorf("Keys() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValues(t *testing.T) {
	type args struct {
		items map[string]int32
	}
	tests := []struct {
		name string
		args args
		want []int32
	}{
		{
			name: "Empty Map",
			args: args{
				items: map[string]int32{},
			},
			want: []int32{},
		},
		{
			name: "One Item Map",
			args: args{
				items: map[string]int32{
					"foo": 0,
				},
			},
			want: []int32{0},
		},
		{
			name: "Two Item Map",
			args: args{
				items: map[string]int32{
					"foo": 0,
					"bar": 1,
				},
			},
			want: []int32{0, 1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Values(tt.args.items)
			slices.Sort(got)
			slices.Sort(tt.want)
			if !apiequality.Semantic.DeepEqual(got, tt.want) {
				t.Errorf("Values() = %v, want %v", got, tt.want)
			}
		})
	}
}
