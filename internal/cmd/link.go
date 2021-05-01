package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/chromedp/chromedp"
	"github.com/google/uuid"
	"github.com/meilisearch/meilisearch-go"
	"github.com/petaki/pylon/internal/meta"
	"github.com/petaki/pylon/internal/models"
	"github.com/petaki/support-go/cli"
)

// LinkAdd command.
func LinkAdd(group *cli.Group, command *cli.Command, arguments []string) int {
	meiliSearchHost, meiliSearchAPIKey := createMeiliSearchFlags(command)
	headlessShellHost := command.FlagSet().String("headless-shell-host", os.Getenv("HEADLESS_SHELL_HOST"), "Headless Shell Host")
	tags := command.FlagSet().String("tags", "", "Link Tags")

	parsed, err := command.Parse(arguments)
	if err != nil {
		return command.PrintHelp(group)
	}

	if !isURL(parsed[0]) {
		return printError(fmt.Errorf("invalid url: %s", parsed[0]))
	}

	fmt.Println("=> Get " + cli.Green("WebSocket Debugger URL"))

	webSocketDebuggerURL, err := getWebSocketDebuggerURL(*headlessShellHost)
	if err != nil {
		return printError(err)
	}

	fmt.Println("=> Get " + cli.Green("URL content"))

	allocatorContext, cancel := chromedp.NewRemoteAllocator(context.Background(), webSocketDebuggerURL)
	defer cancel()

	ctxt, cancel := chromedp.NewContext(allocatorContext)
	defer cancel()

	var body, url string

	err = chromedp.Run(ctxt,
		chromedp.Navigate(parsed[0]),
		chromedp.OuterHTML("html", &body),
		chromedp.Location(&url),
	)
	if err != nil {
		return printError(err)
	}

	fmt.Println("=> Parse " + cli.Green("HTML source"))

	data, err := meta.Parse(strings.NewReader(body))
	if err != nil {
		return printError(err)
	}

	link := (&models.Link{
		ID:   uuid.NewString(),
		URL:  url,
		Tags: strings.Split(*tags, ","),
	}).Fill(data)

	fmt.Println("=> Send data to " + cli.Green("MeiliSearch"))

	meiliSearchClient := meilisearch.NewClient(meilisearch.Config{
		Host:   *meiliSearchHost,
		APIKey: *meiliSearchAPIKey,
	})

	_, err = meiliSearchClient.Documents("pylon").AddOrUpdate([]*models.Link{
		link,
	})
	if err != nil {
		return printError(err)
	}

	return 0
}

// LinkDelete command.
func LinkDelete(group *cli.Group, command *cli.Command, arguments []string) int {
	fmt.Println(arguments)

	return 0
}

func getWebSocketDebuggerURL(headlessShellHost string) (string, error) {
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
