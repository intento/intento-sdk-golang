package intento

import "errors"

var (
	ProviderRelatedError      = errors.New("provider-related error")
	AuthKeyIsMissingError     = errors.New("intento: auth key is missing")
	AuthKeyIsInvalidError     = errors.New("intento: auth key is invalid")
	NotFoundError             = errors.New("intento: intent/provider not found")
	CapabilitiesMismatchError = errors.New("intento: capabilities mismatch for the chosen provider")
	APIRateLimitError         = errors.New("intento: API rate limit exceeded")
	InternalError             = errors.New("intento: internal error")
	NotImplemented            = errors.New("intento: not implemented")
	GatewayTimeoutError       = errors.New("intento: gateway timeout errors")
)

func httpStatusCodeToError(statusCode int) error {
	if statusCode >= 200 && statusCode <= 299 {
		return nil
	}

	switch statusCode {
	case 400:
		return ProviderRelatedError
	case 401:
		return AuthKeyIsMissingError
	case 403:
		return AuthKeyIsInvalidError
	case 404:
		return NotFoundError
	case 413:
		return CapabilitiesMismatchError
	case 429:
		return APIRateLimitError
	case 500:
		return InternalError
	case 501:
		return NotImplemented
	case 502:
		return GatewayTimeoutError
	default:
		return errors.New("unexpected status code")
	}
}
