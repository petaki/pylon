package models

import (
	"fmt"
	"net/url"
	"path"
	"strings"

	"github.com/petaki/pylon/internal/meta"
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
	l.Tags = strings.Split(rawTags, ",")

	return l
}

// Fill function.
func (l *Link) Fill(data *meta.Data) *Link {
	if data.OgTitle != "" {
		l.Title = data.OgTitle
	} else if data.OgSiteName != "" {
		l.Title = data.OgSiteName
	} else {
		l.Title = data.Title
	}

	if data.OgDescription != "" {
		l.Description = data.OgDescription
	} else {
		l.Description = data.Description
	}

	if len(data.OgImages) > 0 {
		l.Image = l.toAbsoluteURL(data.OgImages[0])
	} else if len(data.Images) > 0 {
		l.Image = l.toAbsoluteURL(data.Images[0])
	}

	return l
}

func (l *Link) toAbsoluteURL(relPath string) string {
	base, err := url.Parse(l.URL)
	if err != nil {
		return relPath
	}

	src, err := url.Parse(relPath)
	if err == nil && src.IsAbs() {
		return src.String()
	}

	if strings.HasPrefix(relPath, "/") {
		return fmt.Sprintf("%s://%s%s", base.Scheme, base.Host, relPath)
	}

	return fmt.Sprintf("%s://%s%s", base.Scheme, base.Host, path.Join(base.Path, relPath))
}
