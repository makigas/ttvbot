package hooks

import (
	"fmt"

	"github.com/nicklaw5/helix/v2"
	"gopkg.makigas.es/ttvbot/integrations/discord"
	"gopkg.makigas.es/ttvbot/integrations/ttv/eventsub"
	"gopkg.makigas.es/ttvbot/integrations/ttv/helixapi"
)

func StreamNotifications(
	events *eventsub.EventSubManager,
	hapi *helixapi.HelixApi,
	dc *discord.DiscordClient,
) {
	events.AddEventListener("stream.online", func(payload map[string]interface{}) {
		// Create the client.
		client, err := hapi.NewAppClient()
		if err != nil {
			fmt.Printf("Warning: stream.online handler error: %s", err.Error())
			return
		}

		id := payload["broadcaster_user_id"].(string)
		res, err := client.GetStreams(&helix.StreamsParams{
			UserIDs: []string{id},
		})

		if err != nil {
			fmt.Printf("Warning: stream.online handler error: %s", err.Error())
			return
		}
		if len(res.Data.Streams) != 1 {
			fmt.Printf("Warning: stream.online handler error: did not return stream info")
			return
		}

		title := res.Data.Streams[0].Title
		username := res.Data.Streams[0].UserName
		dc.AnnounceLivestreamStart(title, "https://twitch.tv/"+username)
	})
}
