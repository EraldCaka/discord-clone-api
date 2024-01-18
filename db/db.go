package db

const (
	DBNAME  = "discord-clone-api"
	MONGODB = "mongodb://localhost:27017"
)

type Store struct {
	User    UserStore
	Server  ServerStore
	Channel ChannelStore
	Message MessageStore
}
