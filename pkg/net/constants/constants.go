package constants

const (
	MAX_BODY_SIZE = 5 * 1024 * 1024 //5MB body limit for HTTP request bodies
)

type errorCodes struct {
	AuthInvalidCredentials string
	AuthInvalidScope       string
	BadRequest             string
	InternalServerError    string
	AuthInvalidAuthHeader  string
	AuthInvalidAuthToken   string
	NotFound               string
}

var (
	ErrorCodes errorCodes = errorCodes{
		AuthInvalidCredentials: "AUTH_INVALID_CREDENTIALS",
		AuthInvalidScope:       "AUTH_INVALID_SCOPE",
		BadRequest:             "BAD_REQUEST",
		NotFound:               "NOT_FOUND",
		InternalServerError:    "INTERNAL_SERVER_ERROR",
		AuthInvalidAuthHeader:  "AUTH_NO_AUTHORIZATION_HEADER",
		AuthInvalidAuthToken:   "AUTH_INVALID_TOKEN",
	}
)
