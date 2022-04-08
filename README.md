# tweet_automated_bot

This bot is created to send out automated tweets regarding the [GoWest Conference](https://gowestconf.com). 

## development steps

1. ~set up authentication~
1. ~create tweet endpoint~
1. create tweet storage
1. determine method for selecting tweet for storage
1. ~create scheduler for sending out tweets~
1. ~create input mechanism for new tweets~

## Authentication

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
