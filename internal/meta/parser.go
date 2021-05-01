package meta

import (
	"io"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

// Parse function.
func Parse(buffer io.Reader) (*Data, error) {
	data := &Data{}
	isTitleText := false

	z := html.NewTokenizer(buffer)

	for {
		tt := z.Next()

		switch tt {
		case html.ErrorToken:
			if z.Err() == io.EOF {
				return data, nil
			}

			return nil, z.Err()
		case html.TextToken:
			if isTitleText {
				data.Title = strings.TrimSpace(string(z.Text()))
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
						data.Description = attributes["content"]
					}

					continue
				}

				metaProperty, ok := attributes["property"]
				if ok {
					switch metaProperty {
					case "og:title":
						data.OgTitle = attributes["content"]
					case "og:site_name":
						data.OgSiteName = attributes["content"]
					case "og:description":
						data.OgDescription = attributes["content"]
					case "og:image":
						data.OgImages = append(data.OgImages, attributes["content"])
					}
				}

				continue
			}

			if code == atom.Img {
				imgSrc, ok := attributes["src"]
				if ok {
					data.Images = append(data.Images, imgSrc)
				}
			}
		}
	}
}
