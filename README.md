# tweet_automated_bot

[![Actions Status](https://github.com/forgeutah/tweet_automated_bot/workflows/build/badge.svg)](https://github.com/forgeutah/tweet_automated_bot/actions)
[![codecov](https://codecov.io/gh/forgeutah/tweet_automated_bot/branch/master/graph/badge.svg)](https://codecov.io/gh/forgeutah/tweet_automated_bot)

This bot is created to send out automated tweets regarding the [GoWest Conference](https://gowestconf.com). As well tweet videos from the Forge Utah [youtube channel](https://www.youtube.com/channel/UC7aCz1ur-s48fwm8Zfhjlbg)

## Authentication
WIP: move to 1password

We are using the authentication defined in the [Twitter API Docs](https://developer.twitter.com/en/docs/authentication)

to get the consumer key and consumer secret. store them locally and run
```bash
source .secret
```

## db
![db-relation-diagram.drawio]
[edit here](https://app.diagrams.net/#HSoypete%2Ftweet_automated_bot%2Fsoypete-new-bot%2Fdb-relation-diagram.drawio)

## scheduled tweet flow
![/dev/scheduler_flow.drawio.png]

# dependencies

dependenies are managed via go.mod file. [Dependabot](https://github.com/dependabot) scheduled action looks for updates and vulnerabilites in the dependencies.
- [go-twitter](https://github.com/dghubble/go-twitter)
- [dghubble/oauth1](https://github.com/dghubble/oauth1)
- [discord-go](https://github.com/bwmarrin/discordgo)
