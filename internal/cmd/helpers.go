package cmd

import (
	"net/url"

	"github.com/petaki/pylon/internal/models"
	"github.com/petaki/support-go/cli"
)

func isValidURL(value string) bool {
	u, err := url.ParseRequestURI(value)
	if err != nil {
		return false
	}

	return (u.Scheme == "http" || u.Scheme == "https") && u.Host != ""
}

func printLinkTable(links []*models.Link) int {
	rows := make([][]string, len(links))

	for i, link := range links {
		rows[i] = []string{
			link.ID,
			link.URL,
			link.Title,
		}
	}

	return (&cli.Table{
		Headers: []string{
			"ID",
			"URL",
			"Title",
		},
		Rows: rows,
	}).Print()
}
