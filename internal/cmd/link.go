package cmd

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/meilisearch/meilisearch-go"
	"github.com/petaki/pylon/internal/models"
	"github.com/petaki/support-go/cli"
)

// LinkAdd command.
func LinkAdd(group *cli.Group, command *cli.Command, arguments []string) int {
	search := newMeiliSearch(command)
	linkCrawler := newCrawler(command)
	tags := command.FlagSet().String("tags", "", "Link Tags")

	parsed, err := command.Parse(arguments)
	if err != nil {
		return command.PrintErrorWithHelp(group, err)
	}

	if !isValidURL(parsed[0]) {
		return command.PrintError(fmt.Errorf("%w: %s", ErrInvalidURL, parsed[0]))
	}

	meta, url, err := linkCrawler.Crawl(parsed[0])
	if err != nil {
		return command.PrintError(err)
	}

	link := (&models.Link{
		ID:  uuid.NewString(),
		URL: url,
	}).ParseTags(*tags).Fill(meta)

	fmt.Println("=> Send data to " + cli.Green("MeiliSearch"))

	_, err = search.Index().AddDocuments([]*models.Link{
		link,
	}, nil)
	if err != nil {
		return command.PrintError(err)
	}

	return printLinkTable([]*models.Link{link})
}

// LinkSearch command.
func LinkSearch(group *cli.Group, command *cli.Command, arguments []string) int {
	search := newMeiliSearch(command)

	parsed, err := command.Parse(arguments)
	if err != nil {
		return command.PrintErrorWithHelp(group, err)
	}

	fmt.Println("=> Search " + cli.Green(parsed[0]) + " in " + cli.Green("MeiliSearch"))

	result, err := search.Index().Search(parsed[0], &meilisearch.SearchRequest{
		AttributesToRetrieve: []string{
			"id",
			"url",
			"title",
		},
	})
	if err != nil {
		return command.PrintError(err)
	}

	links := []*models.Link{}

	err = result.Hits.DecodeInto(&links)
	if err != nil {
		return command.PrintError(err)
	}

	return printLinkTable(links)
}

// LinkUpdate command.
func LinkUpdate(group *cli.Group, command *cli.Command, arguments []string) int {
	search := newMeiliSearch(command)
	linkCrawler := newCrawler(command)
	tags := command.FlagSet().String("tags", "", "Link Tags")

	parsed, err := command.Parse(arguments)
	if err != nil {
		return command.PrintErrorWithHelp(group, err)
	}

	fmt.Println("=> Find document " + cli.Green(parsed[0]) + " in " + cli.Green("MeiliSearch"))

	index := search.Index()

	var original models.Link

	err = index.GetDocument(parsed[0], nil, &original)
	if err != nil {
		return command.PrintError(err)
	}

	meta, url, err := linkCrawler.Crawl(original.URL)
	if err != nil {
		return command.PrintError(err)
	}

	link := (&models.Link{
		ID:  original.ID,
		URL: url,
	}).ParseTags(*tags).Fill(meta)

	fmt.Println("=> Send data to " + cli.Green("MeiliSearch"))

	_, err = index.UpdateDocuments([]*models.Link{
		link,
	}, nil)
	if err != nil {
		return command.PrintError(err)
	}

	return printLinkTable([]*models.Link{link})
}

// LinkDelete command.
func LinkDelete(group *cli.Group, command *cli.Command, arguments []string) int {
	search := newMeiliSearch(command)

	parsed, err := command.Parse(arguments)
	if err != nil {
		return command.PrintErrorWithHelp(group, err)
	}

	fmt.Println("=> Delete document " + cli.Green(parsed[0]) + " from " + cli.Green("MeiliSearch"))

	_, err = search.Index().DeleteDocument(parsed[0], nil)
	if err != nil {
		return command.PrintError(err)
	}

	return cli.Success
}

// LinkDeleteAll command.
func LinkDeleteAll(group *cli.Group, command *cli.Command, arguments []string) int {
	search := newMeiliSearch(command)

	_, err := command.Parse(arguments)
	if err != nil {
		return command.PrintErrorWithHelp(group, err)
	}

	fmt.Println("=> Delete all documents from " + cli.Green("MeiliSearch"))

	_, err = search.Index().DeleteAllDocuments(nil)
	if err != nil {
		return command.PrintError(err)
	}

	return cli.Success
}
