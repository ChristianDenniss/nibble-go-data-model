package entity

import (
	"errors"
	"testing"
)

func TestNormalizeImageURL(t *testing.T) {
	cases := []struct {
		raw     string
		want    string
		wantErr bool
	}{
		{raw: "", want: ""},
		{raw: "   ", want: ""},
		{raw: " https://cdn.example.com/a.jpg ", want: "https://cdn.example.com/a.jpg"},
		{raw: "http://localhost:9000/bucket/a.png", want: "http://localhost:9000/bucket/a.png"},
		{raw: "javascript:alert(1)", wantErr: true},
		{raw: "data:image/png;base64,AAAA", wantErr: true},
		{raw: "/images/a.jpg", wantErr: true},
		{raw: "https://", wantErr: true},
	}
	for _, tc := range cases {
		got, err := NormalizeImageURL(tc.raw)
		if tc.wantErr {
			if !errors.Is(err, ErrImageURLInvalid) {
				t.Errorf("NormalizeImageURL(%q) err = %v, want ErrImageURLInvalid", tc.raw, err)
			}
			continue
		}
		if err != nil || got != tc.want {
			t.Errorf("NormalizeImageURL(%q) = %q, %v; want %q", tc.raw, got, err, tc.want)
		}
	}
}
