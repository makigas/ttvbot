package config

import (
	"io"
	"os"

	"github.com/pelletier/go-toml/v2"
)

type httpdConfig struct {
	ServerBind string `toml:"server_bind"` // The IP address where the HTTP server is bound
	ServerPort uint16 `toml:"server_port"` // The port where the HTTP server is listening
}

type redisConfig struct {
	Host     string `toml:"host"`     // The hostname of the Redis database
	Port     int    `toml:"port"`     // The port of the Redis database
	Password string `toml:"password"` // The password of the Redis database
	Database int    `toml:"db"`       // The database index to use
}

type discordConfig struct {
	StreamingWebhookUrl string `toml:"streaming_webhook_url"` // The discord URL to notify about streams
	StreamingPingRole   string `toml:"streaming_ping_role"`   // The role to ping when a stream starts
}

type helixConfig struct {
	AppClientId     string `toml:"app_client_id"`     // The user-less OAuth2 client ID
	AppClientSecret string `toml:"app_client_secret"` // The user-less OAuth2 secret key
}

// A Config is the static configuration of the applciation.
type Config struct {
	Httpd   httpdConfig   `toml:"httpd"`   // HTTP daemon settings
	Redis   redisConfig   `toml:"redis"`   // Redis connection settings
	Discord discordConfig `toml:"discord"` // Discord client settings
	Helix   helixConfig   `toml:"helix"`   // Helix API client settings
}

// ReadFromToml decodes the given reader as a TOML stream and updates the
// configuration based on what is read. The decoder is set to fail if any
// unknown key is provided.
func (cfg *Config) ReadFromToml(r io.Reader) error {
	decoder := toml.NewDecoder(r)
	decoder.DisallowUnknownFields()
	return decoder.Decode(cfg)
}

// ReadFromFile opens the given file path as a TOML file and loads the
// config from it. If it cannot be done, it will fail with an error.
func (cfg *Config) ReadFromFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	if err := cfg.ReadFromToml(file); err != nil {
		return err
	}
	return nil
}
