package constants

type contextKey string

type Role string

const (
	MAX_BODY_SIZE               = 5 * 1024 * 1024 //5MB body limit for HTTP request bodies
	USER_CONTEXT_KEY contextKey = "user"

	UserRole      Role = "user"
	AdminRole     Role = "admin"
	OrganizerRole Role = "organizer"
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
