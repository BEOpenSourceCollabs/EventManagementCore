package google

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rsa"
	"fmt"
	"math/big"
	"time"
)

const (
	es256KeySize      int    = 32
	googleIAPCertsURL string = "https://www.gstatic.com/iap/verify/public_key-jwk"
	googleSACertsURL  string = "https://www.googleapis.com/oauth2/v3/certs"
)

// Validator provides a way to validate Google ID Tokens
type Validator struct {
}

func NewValidator() *Validator {
	return &Validator{}
}

// Validate is used to validate the provided idToken with a known Google cert
// URL. If audience is not empty the audience claim of the Token is validated.
// Upon successful validation a parsed token Payload is returned allowing the
// caller to valdate any additional claims.i
func (v *Validator) ValidateToken(idToken string, audience string) (*GoogleIdTokenPayload, error) {
	return v.validate(idToken, audience)
}

func (v *Validator) validate(idToken, audience string) (*GoogleIdTokenPayload, error) {

	jwt, err := parseJWT(idToken)
	if err != nil {
		return nil, err
	}
	header, err := jwt.parsedHeader()
	if err != nil {
		return nil, err
	}
	payload, err := jwt.parsedPayload()
	if err != nil {
		return nil, err
	}
	sig, err := jwt.decodedSignature()
	if err != nil {
		return nil, err
	}

	if audience != "" && payload.Audience != audience {
		return nil, fmt.Errorf("idtoken: audience provided does not match aud claim in the JWT")
	}

	now := time.Now()

	if now.Unix() > payload.Expires {
		return nil, fmt.Errorf("idtoken: token expired: now=%v, expires=%v", now.Unix(), payload.Expires)
	}

	switch header.Algorithm {
	case "RS256":
		if err := v.validateRS256(header.KeyID, jwt.hashedContent(), sig); err != nil {
			return nil, err
		}
	case "ES256":
		if err := v.validateES256(header.KeyID, jwt.hashedContent(), sig); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("idtoken: expected JWT signed with RS256 or ES256 but found %q", header.Algorithm)
	}

	return payload, nil
}

func (v *Validator) validateRS256(keyID string, hashedContent []byte, sig []byte) error {

	certResp, err := getCert(googleSACertsURL)
	if err != nil {
		return err
	}
	j, err := findMatchingKey(certResp, keyID)
	if err != nil {
		return err
	}
	dn, err := decode(j.N)
	if err != nil {
		return err
	}
	de, err := decode(j.E)
	if err != nil {
		return err
	}

	pk := &rsa.PublicKey{
		N: new(big.Int).SetBytes(dn),
		E: int(new(big.Int).SetBytes(de).Int64()),
	}
	return rsa.VerifyPKCS1v15(pk, crypto.SHA256, hashedContent, sig)
}

func (v *Validator) validateES256(keyID string, hashedContent []byte, sig []byte) error {
	certResp, err := getCert(googleIAPCertsURL)
	if err != nil {
		return err
	}
	j, err := findMatchingKey(certResp, keyID)
	if err != nil {
		return err
	}
	dx, err := decode(j.X)
	if err != nil {
		return err
	}
	dy, err := decode(j.Y)
	if err != nil {
		return err
	}

	pk := &ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X:     new(big.Int).SetBytes(dx),
		Y:     new(big.Int).SetBytes(dy),
	}
	r := big.NewInt(0).SetBytes(sig[:es256KeySize])
	s := big.NewInt(0).SetBytes(sig[es256KeySize:])
	if valid := ecdsa.Verify(pk, hashedContent, r, s); !valid {
		return fmt.Errorf("idtoken: ES256 signature not valid")
	}
	return nil
}
