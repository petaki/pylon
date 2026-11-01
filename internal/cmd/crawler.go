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
	"github.com/chromedp/chromedp/remote"
	"github.com/petaki/pylon/internal/models"
	"github.com/petaki/support-go/cli"
)

const (
	headerUserAgent = "User-Agent"
	requestTimeout  = 10 * time.Second
	crawlTimeout    = 30 * time.Second
	maxBodySize     = 10 << 20
)

type headlessShellVersion struct {
	UserAgent            string `json:"User-Agent"`
	WebSocketDebuggerURL string `json:"webSocketDebuggerUrl"`
}

type crawler struct {
	client            *http.Client
	headlessShellHost *string
}

func newCrawler(command *cli.Command) *crawler {
	return &crawler{
		client:            &http.Client{Timeout: requestTimeout},
		headlessShellHost: command.String("headless-shell-host", "HEADLESS_SHELL_HOST", "http://127.0.0.1:9222", "Headless Shell Host"),
	}
}

// Crawl function.
func (c *crawler) Crawl(rawURL string) (*models.Meta, string, error) {
	fmt.Println("=> Get " + cli.Green("User-Agent") + " and " + cli.Green("WebSocket Debugger URL"))

	version, err := c.version()
	if err != nil {
		return nil, "", err
	}

	fmt.Println("=> Get " + cli.Green("URL content"))

	body, url, err := c.localContent(version.UserAgent, rawURL)
	if err != nil {
		return nil, "", err
	}

	fmt.Println("=> Parse " + cli.Green("HTML source"))

	meta, err := models.ParseMeta(body)
	if err != nil {
		return nil, "", err
	}

	if len(meta.OgImages) > 0 {
		return meta, url, nil
	}

	fmt.Println("=> Retry crawling with " + cli.Green("JavaScript"))

	body, url, err = c.remoteContent(version.WebSocketDebuggerURL, rawURL)
	if err != nil {
		return nil, "", err
	}

	meta, err = models.ParseMeta(body)
	if err != nil {
		return nil, "", err
	}

	return meta, url, nil
}

func (c *crawler) version() (*headlessShellVersion, error) {
	req, err := http.NewRequest(http.MethodGet, *c.headlessShellHost+"/json/version", nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: %s, %d", ErrInvalidStatusCode, *c.headlessShellHost, resp.StatusCode)
	}

	var version headlessShellVersion

	err = json.NewDecoder(resp.Body).Decode(&version)
	if err != nil {
		return nil, err
	}

	if version.UserAgent == "" || version.WebSocketDebuggerURL == "" {
		return nil, fmt.Errorf("%w: %s", ErrInvalidHeadlessShell, *c.headlessShellHost)
	}

	return &version, nil
}

func (c *crawler) localContent(userAgent, rawURL string) (io.Reader, string, error) {
	req, err := http.NewRequest(http.MethodGet, rawURL, nil)
	if err != nil {
		return nil, "", err
	}

	req.Header.Set(headerUserAgent, userAgent)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, "", err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("%w: %s, %d", ErrInvalidStatusCode, rawURL, resp.StatusCode)
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, maxBodySize))
	if err != nil {
		return nil, "", err
	}

	return bytes.NewReader(body), resp.Request.URL.String(), nil
}

func (c *crawler) remoteContent(webSocketDebuggerURL, rawURL string) (io.Reader, string, error) {
	timeoutContext, cancel := context.WithTimeout(context.Background(), crawlTimeout)
	defer cancel()

	allocatorContext, cancel := remote.NewAllocator(timeoutContext, webSocketDebuggerURL)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocatorContext)
	defer cancel()

	err := chromedp.Do(ctx,
		chromedp.Navigate(rawURL),
		chromedp.Sleep(1*time.Second),
	)
	if err != nil {
		return nil, "", err
	}

	body, err := chromedp.Run(ctx, chromedp.OuterHTML("html"))
	if err != nil {
		return nil, "", err
	}

	url, err := chromedp.Run(ctx, chromedp.Location())
	if err != nil {
		return nil, "", err
	}

	return strings.NewReader(body), url, nil
}
