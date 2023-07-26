package hooks

import (
	"gopkg.makigas.es/ttvbot/integrations/discord"
	"gopkg.makigas.es/ttvbot/integrations/ttv/eventsub"
	"gopkg.makigas.es/ttvbot/integrations/ttv/helixapi"
)

func StreamNotifications(
	events *eventsub.EventSubManager,
	helix *helixapi.HelixApi,
	dc *discord.DiscordClient,
) {
	events.AddEventListener("stream.online", func(payload map[string]interface{}) {
		id := payload["broadcaster_user_id"].(string)
		req := helixapi.GetStreamInformationRequest{Id: id}
		if data, err := helix.GetStreamInformation(&req); err == nil {
			dc.AnnounceLivestreamStart(data.Title, "https://twitch.tv/"+data.Username)
		}
	})

	events.AddEventListener("stream.offline", func(payload map[string]interface{}) {
		dc.AnnounceLivestreamEnd("Se acabó el stream, la próxima vez madrugas")
	})
}
