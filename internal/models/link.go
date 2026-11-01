package models

import (
	"net/url"
	"strings"
)

// Link type.
type Link struct {
	ID          string   `json:"id"`
	URL         string   `json:"url"`
	Title       string   `json:"title"`
	Image       string   `json:"image"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
}

// ParseTags function.
func (l *Link) ParseTags(rawTags string) *Link {
	l.Tags = []string{}

	for tag := range strings.SplitSeq(rawTags, ",") {
		tag = strings.TrimSpace(tag)
		if tag != "" {
			l.Tags = append(l.Tags, tag)
		}
	}

	return l
}

// Fill function.
func (l *Link) Fill(meta *Meta) *Link {
	if meta.OgTitle != "" {
		l.Title = meta.OgTitle
	} else if meta.OgSiteName != "" {
		l.Title = meta.OgSiteName
	} else {
		l.Title = meta.Title
	}

	if meta.OgDescription != "" {
		l.Description = meta.OgDescription
	} else {
		l.Description = meta.Description
	}

	if len(meta.OgImages) > 0 {
		l.Image = l.toAbsoluteURL(meta.OgImages[0])
	} else if len(meta.Images) > 0 {
		l.Image = l.toAbsoluteURL(meta.Images[0])
	}

	return l
}

func (l *Link) toAbsoluteURL(rawURL string) string {
	base, err := url.Parse(l.URL)
	if err != nil {
		return rawURL
	}

	ref, err := url.Parse(rawURL)
	if err != nil {
		return rawURL
	}

	return base.ResolveReference(ref).String()
}
