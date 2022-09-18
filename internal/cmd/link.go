package cmd

import (
	"bytes"
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
		return command.PrintError(fmt.Errorf("invalid URL: %s", parsed[0]))
	}

	fmt.Println("=> Get " + cli.Green("User-Agent") + " and " + cli.Green("WebSocket Debugger URL"))

	userAgent, webSocketDebuggerURL, err := getUserAgentAndWebSocketDebuggerURL(*headlessShellHost)
	if err != nil {
		return command.PrintError(err)
	}

	fmt.Println("=> Get " + cli.Green("URL content"))

	body, url, err := getLocalURLContent(userAgent, parsed[0])
	if err != nil {
		return command.PrintError(err)
	}

	fmt.Println("=> Parse " + cli.Green("HTML source"))

	meta, err := models.ParseMeta(body)
	if err != nil {
		return command.PrintError(err)
	}

	if len(meta.OgImages) == 0 {
		fmt.Println("=> Retry crawling with " + cli.Green("JavaScript"))

		body, url, err = getRemoteURLContent(webSocketDebuggerURL, parsed[0])
		if err != nil {
			return command.PrintError(err)
		}

		meta, err = models.ParseMeta(body)
		if err != nil {
			return command.PrintError(err)
		}
	}

	link := (&models.Link{
		ID:  uuid.NewString(),
		URL: url,
	}).ParseTags(*tags).Fill(meta)

	fmt.Println("=> Send data to " + cli.Green("MeiliSearch"))

	meiliSearchClient := meilisearch.NewClient(meilisearch.ClientConfig{
		Host:   *meiliSearchHost,
		APIKey: *meiliSearchAPIKey,
	})

	_, err = meiliSearchClient.Index(*meiliSearchIndex).AddDocuments([]*models.Link{
		link,
	})
	if err != nil {
		return command.PrintError(err)
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

	meiliSearchClient := meilisearch.NewClient(meilisearch.ClientConfig{
		Host:   *meiliSearchHost,
		APIKey: *meiliSearchAPIKey,
	})

	result, err := meiliSearchClient.Index(*meiliSearchIndex).Search(parsed[0], &meilisearch.SearchRequest{
		AttributesToRetrieve: []string{
			"id",
			"url",
			"title",
		},
	})
	if err != nil {
		return command.PrintError(err)
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

	meiliSearchClient := meilisearch.NewClient(meilisearch.ClientConfig{
		Host:   *meiliSearchHost,
		APIKey: *meiliSearchAPIKey,
	})

	var original models.Link

	err = meiliSearchClient.Index(*meiliSearchIndex).GetDocument(parsed[0], nil, &original)
	if err != nil {
		return command.PrintError(err)
	}

	fmt.Println("=> Get " + cli.Green("User-Agent") + " and " + cli.Green("WebSocket Debugger URL"))

	userAgent, webSocketDebuggerURL, err := getUserAgentAndWebSocketDebuggerURL(*headlessShellHost)
	if err != nil {
		return command.PrintError(err)
	}

	fmt.Println("=> Get " + cli.Green("URL content"))

	body, url, err := getLocalURLContent(userAgent, original.URL)
	if err != nil {
		return command.PrintError(err)
	}

	fmt.Println("=> Parse " + cli.Green("HTML source"))

	meta, err := models.ParseMeta(body)
	if err != nil {
		return command.PrintError(err)
	}

	if len(meta.OgImages) == 0 {
		fmt.Println("=> Retry crawling with " + cli.Green("JavaScript"))

		body, url, err = getRemoteURLContent(webSocketDebuggerURL, original.URL)
		if err != nil {
			return command.PrintError(err)
		}

		meta, err = models.ParseMeta(body)
		if err != nil {
			return command.PrintError(err)
		}
	}

	link := (&models.Link{
		ID:  original.ID,
		URL: url,
	}).ParseTags(*tags).Fill(meta)

	fmt.Println("=> Send data to " + cli.Green("MeiliSearch"))

	_, err = meiliSearchClient.Index(*meiliSearchIndex).UpdateDocuments([]*models.Link{
		link,
	})
	if err != nil {
		return command.PrintError(err)
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

	meiliSearchClient := meilisearch.NewClient(meilisearch.ClientConfig{
		Host:   *meiliSearchHost,
		APIKey: *meiliSearchAPIKey,
	})

	_, err = meiliSearchClient.Index(*meiliSearchIndex).DeleteDocument(parsed[0])
	if err != nil {
		return command.PrintError(err)
	}

	return cli.Success
}

// LinkDeleteAll command.
func LinkDeleteAll(group *cli.Group, command *cli.Command, arguments []string) int {
	meiliSearchHost, meiliSearchAPIKey, meiliSearchIndex := createMeiliSearchFlags(command)

	_, err := command.Parse(arguments)
	if err != nil {
		return command.PrintHelp(group)
	}

	fmt.Println("=> Delete all documents from " + cli.Green("MeiliSearch"))

	meiliSearchClient := meilisearch.NewClient(meilisearch.ClientConfig{
		Host:   *meiliSearchHost,
		APIKey: *meiliSearchAPIKey,
	})

	_, err = meiliSearchClient.Index(*meiliSearchIndex).DeleteAllDocuments()
	if err != nil {
		return command.PrintError(err)
	}

	return cli.Success
}

func getUserAgentAndWebSocketDebuggerURL(headlessShellHost string) (string, string, error) {
	req, err := http.NewRequest("GET", headlessShellHost+"/json/version", nil)
	if err != nil {
		return "", "", err
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", "", err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("headless shell status code: %d", resp.StatusCode)
	}

	var data map[string]interface{}

	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return "", "", err
	}

	return data["User-Agent"].(string), data["webSocketDebuggerUrl"].(string), nil
}

func getLocalURLContent(userAgent, rawURL string) (io.Reader, string, error) {
	req, err := http.NewRequest("GET", rawURL, nil)
	if err != nil {
		return nil, "", err
	}

	req.Header.Set("User-Agent", userAgent)

	resp, err := client.Do(req)
	if err != nil {
		return nil, "", err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}

	return bytes.NewReader(body), resp.Request.URL.String(), nil
}

func getRemoteURLContent(webSocketDebuggerURL, rawURL string) (io.Reader, string, error) {
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
