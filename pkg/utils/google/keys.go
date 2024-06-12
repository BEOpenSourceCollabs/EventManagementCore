package google

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// gcertResponse represents a list jwks. It is the format returned from known
// Google cert endpoints.
type gcertResponse struct {
	Keys []gjwk `json:"keys"`
}

// jwk is a simplified representation of a standard jwk. It only includes the
// fields used by Google's cert endpoints.
type gjwk struct {
	Alg string `json:"alg"`
	Crv string `json:"crv"`
	Kid string `json:"kid"`
	Kty string `json:"kty"`
	Use string `json:"use"`
	E   string `json:"e"`
	N   string `json:"n"`
	X   string `json:"x"`
	Y   string `json:"y"`
}

// findMatchingKey extracts matching key from certificates/keys for given keyID
func findMatchingKey(response *gcertResponse, keyID string) (*gjwk, error) {
	if response == nil {
		return nil, fmt.Errorf("idtoken: cert response is nil")
	}
	for _, v := range response.Keys {
		if v.Kid == keyID {
			return &v, nil
		}
	}
	return nil, fmt.Errorf("idtoken: could not find matching cert keyId for the token provided")
}

// getCert fetches keys/certificates from google servers
func getCert(url string) (*gcertResponse, error) {

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("idtoken: unable to retrieve cert, got status code %d", resp.StatusCode)
	}

	certResp := &gcertResponse{}
	if err := json.NewDecoder(resp.Body).Decode(certResp); err != nil {
		return nil, err

	}

	return certResp, nil
}
