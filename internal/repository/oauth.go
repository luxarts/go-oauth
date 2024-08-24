package repository

import (
	"encoding/json"
	"errors"
	"github.com/go-resty/resty/v2"
	"github.com/luxarts/go-oauth/internal/defines"
	"github.com/luxarts/go-oauth/internal/domain"
	"log"
	"net/url"
	"os"
)

type OAuthRepository interface {
	GetAuthorizeURL() (string, error)
	ExchangeCode(code string, codeVerifier string, clientID string, clientSecret string, redirectURI string) (*domain.Token, error)
}

type oauthRepository struct {
	rc     *resty.Client
	oidCfg domain.OpenIDConfig
}

func NewOAuthRepository(rc *resty.Client) OAuthRepository {
	configURL := os.Getenv(defines.EnvOpenIDConfigURL)
	resp, err := rc.R().Get(configURL)
	if err != nil {
		log.Fatalln("failed to get openid config", err)
	}
	var oidCfg domain.OpenIDConfig
	err = json.Unmarshal(resp.Body(), &oidCfg)
	if err != nil {
		log.Fatalln("failed to unmarshal openid config", err)
	}

	return &oauthRepository{
		oidCfg: oidCfg,
		rc:     rc,
	}
}

func (repo *oauthRepository) GetAuthorizeURL() (string, error) {
	if repo.oidCfg.AuthorizationEndpoint == "" {
		return "", errors.New("no authorize url")
	}

	return repo.oidCfg.AuthorizationEndpoint, nil
}
func (repo *oauthRepository) ExchangeCode(code string, codeVerifier string, clientID string, clientSecret string, redirectURI string) (*domain.Token, error) {
	uv := make(url.Values)
	uv.Add("client_id", clientID)
	uv.Add("client_secret", clientSecret)
	uv.Add("redirect_uri", redirectURI)
	uv.Add("code", code)
	uv.Add("grant_type", "authorization_code")
	uv.Add("code_verifier", codeVerifier)

	resp, err := repo.rc.R().
		SetFormDataFromValues(uv).
		Post(repo.oidCfg.TokenEndpoint)
	if err != nil {
		return nil, err
	}

	var token domain.Token
	err = json.Unmarshal(resp.Body(), &token)
	if err != nil {
		return nil, err
	}

	return &token, nil
}
