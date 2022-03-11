package intento

type TextFormat string

const (
	FormatHTML TextFormat = "html"
)

// TranslationWithSourceTextFormat specifies the format of the source text.
func TranslationWithSourceTextFormat(format TextFormat) TranslationOption {
	return newFuncTranslationOption(func(o *translationOptions) {
		o.Context.Format = format
	})
}

// TranslationWithTrace sets trace flag to provide more information to support team.
//
// By default, Intento API works in the “no trace” mode, with no payload stored on
// the Intento side. Even temporary translation results are not accessible by
// Intento employees.
//
// If you have a reproducible error when using the Intento API and require
// technical support, we recommend you enable the payload logging to provide
// more information to our support team.
func TranslationWithTrace() TranslationOption {
	return newFuncTranslationOption(func(o *translationOptions) {
		o.Service.Trace = true
	})
}

// TranslationWithProvider sets a translation provider.
func TranslationWithProvider(providerID string) TranslationOption {
	return newFuncTranslationOption(func(o *translationOptions) {
		o.Service.Provider = providerID
	})
}

// TranslationWithSmartRouting sets a smart routing.
//
// The publicly available smart routing feature routes user's translation
// requests to the best MT provider. Intento selects the MT model for the text
// and language pair based on the following information:
//
// - apriori benchmark on the standard dataset
// - provider usage statistics
//
// Besides publicly available schemes, a user can have custom routing schemes.
func TranslationWithSmartRouting(routing string) TranslationOption {
	return newFuncTranslationOption(func(o *translationOptions) {
		o.Service.Routing = routing
	})
}

// TranslationWithCache sets the rules for using the TM-cache.
//
// The Intento MT Cache feature enables you to cache previously translated
// strings on the Intento side and reuse them in the subsequent requests. In
// order to enable MT Cache, you need to alter the translation requests with the
// following attributes:
// apply = true - pulls the translation from the cache if it’s previously cached.
// update = true - updates the translation in the cache if it’s previously cached,
// or puts the translation to the cache if it isn’t.
//
// You can specify both attributes so the request will use cached translation
// before the translation is sent to MT Provider, and put the translation to the
// cache if previously not cached after the translation with MT Provider.
//
func TranslationWithCache(apply, update bool) TranslationOption {
	return newFuncTranslationOption(func(o *translationOptions) {
		o.Service.Cache.Apply = apply
		o.Service.Cache.Update = update
	})
}

// TranslationWithNoTranslateProtection is NOTRANSLATE protection.
//
// Intento API allows you to screen a word or parts of sentences to protect
// against translation. In order to do that one needs to modify the source text so
// the desired word or words are surrounded with NOTRANSLATE suffix and
// prefix, and used suffix and prefix are mentioned in the service section of the
// translation request.
//
// Parameter removeMarkup defines whether the prefix and suffix tags should
// remain in the translation result.
//
// Please note that NOTRANSLATE protection only works if html format is specified.
//
// Example:
// ```
// result, err := client.Translate(
//     ctx,
//     "en",
//     "es",
//     []string{"Hello <span class=\"notranslate\">Old Friend</span>"},
//     TranslationWithNoTranslateProtection(
//         "<span class=\"notranslate\">",
//         "</span>",
//         true,
//     ),
// )
// ```
// Output: "Hola Old Friend"
//
func TranslationWithNoTranslateProtection(prefix, suffix string, removeMarkup bool) TranslationOption {
	return newFuncTranslationOption(func(o *translationOptions) {
		o.Service.NoTranslate.Prefix = prefix
		o.Service.NoTranslate.Suffix = suffix
		o.Service.NoTranslate.RemoveMarkup = removeMarkup
	})
}

// TranslationWithProfanityDetection runs a profanity check of the translated content.
//
// The content parameter is a list of unwanted content types that are detected.
// For now, only profanity is implemented, planned to be extended to violence, racism, etc.
//
// Please note that profanity detection does not modify the translation result.
//
func TranslationWithProfanityDetection(content []string) TranslationOption {
	return newFuncTranslationOption(func(o *translationOptions) {
		o.Service.Moderation.Used = true
		o.Service.Moderation.Content = content
		o.Service.Moderation.Action = "inform"
	})
}

// TranslationOption configures how we set up the connection.
type TranslationOption interface {
	apply(*translationOptions)
}

// translationOptions configure a translation process.
type translationOptions struct {
	Context struct {
		From   string     `json:"from,omitempty"`
		To     string     `json:"to,omitempty"`
		Text   []string   `json:"text,omitempty"`
		Format TextFormat `json:"format,omitempty"`
	} `json:"context"`
	Service struct {
		Async    bool   `json:"async,omitempty"`
		Trace    bool   `json:"trace,omitempty"`
		Provider string `json:"provider,omitempty"`
		Routing  string `json:"routing,omitempty"`
		Cache    struct {
			Apply  bool `json:"apply,omitempty"`
			Update bool `json:"update,omitempty"`
		} `json:"cache,omitempty"`
		NoTranslate struct {
			Prefix       string `json:"prefix,omitempty"`
			Suffix       string `json:"suffix,omitempty"`
			RemoveMarkup bool   `json:"remove_markup,omitempty"`
		} `json:"notranslate,omitempty"`
		Moderation struct {
			Action  string   `json:"action,omitempty"`
			Used    bool     `json:"used,omitempty"`
			Content []string `json:"content,omitempty"`
		} `json:"moderation,omitempty"`
	} `json:"service"`
}

// funcTranslationOption wraps a function that modifies clientOptions into an implementation of the ClientOption interface.
type funcTranslationOption struct {
	fn func(*translationOptions)
}

func (fco *funcTranslationOption) apply(do *translationOptions) {
	fco.fn(do)
}

func newFuncTranslationOption(fn func(*translationOptions)) *funcTranslationOption {
	return &funcTranslationOption{
		fn: fn,
	}
}
