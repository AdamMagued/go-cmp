// Copyright 2026, The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmpopts_test

import (
	"math"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestEquateApproxDuration(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		x, y      time.Duration
		margin    time.Duration
		wantEqual bool
	}{
		{
			name:      "identical zero durations with zero margin",
			x:         0,
			y:         0,
			margin:    0,
			wantEqual: true,
		},
		{
			name:      "identical positive durations with zero margin",
			x:         5 * time.Second,
			y:         5 * time.Second,
			margin:    0,
			wantEqual: true,
		},
		{
			name:      "identical negative durations with zero margin",
			x:         -5 * time.Second,
			y:         -5 * time.Second,
			margin:    0,
			wantEqual: true,
		},
		{
			name:      "different durations with zero margin",
			x:         5 * time.Second,
			y:         5*time.Second + 1,
			margin:    0,
			wantEqual: false,
		},
		{
			name:      "difference strictly within margin",
			x:         10 * time.Millisecond,
			y:         12 * time.Millisecond,
			margin:    5 * time.Millisecond,
			wantEqual: true,
		},
		{
			name:      "difference exactly at positive margin",
			x:         10 * time.Millisecond,
			y:         15 * time.Millisecond,
			margin:    5 * time.Millisecond,
			wantEqual: true,
		},
		{
			name:      "difference exactly at negative margin",
			x:         15 * time.Millisecond,
			y:         10 * time.Millisecond,
			margin:    5 * time.Millisecond,
			wantEqual: true,
		},
		{
			name:      "difference strictly exceeding margin positive",
			x:         10 * time.Millisecond,
			y:         15*time.Millisecond + 1,
			margin:    5 * time.Millisecond,
			wantEqual: false,
		},
		{
			name:      "difference strictly exceeding margin negative",
			x:         15*time.Millisecond + 1,
			y:         10 * time.Millisecond,
			margin:    5 * time.Millisecond,
			wantEqual: false,
		},
		{
			name:      "cross zero within margin",
			x:         -2 * time.Second,
			y:         3 * time.Second,
			margin:    5 * time.Second,
			wantEqual: true,
		},
		{
			name:      "cross zero exceeding margin",
			x:         -2 * time.Second,
			y:         3 * time.Second,
			margin:    5*time.Second - 1,
			wantEqual: false,
		},
		{
			name:      "boundary MinInt64 equal to itself",
			x:         time.Duration(math.MinInt64),
			y:         time.Duration(math.MinInt64),
			margin:    0,
			wantEqual: true,
		},
		{
			name:      "boundary MaxInt64 equal to itself",
			x:         time.Duration(math.MaxInt64),
			y:         time.Duration(math.MaxInt64),
			margin:    0,
			wantEqual: true,
		},
		{
			name:      "boundary MinInt64 offset within margin",
			x:         time.Duration(math.MinInt64),
			y:         time.Duration(math.MinInt64) + 10*time.Second,
			margin:    10 * time.Second,
			wantEqual: true,
		},
		{
			name:      "boundary MinInt64 offset exceeding margin",
			x:         time.Duration(math.MinInt64),
			y:         time.Duration(math.MinInt64) + 10*time.Second + 1,
			margin:    10 * time.Second,
			wantEqual: false,
		},
		{
			name:      "boundary MaxInt64 offset within margin",
			x:         time.Duration(math.MaxInt64) - 10*time.Second,
			y:         time.Duration(math.MaxInt64),
			margin:    10 * time.Second,
			wantEqual: true,
		},
		{
			name:      "boundary MaxInt64 offset exceeding margin",
			x:         time.Duration(math.MaxInt64) - 10*time.Second - 1,
			y:         time.Duration(math.MaxInt64),
			margin:    10 * time.Second,
			wantEqual: false,
		},
		{
			name:      "boundary MinInt64 and MaxInt64 do not overflow to equal",
			x:         time.Duration(math.MinInt64),
			y:         time.Duration(math.MaxInt64),
			margin:    time.Second,
			wantEqual: false,
		},
		{
			name:      "boundary MinInt64 and 0 do not overflow to equal",
			x:         time.Duration(math.MinInt64),
			y:         0,
			margin:    time.Second,
			wantEqual: false,
		},
		{
			name:      "boundary 0 and MaxInt64 do not overflow to equal",
			x:         0,
			y:         time.Duration(math.MaxInt64),
			margin:    time.Second,
			wantEqual: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			opt := cmpopts.EquateApproxDuration(tc.margin)
			gotEqual := cmp.Equal(tc.x, tc.y, opt)
			if gotEqual != tc.wantEqual {
				t.Errorf("cmp.Equal(%v, %v, margin=%v) = %v; want %v", tc.x, tc.y, tc.margin, gotEqual, tc.wantEqual)
			}

			// Verify symmetry: Equal(x, y) == Equal(y, x)
			gotSymmetric := cmp.Equal(tc.y, tc.x, opt)
			if gotSymmetric != gotEqual {
				t.Errorf("symmetry violated: cmp.Equal(%v, %v) = %v, but cmp.Equal(%v, %v) = %v",
					tc.x, tc.y, gotEqual, tc.y, tc.x, gotSymmetric)
			}
		})
	}
}

func TestEquateApproxDuration_Panic(t *testing.T) {
	t.Parallel()

	negativeMargins := []time.Duration{
		-1,
		-1 * time.Second,
		time.Duration(math.MinInt64),
	}

	for _, margin := range negativeMargins {
		t.Run(margin.String(), func(t *testing.T) {
			defer func() {
				r := recover()
				if r == nil {
					t.Fatalf("EquateApproxDuration(%v) did not panic; want panic", margin)
				}
				const wantMsg = "margin must be a non-negative number"
				if r != wantMsg {
					t.Fatalf("EquateApproxDuration(%v) panic = %q; want %q", margin, r, wantMsg)
				}
			}()
			_ = cmpopts.EquateApproxDuration(margin)
		})
	}
}

func TestEquateApproxDuration_StructField(t *testing.T) {
	t.Parallel()

	type Task struct {
		Name    string
		Elapsed time.Duration
	}

	t1 := Task{Name: "task-a", Elapsed: 100 * time.Millisecond}
	t2 := Task{Name: "task-a", Elapsed: 105 * time.Millisecond}
	t3 := Task{Name: "task-a", Elapsed: 120 * time.Millisecond}

	opt := cmpopts.EquateApproxDuration(10 * time.Millisecond)

	if diff := cmp.Diff(t1, t2, opt); diff != "" {
		t.Errorf("unexpected diff for t1 and t2 (-want +got):\n%s", diff)
	}

	if diff := cmp.Diff(t1, t3, opt); diff == "" {
		t.Errorf("expected diff for t1 and t3, but got none")
	}
}
