package gormx

import (
	"github.com/AlekSi/pointer"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUnpointer(t *testing.T) {
	tests := []struct {
		msg    string
		limit  *int64
		offset *int64
		want   *LimitAndOffset
	}{
		{
			msg:    "both zero",
			limit:  pointer.ToInt64(0),
			offset: pointer.ToInt64(0),
			want:   &LimitAndOffset{Limit: 0, Offset: 0},
		},
		{
			msg:    "both nil",
			limit:  nil,
			offset: nil,
			want:   &LimitAndOffset{Limit: DefaultLimit, Offset: 0},
		},
		{
			msg:    "both not zero",
			limit:  pointer.ToInt64(13),
			offset: pointer.ToInt64(13),
			want:   &LimitAndOffset{Limit: 13, Offset: 13},
		},
		{
			msg:    "limit nil",
			limit:  nil,
			offset: pointer.ToInt64(13),
			want:   &LimitAndOffset{Limit: DefaultLimit, Offset: 13},
		},
		{
			msg:    "offset nil",
			limit:  pointer.ToInt64(13),
			offset: nil,
			want:   &LimitAndOffset{Offset: 0, Limit: 13},
		},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.want, Unpointer(tt.offset, tt.limit), tt.msg)
	}

}
