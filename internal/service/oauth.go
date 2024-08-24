package service

import (
	"crypto/sha256"
	"encoding/base64"
	"go-oauth/internal/defines"
	"go-oauth/internal/domain"
	"go-oauth/internal/repository"
	"math/rand"
	"net/url"
	"os"
	"time"
)

type OAuthService interface {
	GetLoginURL() (*string, error)
	Callback(state string, code string) (*domain.Token, error)
}

type oauthService struct {
	repo      repository.OAuthRepository
	stateRepo repository.StateRepository
}

func NewOAuthService(r repository.OAuthRepository, sr repository.StateRepository) OAuthService {
	return &oauthService{
		repo:      r,
		stateRepo: sr,
	}
}

func (svc *oauthService) GetLoginURL() (*string, error) {
	authURL, err := svc.repo.GetAuthorizeURL()
	if err != nil {
		return nil, err
	}

	u, err := url.Parse(authURL)
	if err != nil {
		return nil, err
	}

	challenge, verifier, err := svc.generatePKCE()
	if err != nil {
		return nil, err
	}

	state := svc.stateRepo.Create(verifier)

	uv := make(url.Values)
	uv.Add("client_id", os.Getenv(defines.EnvClientID))
	uv.Add("redirect_uri", os.Getenv(defines.EnvRedirectURI))
	uv.Add("response_type", "code")
	uv.Add("scope", "openid")
	uv.Add("code_challenge_method", "S256")
	uv.Add("code_challenge", challenge)
	uv.Add("state", state)

	u.RawQuery = uv.Encode()

	uStr := u.String()
	return &uStr, nil
}
func (svc *oauthService) Callback(state string, code string) (*domain.Token, error) {
	codeVerifier, err := svc.stateRepo.Get(state)
	if err != nil {
		return nil, err
	}
	clientID := os.Getenv(defines.EnvClientID)
	clientSecret := os.Getenv(defines.EnvClientSecret)
	redirectURI := os.Getenv(defines.EnvRedirectURI)

	return svc.repo.ExchangeCode(code, codeVerifier, clientID, clientSecret, redirectURI)
}

func (svc *oauthService) generatePKCE() (challenge string, verifier string, err error) {
	cv := make([]byte, 32)
	_, err = rand.New(rand.NewSource(time.Now().UnixNano())).Read(cv)
	if err != nil {
		return "", "", err
	}

	verifier = base64.RawURLEncoding.EncodeToString(cv)

	cvhash := sha256.Sum256([]byte(verifier))
	challenge = base64.RawURLEncoding.EncodeToString(cvhash[:])

	return challenge, verifier, nil
}
