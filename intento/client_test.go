package intento_test

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"intento-golang/intento"
)

var apiKey = readApiKey()

var text = []string{
	"Hello World!",
}

func ExampleClient_Translate_noOptions() {
	ctx := context.Background()

	client := intento.New(apiKey)

	result, err := client.Translate(ctx, text, "en", "es")
	if err != nil {
		log.Fatalf("translate: %v", err)
	}

	fmt.Println(result.Results[0])

	// Output:
	// Hola, mundo.
}

func ExampleClient_Translate_allOptions() {
	ctx := context.Background()

	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}

	logger := func(ctx context.Context, format string, args ...interface{}) {
		log.Printf(format, args...)
	}

	client := intento.New(
		apiKey,
		intento.ClientWithHttpClient(httpClient),
		intento.ClientWithLogger(logger),
	)

	providerID := "ai.text.translate.tencent.machine_translation_api"

	result, err := client.Translate(
		ctx,
		text,
		intento.AutoDetectSourceLanguage,
		"es",
		intento.TranslationWithSourceTextFormat(intento.FormatHTML),
		intento.TranslationWithTrace(),
		intento.TranslationWithProvider(providerID),
		intento.TranslationWithSmartRouting("best"),
		intento.TranslationWithCache(true, true),
		intento.TranslationWithNoTranslateProtection(
			`<span class="notranslate">`,
			`</span>`,
			true,
		),
		intento.TranslationWithProfanityDetection([]string{
			"profanity",
		}),
	)
	if err != nil {
		log.Fatalf("translate: %v", err)
	}

	fmt.Println(result.Results[0])

	// Output:
	// ¡Hola mundo!
}

func ExampleClient_Translate_errorsHandling() {
	ctx := context.Background()

	client := intento.New("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")

	_, err := client.Translate(ctx, text, "en", "es")
	if err != nil {
		var authKeyIsMissingError *intento.AuthKeyIsMissingError

		if errors.As(err, &authKeyIsMissingError) {
			fmt.Println("Please update your authorization key!")
			return
		}

		log.Fatalf("translate: %v", err)
	}

	// Output:
	// Please update your authorization key!
}

func ExampleClient_AvailableProviders() {
	ctx := context.Background()

	client := intento.New(apiKey)

	providers, err := client.AvailableProviders(ctx)
	if err != nil {
		log.Fatalf("get available providers: %v", err)
	}

	var providerIDs []string

	for _, provider := range providers {
		providerIDs = append(providerIDs, provider.ID)
	}

	sort.Strings(providerIDs)

	for _, providerID := range providerIDs {
		log.Println(providerID)
	}

	// Output:
}

func ExampleClient_AvailableLanguages() {
	ctx := context.Background()

	client := intento.New(apiKey)

	languages, err := client.AvailableLanguages(ctx)
	if err != nil {
		log.Fatalf("get available languages: %v", err)
	}

	var languageCodes []string

	for _, language := range languages {
		languageCodes = append(languageCodes, fmt.Sprintf("%-7s - %s", language.IntentoCode, language.ISOName))
	}

	sort.Strings(languageCodes)

	for _, languageCode := range languageCodes {
		log.Println(languageCode)
	}

	// Output:
}

func ExampleClient_SmartRoutingList() {
	ctx := context.Background()

	client := intento.New(apiKey)

	smartRoutingList, err := client.SmartRoutingList(ctx)
	if err != nil {
		log.Fatalf("get smart routing: %v", err)
	}

	var list []string

	for _, smartRouting := range smartRoutingList {
		list = append(list, fmt.Sprintf("%s - %s", smartRouting.Name, smartRouting.Description))
	}

	sort.Strings(list)

	for _, item := range list {
		log.Println(item)
	}

	// Output:
}

func TestClient_AvailableProviders(t *testing.T) {
	ctx := context.Background()

	mockHttpClient := &HttpClientMock{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(strings.NewReader("[{}]")),
			}, nil
		},
	}

	mockLogger := func(ctx context.Context, format string, args ...interface{}) {
		t.Logf(format, args)
	}

	const (
		mockApiKey = "api_key_1"
	)

	client := intento.New(
		mockApiKey,
		intento.ClientWithHttpClient(mockHttpClient),
		intento.ClientWithLogger(mockLogger),
	)

	providers, err := client.AvailableProviders(ctx)

	assert.NoError(t, err)
	assert.NotEmpty(t, providers)
}

func readApiKey() string {
	data, err := ioutil.ReadFile("../api_key.txt")
	if err != nil {
		log.Fatalf("read file: %v", err)
	}

	return strings.TrimSpace(string(data))
}
