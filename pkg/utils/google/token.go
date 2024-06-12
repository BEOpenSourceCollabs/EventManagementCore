package google

type GoogleIdTokenPayload struct {
	Issuer   string                 `json:"iss"`
	Audience string                 `json:"aud"`
	Expires  int64                  `json:"exp"`
	IssuedAt int64                  `json:"iat"`
	Subject  string                 `json:"sub,omitempty"`
	Claims   map[string]interface{} `json:"-"`
}

type GoogleUserClaims struct {
	Email         string
	EmailVerified bool
	Id            string
	Name          string
	Picture       string
	FirstName     string
	LastName      string
}

func (p *GoogleIdTokenPayload) GetClaims() *GoogleUserClaims {

	claims := &GoogleUserClaims{}

	if v, ok := p.Claims["email"].(string); ok {
		claims.Email = v
	}

	if v, ok := p.Claims["email_verified"].(bool); ok {
		claims.EmailVerified = v
	}

	if v, ok := p.Claims["sub"].(string); ok {
		claims.Id = v
	}

	if v, ok := p.Claims["name"].(string); ok {
		claims.Name = v
	}

	if v, ok := p.Claims["given_name"].(string); ok {
		claims.FirstName = v
	}

	if v, ok := p.Claims["family_name"].(string); ok {
		claims.LastName = v
	}

	if v, ok := p.Claims["picture"].(string); ok {
		claims.LastName = v
	}

	return claims
}
