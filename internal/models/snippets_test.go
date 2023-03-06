package models

import (
	"testing"
	"time"
)

func TestSnippet(t *testing.T) {
	tests := []struct {
		name string
		tm   time.Time
		want string
	}{{
		name: "UTC",
		tm:   time.Date(2023, 3, 6, 10, 15, 0, 0, time.UTC),
		want: "Mar 06, 2023 at 10:15",
	},
		{
			name: "Empty",
			tm:   time.Time{},
			want: "",
		},
		{
			name: "CET",
			tm:   time.Date(2023, 3, 6, 10, 15, 0, 0, time.FixedZone("CET", 1*60*60)),
			want: "Mar 06, 2023 at 09:15",
		},
	}

	snippet := Snippet{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hd := snippet.ReadableDate(tt.tm)

			if hd != tt.want {
				t.Errorf("got %q; want %q", hd, tt.want)
			}
		})
	}
}
