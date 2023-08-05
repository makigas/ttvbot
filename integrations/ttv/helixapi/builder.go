package helixapi

import (
	"errors"

	"github.com/nicklaw5/helix/v2"
)

// A Builder allows to build different clients depending on the authentication
// strategy required to login against the application. It returns fully
// initialized helix.Client instances, which can be used directly.
type Builder interface {
	// WithTokenProvider changes the token provider that will get or issue tokens
	// when the client is built. It is important to provide the proper token
	// provider strategy depending on the operation that the client will be used
	// for. Actions on behalf on an user will require an user access token.
	WithTokenProvider(tp TokenProvider)

	// Build crafts the proper client, returning an error if it cannot be built
	// due to an invalid state. If the client cannot be built, it errors.
	Build() (*helix.Client, error)
}

type helixBuilder struct {
	clientId     string
	clientSecret string
	provider     TokenProvider
}

func NewHelixBuilder(clientId, clientSecret string) Builder {
	return &helixBuilder{
		clientId:     clientId,
		clientSecret: clientSecret,
	}
}

func (builder *helixBuilder) WithTokenProvider(tp TokenProvider) {
	builder.provider = tp
	builder.provider.SetAppCredentials(builder.clientId, builder.clientSecret)
}

func (builder *helixBuilder) Build() (*helix.Client, error) {
	// A token provider is required to build the client.
	if builder.provider == nil {
		return nil, errors.New("missing TokenProvider")
	}

	// Build and prepare the helix client.
	client, err := helix.NewClient(&helix.Options{
		ClientID:     builder.clientId,
		ClientSecret: builder.clientSecret,
	})
	if err != nil {
		return nil, err
	}

	// Authenticate the client.
	if err := builder.provider.PlaceToken(client); err != nil {
		return nil, err
	}
	return client, nil
}
