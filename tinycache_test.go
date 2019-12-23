package tinycache

import (
	"github.com/c2h5oh/datasize"
	"testing"
	"time"
)

func TestTinyCache_SetMaxMemory(t *testing.T) {
	type wantResult struct {
		ret  bool
		size datasize.ByteSize
	}
	tests := []struct {
		size string
		want wantResult
	}{
		{"1KB", wantResult{true, 1 * datasize.KB}},
		{"1kB", wantResult{true, 1 * datasize.KB}},
		{"100KB", wantResult{true, 100 * datasize.KB}},
		{"1MB", wantResult{true, 1 * datasize.MB}},
		{"2MB", wantResult{true, 2 * datasize.MB}},
		{"1GB", wantResult{true, 1 * datasize.GB}},
		{"1TB", wantResult{true, 1 * datasize.TB}},
		{"0KB", wantResult{false, 0}},
	}
	for _, tt := range tests {
		t.Run(tt.size, func(t *testing.T) {
			c := &tinyCache{}
			if got := c.SetMaxMemory(tt.size); got != tt.want.ret || c.maxMem != tt.want.size {
				t.Errorf("SetMaxMemory() = %v, ret %v", got, true)
			}
		})
	}
}

func TestTinyCache_Set(t *testing.T) {
	type fields struct {
		maxMem datasize.ByteSize
		items  map[string]*cacheItem
	}
	type args struct {
		key    string
		val    interface{}
		expire time.Duration
	}
	tests := []struct {
		fields fields
		args   args
		want   bool
	}{
		{
			args: args{
				key: "string",
				val: "string value",
			},
			want: true,
		},
		{
			args: args{
				key: "slice",
				val: &[]string{
					"string1",
					"string2",
				},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.args.key, func(t *testing.T) {
			c := New()
			if c.maxMem == 0 {
				c.maxMem = 100 * datasize.MB
			}
			if c.items == nil {
				c.items = make(map[string]*cacheItem)
			}

			c.Set(tt.args.key, tt.args.val, tt.args.expire)
			item, ok := c.items[tt.args.key]
			if ok != tt.want {
				t.Errorf("Set(%v) failed", tt.args)
				return
			}
			if item.data != tt.args.val {
				t.Errorf("Set(%v) error", tt.args)
			}
		})
	}
}

func TestTinyCache_Expiration(t *testing.T) {
	cache := New()
	cache.Set("key", "value", 1*time.Millisecond)
	time.Sleep(2 * time.Millisecond)
	if v, ok := cache.Get("key"); ok {
		t.Errorf("expire item(%v) = %v failed", "key", v)
	}
}
