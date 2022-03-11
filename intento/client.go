package intento

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

//go:generate moq -pkg intento_test -out mocks_test.go . HttpClient

// HttpClient is an interface for sending an HTTP request.
type HttpClient interface {
	Do(req *http.Request) (resp *http.Response, err error)
}

// Logger is a function for writing to the log.
type Logger func(ctx context.Context, format string, args ...interface{})

// Client is the client for interacting with the Intento API.
type Client struct {
	clientOptions
	apiKey string
}

// New creates an instance of Client.
func New(
	apiKey string,
	options ...ClientOption,
) *Client {
	client := &Client{
		clientOptions: defaultClientOptions(),
		apiKey:        apiKey,
	}

	for _, opt := range options {
		opt.apply(&client.clientOptions)
	}

	return client
}

// Provider describes a translation service provider.
type Provider struct {
	Production           bool     `json:"production"`
	Integrated           bool     `json:"integrated"`
	Billable             bool     `json:"billable"`
	OwnAuth              bool     `json:"own_auth"`
	StockModel           bool     `json:"stock_model"`
	CustomModel          bool     `json:"custom_model"`
	DelegatedCredentials bool     `json:"delegated_credentials"`
	AsyncOnly            bool     `json:"async_only"`
	ID                   string   `json:"id"`
	Name                 string   `json:"name"`
	Vendor               string   `json:"vendor"`
	Score                int      `json:"score"`
	Price                int      `json:"price"`
	ApiID                string   `json:"api_id"`
	Picture              string   `json:"picture"`
	Type                 string   `json:"type"`
	Description          string   `json:"description"`
	Tone                 []string `json:"tone"`
	Symmetric            []string `json:"symmetric"`
	Pairs                []struct {
		From string `json:"from"`
		To   string `json:"to"`
	} `json:"pairs"`
}

// AvailableProviders returns the list of available providers.
func (c *Client) AvailableProviders(ctx context.Context) ([]Provider, error) {
	var providers []Provider

	err := c.apiGetRequest(ctx, "https://syncwrapper.inten.to/ai/text/translate", &providers)
	if err != nil {
		return nil, err
	}

	return providers, nil
}

const AutoDetectSourceLanguage = ""

// TranslationResult describes a result of translation.
type TranslationResult struct {
	ID      string   `json:"id"`
	Results []string `json:"results"`
	Meta    struct {
		DetectedSourceLanguage []string `json:"detected_source_language"`
	} `json:"meta"`
	Service struct {
		Provider struct {
			ID          string `json:"id"`
			Name        string `json:"name"`
			Vendor      string `json:"vendor"`
			Description string `json:"description"`
			Logo        string `json:"logo"`
		} `json:"provider"`
	} `json:"service"`
}

// Translate text with given settings.
func (c *Client) Translate(
	ctx context.Context,
	text []string,
	from string,
	to string,
	options ...TranslationOption,
) (TranslationResult, error) {
	params := translationOptions{}
	params.Context.Text = text
	params.Context.From = from
	params.Context.To = to

	for _, opt := range options {
		opt.apply(&params)
	}

	var result TranslationResult

	err := c.apiPostRequest(ctx, "https://syncwrapper.inten.to/ai/text/translate", &params, &result)
	if err != nil {
		return TranslationResult{}, err
	}

	return result, nil
}

func (c *Client) apiGetRequest(ctx context.Context, url string, result interface{}) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("create http request: %w", err)
	}

	return c.apiRequest(ctx, req, result)
}

func (c *Client) apiPostRequest(ctx context.Context, url string, params interface{}, result interface{}) error {
	requestBody, err := json.Marshal(params)
	if err != nil {
		return fmt.Errorf("marshal json: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(requestBody))
	if err != nil {
		return fmt.Errorf("create http request: %w", err)
	}

	return c.apiRequest(ctx, req, result)
}

func (c *Client) apiRequest(ctx context.Context, req *http.Request, result interface{}) error {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("apikey", c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("http do request: %w", err)
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			c.logger(ctx, "close request body: %v", err)
		}
	}()

	err = httpStatusCodeToError(resp.StatusCode)
	if err != nil {
		return fmt.Errorf("check http status code: %w", err)
	}

	err = json.NewDecoder(resp.Body).Decode(result)
	if err != nil {
		return fmt.Errorf("unmarshal response json: %w", err)
	}

	return nil
}
