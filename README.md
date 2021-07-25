## bassface

`Bassface` is a general purpose irc bot, written in GoLang, based on [`whyrusleeping/hellabot`](https://github.com/whyrusleeping/hellabot), for use in the [Jungletrain.net chat](https://jungletrain.net/chat). Interested in Jungle/Drum and Bass music? Come say hi!

A predecessor to `Bassface` [used to fart](https://www.reddit.com/r/devops/comments/orctqs/what_do_you_do_with_python/h6hl3wa/). But now it does ... other things.

## capabilities

At the time of this writing, `bassface` can:

* if given permission, kick/ban users for using words that match a static list, and report the action to a list of users (responds to channel messages with specific content)
* hello world (`!bassface`)
* randomly respond to other bots in the channel (responds to channel messages with specific content)
* send random ascii boobs based on word matches from other users (responds to channel messages with specific content)
* query the [discogs.com](https://www.discogs.com) database in various ways (!discogs)
* register itself to `nickserv` (responds to specific direct messages from specific users)
* post links to listen to the [jungletrain.net](https://www.discogs.com) radio stream (`!streams` or `!pls`)
* respond with some text when some other text matches (`whagwarn?`)

## usage

The `Makefile` describes everything you need to do here - including which secrets to mount in as environment variables to run. 

* Compile bassface, build container, tag `latest`, push to `mars64/bassface:latest`
```
make all
```

* delete local compiled binaries, clean docker images
```
make clean
```

When the container is staged, use the helm deployment in `mars64.io/linode/helm/bassface` to deploy.

Use the helm templates to deploy to a Kubernetes cluster. All secrets are handled as envvars -- set these before you run the commands above for great fun and happiness.

Once the bot is connected to the channel of choice, you can use the commands listed in the `hbot.Trigger` sections. See [`whyrysleeping/hellabot`](https://github.com/whyrusleeping/hellabot) for more info on the framework.