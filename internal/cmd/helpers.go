package cmd

import (
	"net/url"
	"os"

	"github.com/petaki/support-go/cli"
)

func createMeiliSearchFlags(command *cli.Command) (*string, *string, *string) {
	meiliSearchHost := command.FlagSet().String("meilisearch-host", os.Getenv("MEILISEARCH_HOST"), "MeiliSearch Host")
	meiliSearchAPIKey := command.FlagSet().String("meilisearch-api-key", os.Getenv("MEILISEARCH_API_KEY"), "MeiliSearch API Key")
	meiliSearchIndex := command.FlagSet().String("meilisearch-index", os.Getenv("MEILISEARCH_INDEX"), "MeiliSearch Index")

	return meiliSearchHost, meiliSearchAPIKey, meiliSearchIndex
}

func createLinkFlags(command *cli.Command) (*string, *string) {
	headlessShellHost := command.FlagSet().String("headless-shell-host", os.Getenv("HEADLESS_SHELL_HOST"), "Headless Shell Host")
	tags := command.FlagSet().String("tags", "", "Link Tags")

	return headlessShellHost, tags
}

func isValidURL(value string) bool {
	_, err := url.ParseRequestURI(value)
	if err != nil {
		return false
	}

	u, err := url.Parse(value)
	if err != nil {
		return false
	}

	if u.Scheme == "" {
		return false
	}

	if u.Host == "" {
		return false
	}

	return true
}
