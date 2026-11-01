package models

import (
	"slices"
	"testing"
)

func TestLinkParseTags(t *testing.T) {
	tests := []struct {
		rawTags string
		want    []string
	}{
		{rawTags: "", want: []string{}},
		{rawTags: "go, cli,,web ", want: []string{"go", "cli", "web"}},
	}

	for _, tt := range tests {
		got := (&Link{}).ParseTags(tt.rawTags).Tags
		if !slices.Equal(got, tt.want) || got == nil {
			t.Errorf("ParseTags(%q) = %#v, want %#v", tt.rawTags, got, tt.want)
		}
	}
}

func TestLinkFillImage(t *testing.T) {
	tests := []struct {
		url   string
		image string
		want  string
	}{
		{url: "https://example.com", image: "logo.png", want: "https://example.com/logo.png"},
		{url: "https://example.com/a/b.html", image: "c.png", want: "https://example.com/a/c.png"},
		{url: "https://example.com/a/", image: "/c.png", want: "https://example.com/c.png"},
		{url: "https://example.com", image: "//cdn.example.com/c.png", want: "https://cdn.example.com/c.png"},
	}

	for _, tt := range tests {
		got := (&Link{URL: tt.url}).Fill(&Meta{OgImages: []string{tt.image}}).Image
		if got != tt.want {
			t.Errorf("Fill(%q, %q) = %q, want %q", tt.url, tt.image, got, tt.want)
		}
	}
}
