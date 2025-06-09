// SPDX-FileCopyrightText: Copyright (C) SchedMD LLC.
// SPDX-License-Identifier: Apache-2.0

package placeholderinfo

import (
	"testing"

	apiequality "k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/utils/ptr"
)

func TestPlaceholderInfo_Equal(t *testing.T) {
	type fields struct {
		Pods []string
	}
	type args struct {
		cmp PlaceholderInfo
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "Empty",
			fields: fields{
				Pods: []string{},
			},
			args: args{
				cmp: PlaceholderInfo{
					Pods: []string{},
				},
			},
			want: true,
		},
		{
			name: "Equal",
			fields: fields{
				Pods: []string{"bar/foo"},
			},
			args: args{
				cmp: PlaceholderInfo{
					Pods: []string{"bar/foo"},
				},
			},
			want: true,
		},
		{
			name: "Not equal",
			fields: fields{
				Pods: []string{"bar/foo"},
			},
			args: args{
				cmp: PlaceholderInfo{
					Pods: []string{"buz/biz"},
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			phInfo := &PlaceholderInfo{
				Pods: tt.fields.Pods,
			}
			if got := phInfo.Equal(tt.args.cmp); got != tt.want {
				t.Errorf("PlaceholderInfo.Equal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPlaceholderInfo_ToString(t *testing.T) {
	type fields struct {
		Pods []string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Empty",
			fields: fields{
				Pods: []string{},
			},
			want: `{"pods":[]}`,
		},
		{
			name: "With a pod",
			fields: fields{
				Pods: []string{"bar/foo"},
			},
			want: `{"pods":["bar/foo"]}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			podInfo := &PlaceholderInfo{
				Pods: tt.fields.Pods,
			}
			if got := podInfo.ToString(); got != tt.want {
				t.Errorf("PlaceholderInfo.ToString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseIntoPlaceholderInfo(t *testing.T) {
	type args struct {
		str *string
		out *PlaceholderInfo
	}
	tests := []struct {
		name    string
		args    args
		want    *PlaceholderInfo
		wantErr bool
	}{
		{
			name: "Empty string",
			args: args{
				str: ptr.To(""),
				out: &PlaceholderInfo{},
			},
			want:    &PlaceholderInfo{},
			wantErr: true,
		},
		{
			name: "Empty values",
			args: args{
				str: ptr.To(`{"pods":[]}`),
				out: &PlaceholderInfo{},
			},
			want:    &PlaceholderInfo{Pods: []string{}},
			wantErr: false,
		},
		{
			name: "Single pod",
			args: args{
				str: ptr.To(`{"pods":["bar/foo"]}`),
				out: &PlaceholderInfo{},
			},
			want:    &PlaceholderInfo{Pods: []string{"bar/foo"}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ParseIntoPlaceholderInfo(tt.args.str, tt.args.out); (err != nil) != tt.wantErr {
				t.Errorf("ParseIntoPlaceholderInfo() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got := tt.args.out; !apiequality.Semantic.DeepEqual(got, tt.want) {
				t.Errorf("ParseIntoPlaceholderInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}
