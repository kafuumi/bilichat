package bilichat

import (
	"testing"
	"time"
)

func TestBuffer_Put(t *testing.T) {
	tests := []struct {
		name string
		caps int
		freq time.Duration
	}{
		{"by_time", 20, time.Second},
		{"by_caps", 5, time.Second},
		{"caps_or_time", 10, time.Second},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			item := 0
			buf := newBuffer[int](test.caps, test.freq, false, func(items []int) {
				t.Logf("items len:%d, items:%v, item:%d", len(items), items, item)
			})
			for {
				if item >= 30 {
					buf.MustFlush()
					buf.Free()
					break
				}
				buf.Put(item)
				item++
				time.Sleep(100 * time.Millisecond)
			}
		})
	}
}
