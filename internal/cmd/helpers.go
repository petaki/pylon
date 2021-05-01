package cmd

import (
	"fmt"
	"net/url"
	"os"

	"github.com/petaki/support-go/cli"
)

func createMeiliSearchFlags(command *cli.Command) (*string, *string) {
	meiliSearchHost := command.FlagSet().String("meilisearch-host", os.Getenv("MEILISEARCH_HOST"), "MeiliSearch Host")
	meiliSearchAPIKey := command.FlagSet().String("meilisearch-api-key", os.Getenv("MEILISEARCH_API_KEY"), "MeiliSearch Api Key")

	return meiliSearchHost, meiliSearchAPIKey
}

func printError(err error) int {
	fmt.Println(err)

	return 1
}

func isURL(value string) bool {
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
