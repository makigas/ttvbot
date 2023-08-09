# ttvbot

This is the main repository for **ttvbot**.

**ttvbot** is a custom Twitch integration. It offers end features like:

- Notifying via Discord whenever the stream starts or ends.
- (To do) Reply to commands via the Twitch chat.
- (To do, pending rewrite) Detect and ban Twitch bots that join the chat.
- (To do) Manage VIP lists, which is a VERY specific reward in my chat.

# Package organization

* cmd/: ttvbot commands.
  * ttvbotd/: the ttvbot daemon; this is the main program.
* domain/: internal domains.
  * commands/: command parser and chatbot command execution.
  * config/: global configuration and application state.
* hooks/: global callbacks that wish to be managed by go-fx.
* httpd/: the HTTP gateway manager.
  * server/: the HTTPD that ttvbotd spawns.
* integrations/: third party integrations
  * discord/: integrations related to Discord.
  * ttv/: integrations related to Twitch.tv
    * eventsub/: EventSub handler.
    * helixapi/: Helix client.
* persistence/: packages related to data storage.
  * dbcommands/: Stores commands to use in the Twitch chatbot.
  * redis/: Redis client, exposed as a go-fx module.

# Commands

The application is designed around the idea that eventually at some point more
than one executable may be needed. For instance, ttvbotctl could be a command
line client that interacts with the ttvbotd daemon, because it will be easier
to issue a command line request like `ttvbotd set-title New title` rather than
frenetically issuing a payload to http://bot.example.com/stream/info.

The current list of commands is:

* ttvbotd (`gopkg.makigas.es/ttvbot/cmd/ttvbotd`): the main ttvbot daemon.

Review the README on each subdirectory of the cmd/ package for information.

# System architecture

I use go-fx to manage the different modules of the application as independent
things. For instance, the Config is exposed as a module.

go-fx is a dependency injection framework. Once the high level components are
exposed as FxModules, they can be injected into other high level components
that also wish to be exposed as FxModule.

As a practical example, the Config object is also an FxModule that can be
injected into other FxModules that require to be configured. The Redis
component depends on the Config and then exposes its own Redis client so that
other components can depend on it (to get or put keys).

# End to end architecture

ttvbot will be controlled by the broadcaster using different transports:

- Either the HTTP server daemon, which will at some point will gain the
  ability to expose the internal daemon state or to modify it via REST API,
  as well as to allow the broadcaster to trigger actions via the HTTP server
  itself.
- Or using specific Twitch chat commands that allow the broadcaster or a
  moderator to trigger the same actions or to fetch or update the internal
  daemon state.

As a practical example, "update stream title" is an action that may be
performed either by issuing a PUT request to /stream/info with the title as
a payload, or via !settitle command in the Twitch chat.

Once the HTTP server gains private functionality, it will be necessary to also
have an authentication and authorization mechanism so that only authorized
users can issue these commands.

# Open Source Policy

This package has been made open source in the hope that it is useful for people
studying the behaviour of this software or the programming language or library
set.

However, this is not an open effort. Therefore, issues and pull requests may be
ignored. This program was designed to fulfill some specific requirements that
may not fit the requirements of other people. If other people is reading this
and considering that the application does not behave as expected, they are free
to write their own integrations.
