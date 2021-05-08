package models

import (
	"io"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

// Meta type.
type Meta struct {
	Title         string
	Description   string
	Images        []string
	OgTitle       string
	OgSiteName    string
	OgDescription string
	OgImages      []string
}

// ParseMeta function.
func ParseMeta(buffer io.Reader) (*Meta, error) {
	m := new(Meta)
	isTitleText := false

	z := html.NewTokenizer(buffer)

	for {
		tt := z.Next()

		switch tt {
		case html.ErrorToken:
			if z.Err() == io.EOF {
				return m, nil
			}

			return nil, z.Err()
		case html.TextToken:
			if isTitleText {
				m.Title = strings.TrimSpace(string(z.Text()))
			}
		case html.StartTagToken, html.SelfClosingTagToken, html.EndTagToken:
			name, hasAttr := z.TagName()
			code := atom.Lookup(name)

			if code == atom.Title {
				isTitleText = tt == html.StartTagToken

				continue
			}

			if !hasAttr {
				continue
			}

			if code != atom.Meta && code != atom.Img {
				continue
			}

			attributes := make(map[string]string)

			var key, val []byte
			for hasAttr {
				key, val, hasAttr = z.TagAttr()
				attributes[atom.String(key)] = string(val)
			}

			if code == atom.Meta {
				metaName, ok := attributes["name"]
				if ok {
					switch metaName {
					case "description":
						m.Description = attributes["content"]
					}

					continue
				}

				metaProperty, ok := attributes["property"]
				if ok {
					switch metaProperty {
					case "og:title":
						m.OgTitle = attributes["content"]
					case "og:site_name":
						m.OgSiteName = attributes["content"]
					case "og:description":
						m.OgDescription = attributes["content"]
					case "og:image":
						m.OgImages = append(m.OgImages, attributes["content"])
					}
				}

				continue
			}

			if code == atom.Img {
				imgSrc, ok := attributes["src"]
				if ok {
					m.Images = append(m.Images, imgSrc)
				}
			}
		}
	}
}
