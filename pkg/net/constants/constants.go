package constants

const (
	MAX_BODY_SIZE = 5 * 1024 * 1024 //5MB body limit for HTTP request bodies
)

type errorCodes struct {
	AuthInvalidCredentials string
	BadRequest             string
	InternalServerError    string
}

var (
	ErrorCodes errorCodes = errorCodes{
		AuthInvalidCredentials: "AUTH_INVALID_CREDENTIALS",
		BadRequest:             "BAD_REQUEST",
		InternalServerError:    "INTERNAL_SERVER_ERROR",
	}
)
