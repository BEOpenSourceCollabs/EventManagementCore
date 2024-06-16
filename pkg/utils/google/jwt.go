package google

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
)

// gjwt represents the segments of a jwt and exposes convenience methods for
// working with the different segments.
type gjwt struct {
	header    string
	payload   string
	signature string
}

// gJwtHeader represents a parted jwt's header segment.
type gJwtHeader struct {
	Algorithm string `json:"alg"`
	Type      string `json:"typ"`
	KeyID     string `json:"kid"`
}

// decodedHeader base64 decodes the header segment.
func (j *gjwt) decodedHeader() ([]byte, error) {
	dh, err := decode(j.header)
	if err != nil {
		return nil, fmt.Errorf("idtoken: unable to decode JWT header: %v", err)
	}
	return dh, nil
}

// decodedPayload base64 payload the header segment.
func (j *gjwt) decodedPayload() ([]byte, error) {
	p, err := decode(j.payload)
	if err != nil {
		return nil, fmt.Errorf("idtoken: unable to decode JWT payload: %v", err)
	}

	return p, nil
}

// decodedPayload base64 payload the header segment.
func (j *gjwt) decodedSignature() ([]byte, error) {
	p, err := decode(j.signature)
	if err != nil {
		return nil, fmt.Errorf("idtoken: unable to decode JWT signature: %v", err)
	}
	return p, nil
}

// parsedHeader returns a struct representing a JWT header.
func (j *gjwt) parsedHeader() (gJwtHeader, error) {
	var h gJwtHeader
	dh, err := j.decodedHeader()
	if err != nil {
		return h, err
	}
	err = json.Unmarshal(dh, &h)
	if err != nil {
		return h, fmt.Errorf("idtoken: unable to unmarshal JWT header: %v", err)
	}
	return h, nil
}

// parsedPayload returns a struct representing a JWT payload.
func (j *gjwt) parsedPayload() (*GoogleIdTokenPayload, error) {
	var p GoogleIdTokenPayload
	dp, err := j.decodedPayload()
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(dp, &p); err != nil {
		return nil, fmt.Errorf("idtoken: unable to unmarshal JWT payload: %v", err)
	}
	if err := json.Unmarshal(dp, &p.Claims); err != nil {
		return nil, fmt.Errorf("idtoken: unable to unmarshal JWT payload claims: %v", err)
	}
	return &p, nil
}

// hashedContent gets the SHA256 checksum for verification of the JWT.
func (j *gjwt) hashedContent() []byte {
	signedContent := j.header + "." + j.payload
	hashed := sha256.Sum256([]byte(signedContent))
	return hashed[:]
}

func (j *gjwt) String() string {
	return fmt.Sprintf("%s.%s.%s", j.header, j.payload, j.signature)
}
