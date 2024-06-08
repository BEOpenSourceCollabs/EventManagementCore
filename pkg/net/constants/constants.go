package constants

type contextKey string

const (
	MAX_BODY_SIZE               = 5 * 1024 * 1024 //5MB body limit for HTTP request bodies
	USER_CONTEXT_KEY contextKey = "user"
)

type errorCodes struct {
	AuthInvalidCredentials string
	BadRequest             string
	InternalServerError    string
	AuthInvalidAuthHeader  string
	AuthInvalidAuthToken   string
	NotFound               string
}

var (
	ErrorCodes errorCodes = errorCodes{
		AuthInvalidCredentials: "AUTH_INVALID_CREDENTIALS",
		BadRequest:             "BAD_REQUEST",
		NotFound:               "NOT_FOUND",
		InternalServerError:    "INTERNAL_SERVER_ERROR",
		AuthInvalidAuthHeader:  "AUTH_NO_AUTHORIZATION_HEADER",
		AuthInvalidAuthToken:   "AUTH_INVALID_TOKEN",
	}
)
