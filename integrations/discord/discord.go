package discord

import (
	"fmt"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/webhook"
	"go.uber.org/fx"
	"gopkg.makigas.es/ttvbot/domain/config"
)

var Module = fx.Module("Discord", fx.Provide(NewDiscordClient))

type DiscordClient struct {
	livestreamWebhookUrl string
	livestreamWebookRole string
}

func NewDiscordClient(config *config.Config) *DiscordClient {
	return &DiscordClient{
		livestreamWebhookUrl: config.Discord.StreamingWebhookUrl,
		livestreamWebookRole: config.Discord.StreamingPingRole,
	}
}

func (dc *DiscordClient) AnnounceLivestreamStart(message, url string) {
	client, err := webhook.NewWithURL(dc.livestreamWebhookUrl)
	if err == nil {
		client.CreateMessage(discord.WebhookMessageCreate{
			Content: fmt.Sprintf("%s <%s> <@&%s>", message, url, dc.livestreamWebookRole),
		})
	}
}

func (dc *DiscordClient) AnnounceLivestreamEnd(message string) {
	client, err := webhook.NewWithURL(dc.livestreamWebhookUrl)
	if err == nil {
		client.CreateMessage(discord.WebhookMessageCreate{
			Content: message,
		})
	}

}
