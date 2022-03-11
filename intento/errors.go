package intento

import "errors"

type ProviderRelatedError struct{}

func (e *ProviderRelatedError) Error() string {
	return "provider-related error"
}

type AuthKeyIsMissingError struct{}

func (e *AuthKeyIsMissingError) Error() string {
	return "intento: auth key is missing"
}

type AuthKeyIsInvalidError struct{}

func (e *AuthKeyIsInvalidError) Error() string {
	return "intento: auth key is invalid"
}

type NotFoundError struct{}

func (e *NotFoundError) Error() string {
	return "intento: intent/provider not found"
}

type CapabilitiesMismatchError struct{}

func (e *CapabilitiesMismatchError) Error() string {
	return "intento: capabilities mismatch for the chosen provider"
}

type APIRateLimitError struct{}

func (e *APIRateLimitError) Error() string {
	return "intento: API rate limit exceeded"
}

type InternalError struct{}

func (e *InternalError) Error() string {
	return "intento: internal error"
}

type NotImplemented struct{}

func (e *NotImplemented) Error() string {
	return "intento: not implemented"
}

type GatewayTimeoutError struct{}

func (e *GatewayTimeoutError) Error() string {
	return "intento: gateway timeout errors"
}

func httpStatusCodeToError(statusCode int) error {
	if statusCode >= 200 && statusCode <= 299 {
		return nil
	}

	switch statusCode {
	case 400:
		return &ProviderRelatedError{}
	case 401:
		return &AuthKeyIsMissingError{}
	case 403:
		return &AuthKeyIsInvalidError{}
	case 404:
		return &NotFoundError{}
	case 413:
		return &CapabilitiesMismatchError{}
	case 429:
		return &APIRateLimitError{}
	case 500:
		return &InternalError{}
	case 501:
		return &NotImplemented{}
	case 502:
		return &GatewayTimeoutError{}
	default:
		return errors.New("unexpected status code")
	}
}
