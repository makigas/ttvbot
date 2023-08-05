package helixapi

import (
	"errors"

	"github.com/nicklaw5/helix/v2"
)

// A TokenProvider issues tokens when requested.
type TokenProvider interface {
	// SetAppCredentials changes the internal ClientID and ClientSecret required to
	// communicate against the Helix API in order to fetch a fresh token. This
	// method is meant to be used by the builder when creating a client token so
	// that it can provide the internal clientId and clientSecret.
	SetAppCredentials(clientId, clientSecret string)

	// Token should issue a query against the Helix API in order to get a new
	// session token based on the strategy that is being implemented. In case of
	// error, it errs.
	Token() (string, error)

	// PlaceToken calls the proper method of the helix Client in order to set the
	// session token. This method exists because depending on the kind of token
	// being used, one has to call either the SetAppAccessToken or the
	// SetUserAccessToken method of the client.
	PlaceToken(*helix.Client) error
}

func NewAppAccessTokenProvider(scopes []string) TokenProvider {
	return &appAccessTokenProvider{scopes: scopes}
}

type appAccessTokenProvider struct {
	clientId     string
	clientSecret string
	scopes       []string
}

func (atp *appAccessTokenProvider) SetAppCredentials(clientId, clientSecret string) {
	atp.clientId = clientId
	atp.clientSecret = clientSecret
}

func (atp *appAccessTokenProvider) Token() (string, error) {
	// Spawn a new Helix client.
	client, err := helix.NewClient(&helix.Options{
		ClientID:     atp.clientId,
		ClientSecret: atp.clientSecret,
	})
	if err != nil {
		return "", err
	}

	// Use the endpoint to fetch an application token.
	resp, err := client.RequestAppAccessToken(atp.scopes)
	if err != nil {
		return "", nil
	}
	if token := resp.Data.AccessToken; token != "" {
		return token, nil
	}
	return "", errors.New("empty access token")
}

func (atp *appAccessTokenProvider) PlaceToken(client *helix.Client) error {
	token, err := atp.Token()
	if err != nil {
		return err
	}
	client.SetAppAccessToken(token)
	return nil
}
