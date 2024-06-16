package google

import (
	"encoding/base64"
	"fmt"
	"strings"
)

// decode takes base64 encoded string and returns deconded string
func decode(s string) ([]byte, error) {
	return base64.RawURLEncoding.DecodeString(s)
}

// parseJWT parses JWT token
func parseJWT(idToken string) (*gjwt, error) {
	segments := strings.Split(idToken, ".")
	if len(segments) != 3 {
		return nil, fmt.Errorf("idtoken: invalid token, token must have three segments; found %d", len(segments))
	}
	return &gjwt{
		header:    segments[0],
		payload:   segments[1],
		signature: segments[2],
	}, nil
}

// ParsePayload parses the given token and returns its payload.
//
// Warning: This function does not validate the token prior to parsing it.
//
// ParsePayload is primarily meant to be used to inspect a token's payload. This is
// useful when validation fails and the payload needs to be inspected.
//
// Note: A successful Validate() invocation with the same token will return an
// identical payload.
func ParsePayload(idToken string) (*GoogleIdTokenPayload, error) {
	jwt, err := parseJWT(idToken)
	if err != nil {
		return nil, err
	}
	return jwt.parsedPayload()
}
