package intento

import (
	"context"
	"log"
	"net/http"
)

// ClientWithHttpClient sets HttpClient.
func ClientWithHttpClient(httpClient HttpClient) ClientOption {
	return newFuncClientOption(func(o *clientOptions) {
		o.httpClient = httpClient
	})
}

// ClientWithLogger sets Logger.
func ClientWithLogger(logger Logger) ClientOption {
	return newFuncClientOption(func(o *clientOptions) {
		o.logger = logger
	})
}

// ClientWithServerURL sets server url.
func ClientWithServerURL(serverURL string) ClientOption {
	return newFuncClientOption(func(o *clientOptions) {
		o.apiServerURL = serverURL
	})
}

// ClientOption configures how we set up the connection.
type ClientOption interface {
	apply(*clientOptions)
}

// clientOptions configure a Client.
type clientOptions struct {
	httpClient   HttpClient
	logger       Logger
	apiServerURL string
}

func defaultClientOptions() clientOptions {
	return clientOptions{
		httpClient:   http.DefaultClient,
		logger:       func(ctx context.Context, format string, args ...interface{}) { log.Printf(format, args...) },
		apiServerURL: "https://syncwrapper.inten.to", // "https://api.inten.to"
	}
}

// funcClientOption wraps a function that modifies clientOptions into an implementation of the ClientOption interface.
type funcClientOption struct {
	fn func(*clientOptions)
}

func (fco *funcClientOption) apply(do *clientOptions) {
	fco.fn(do)
}

func newFuncClientOption(fn func(*clientOptions)) *funcClientOption {
	return &funcClientOption{
		fn: fn,
	}
}
