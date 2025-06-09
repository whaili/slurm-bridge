// SPDX-FileCopyrightText: Copyright (C) SchedMD LLC.
// SPDX-License-Identifier: Apache-2.0

package node

import (
	"context"
	"testing"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/util/workqueue"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func newQueue() workqueue.TypedRateLimitingInterface[reconcile.Request] {
	return workqueue.NewTypedRateLimitingQueue(workqueue.DefaultTypedControllerRateLimiter[reconcile.Request]())
}

func Test_nodeEventHandler_Create(t *testing.T) {
	type fields struct {
		Reader client.Reader
	}
	type args struct {
		ctx context.Context
		evt event.CreateEvent
		q   workqueue.TypedRateLimitingInterface[reconcile.Request]
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
	}{
		{
			name: "Empty",
			fields: fields{
				Reader: fake.NewFakeClient(),
			},
			args: args{
				ctx: context.TODO(),
				evt: event.CreateEvent{},
				q:   newQueue(),
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &nodeEventHandler{
				Reader: tt.fields.Reader,
			}
			h.Create(tt.args.ctx, tt.args.evt, tt.args.q)
			if got := tt.args.q.Len(); got > tt.want {
				t.Errorf("Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_nodeEventHandler_Delete(t *testing.T) {
	type fields struct {
		Reader client.Reader
	}
	type args struct {
		ctx context.Context
		evt event.DeleteEvent
		q   workqueue.TypedRateLimitingInterface[reconcile.Request]
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
	}{
		{
			name: "Empty",
			fields: fields{
				Reader: fake.NewFakeClient(),
			},
			args: args{
				ctx: context.TODO(),
				evt: event.DeleteEvent{},
				q:   newQueue(),
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &nodeEventHandler{
				Reader: tt.fields.Reader,
			}
			h.Delete(tt.args.ctx, tt.args.evt, tt.args.q)
			if got := tt.args.q.Len(); got > tt.want {
				t.Errorf("Delete() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_nodeEventHandler_Generic(t *testing.T) {
	type fields struct {
		Reader client.Reader
	}
	type args struct {
		ctx context.Context
		evt event.GenericEvent
		q   workqueue.TypedRateLimitingInterface[reconcile.Request]
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
	}{
		{
			name: "Empty",
			fields: fields{
				Reader: fake.NewFakeClient(),
			},
			args: args{
				ctx: context.TODO(),
				evt: event.GenericEvent{},
				q:   newQueue(),
			},
			want: 0,
		},
		{
			name: "Populated",
			fields: fields{
				Reader: fake.NewFakeClient(),
			},
			args: args{
				ctx: context.TODO(),
				evt: event.GenericEvent{
					Object: &corev1.Node{
						ObjectMeta: metav1.ObjectMeta{
							Name: "node-0",
						},
					},
				},
				q: newQueue(),
			},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &nodeEventHandler{
				Reader: tt.fields.Reader,
			}
			h.Generic(tt.args.ctx, tt.args.evt, tt.args.q)
			if got := tt.args.q.Len(); got > tt.want {
				t.Errorf("Generic() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_nodeEventHandler_Update(t *testing.T) {
	type fields struct {
		Reader client.Reader
	}
	type args struct {
		ctx context.Context
		evt event.UpdateEvent
		q   workqueue.TypedRateLimitingInterface[reconcile.Request]
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
	}{
		{
			name: "Empty",
			fields: fields{
				Reader: fake.NewFakeClient(),
			},
			args: args{
				ctx: context.TODO(),
				evt: event.UpdateEvent{},
				q:   newQueue(),
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &nodeEventHandler{
				Reader: tt.fields.Reader,
			}
			h.Update(tt.args.ctx, tt.args.evt, tt.args.q)
			if got := tt.args.q.Len(); got > tt.want {
				t.Errorf("Update() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_enqueueNode(t *testing.T) {
	type args struct {
		q    workqueue.TypedRateLimitingInterface[reconcile.Request]
		node *corev1.Node
	}
	tests := []struct {
		name    string
		args    args
		enqueue bool
	}{
		{
			name: "do nothing",
			args: args{
				q:    newQueue(),
				node: nil,
			},
			enqueue: false,
		},
		{
			name: "enqueue node",
			args: args{
				q: newQueue(),
				node: &corev1.Node{
					ObjectMeta: metav1.ObjectMeta{
						Name: "node-0",
					},
				},
			},
			enqueue: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			enqueueNode(tt.args.q, tt.args.node)
			if tt.args.q.Len() > 0 != tt.enqueue {
				t.Errorf("enqueueNode() have = %v, enqueue %v", tt.args.q.Len(), tt.enqueue)
			}
		})
	}
}

func Test_enqueueNodeAfter(t *testing.T) {
	type args struct {
		q        workqueue.TypedRateLimitingInterface[reconcile.Request]
		node     *corev1.Node
		duration time.Duration
	}
	tests := []struct {
		name    string
		args    args
		enqueue bool
	}{
		{
			name: "do nothing",
			args: args{q: newQueue(),
				node: nil,
			},
			enqueue: false,
		},
		{
			name: "enqueue node",
			args: args{
				q: newQueue(),
				node: &corev1.Node{
					ObjectMeta: metav1.ObjectMeta{
						Name: "node-0",
					},
				},
			},
			enqueue: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			enqueueNodeAfter(tt.args.q, tt.args.node, tt.args.duration)
			if tt.args.q.Len() > 0 != tt.enqueue {
				t.Errorf("enqueueNode() have = %v, enqueue %v", tt.args.q.Len(), tt.enqueue)
			}
		})
	}
}
