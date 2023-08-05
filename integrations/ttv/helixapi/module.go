package helixapi

import (
	"github.com/nicklaw5/helix/v2"
	"go.uber.org/fx"
	"gopkg.makigas.es/ttvbot/domain/config"
)

type HelixApi struct {
	clientId     string
	clientSecret string
}

func (hapi *HelixApi) NewAppClient() (*helix.Client, error) {
	builder := NewHelixBuilder(hapi.clientId, hapi.clientSecret)
	builder.WithTokenProvider(NewAppAccessTokenProvider([]string{}))
	return builder.Build()
}

type HelixApiParams struct {
	fx.In
	AppConf *config.Config
}

type HelixApiResult struct {
	fx.Out
	Builder *HelixApi
}

func NewHelixApi(params HelixApiParams) HelixApiResult {
	clientId := params.AppConf.Helix.AppClientId
	clientSecret := params.AppConf.Helix.AppClientSecret
	return HelixApiResult{
		Builder: &HelixApi{clientId: clientId, clientSecret: clientSecret},
	}
}

var Module = fx.Module("helixapi", fx.Provide(NewHelixApi))
