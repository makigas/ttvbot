package helixapi

import (
	"errors"

	"github.com/nicklaw5/helix/v2"
	"go.uber.org/fx"
	"gopkg.makigas.es/ttvbot/domain/config"
)

var Module = fx.Module("HelixApi", fx.Provide(NewHelixApi))

var (
	ErrUnexpectedResponse = errors.New("unexpected response from Twitch")
	ErrAuthentication     = errors.New("authentication error")
)

type HelixApi struct {
	client *helix.Client
}

func NewHelixApi(cfg *config.Config) (*HelixApi, error) {
	client, err := helix.NewClient(&helix.Options{
		ClientID:     cfg.Helix.AppClientId,
		ClientSecret: cfg.Helix.AppClientSecret,
	})
	if err != nil {
		return nil, err
	}
	return &HelixApi{client: client}, nil
}

func (ha *HelixApi) refreshApplicationToken() error {
	token, err := ha.client.RequestAppAccessToken([]string{})
	if err != nil {
		return err
	}
	ha.client.SetAppAccessToken(token.Data.AccessToken)
	return nil
}
