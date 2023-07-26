/*
Ttvbotd is the daemon that interacts with the Twitch.tv API. At the moment,
it is used to handle livestream events via Twitch's EventSub API, but in the
future the IRC bot will also be handled by this program.

# Synopsis

	ttvbotd -config <config-file>

Parameters:

	-config <config-file>
	  The path to the config file to use. Checkout the Configruation section
	  for more information on how to provide this parameter.

# What can the daemon bot?

The current use case for ttvbotd is creating an HTTP server that can listen
to EventSub notifications coming from Twitch. The daemon will only actively
interact with stream.online and stream.offline events, and once a notification
is received, it will forward the notification to a Discord webhook so that
other people can get the news.

# Configuration

So many configuration keys are needed to operate ttvbotd that it was found
that using a config file was more flexible. The bot needs to be given a
config file with all the settings. The config file is a TOML file with the
following sections:

- httpd: settings for the HTTP server
  - server_bind: the hostname where to bind the HTTP server. For running in
    development mode, '127.0.0.1' is recommended. To serve the bot in
    production, '0.0.0.0' is suggested unless there is a reverse proxy in
    front of the application that could forward traffic.
  - server_port: the port where the HTTP server should listen.

- redis: settings for the Redis persistence database.
  - host: the hostname where to connect.
  - port: the port where to connect.
  - password: if given, a password to use during authentication.
  - db: if given, a database index. Otherwise, 0.

- discord: settings for the Discord integration.
  - streaming_webhook_url: the webhook URL to announce streams as they start.
  - streaming_ping_role: the role to ping about the stream notifications.

- helix: settings for the Twitch API integration.
  - app_client_id: the OAuth application client ID.
  - app_client_secret: the OAuth application client secret.
*/
package main
