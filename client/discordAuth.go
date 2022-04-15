package client

type discordAuth interface {
	Creds() []interface{}
}

type BearerAuth struct {
	Token string
}

func (b BearerAuth) Creds() []interface{} {
	return []interface{}{"Bearer " + b.Token}
}

type BotAuth struct {
	Token string
}

func (b BotAuth) Creds() []interface{} {
	return []interface{}{"Bot " + b.Token}
}

type UserPassAuth struct {
	User, Pass string
}

func (b UserPassAuth) Creds() []interface{} {
	return []interface{}{b.User, b.Pass}
}

type UserPassTokenAuth struct {
	User, Pass, Token string
}

func (b UserPassTokenAuth) Creds() []interface{} {
	return []interface{}{b.User, b.Pass, b.Token}
}
