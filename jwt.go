package jwt

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"github.com/go-resty/resty/v2"
	"github.com/luxarts/go-oauth/internal/defines"
	"log"
	"math/big"
	"os"
	"strings"
)

type key struct {
	Kid     string   `json:"kid"`
	Kty     string   `json:"kty"`
	Alg     string   `json:"alg"`
	Use     string   `json:"use"`
	N       string   `json:"n"`
	E       string   `json:"e"`
	X5C     []string `json:"x5c"`
	X5T     string   `json:"x5t"`
	X5TS256 string   `json:"x5t#S256"`
}

type certs struct {
	Keys []key `json:"keys"`
}

type Validator struct {
	jwks certs
}

func NewValidator() *Validator {
	var validator Validator

	jwksURL := os.Getenv(defines.EnvJwksURL)
	if jwksURL == "" {
		log.Fatalln(defines.EnvJwksURL, "is empty")
	}
	rc := resty.New()
	resp, err := rc.R().Get(jwksURL)
	if err != nil {
		log.Fatalln("failed to get jwks", err)
	}

	err = json.Unmarshal(resp.Body(), &validator.jwks)
	if err != nil {
		log.Fatalln("failed to unmarshal jwks", err)
	}

	return &validator
}

func (v *Validator) IsValid(token string) bool {
	parts := strings.Split(token, ".")

	if len(parts) != 3 {
		return false
	}

	headerBytes, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return false
	}

	var header map[string]interface{}
	err = json.Unmarshal(headerBytes, &header)
	if err != nil {
		return false
	}

	kid, ok := header["kid"].(string)
	if !ok {
		return false
	}

	var foundKey key
	for _, k := range v.jwks.Keys {
		if k.Kid == kid {
			foundKey = k
			break
		}
	}

	if foundKey.Kid == "" {
		return false
	}

	modulusBytes, err := base64.RawURLEncoding.DecodeString(foundKey.N)
	if err != nil {
		return false
	}
	exponentBytes, err := base64.RawURLEncoding.DecodeString(foundKey.E)
	if err != nil {
		return false
	}

	exponent := 0
	for _, b := range exponentBytes {
		exponent = exponent<<8 + int(b)
	}

	rsaPublicKey := &rsa.PublicKey{
		N: new(big.Int).SetBytes(modulusBytes),
		E: exponent,
	}

	signatureBytes, err := base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil {
		return false
	}

	hashed := sha256.Sum256([]byte(parts[0] + "." + parts[1]))

	err = rsa.VerifyPKCS1v15(rsaPublicKey, crypto.SHA256, hashed[:], signatureBytes)
	if err != nil {
		return false
	}

	return true
}
