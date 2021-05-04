package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/google/uuid"
	"github.com/meilisearch/meilisearch-go"
	"github.com/petaki/pylon/internal/meta"
	"github.com/petaki/pylon/internal/models"
	"github.com/petaki/support-go/cli"
)

// LinkAdd command.
func LinkAdd(group *cli.Group, command *cli.Command, arguments []string) int {
	meiliSearchHost, meiliSearchAPIKey, meiliSearchIndex := createMeiliSearchFlags(command)
	headlessShellHost, tags := createLinkFlags(command)

	parsed, err := command.Parse(arguments)
	if err != nil {
		return command.PrintHelp(group)
	}

	if !isValidURL(parsed[0]) {
		return printError(fmt.Errorf("invalid URL: %s", parsed[0]))
	}

	fmt.Println("=> Get " + cli.Green("WebSocket Debugger URL"))

	webSocketDebuggerURL, err := getWebSocketDebuggerURL(*headlessShellHost)
	if err != nil {
		return printError(err)
	}

	fmt.Println("=> Get " + cli.Green("URL content"))

	body, url, err := getURLContent(webSocketDebuggerURL, parsed[0])
	if err != nil {
		return printError(err)
	}

	fmt.Println("=> Parse " + cli.Green("HTML source"))

	data, err := meta.Parse(body)
	if err != nil {
		return printError(err)
	}

	link := (&models.Link{
		ID:  uuid.NewString(),
		URL: url,
	}).ParseTags(*tags).Fill(data)

	fmt.Println("=> Send data to " + cli.Green("MeiliSearch"))

	meiliSearchClient := meilisearch.NewClient(meilisearch.Config{
		Host:   *meiliSearchHost,
		APIKey: *meiliSearchAPIKey,
	})

	_, err = meiliSearchClient.Documents(*meiliSearchIndex).AddOrUpdate([]*models.Link{
		link,
	})
	if err != nil {
		return printError(err)
	}

	return printLinkTable([]*models.Link{link})
}

// LinkSearch command.
func LinkSearch(group *cli.Group, command *cli.Command, arguments []string) int {
	meiliSearchHost, meiliSearchAPIKey, meiliSearchIndex := createMeiliSearchFlags(command)

	parsed, err := command.Parse(arguments)
	if err != nil {
		return command.PrintHelp(group)
	}

	fmt.Println("=> Search " + cli.Green(parsed[0]) + " in " + cli.Green("MeiliSearch"))

	meiliSearchClient := meilisearch.NewClient(meilisearch.Config{
		Host:   *meiliSearchHost,
		APIKey: *meiliSearchAPIKey,
	})

	result, err := meiliSearchClient.Search(*meiliSearchIndex).Search(meilisearch.SearchRequest{
		Query: parsed[0],
		AttributesToRetrieve: []string{
			"id",
			"url",
			"title",
		},
	})
	if err != nil {
		return printError(err)
	}

	links := make([]*models.Link, len(result.Hits))

	for i, hit := range result.Hits {
		links[i] = &models.Link{
			ID:    hit.(map[string]interface{})["id"].(string),
			URL:   hit.(map[string]interface{})["url"].(string),
			Title: hit.(map[string]interface{})["title"].(string),
		}
	}

	return printLinkTable(links)
}

// LinkUpdate command.
func LinkUpdate(group *cli.Group, command *cli.Command, arguments []string) int {
	meiliSearchHost, meiliSearchAPIKey, meiliSearchIndex := createMeiliSearchFlags(command)
	headlessShellHost, tags := createLinkFlags(command)

	parsed, err := command.Parse(arguments)
	if err != nil {
		return command.PrintHelp(group)
	}

	fmt.Println("=> Find document " + cli.Green(parsed[0]) + " in " + cli.Green("MeiliSearch"))

	meiliSearchClient := meilisearch.NewClient(meilisearch.Config{
		Host:   *meiliSearchHost,
		APIKey: *meiliSearchAPIKey,
	})

	var original models.Link

	err = meiliSearchClient.Documents(*meiliSearchIndex).Get(parsed[0], &original)
	if err != nil {
		return printError(err)
	}

	fmt.Println("=> Get " + cli.Green("WebSocket Debugger URL"))

	webSocketDebuggerURL, err := getWebSocketDebuggerURL(*headlessShellHost)
	if err != nil {
		return printError(err)
	}

	fmt.Println("=> Get " + cli.Green("URL content"))

	body, url, err := getURLContent(webSocketDebuggerURL, original.URL)
	if err != nil {
		return printError(err)
	}

	fmt.Println("=> Parse " + cli.Green("HTML source"))

	data, err := meta.Parse(body)
	if err != nil {
		return printError(err)
	}

	link := (&models.Link{
		ID:  original.ID,
		URL: url,
	}).ParseTags(*tags).Fill(data)

	fmt.Println("=> Send data to " + cli.Green("MeiliSearch"))

	_, err = meiliSearchClient.Documents(*meiliSearchIndex).AddOrUpdate([]*models.Link{
		link,
	})
	if err != nil {
		return printError(err)
	}

	return printLinkTable([]*models.Link{link})
}

// LinkDelete command.
func LinkDelete(group *cli.Group, command *cli.Command, arguments []string) int {
	meiliSearchHost, meiliSearchAPIKey, meiliSearchIndex := createMeiliSearchFlags(command)

	parsed, err := command.Parse(arguments)
	if err != nil {
		return command.PrintHelp(group)
	}

	fmt.Println("=> Delete document " + cli.Green(parsed[0]) + " from " + cli.Green("MeiliSearch"))

	meiliSearchClient := meilisearch.NewClient(meilisearch.Config{
		Host:   *meiliSearchHost,
		APIKey: *meiliSearchAPIKey,
	})

	_, err = meiliSearchClient.Documents(*meiliSearchIndex).Delete(parsed[0])
	if err != nil {
		return printError(err)
	}

	return 0
}

// LinkDeleteAll command.
func LinkDeleteAll(group *cli.Group, command *cli.Command, arguments []string) int {
	meiliSearchHost, meiliSearchAPIKey, meiliSearchIndex := createMeiliSearchFlags(command)

	_, err := command.Parse(arguments)
	if err != nil {
		return command.PrintHelp(group)
	}

	fmt.Println("=> Delete all documents from " + cli.Green("MeiliSearch"))

	meiliSearchClient := meilisearch.NewClient(meilisearch.Config{
		Host:   *meiliSearchHost,
		APIKey: *meiliSearchAPIKey,
	})

	_, err = meiliSearchClient.Documents(*meiliSearchIndex).DeleteAllDocuments()
	if err != nil {
		return printError(err)
	}

	return 0
}

func getWebSocketDebuggerURL(headlessShellHost string) (string, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("GET", headlessShellHost+"/json/version", nil)
	if err != nil {
		return "", err
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("headless shell status code: %d", resp.StatusCode)
	}

	var data map[string]interface{}

	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return "", err
	}

	return data["webSocketDebuggerUrl"].(string), nil
}

func getURLContent(webSocketDebuggerURL, rawURL string) (io.Reader, string, error) {
	allocatorContext, cancel := chromedp.NewRemoteAllocator(context.Background(), webSocketDebuggerURL)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocatorContext)
	defer cancel()

	var body, url string

	err := chromedp.Run(ctx,
		chromedp.Navigate(rawURL),
		chromedp.Sleep(1*time.Second),
		chromedp.OuterHTML("html", &body),
		chromedp.Location(&url),
	)
	if err != nil {
		return nil, "", err
	}

	return strings.NewReader(body), url, nil
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
