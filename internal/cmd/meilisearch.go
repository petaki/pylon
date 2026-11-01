package cmd

import (
	"net/http"

	"github.com/meilisearch/meilisearch-go"
	"github.com/petaki/support-go/cli"
)

type meiliSearch struct {
	host   *string
	apiKey *string
	index  *string
}

func newMeiliSearch(command *cli.Command) *meiliSearch {
	return &meiliSearch{
		host:   command.String("meilisearch-host", "MEILISEARCH_HOST", "http://127.0.0.1:7700", "MeiliSearch Host"),
		apiKey: command.Secret("meilisearch-api-key", "MEILISEARCH_API_KEY", "MeiliSearch API Key"),
		index:  command.String("meilisearch-index", "MEILISEARCH_INDEX", "pylon", "MeiliSearch Index"),
	}
}

// Index function.
func (m *meiliSearch) Index() meilisearch.IndexManager {
	return meilisearch.New(
		*m.host,
		meilisearch.WithAPIKey(*m.apiKey),
		meilisearch.WithCustomClient(&http.Client{Timeout: requestTimeout}),
	).Index(*m.index)
}
