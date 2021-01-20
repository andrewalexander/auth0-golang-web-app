package auth

import (
	"context"
	"log"
	"net/url"
	"os"

	"golang.org/x/oauth2"

	oidc "github.com/coreos/go-oidc"
)

type Authenticator struct {
	Provider *oidc.Provider
	Config   oauth2.Config
	Ctx      context.Context
}

func NewAuthenticator() (*Authenticator, error) {
	ctx := context.Background()
	providerUrl, err := url.Parse(os.Getenv("AUTH_PROVIDER_URL"))
	if err != nil {
		log.Printf("failed to parse url from env. %s", err)
		return nil, err
	}
	provider, err := oidc.NewProvider(ctx, providerUrl.String())
	if err != nil {
		log.Printf("failed to get provider: %v", err)
		return nil, err
	}
	conf := oauth2.Config{
		ClientID:     os.Getenv("AUTH_CLIENT_ID"),
		ClientSecret: os.Getenv("AUTH_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("AUTH_REDIRECT_URI"),
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, "profile"},
	}

	return &Authenticator{
		Provider: provider,
		Config:   conf,
		Ctx:      ctx,
	}, nil
}
